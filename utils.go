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
	classes          []*SchoolClass
	gradeSubsections []*GradeSection
	gpa              float64
	credits          int64
}

type SchoolClass struct {
	credits       int64
	grade         float64
	totalWeight   float64
	gradeParts    []*GradePart
	name          string
	explicitGrade string
	desiredGrade  float64
}

type GradePart struct {
	weight         float64
	name           string
	dropLowest     int64
	pointsRecieved float64
	pointsTotal    float64
	scores         []*Score
	dropped        []*Score
}

type Score struct {
	pointsRecieved float64
	pointsTotal    float64
}

func printError(errLog *log.Logger, errMsg string) {
	errLog.Println("error " + errMsg)
}

func printLineError(errLog *log.Logger, fileName string, lineIndex int, errMsg string) {
	printError(errLog, fmt.Sprintf("[%s:%d]: %s", fileName, lineIndex+1, errMsg))
}

func parseOptionLine(fileName string, line string) (string, string, error) {
	commentPrefixIndex := strings.Index(line, "#")
	if commentPrefixIndex != -1 {
		line = line[:commentPrefixIndex]
	}

	field_key, field_value, found := strings.Cut(line, "=")
	field_key = strings.Trim(strings.TrimSpace(field_key), "\"")

	if found {
		field_value = strings.Trim(strings.TrimSpace(field_value), "\"")
	}

	return field_key, field_value, nil
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
