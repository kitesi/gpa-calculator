package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// handle a directory recursively and return the GradeSection and a status, 0=good, 1=error
func handleDirectory(errLog *log.Logger, dirName string, gradeSection GradeSection) (*GradeSection, int) {
	dir, err := os.ReadDir(dirName)

	if err != nil {
		printError(errLog, fmt.Sprintf("could not read directory '%s'", dirName))
		return &GradeSection{}, 1
	}

	for _, file := range dir {
		nextPath := filepath.Join(dirName, file.Name())
		if file.IsDir() {
			child, status := handleDirectory(errLog, nextPath, GradeSection{name: file.Name()})

			// if there's an error reading the directory, we still want to compute the gpa
			if status == 1 {
				continue
			}

			gradeSection.gradeSubsections = append(gradeSection.gradeSubsections, child)
			gradeSection.credits += child.credits
		} else {
			schoolClass, status := handleFile(errLog, nextPath)

			// even if a file has an error, we still want to compute the gpa
			if status == 2 || status == 1 {
				continue
			}

			gradeSection.classes = append(gradeSection.classes, schoolClass)

			// don't add to total credits if no grade is determined yet
			if schoolClass.grade != -1 {
				gradeSection.credits += schoolClass.credits
			}
		}
	}

	return &gradeSection, 0
}

