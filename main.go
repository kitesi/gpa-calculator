package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
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
	gradeParts := map[string]*GradePart{}
	scanner := bufio.NewScanner(strings.NewReader(fileContent))
	inMetaOptions := false
	var credits int64 = 4
	var current_grade_part_name string

	fileName = filepath.Base(fileName)

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "~ Meta") {
			if inMetaOptions {
				errLog.Printf("error [%s]: recieved more than one meta headers", fileName)
				os.Exit(1)
			}

			inMetaOptions = true
			continue
		}

		if inMetaOptions && !strings.HasPrefix(line, ">") {
			field_name, field_value := parseOptionLine(errLog, fileName, line)

			if field_name == "credit" {
				c, err := strconv.ParseInt(field_value, 10, 64)
				checkErr(errLog, err)
				credits = c
			}

			continue
		}

		if strings.HasPrefix(line, ">") {
			inMetaOptions = false
			current_grade_part_name = strings.TrimSpace(trimFirstRune(line))
			gradeParts[current_grade_part_name] = &GradePart{}
		} else if current_grade_part_name == "" {
			errLog.Printf("error [%s]: recieved a line that is not under a grade part\n\t$ %s\n", fileName, line)
			os.Exit(1)
		} else {
			field_name, field_value := parseOptionLine(errLog, fileName, line)

			if field_name == "weight" {
				field_value_float, err := strconv.ParseFloat(field_value, 32)

				if err != nil {
					errLog.Printf("error [%s]: the value for weight did not compile to a float\n\t$ %s\n", fileName, field_value)
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
						errLog.Printf("error [%s]: one of the scores did not follow the x/y format\n\t$%s", fileName, score)
						os.Exit(1)
					}

					numerator, err := strconv.ParseFloat(score_fractions[0], 32)

					if err != nil {
						errLog.Printf("error [%s]: the numerator in one of the scores did not compile to a float\n\t$ %s", fileName, score)
						os.Exit(1)
					}

					denominator, err := strconv.ParseFloat(score_fractions[1], 32)

					if err != nil {
						errLog.Printf("error [%s]: the denominator in one of the scores did not compile to a float\n\t$ %s", fileName, score)
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
				errLog.Printf("error [%s]: recieved an invalid field name: %s", fileName, field_name)
				os.Exit(1)
			}
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

	return &SchoolClass{grade: grade, gradeParts: gradeParts, credits: credits}
}

func printGrades(gs *GradeSection, prefix string, verbose bool) {
	if len(gs.gradeSubsections) > 0 {
		for i, gSubsection := range gs.gradeSubsections {
			if i == len(gs.gradeSubsections)-1 {
				fmt.Printf("%s└── %s (%.2f)\n", prefix, gSubsection.name, gSubsection.gpa)
			} else if i == 0 {
				fmt.Printf("%s├── %s (%.2f)\n", prefix, gSubsection.name, gSubsection.gpa)
			} else {
				fmt.Printf("%s│  %s (%.2f)\n", prefix, gSubsection.name, gSubsection.gpa)
			}

			if i == len(gs.gradeSubsections)-1 {
				printGrades(gSubsection, prefix+"    ", verbose)
			} else {
				printGrades(gSubsection, prefix+"│   ", verbose)
			}
		}
	}

	if len(gs.classes) > 0 {
		i := 0

		for sClassName, sClass := range gs.classes {
			connecter := "├──"

			if i == len(gs.classes)-1 {
				connecter = "└──"
			}

			if sClass.grade == -1 {
				fmt.Printf("%s %s (unset)\n", prefix+connecter, sClassName)
			} else {
				fmt.Printf("%s %s (%.2f) [%s]\n", prefix+connecter, sClassName, sClass.grade*100, getGradeLetter(sClass.grade))
			}

			if verbose {
				j := 0

				for gradePartName, gradePart := range sClass.gradeParts {
					subconnector := "├──"
					additionalPrefix := "│"

					if j == len(sClass.gradeParts)-1 {
						subconnector = "└──"
					}

					if i == len(gs.classes)-1 {
						additionalPrefix = " "
					}

					if gradePart.pointsTotal == 0 {
						fmt.Printf("%s    %s %s (unset)\n", prefix+additionalPrefix, subconnector, gradePartName)
					} else {
						fmt.Printf("%s    %s %s (%.2f) [%s]\n", prefix+additionalPrefix, subconnector, gradePartName, (gradePart.pointsRecieved/gradePart.pointsTotal)*100, getGradeLetter(gradePart.pointsRecieved/gradePart.pointsTotal))
					}

					j += 1
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

			totalCreditsAdded += getGradeGPA(sClass.grade) * float64(sClass.credits)
			gs.gpa += getGradeGPA(sClass.grade) * float64(sClass.credits) / float64(gs.totalCredits)
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
			fmt.Println("Usage: gpa <grades_directory> [-h|--help] [-v|--verbose]\ngrades_directory: a path to examine the grade(s), it can be a file or directory")
			os.Exit(0)
		} else if arg == "-v" || arg == "--verbose" {
			verbose = true
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
		fmt.Printf("%s (%.2f)\n", d.name, d.gpa)
		printGrades(d, "", verbose)

	} else {
		f := handleFile(errLog, fileName)
		d := &GradeSection{name: "", classes: map[string]*SchoolClass{filepath.Base(fileName): f}}

		printGrades(d, "", verbose)
	}
}
