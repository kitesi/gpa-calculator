package main

import (
	"log"
	"os"
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
	credits    int64
	grade      float64
	gradeParts map[string]*GradePart
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

func printLineError(errLog *log.Logger, fileName string, lineIndex int, errMsg string) {
	errLog.Printf("error [%s:%d]: %s\n", fileName, lineIndex+1, errMsg)
	os.Exit(1)
}

func parseOptionLine(errLog *log.Logger, fileName string, line string) (string, string) {
	fields := strings.Split(line, "=")

	if len(fields) != 2 {
		errLog.Printf("error [%s]: recieved a line that does not follow the x = y format\n\t$ %s\n", fileName, line)
		os.Exit(1)
	}

	return strings.TrimSpace(fields[0]), strings.TrimSpace(fields[1])
}

// https://stackoverflow.com/a/48801414
func trimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}
