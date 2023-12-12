package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

type GradeSection struct {
	name             string
	classes          map[string]*SchoolClass
	gradeSubsections []*GradeSection
	gpa              float64
	totalCredits     int64
}

type SchoolClass struct {
	credits       int64
	grade         float64
	totalWeight   float64
	gradeParts    map[string]*GradePart
	name          string
	explicitGrade string
	desiredGrade  float64
}

type GradePart struct {
	weight         float64
	pointsRecieved float64
	pointsTotal    float64
}

func checkErr(errLog *log.Logger, err error) {
	if err != nil {
		errLog.Println(err)
		os.Exit(1)
	}
}

func printError(errLog *log.Logger, errMsg string) {
	errLog.Println("error " + errMsg)
	os.Exit(1)
}

func printLineError(errLog *log.Logger, fileName string, lineIndex int, errMsg string) {
	printError(errLog, fmt.Sprintf("[%s:%d]: %s\n", fileName, lineIndex+1, errMsg))
}

func parseOptionLine(errLog *log.Logger, fileName string, line string, lineIndex int) (string, string) {
	fields := strings.Split(line, "=")

	if len(fields) != 2 {
		printError(errLog, fmt.Sprintf("[%s:%d]: recieved a line that does not follow the x = y format", fileName, lineIndex+1))

		os.Exit(1)
	}

	commentPrefixIndex := strings.Index(fields[1], "#")

	if commentPrefixIndex != -1 {
		fields[1] = fields[1][:commentPrefixIndex]
	}

	return strings.TrimSpace(fields[0]), strings.TrimSpace(fields[1])
}

func fuzzyFindFile(dir string, search string) string {
	files, err := os.ReadDir(dir)

	if err != nil {
		return ""
	}

	for _, file := range files {
		nextPath := filepath.Join(dir, file.Name())

		if file.IsDir() {
			innerDirSearchTest := fuzzyFindFile(nextPath, search)

			if innerDirSearchTest != "" {
				return innerDirSearchTest
			}
		}

		if strings.Contains(file.Name(), search) {
			fmt.Println("found file: " + nextPath)
			return nextPath
		}
	}

	return ""
}

// https://stackoverflow.com/a/48801414
func trimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}