// returns the SchoolClass and a status, 0=good, 1=error, 2=specified file ignore
func handleFile(errLog *log.Logger, fileName string) (*SchoolClass, int) {
	fileContentBuffer, err := os.ReadFile(fileName)

	if err != nil {
		printError(errLog, fmt.Sprintf("could not read file '%s'", fileName))
		return &SchoolClass{}, 1
	}

	fileContent := string(fileContentBuffer)
	lines := regexp.MustCompile(`\r?\n`).Split(fileContent, -1)
	gradeParts := []*GradePart{}

	inMetaOptions := false
	status := 0
	userExplicitGrade := ""
	currentGradePartIndex := -1
	desiredGrade := -1.0

	var credits int64 = 4

	fileName = filepath.Base(fileName)
	className := fileName

	lineIndex := -1

	for lineIndex+1 < len(lines) {
		lineIndex += 1
		line := lines[lineIndex]

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "~ Meta") {
			if inMetaOptions {
				printLineError(errLog, fileName, lineIndex, "recieved more than one meta headers")
				return &SchoolClass{}, 1
			}

			inMetaOptions = true
			continue
		}

		if inMetaOptions && !strings.HasPrefix(line, ">") {
			field_name, field_value, err := parseOptionLine(fileName, line)

			if err != nil {
				printLineError(errLog, fileName, lineIndex, err.Error())
				return &SchoolClass{}, 1
			}

			if field_name == "credits" {
				c, err := strconv.ParseInt(field_value, 10, 64)

				if err != nil {
					printLineError(errLog, fileName, lineIndex, fmt.Sprintf("the value for credits did not compile to an int: '%s'", field_value))
					return &SchoolClass{}, 1
				}

				credits = c
			}

			if field_name == "name" {
				className = field_value
			}

			if field_name == "grade" {
				userExplicitGrade = field_value

				if getGradeGPA(userExplicitGrade) == -1 {
					printLineError(errLog, fileName, lineIndex, fmt.Sprintf("recieved an invalid grade: '%s'", userExplicitGrade))
					return &SchoolClass{}, 1
				}
			}

			if field_name == "desired_grade" {
				desiredGrade, err = strconv.ParseFloat(field_value, 64)

				if err != nil {
					printLineError(errLog, fileName, lineIndex, fmt.Sprintf("the value for desired_grade did not compile to a float: '%s'", field_value))
					return &SchoolClass{}, 1
				}

				if desiredGrade < 0 || desiredGrade > 100 {
					printLineError(errLog, fileName, lineIndex, fmt.Sprintf("the value for desired_grade is not between 0 and 100: '%s'", field_value))
					return &SchoolClass{}, 1
				}

				desiredGrade /= 100
			}

			if field_name == "ignore" {
				if field_value == "true" {
					status = 2
				} else {
					printLineError(errLog, fileName, lineIndex, fmt.Sprintf("the value for ignore can only be 'true': '%s'", field_value))
					return &SchoolClass{}, 1
				}
			}

			continue
		}

		if strings.HasPrefix(line, ">") {
			inMetaOptions = false
			gradePartName := strings.TrimSpace(trimFirstRune(strings.TrimSpace(line)))

			if gradePartName == "" {
				printLineError(errLog, fileName, lineIndex, "recieved a grade part with no name")
				return &SchoolClass{}, 1
			}

			for _, gradePart := range gradeParts {
				if gradePart.name == gradePartName {
					printLineError(errLog, fileName, lineIndex, fmt.Sprintf("recieved a duplicate grade part name: '%s'", gradePartName))
					return &SchoolClass{}, 1
				}
			}

			currentGradePartIndex += 1
			gradeParts = append(gradeParts, &GradePart{name: gradePartName})
		} else if currentGradePartIndex == -1 {
			printLineError(errLog, fileName, lineIndex, "recieved a line that is not under a grade part")
			return &SchoolClass{}, 1
		} else {
			option_string := strings.TrimSpace(line)
			nextLineIndex := lineIndex + 1

			for {
				if nextLineIndex >= len(lines) {
					break
				}

				nextLine := lines[nextLineIndex]
				commentPrefixIndex := strings.Index(nextLine, "#")

				if commentPrefixIndex != -1 {
					nextLine = nextLine[:commentPrefixIndex]
				}

				if strings.HasPrefix(nextLine, ">") || strings.HasPrefix(nextLine, "~") || strings.Contains(nextLine, "=") {
					break
				}

				option_string += " " + strings.TrimSpace(nextLine)
				nextLineIndex += 1
			}

			field_name, field_value, err := parseOptionLine(fileName, option_string)

			if err != nil {
				printLineError(errLog, fileName, lineIndex, err.Error())
				return &SchoolClass{}, 1
			}

			if field_name == "weight" {
				field_value_float, err := strconv.ParseFloat(field_value, 32)

				if err != nil {
					printLineError(errLog, fileName, lineIndex, fmt.Sprintf("the value for weight did not compile to a float: '%s'", field_value))
					return &SchoolClass{}, 1
				}

				gradeParts[currentGradePartIndex].weight = field_value_float
			} else if field_name == "data" {
				for _, score := range strings.Split(strings.TrimSpace(field_value), ",") {

					// trailling commas are allowed
					if strings.TrimSpace(score) == "" {
						continue
					}

					score_fractions := strings.Split(strings.TrimSpace(score), "/")

					if len(score_fractions) != 2 {
						printLineError(errLog, fileName, lineIndex, fmt.Sprintf("one of the scores did not follow the x/y format: '%s'", score))
						return &SchoolClass{}, 1
					}

					numerator, err := strconv.ParseFloat(score_fractions[0], 32)

					if err != nil {
						printLineError(errLog, fileName, lineIndex, fmt.Sprintf("the numerator in one of the scores did not compile to a float: '%s'", score))
						return &SchoolClass{}, 1
					}

					denominator, err := strconv.ParseFloat(score_fractions[1], 32)

					if err != nil {
						printLineError(errLog, fileName, lineIndex, fmt.Sprintf("the denominator in one of the scores did not compile to a float: '%s'", score))
						return &SchoolClass{}, 1
					}

					gradeParts[currentGradePartIndex].pointsRecieved += numerator
					gradeParts[currentGradePartIndex].pointsTotal += denominator
				}

			} else {
				printLineError(errLog, fileName, lineIndex, fmt.Sprintf("recieved an invalid field name: '%s'", field_name))
				return &SchoolClass{}, 1
			}

			lineIndex = nextLineIndex - 1
		}
	}

	totalWeight := 0.0
	totalGrades := 0.0

	for _, gradePart := range gradeParts {
		if gradePart.pointsTotal == 0 {
			continue
		}

		totalWeight += gradePart.weight
		totalGrades += (gradePart.pointsRecieved / gradePart.pointsTotal) * gradePart.weight
	}

	var grade float64 = -1

	if totalGrades != 0 && totalWeight != 0 {
		grade = totalGrades / totalWeight
	}

	return &SchoolClass{grade: grade, totalWeight: totalWeight, gradeParts: gradeParts, credits: credits, name: className, explicitGrade: userExplicitGrade, desiredGrade: desiredGrade}, status
}

