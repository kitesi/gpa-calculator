package main

import (
	"bytes"
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

	originalWriter := errLog.Writer()
	temporaryBuffer := new(bytes.Buffer)
	errLog.SetOutput(temporaryBuffer)

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
				continue
			}

			inMetaOptions = true
			continue
		}

		if inMetaOptions && !strings.HasPrefix(line, ">") {
			field_name, field_value, err := parseOptionLine(fileName, line)

			if err != nil {
				printLineError(errLog, fileName, lineIndex, err.Error())
				continue
			}

			if field_name == "ignore" {
				if field_value == "" || field_value == "true" || field_value == "1" {
					status = 2
				} else {
					printLineError(errLog, fileName, lineIndex, fmt.Sprintf("the value for ignore can only be 'true', '1', or nothing: '%s'", field_value))
				}
				continue
			}

			// past this point requires field value
			if field_value == "" {
				printLineError(errLog, fileName, lineIndex, fmt.Sprintf("recieved a field \"%s\" with no value (x = y)", field_name))
				continue
			}

			if field_name == "credits" {
				c, err := strconv.ParseInt(field_value, 10, 64)

				if err != nil {
					printLineError(errLog, fileName, lineIndex, fmt.Sprintf("the value for credits did not compile to an int: '%s'", field_value))
					continue
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
					continue
				}
			}

			if field_name == "desired_grade" {
				desiredGrade, err = strconv.ParseFloat(field_value, 64)

				if err != nil {
					printLineError(errLog, fileName, lineIndex, fmt.Sprintf("the value for desired_grade did not compile to a float: '%s'", field_value))
					continue
				}

				if desiredGrade < 0 || desiredGrade > 100 {
					printLineError(errLog, fileName, lineIndex, fmt.Sprintf("the value for desired_grade is not between 0 and 100: '%s'", field_value))
					continue
				}

				desiredGrade /= 100
			}

			continue
		}

		if strings.HasPrefix(line, ">") {
			inMetaOptions = false
			gradePartName := strings.TrimSpace(trimFirstRune(strings.TrimSpace(line)))

			if gradePartName == "" {
				printLineError(errLog, fileName, lineIndex, "recieved a grade part with no name")
				continue
			}

			for _, gradePart := range gradeParts {
				if gradePart.name == gradePartName {
					printLineError(errLog, fileName, lineIndex, fmt.Sprintf("recieved a duplicate grade part name: '%s'", gradePartName))
					continue
				}
			}

			currentGradePartIndex += 1
			gradeParts = append(gradeParts, &GradePart{name: gradePartName})
		} else if currentGradePartIndex == -1 {
			printLineError(errLog, fileName, lineIndex, "recieved a line that is not under a grade part")
			continue
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

				if strings.HasPrefix(nextLine, ">") || strings.HasPrefix(nextLine, "~") || strings.Contains(nextLine, "=") || strings.Contains(nextLine, "drop_lowest") {
					break
				}

				option_string += " " + strings.TrimSpace(nextLine)
				nextLineIndex += 1
			}

			field_name, field_value, err := parseOptionLine(fileName, option_string)

			if err != nil {
				printLineError(errLog, fileName, lineIndex, err.Error())
				continue
			}

			var drop_lowest int64 = 0

			if field_name == "drop_lowest" {
				if field_value != "" {
					drop_lowest_int, err := strconv.ParseInt(field_value, 10, 64)

					if err != nil {
						printLineError(errLog, fileName, lineIndex, fmt.Sprintf("the value for drop_lowest did not compile to an int: '%s'", field_value))
						continue
					}

					drop_lowest = drop_lowest_int
				} else {
					drop_lowest = 1
				}
				gradeParts[currentGradePartIndex].dropLowest = drop_lowest
				continue
			}

			// similarily, past this point requires field value
			if field_value == "" {
				printLineError(errLog, fileName, lineIndex, fmt.Sprintf("recieved a field \"%s\" with no value (x = y)", field_name))
				continue
			}

			if field_name == "weight" {
				field_value_float, err := strconv.ParseFloat(field_value, 32)

				if err != nil {
					printLineError(errLog, fileName, lineIndex, fmt.Sprintf("the value for weight did not compile to a float: '%s'", field_value))
					continue
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
						continue
					}

					numerator, err := strconv.ParseFloat(score_fractions[0], 32)

					if err != nil {
						printLineError(errLog, fileName, lineIndex, fmt.Sprintf("the numerator in one of the scores did not compile to a float: '%s'", score))
						continue
					}

					denominator, err := strconv.ParseFloat(score_fractions[1], 32)

					if err != nil {
						printLineError(errLog, fileName, lineIndex, fmt.Sprintf("the denominator in one of the scores did not compile to a float: '%s'", score))
						continue
					}

					// todo: validate numerator/denominator

					gradeParts[currentGradePartIndex].pointsRecieved += numerator
					gradeParts[currentGradePartIndex].pointsTotal += denominator
					gradeParts[currentGradePartIndex].scores = append(gradeParts[currentGradePartIndex].scores, &Score{pointsRecieved: numerator, pointsTotal: denominator})

				}

			} else {
				printLineError(errLog, fileName, lineIndex, fmt.Sprintf("recieved an invalid field name: '%s'", field_name))
				continue
			}

			lineIndex = nextLineIndex - 1
		}
	}

	totalWeight := 0.0
	totalGrades := 0.0

	for _, gradePart := range gradeParts {
		amountScores := len(gradePart.scores)
		if amountScores == 0 {
			continue
		}

		totalR := gradePart.pointsRecieved
		totalT := gradePart.pointsTotal

		// sort the scores from lowest to highest
		sort.Slice(gradePart.scores, func(i, j int) bool {
			scoreI := gradePart.scores[i]
			scoreJ := gradePart.scores[j]

			return (totalR-scoreI.pointsRecieved)/(totalT-scoreI.pointsTotal) > (totalR-scoreJ.pointsRecieved)/(totalT-scoreJ.pointsTotal)
		})

		dropAmount := int(gradePart.dropLowest)

		if dropAmount >= amountScores {
			dropAmount = amountScores - 1
		}

		for i := 0; i < dropAmount; i++ {
			score := gradePart.scores[i]
			gradePart.dropped = append(gradePart.dropped, score)
			gradePart.pointsRecieved -= score.pointsRecieved
			gradePart.pointsTotal -= score.pointsTotal
		}

		totalWeight += gradePart.weight
		totalGrades += (gradePart.pointsRecieved / gradePart.pointsTotal) * gradePart.weight
	}

	var grade float64 = -1

	if totalGrades != 0 && totalWeight != 0 {
		grade = totalGrades / totalWeight
	}

	errLog.SetOutput(originalWriter)

	if temporaryBuffer.String() != "" {
		errLog.Print(temporaryBuffer.String())
		return &SchoolClass{}, 1
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

			if i == len(gs.gradeSubsections)-1 && len(gs.classes) == 0 {
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

						p := gradePart.pointsRecieved / gradePart.pointsTotal
						percentLossed := (1 - p) * gradePart.weight

						if percentLossed > 0.0001 {
							fmt.Printf("%s    %s %s (%.2f) (%s) (-%.2f%%)\n", prefix+additionalPrefix, subconnector, gradePart.name, (p)*100, getGradeLetter(p), percentLossed*100)
						} else {
							fmt.Printf("%s    %s %s (%.2f) (%s)\n", prefix+additionalPrefix, subconnector, gradePart.name, (p)*100, getGradeLetter(p))
						}
					}

					// todo: test all of this drop stuff
					if gradePart.dropLowest > 0 {
						droppedStr := ""
						anotherAdditionalPrefix := "│"

						if j == len(sClass.gradeParts)-1 && sClass.desiredGrade == -1 {
							anotherAdditionalPrefix = " "
						}

						for k, score := range gradePart.dropped {
							droppedStr += fmt.Sprintf("%.2f/%.2f ", score.pointsRecieved, score.pointsTotal)

							if k != len(gradePart.dropped)-1 {
								droppedStr += ", "
							}
						}
						fmt.Printf("%s    %s   └── dropped (%d) %s\n", prefix+additionalPrefix, anotherAdditionalPrefix, len(gradePart.dropped), droppedStr)
					}

					if sClass.desiredGrade != -1 && strings.HasPrefix(strings.ToLower(gradePart.name), "final") {
						finalGradePart = gradePart
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

func calculateGPA(gs *GradeSection, unweighted bool) (float64, float64) {
	totalCreditsAdded := 0.0

	if len(gs.gradeSubsections) != 0 {
		for _, gSubsection := range gs.gradeSubsections {
			childGpa, childTotalCreditsAdded := calculateGPA(gSubsection, unweighted)
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

			if unweighted && correspondingGPA > 4.0 {
				correspondingGPA = 4.0
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
	unweighted := false
	edit := false
	posArgs := []string{}

	for _, arg := range args {
		if arg == "-h" || arg == "--help" {
			fmt.Println("Usage: gpa [file] [-e|--edit] [-h|--help] [-v|--verbose] [--unweighted|-u] [--version]\nfile: a path to examine, it can be a file or directory")
			return 0
		} else if arg == "-v" || arg == "--verbose" {
			verbose = true
		} else if arg == "--version" {
			fmt.Println("gpa-calculator version 1.1.0")
			return 0
		} else if arg == "-e" || arg == "--edit" {
			edit = true
		} else if arg == "-u" || arg == "--unweighted" {
			unweighted = true
		} else if strings.HasPrefix(arg, "-") {
			printError(errLog, fmt.Sprintf("unknown flag: '%s'", arg))
			return 1
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

		calculateGPA(d, unweighted)

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

		fmt.Println(fileName)
		d := &GradeSection{name: "", classes: []*SchoolClass{f}}
		printGrades(errLog, d, "", verbose)
	}

	return 0
}

func main() {
	os.Exit(run(os.Args[1:]))
}
