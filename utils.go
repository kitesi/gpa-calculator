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
	pointsRecieved float64
	pointsTotal    float64
	name           string
}

func printError(errLog *log.Logger, errMsg string) {
	errLog.Println("error " + errMsg)
}

func printLineError(errLog *log.Logger, fileName string, lineIndex int, errMsg string) {
	printError(errLog, fmt.Sprintf("[%s:%d]: %s\n", fileName, lineIndex+1, errMsg))
}

func parseOptionLine(fileName string, line string) (string, string, error) {
	fields := strings.Split(line, "=")

	if len(fields) != 2 {
		return "", "", fmt.Errorf("recieved a line that does not follow the x = y format")
	}

	commentPrefixIndex := strings.Index(fields[1], "#")

	if commentPrefixIndex != -1 {
		fields[1] = fields[1][:commentPrefixIndex]
	}

	s1 := strings.Trim(strings.TrimSpace(fields[0]), "\"")
	s2 := strings.Trim(strings.TrimSpace(fields[1]), "\"")

	return s1, s2, nil
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