func printGrades(errLog *log.Logger, gs *GradeSection, prefix string, verbose bool) {
	if len(gs.gradeSubsections) > 0 {
		for i, gSubsection := range gs.gradeSubsections {
			if i == len(gs.gradeSubsections)-1 && len(gs.classes) == 0 {
				fmt.Printf("%s└── %s (%.2f)\n", prefix, gSubsection.name, gSubsection.gpa)
			} else {
				fmt.Printf("%s├── %s (%.2f)\n", prefix, gSubsection.name, gSubsection.gpa)
			}

			if i == len(gs.gradeSubsections)-1 {
				printGrades(errLog, gSubsection, prefix+"    ", verbose)
			} else {
				printGrades(errLog, gSubsection, prefix+"│   ", verbose)
			}
		}
	}

	if len(gs.classes) > 0 {
		i := 0

		sort.Slice(gs.classes, func(i, j int) bool {
			return gs.classes[i].name < gs.classes[j].name
		})

		for _, sClass := range gs.classes {
			connecter := "├──"

			if i == len(gs.classes)-1 {
				connecter = "└──"
			}

			if sClass.grade == -1 {
				fmt.Printf("%s %s (unset)\n", prefix+connecter, sClass.name)
			} else {
				gradeLetter := sClass.explicitGrade

				if gradeLetter == "" {
					gradeLetter = getGradeLetter(sClass.grade)
				}

				fmt.Printf("%s %s (%.2f) (%s)\n", prefix+connecter, sClass.name, sClass.grade*100, gradeLetter)
			}

			if verbose {
				j := 0
				additionalPrefix := "│"

				if i == len(gs.classes)-1 {
					additionalPrefix = " "
				}

				var finalGradePart *GradePart

				for _, gradePart := range sClass.gradeParts {
					subconnector := "├──"

					if j == len(sClass.gradeParts)-1 && sClass.desiredGrade == -1 {
						subconnector = "└──"
					}

					if gradePart.pointsTotal == 0 {
						fmt.Printf("%s    %s %s (unset)\n", prefix+additionalPrefix, subconnector, gradePart.name)
					} else {
						fmt.Printf("%s    %s %s (%.2f) (%s)\n", prefix+additionalPrefix, subconnector, gradePart.name, (gradePart.pointsRecieved/gradePart.pointsTotal)*100, getGradeLetter(gradePart.pointsRecieved/gradePart.pointsTotal))
					}

					if sClass.desiredGrade != -1 && strings.HasPrefix(strings.ToLower(gradePart.name), "final") {
						finalGradePart = gradePart
						break
					}

					j += 1
				}

				if sClass.desiredGrade != -1 {
					if finalGradePart == nil {
						printError(errLog, fmt.Sprintf("[%s]: could not find a grade part that starts with 'final'", sClass.name))
						continue
					}

					// if the final grade is already set, we want to remove it from the calculation
					gradeWithoutFinal := sClass.grade

					if finalGradePart.pointsTotal != 0 {
						gradeWithoutFinal = (gradeWithoutFinal*sClass.totalWeight - (finalGradePart.pointsRecieved/finalGradePart.pointsTotal)*finalGradePart.weight) / (sClass.totalWeight - finalGradePart.weight)
					}

					// if no grade is set without the final, then to get the desired grade you need to get that desired grade on the final
					if gradeWithoutFinal == -1 {
						gradeWithoutFinal = sClass.desiredGrade
					}

					fmt.Printf("%s    └── to get a %.2f%% you need at least a %.2f%% on the final\n", prefix+additionalPrefix, sClass.desiredGrade*100, (sClass.desiredGrade-gradeWithoutFinal*(1-finalGradePart.weight))*100/finalGradePart.weight)
				}

			}

			i += 1
		}
	}
}

