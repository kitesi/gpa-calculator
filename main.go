package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func handleDirectory(errLog *log.Logger, dirName string, gradeSection GradeSection) *GradeSection {
	dir, err := os.ReadDir(dirName)
	checkErr(errLog, err)

	for _, file := range dir {
		nextPath := filepath.Join(dirName, file.Name())
		if file.IsDir() {
			child := handleDirectory(errLog, nextPath, GradeSection{name: file.Name(), classes: make(map[string]*SchoolClass)})
			gradeSection.gradeSubsections = append(gradeSection.gradeSubsections, child)
			gradeSection.totalCredits += child.totalCredits
		} else {
			schoolClass := handleFile(errLog, nextPath)
			gradeSection.classes[file.Name()] = schoolClass

			// don't add to total credits if no grade is determined yet
			if schoolClass.grade != -1 {
				gradeSection.totalCredits += gradeSection.classes[file.Name()].credits
			}
		}
	}

	return &gradeSection
}

func handleFile(errLog *log.Logger, fileName string) *SchoolClass {
	fileContentBuffer, err := os.ReadFile(fileName)
	checkErr(errLog, err)

	fileContent := string(fileContentBuffer)
	lines := regexp.MustCompile(`\r?\n`).Split(fileContent, -1)
	gradeParts := map[string]*GradePart{}

	inMetaOptions := false
	userExplicitGrade := ""
	current_grade_part_name := ""
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
				os.Exit(1)
			}

			inMetaOptions = true
			continue
		}

		if inMetaOptions && !strings.HasPrefix(line, ">") {
			field_name, field_value := parseOptionLine(errLog, fileName, line, lineIndex)

			if field_name == "credits" {
				c, err := strconv.ParseInt(field_value, 10, 64)

				if err != nil {
					printLineError(errLog, fileName, lineIndex, fmt.Sprintf("the value for credits did not compile to an int: %s", field_value))
				}

				credits = c
			}

			if field_name == "name" {
				className = field_value
			}

			if field_name == "grade" {
				userExplicitGrade = field_value

				if getGradeGPA(userExplicitGrade) == -1 {
					printLineError(errLog, fileName, lineIndex, fmt.Sprintf("recieved an invalid grade: %s", userExplicitGrade))
					os.Exit(1)
				}
			}

			if field_name == "desired_grade" {
				desiredGrade, err = strconv.ParseFloat(field_value, 64)

				if err != nil {
					printLineError(errLog, fileName, lineIndex, fmt.Sprintf("the value for desired_grade did not compile to a float: %s", field_value))
				}

				if desiredGrade < 0 || desiredGrade > 100 {
					printLineError(errLog, fileName, lineIndex, fmt.Sprintf("the value for desired_grade is not between 0 and 100: %s", field_value))
				}

				desiredGrade /= 100
			}

			continue
		}

		if strings.HasPrefix(line, ">") {
			inMetaOptions = false
			current_grade_part_name = strings.TrimSpace(trimFirstRune(strings.TrimSpace(line)))

			if _, ok := gradeParts[current_grade_part_name]; ok {
				printLineError(errLog, fileName, lineIndex, fmt.Sprintf("recieved a duplicate grade part name: %s", current_grade_part_name))
				os.Exit(1)
			}

			gradeParts[current_grade_part_name] = &GradePart{}
		} else if current_grade_part_name == "" {
			printLineError(errLog, fileName, lineIndex, "recieved a line that is not under a grade part")
			os.Exit(1)
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

			field_name, field_value := parseOptionLine(errLog, fileName, option_string, lineIndex)

			if field_name == "weight" {
				field_value_float, err := strconv.ParseFloat(field_value, 32)

				if err != nil {
					printLineError(errLog, fileName, lineIndex, fmt.Sprintf("the value for weight did not compile to a float: %s", field_value))
					errLog.Println(err)
					os.Exit(1)
				}

				if entry, ok := gradeParts[current_grade_part_name]; ok {
					entry.weight = field_value_float
					gradeParts[current_grade_part_name] = entry
				} else {
					errLog.Println("was denied entry to grade_parts map ?")
				}
			} else if field_name == "data" {
				for _, score := range strings.Split(strings.TrimSpace(field_value), ",") {
					score_fractions := strings.Split(strings.TrimSpace(score), "/")

					if len(score_fractions) != 2 {
						printLineError(errLog, fileName, lineIndex, fmt.Sprintf("one of the scores did not follow the x/y format: %s", score))
						os.Exit(1)
					}

					numerator, err := strconv.ParseFloat(score_fractions[0], 32)

					if err != nil {
						printLineError(errLog, fileName, lineIndex, fmt.Sprintf("the numerator in one of the scores did not compile to a float: %s", score))
						os.Exit(1)
					}

					denominator, err := strconv.ParseFloat(score_fractions[1], 32)

					if err != nil {
						printLineError(errLog, fileName, lineIndex, fmt.Sprintf("the denominator in one of the scores did not compile to a float: %s", score))
						os.Exit(1)
					}

					if entry, ok := gradeParts[current_grade_part_name]; ok {
						entry.pointsRecieved += (numerator)
						entry.pointsTotal += (denominator)

						gradeParts[current_grade_part_name] = entry
					} else {
						errLog.Println("was denied entry to grade_parts map ?")
					}
				}

			} else {
				printLineError(errLog, fileName, lineIndex, fmt.Sprintf("recieved an invalid field name: %s", field_name))
				os.Exit(1)
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

	return &SchoolClass{grade: grade, gradeParts: gradeParts, credits: credits, name: className, explicitGrade: userExplicitGrade, desiredGrade: desiredGrade}
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

				for gradePartName, gradePart := range sClass.gradeParts {
					subconnector := "├──"

					if j == len(sClass.gradeParts)-1 && sClass.desiredGrade == -1 {
						subconnector = "└──"
					}

					if gradePart.pointsTotal == 0 {
						fmt.Printf("%s    %s %s (unset)\n", prefix+additionalPrefix, subconnector, gradePartName)
					} else {
						fmt.Printf("%s    %s %s (%.2f) (%s)\n", prefix+additionalPrefix, subconnector, gradePartName, (gradePart.pointsRecieved/gradePart.pointsTotal)*100, getGradeLetter(gradePart.pointsRecieved/gradePart.pointsTotal))
					}

					j += 1
				}

				if sClass.desiredGrade != -1 {
					var finalGradePart *GradePart

					for gradePartName := range sClass.gradeParts {
						if strings.HasPrefix(strings.ToLower(gradePartName), "final") {
							finalGradePart = sClass.gradeParts[gradePartName]
							break
						}
					}

					if finalGradePart == nil {
						printError(errLog, fmt.Sprintf("[%s]: could not find a grade part that starts with 'final'", sClass.name))
					}

					fmt.Printf("%s    └── to get a %.2f%% you need at least a %.2f%% on the final\n", prefix+additionalPrefix, sClass.desiredGrade*100, (sClass.desiredGrade-sClass.grade*(1-finalGradePart.weight))*100/finalGradePart.weight)
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
			gs.gpa += childTotalCreditsAdded / float64(gs.totalCredits)
			totalCreditsAdded += childTotalCreditsAdded
		}
	}

	if len(gs.classes) > 0 && gs.totalCredits > 0 {
		for _, sClass := range gs.classes {
			if sClass.grade == -1 {
				continue
			}

			correspondingGPA := getGradeGPA(getGradeLetter(sClass.grade))

			if sClass.explicitGrade != "" {
				correspondingGPA = getGradeGPA(sClass.explicitGrade)
			}

			totalCreditsAdded += correspondingGPA * float64(sClass.credits)
			gs.gpa += correspondingGPA * float64(sClass.credits) / float64(gs.totalCredits)
		}
	}

	return gs.gpa, totalCreditsAdded
}