func calculateGPA(gs *GradeSection) (float64, float64) {
	totalCreditsAdded := 0.0

	if len(gs.gradeSubsections) != 0 {
		for _, gSubsection := range gs.gradeSubsections {
			childGpa, childTotalCreditsAdded := calculateGPA(gSubsection)
			gSubsection.gpa = childGpa
			gs.gpa += childTotalCreditsAdded / float64(gs.credits)
			totalCreditsAdded += childTotalCreditsAdded
		}
	}

	if len(gs.classes) > 0 && gs.credits > 0 {
		for _, sClass := range gs.classes {
			if sClass.grade == -1 {
				continue
			}

			correspondingGPA := getGradeGPA(getGradeLetter(sClass.grade))

			if sClass.explicitGrade != "" {
				correspondingGPA = getGradeGPA(sClass.explicitGrade)
			}

			totalCreditsAdded += correspondingGPA * float64(sClass.credits)
			gs.gpa += correspondingGPA * float64(sClass.credits) / float64(gs.credits)
		}
	}

	return gs.gpa, totalCreditsAdded
}

func run(args []string) int {
	errLog := log.New(os.Stderr, "", 0)

	verbose := false
	edit := false
	posArgs := []string{}

	for _, arg := range args {
		if arg == "-h" || arg == "--help" {
			fmt.Println("Usage: gpa [file] [-e|--edit] [-h|--help] [-v|--verbose] [--version]\nfile: a path to examine, it can be a file or directory")
			return 0
		} else if arg == "-v" || arg == "--verbose" {
			verbose = true
		} else if arg == "--version" {
			fmt.Println("gpa-calculator version 1.0.0")
			return 0
		} else if arg == "-e" || arg == "--edit" {
			edit = true
		} else {
			posArgs = append(posArgs, arg)
		}
	}

	if len(posArgs) > 1 {
		printError(errLog, fmt.Sprintf("expected 0-1 positional argument, recieved %d", len(posArgs)))
		return 1
	}

	fileName := ""

	if len(posArgs) == 0 {
		if os.Getenv("GRADES_DIR") != "" {
			fileName = os.Getenv("GRADES_DIR")
		} else {
			printError(errLog, "did not recieve a positional argument, and GRADES_DIR is not set")
			return 1
		}
	} else {
		fileName = posArgs[0]
	}

	fileInfo, err := os.Stat(fileName)

	if errors.Is(err, os.ErrNotExist) {
		gradesDir := os.Getenv("GRADES_DIR")

		if gradesDir == "" {
			printError(errLog, fmt.Sprintf("could not find file or directory '%s', no fuzzy find search occured as $GRADES_DIR is not set", fileName))
			return 1
		}

		fuzzyResult := fuzzyFindFile(gradesDir, fileName)

		if fuzzyResult == "" {
			printError(errLog, fmt.Sprintf("could not find file or directory '%s', even with fuzzy find", fileName))
			return 1
		} else {
			fileName = fuzzyResult
		}

		fileInfo, err = os.Stat(fileName)
	}

	// if not os.ErrNotExist or the fuzzy find worked but there's an error accessing it, then just generic print
	if err != nil {
		printError(errLog, err.Error())
		return 1
	}

	if fileInfo.IsDir() {
		if edit {
			printError(errLog, "editing directories is not supported")
			return 1
		}

		d, status := handleDirectory(errLog, fileName, GradeSection{name: filepath.Base(fileName)})

		if status == 1 {
			return 1
		}

		calculateGPA(d)

		fmt.Printf("%s (%.2f)\n", fileName, d.gpa)
		printGrades(errLog, d, "", verbose)
	} else {
		if edit {
			editor := os.Getenv("EDITOR")

			if editor == "" {
				printError(errLog, "$EDITOR is not set")
				return 1
			}

			cmd := exec.Command(editor, fileName)

			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				printError(errLog, fmt.Sprintf("could not open file %s in $EDITOR", fileName))
				return 1
			}
		}

		f, status := handleFile(errLog, fileName)

		if status == 1 {
			return 1
		} else if status == 2 {
			// ignore the ignore flag because if a user is specifying that file directly, they probably want to see it
		}

		d := &GradeSection{name: "", classes: []*SchoolClass{f}}
		printGrades(errLog, d, "", verbose)
	}

	return 0
}

func main() {
	os.Exit(run(os.Args[1:]))
}