func main() {
	errLog := log.New(os.Stderr, "", 0)

	verbose := false
	posArgs := []string{}

	for _, arg := range os.Args {
		if arg == "-h" || arg == "--help" {
			fmt.Println("Usage: gpa <grades_directory> [-h|--help] [-v|--verbose] [--version]\ngrades_directory: a path to examine the grade(s), it can be a file or directory")
			os.Exit(0)
		} else if arg == "-v" || arg == "--verbose" {
			verbose = true
		} else if arg == "--version" {
			fmt.Println("gpa-calculator version 1.0.0")
			os.Exit(0)
		} else {
			posArgs = append(posArgs, arg)
		}
	}

	if len(posArgs) != 2 {
		errLog.Printf("error: expected 1 positional argument, recieved %d\n", len(posArgs)-1)
		os.Exit(1)
	}

	fileName := posArgs[1]
	fileInfo, err := os.Stat(fileName)
	checkErr(errLog, err)

	if fileInfo.IsDir() {
		d := handleDirectory(errLog, fileName, GradeSection{name: filepath.Base(fileName), classes: make(map[string]*SchoolClass)})

		calculateGPA(d)
		fmt.Printf("%s (%.2f)\n", fileName, d.gpa)
		printGrades(errLog, d, "", verbose)

	} else {
		f := handleFile(errLog, fileName)
		d := &GradeSection{name: "", classes: map[string]*SchoolClass{filepath.Base(fileName): f}}

		printGrades(errLog, d, "", verbose)
	}
}
