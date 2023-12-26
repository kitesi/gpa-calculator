package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"testing"

	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MainTestSuite struct {
	suite.Suite
	directory     string
	inner_log_buf *bytes.Buffer
	inner_log     *log.Logger
}

// test handleDirectory() and calculateGPA()
func (suite *MainTestSuite) TestHandleDirectory() {
	t := suite.T()
	d, status := handleDirectory(suite.inner_log, suite.directory, GradeSection{name: "test_files"})
	assert.Equal(t, 0, status)
	calculateGPA(d)

	deep.CompareUnexportedFields = true

	if diff := deep.Equal(*d, expected_gs_from_test_directory); diff != nil {
		t.Error(diff)
	}

	output := suite.inner_log_buf.String()
	assert.Equal(t, "", output, "output should be empty")

	suite.inner_log_buf.Reset()
}

func (suite *MainTestSuite) TestHandleFile() {
	t := suite.T()

	t.Run("test file with duplicate meta header", func(tx *testing.T) {
		_, status := handleFile(suite.inner_log, "test_files/duplicate_meta.grade")
		output := suite.inner_log_buf.String()
		assert.Equal(tx, "error [duplicate_meta.grade:4]: recieved more than one meta headers\n", output)
		assert.Equal(tx, 1, status)
	})

	suite.inner_log_buf.Reset()

	t.Run("test file with meta.credits", func(tx *testing.T) {
		d, status := handleFile(suite.inner_log, "test_files/credits.grade")
		output := suite.inner_log_buf.String()

		assert.Equal(tx, "credits.grade", d.name)
		assert.Equal(tx, int64(2), d.credits)
		assert.Equal(tx, "", output)
		assert.Equal(tx, 0, status)
	})

	suite.inner_log_buf.Reset()

	t.Run("test file with invalid credits", func(tx *testing.T) {
		_, status := handleFile(suite.inner_log, "test_files/invalid_credits.grade")
		output := suite.inner_log_buf.String()
		assert.Equal(tx, "error [invalid_credits.grade:4]: the value for credits did not compile to an int: 'a23'\n", output)
		assert.Equal(tx, 1, status)
	})

	suite.inner_log_buf.Reset()

	t.Run("test file with meta.name", func(tx *testing.T) {
		d, status := handleFile(suite.inner_log, "test_files/name.grade")
		output := suite.inner_log_buf.String()

		assert.Equal(tx, "CSC 110", d.name)
		assert.Equal(tx, "", output)
		assert.Equal(tx, 0, status)
	})

	suite.inner_log_buf.Reset()

	t.Run("test file with set grade", func(tx *testing.T) {
		d, status := handleFile(suite.inner_log, "test_files/set_grade.grade")
		output := suite.inner_log_buf.String()

		assert.Equal(tx, "set_grade.grade", d.name)
		assert.Equal(tx, "A", d.explicitGrade)
		assert.Equal(tx, "", output)
		assert.Equal(tx, 0, status)
	})

	suite.inner_log_buf.Reset()

	t.Run("test file with invalid set grade", func(tx *testing.T) {
		_, status := handleFile(suite.inner_log, "test_files/set_invalid_grade.grade")
		output := suite.inner_log_buf.String()
		assert.Equal(tx, "error [set_invalid_grade.grade:3]: recieved an invalid grade: 'Z'\n", output)
		assert.Equal(tx, 1, status)
	})

	suite.inner_log_buf.Reset()

	// TODO: test file with desired grade

	t.Run("test file with invalid desired grade", func(tx *testing.T) {
		_, status := handleFile(suite.inner_log, "test_files/invalid_desired_grade.grade")
		output := suite.inner_log_buf.String()
		assert.Equal(tx, "error [invalid_desired_grade.grade:3]: the value for desired_grade did not compile to a float: 'a9'\n", output)
		assert.Equal(tx, 1, status)
	})

	suite.inner_log_buf.Reset()

	t.Run("test file with invalid desired grade (>100 or <0)", func(tx *testing.T) {
		_, status := handleFile(suite.inner_log, "test_files/invalid_desired_grade2.grade")
		output := suite.inner_log_buf.String()
		assert.Equal(tx, "error [invalid_desired_grade2.grade:3]: the value for desired_grade is not between 0 and 100: '200'\n", output)
		assert.Equal(tx, 1, status)
	})

	suite.inner_log_buf.Reset()

	t.Run("test file with ignore", func(tx *testing.T) {
		d, status := handleFile(suite.inner_log, "test_files/ignore.grade")
		output := suite.inner_log_buf.String()

		if diff := deep.Equal(*d, expected_ignore_class); diff != nil {
			t.Error(diff)
		}

		assert.Equal(tx, "", output)
		assert.Equal(tx, 2, status)
	})

	suite.inner_log_buf.Reset()

	t.Run("test file with invalid ignore", func(tx *testing.T) {
		_, status := handleFile(suite.inner_log, "test_files/invalid_ignore.grade")
		output := suite.inner_log_buf.String()
		assert.Equal(tx, "error [invalid_ignore.grade:5]: the value for ignore can only be 'true': 'false'\n", output)
		assert.Equal(tx, 1, status)
	})

	suite.inner_log_buf.Reset()

	t.Run("test file with grade part that has no name", func(tx *testing.T) {
		_, status := handleFile(suite.inner_log, "test_files/empty_grade_part.grade")
		output := suite.inner_log_buf.String()
		assert.Equal(tx, "error [empty_grade_part.grade:5]: recieved a grade part with no name\nerror [empty_grade_part.grade:6]: recieved a line that is not under a grade part\nerror [empty_grade_part.grade:7]: recieved a line that is not under a grade part\n", output)
		assert.Equal(tx, 1, status)
	})

	suite.inner_log_buf.Reset()

	t.Run("test file with duplicate grade part", func(tx *testing.T) {
		_, status := handleFile(suite.inner_log, "test_files/duplicate_grade_part.grade")
		output := suite.inner_log_buf.String()
		assert.Equal(tx, "error [duplicate_grade_part.grade:9]: recieved a duplicate grade part name: 'Homework'\n", output)
		assert.Equal(tx, 1, status)
	})

	suite.inner_log_buf.Reset()

	t.Run("test file with an option that has no parent grade part", func(tx *testing.T) {
		_, status := handleFile(suite.inner_log, "test_files/unknown_parent.grade")
		output := suite.inner_log_buf.String()
		assert.Equal(tx, "error [unknown_parent.grade:1]: recieved a line that is not under a grade part\nerror [unknown_parent.grade:7]: recieved a duplicate grade part name: 'Homework'\n", output)
		assert.Equal(tx, 1, status)
	})

	suite.inner_log_buf.Reset()

	// TODO: test that whitespace does not matter

	t.Run("test file with comments", func(tx *testing.T) {
		d, status := handleFile(suite.inner_log, "test_files/comments.grade")
		output := suite.inner_log_buf.String()

		assert.Equal(tx, "comments.grade", d.name)
		assert.Equal(tx, int64(5), d.credits)
		assert.Equal(tx, 3, len(d.gradeParts))
		assert.Equal(tx, "", output)
		assert.Equal(tx, 0, status)
	})

	suite.inner_log_buf.Reset()

	t.Run("test file with multiline options", func(tx *testing.T) {
		d, status := handleFile(log.New(os.Stderr, "", 0), "test_files/multiline.grade")
		output := suite.inner_log_buf.String()

		assert.Equal(tx, "multiline.grade", d.name)
		assert.Equal(tx, "", output)
		assert.Equal(tx, 0, status)

		d2, status2 := handleFile(log.New(os.Stderr, "", 0), "test_files/non_multiline.grade")
		output2 := suite.inner_log_buf.String()

		assert.Equal(tx, "non_multiline.grade", d2.name)
		assert.Equal(tx, "", output2)
		assert.Equal(tx, 0, status2)

		// change name just so we can compare
		d2.name = d.name

		if diff := deep.Equal(d, d2); diff != nil {
			t.Error(diff)
		}
	})

	suite.inner_log_buf.Reset()

	t.Run("test file with meta option that does not follow x = y", func(tx *testing.T) {
		_, status := handleFile(suite.inner_log, "test_files/meta_line_non_assignment.grade")
		output := suite.inner_log_buf.String()
		assert.Equal(tx, "error [meta_line_non_assignment.grade:4]: recieved a line that does not follow the x = y format\n", output)
		assert.Equal(tx, 1, status)
	})

	suite.inner_log_buf.Reset()

	t.Run("test file with grade part option that does not follow x = y", func(tx *testing.T) {
		_, status := handleFile(suite.inner_log, "test_files/grade_part_line_non_assignment.grade")
		output := suite.inner_log_buf.String()
		assert.Equal(tx, "error [grade_part_line_non_assignment.grade:6]: recieved a line that does not follow the x = y format\n", output)
		assert.Equal(tx, 1, status)
	})

	suite.inner_log_buf.Reset()

	t.Run("test file with invalid weight option", func(tx *testing.T) {
		_, status := handleFile(suite.inner_log, "test_files/invalid_weight.grade")
		output := suite.inner_log_buf.String()
		assert.Equal(tx, "error [invalid_weight.grade:4]: the value for weight did not compile to a float: '0.2a'\n", output)
		assert.Equal(tx, 1, status)
	})

	suite.inner_log_buf.Reset()

	t.Run("test file with trailing comma in data option", func(tx *testing.T) {
		d, status := handleFile(suite.inner_log, "test_files/trailing_comma.grade")
		output := suite.inner_log_buf.String()
		assert.Equal(tx, "", output)
		assert.Equal(tx, 0, status)

		// trailing_comma.grade is derived from name.grade
		d2, _ := handleFile(suite.inner_log, "test_files/name.grade")
		d.name = d2.name

		if diff := deep.Equal(d, d2); diff != nil {
			t.Error(diff)
		}
	})

	suite.inner_log_buf.Reset()

	t.Run("test file with invalid data scores", func(tx *testing.T) {
		_, status := handleFile(suite.inner_log, "test_files/invalid_data_scores.grade")
		output := suite.inner_log_buf.String()
		assert.Equal(tx, "error [invalid_data_scores.grade:5]: one of the scores did not follow the x/y format: '20'\n", output)
		assert.Equal(tx, 1, status)
	})

	suite.inner_log_buf.Reset()

	t.Run("test file with invalid numerator", func(tx *testing.T) {
		_, status := handleFile(suite.inner_log, "test_files/invalid_numerator.grade")
		output := suite.inner_log_buf.String()
		assert.Equal(tx, "error [invalid_numerator.grade:5]: the numerator in one of the scores did not compile to a float: '20a/21'\n", output)
		assert.Equal(tx, 1, status)
	})

	suite.inner_log_buf.Reset()

	t.Run("test file with invalid denominator", func(tx *testing.T) {
		_, status := handleFile(suite.inner_log, "test_files/invalid_denominator.grade")
		output := suite.inner_log_buf.String()
		assert.Equal(tx, "error [invalid_denominator.grade:5]: the denominator in one of the scores did not compile to a float: '20/21a'\n", output)
		assert.Equal(tx, 1, status)
	})

	suite.inner_log_buf.Reset()

	t.Run("test file with invalid/misspelt field", func(tx *testing.T) {
		_, status := handleFile(suite.inner_log, "test_files/misspelt_field.grade")
		output := suite.inner_log_buf.String()
		assert.Equal(tx, "error [misspelt_field.grade:7]: recieved an invalid field name: 'weigth'\n", output)
		assert.Equal(tx, 1, status)
	})

	suite.inner_log_buf.Reset()

	t.Run("test file with double quotes in options", func(tx *testing.T) {
		d, status := handleFile(suite.inner_log, "test_files/quotes.grade")
		output := suite.inner_log_buf.String()

		assert.Equal(tx, "CSC 110", d.name)
		assert.Equal(tx, "A+", d.explicitGrade)
		assert.Equal(tx, "", output)
		assert.Equal(tx, 0, status)
	})

	// TODO: single quotes? idk

	suite.inner_log_buf.Reset()

	t.Run("test that the program prints out all the errors, not just one", func(tx *testing.T) {
		_, status := handleFile(suite.inner_log, "test_files/lot-wrong.grade")
		output := suite.inner_log_buf.String()

		assert.Equal(tx, "error [lot-wrong.grade:2]: the value for desired_grade is not between 0 and 100: '-20'\nerror [lot-wrong.grade:4]: recieved more than one meta headers\nerror [lot-wrong.grade:7]: recieved a grade part with no name\nerror [lot-wrong.grade:8]: recieved a line that is not under a grade part\nerror [lot-wrong.grade:9]: recieved a line that is not under a grade part\nerror [lot-wrong.grade:12]: the value for weight did not compile to a float: '0a'\nerror [lot-wrong.grade:15]: recieved a duplicate grade part name: 'Quizzes'\nerror [lot-wrong.grade:16]: the value for weight did not compile to a float: '0k'\n", output)
		assert.Equal(tx, 1, status)
	})

	suite.inner_log_buf.Reset()
}

func (suite *MainTestSuite) TestPrintGrades() {
	t := suite.T()

	t.Run("test print grades on directory", func(tx *testing.T) {
		output := captureCombined(func() {
			d, status := handleDirectory(suite.inner_log, suite.directory, GradeSection{name: "test_files"})
			assert.Equal(t, 0, status)
			calculateGPA(d)

			printGrades(suite.inner_log, d, "", false)
		})

		assert.Equal(t, expected_printed_output_non_verbose, output)

		output = captureCombined(func() {
			d, status := handleDirectory(suite.inner_log, suite.directory, GradeSection{name: "test_files"})
			assert.Equal(t, 0, status)
			calculateGPA(d)

			printGrades(suite.inner_log, d, "", true)
		})

		assert.Equal(t, expected_printed_output_verbose, output)
		assert.Equal(t, "", suite.inner_log_buf.String(), "log should be empty")
	})

	t.Run("test print grades on directory with child files and child dirs side by side", func(tx *testing.T) {
		output := captureCombined(func() {
			d, status := handleDirectory(suite.inner_log, "test_files", GradeSection{name: "test_files"})
			assert.Equal(t, 0, status)
			calculateGPA(d)

			printGrades(suite.inner_log, d, "", false)
		})

		assert.Equal(t, expected_files_with_dirs_output, output)
		assert.NotEqual(t, "", suite.inner_log_buf.String(), "log should not be empty")
	})

}

func (suite *MainTestSuite) TestRun() {
	t := suite.T()

	t.Run("test run with --help and -h", func(tx *testing.T) {
		output := captureCombined(func() {
			status := run([]string{"--help"})
			assert.Equal(t, 0, status)
		})

		assert.Equal(t, "Usage: gpa [file] [-e|--edit] [-h|--help] [-v|--verbose] [--version]\nfile: a path to examine, it can be a file or directory\n", output)

		output = captureCombined(func() {
			status := run([]string{"-h"})
			assert.Equal(t, 0, status)
		})

		assert.Equal(t, "Usage: gpa [file] [-e|--edit] [-h|--help] [-v|--verbose] [--version]\nfile: a path to examine, it can be a file or directory\n", output)
	})

	t.Run("test run with --version", func(tx *testing.T) {
		output := captureCombined(func() {
			status := run([]string{"--version"})
			assert.Equal(t, 0, status)
		})

		assert.Equal(t, "gpa-calculator version 1.1.0\n", output)
	})

	t.Run("test run on directory with verbose", func(tx *testing.T) {
		output := captureCombined(func() {
			status := run([]string{"-v", "test_files/grades/"})
			assert.Equal(t, 0, status)
		})
		assert.Equal(t, "test_files/grades/ (3.19)\n"+expected_printed_output_verbose, output)
	})

	t.Run("test run on directory without verbose", func(tx *testing.T) {
		output := captureCombined(func() {
			status := run([]string{"test_files/grades/"})
			assert.Equal(t, 0, status)
		})
		assert.Equal(t, "test_files/grades/ (3.19)\n"+expected_printed_output_non_verbose, output)
	})

	t.Run("test run on file with verbose", func(tx *testing.T) {
		output := captureCombined(func() {
			status := run([]string{"-v", "test_files/grades/2022/fall/ma100.grade"})
			assert.Equal(t, 0, status)
		})

		expected_output := `test_files/grades/2022/fall/ma100.grade
└── ma100.grade (78.31) (C+)
     ├── Homework (87.36) (B+)
     ├── Quizzes (83.80) (B-)
     ├── Midterm (60.37) (D-)
     └── Final Exam (85.85) (B)
`
		assert.Equal(t, expected_output, output)
	})

	t.Run("test run on file without verbose", func(tx *testing.T) {
		output := captureCombined(func() {
			status := run([]string{"test_files/grades/2022/fall/ma100.grade"})
			assert.Equal(t, 0, status)
		})

		assert.Equal(t, "test_files/grades/2022/fall/ma100.grade\n└── ma100.grade (78.31) (C+)\n", output)
	})

	t.Run("test run on file with more than 1 positional arguments", func(tx *testing.T) {
		output := captureCombined(func() {
			status := run([]string{"test_files/grades/2022/fall/ma100.grade", "test_files/grades/2022/fall/ma100.grade"})
			assert.Equal(t, 1, status)
		})

		assert.Equal(t, "error expected 0-1 positional argument, recieved 2\n", output)
	})

	t.Run("test run with no file provided and no GRADES_DIR", func(tx *testing.T) {
		os.Setenv("GRADES_DIR", "")

		output := captureCombined(func() {
			status := run([]string{})
			assert.Equal(t, 1, status)
		})

		assert.Equal(t, "error did not recieve a positional argument, and GRADES_DIR is not set\n", string(output))
	})

	t.Run("test run with no file provided and GRADES_DIR set", func(tx *testing.T) {
		os.Setenv("GRADES_DIR", "test_files/grades/")
		output := captureCombined(func() {
			status := run([]string{})
			assert.Equal(t, 0, status)
		})

		assert.Equal(t, "test_files/grades/ (3.19)\n"+expected_printed_output_non_verbose, string(output))
	})

	t.Run("test run with invalid file and no GRADES_DIR", func(tx *testing.T) {
		os.Setenv("GRADES_DIR", "")
		output := captureCombined(func() {
			status := run([]string{"test_files/this_file_does_not_exist.grade"})
			assert.Equal(t, 1, status)
		})

		assert.Equal(t, "error could not find file or directory 'test_files/this_file_does_not_exist.grade', no fuzzy find search occured as $GRADES_DIR is not set\n", string(output))
	})

	t.Run("test run with invalid file and GRADES_DIR set", func(tx *testing.T) {
		os.Setenv("GRADES_DIR", "test_files/grades/")
		output := captureCombined(func() {
			status := run([]string{"test_files/this_file_does_not_exist.grade"})
			assert.Equal(t, 1, status)
		})

		assert.Equal(t, "error could not find file or directory 'test_files/this_file_does_not_exist.grade', even with fuzzy find\n", string(output))
	})

	t.Run("test run with valid fuzzy file and GRADES_DIR set", func(tx *testing.T) {
		os.Setenv("GRADES_DIR", "test_files/grades/")
		output := captureCombined(func() {
			status := run([]string{"ma100.grade"})
			assert.Equal(t, 0, status)
		})

		assert.Equal(t, "test_files/grades/2022/fall/ma100.grade\n└── ma100.grade (78.31) (C+)\n", output)
	})

	t.Run("test run with valid fuzzy directory and GRADES_DIR set", func(tx *testing.T) {
		// also test that fuzzy directory search does not go into directories that do not have read permissions
		err := os.Mkdir("test_files/grades/2022/autumn", 0000)

		if err != nil {
			t.Error("could not create directory 'test_files/grades/2022/autumn'")
		}

		os.Setenv("GRADES_DIR", "test_files/grades/")
		output := captureCombined(func() {
			status := run([]string{"fall"})
			assert.Equal(t, 0, status)
		})

		assert.Equal(t, expected_fuzzy_directory_output, output)
		os.Remove("test_files/grades/2022/autumn")
	})

	t.Run("test run with valid fuzzy but file no permission and GRADES_DIR set", func(tx *testing.T) {
		os.Setenv("GRADES_DIR", "test_files/grades/")
		f, err := os.Create("test_files/grades/2022/fall/eng300.grade")

		if err != nil {
			t.Error("could not create file 'test_files/eng300.grade'")
		}

		f.Chmod(0000)
		output := captureCombined(func() {
			status := run([]string{"eng300"})
			assert.Equal(t, 1, status)
		})

		assert.Equal(t, "error could not read file 'test_files/grades/2022/fall/eng300.grade'\n", string(output))
		os.Remove(f.Name())
		f.Close()
	})

	t.Run("test run with valid file but no permission", func(tx *testing.T) {
		f, err := os.Create("test_files/no_permissions.grade")

		if err != nil {
			t.Error("could not create file 'test_files/no_permissions.grade'")
		}

		f.Chmod(0000)
		output := captureCombined(func() {
			status := run([]string{"test_files/no_permissions.grade"})
			assert.Equal(t, 1, status)
		})

		assert.Equal(t, "error could not read file 'test_files/no_permissions.grade'\n", string(output))
		os.Remove(f.Name())
		f.Close()
	})

	t.Run("test run with valid directory but no permission", func(tx *testing.T) {
		err := os.Mkdir("test_files/no_permissions/", 0000)

		if err != nil {
			t.Error("could not create directory 'test_files/no_permissions/'")
		}

		output := captureCombined(func() {
			status := run([]string{"test_files/no_permissions/"})
			assert.Equal(t, 1, status)
		})

		assert.Equal(t, "error could not read directory 'test_files/no_permissions/'\n", string(output))
		os.Remove("test_files/no_permissions/")
	})

	t.Run("test run on directory with one child directory no permissions", func(tx *testing.T) {
		if err := os.Mkdir("test_files/grades/2022/winter", 0000); err != nil {
			t.Error("could not create directory 'test_files/grades/2022/winter'")
		}

		output := captureCombined(func() {
			status := run([]string{"test_files/grades/"})
			assert.Equal(t, 0, status)
		})

		assert.Equal(t, "error could not read directory 'test_files/grades/2022/winter'\ntest_files/grades/ (3.19)\n"+expected_printed_output_non_verbose, string(output))

		if err := os.Remove("test_files/grades/2022/winter"); err != nil {
			t.Error("could not remove directory 'test_files/grades/2022/winter'")
		}

	})

	t.Run("test run with --edit without EDITOR", func(tx *testing.T) {
		os.Setenv("EDITOR", "")
		output := captureCombined(func() {
			status := run([]string{"--edit", "test_files/grades/2022/fall/ma100.grade"})
			assert.Equal(t, 1, status)
		})

		assert.Equal(t, "error $EDITOR is not set\n", output)
	})

	t.Run("test run with --edit on directory", func(tx *testing.T) {
		os.Setenv("EDITOR", "vi")
		output := captureCombined(func() {
			status := run([]string{"--edit", "test_files/grades/"})
			assert.Equal(t, 1, status)
		})

		assert.Equal(t, "error editing directories is not supported\n", output)
	})

	t.Run("test run with --edit on invalid editor", func(tx *testing.T) {
		os.Setenv("EDITOR", "vi-this-editor-does-not-exist")
		output := captureCombined(func() {
			status := run([]string{"--edit", "test_files/grades/2022/fall/ma100.grade"})
			assert.Equal(t, 1, status)
		})

		assert.Equal(t, "error could not open file test_files/grades/2022/fall/ma100.grade in $EDITOR\n", output)
	})

	t.Run("test file with no data", func(tx *testing.T) {
		output := captureCombined(func() {
			status := run([]string{"test_files/no_data.grade"})
			assert.Equal(t, 0, status)
		})

		assert.Equal(t, "test_files/no_data.grade\n└── no_data.grade (unset)\n", output)
	})

	t.Run("test file with no data verbose", func(tx *testing.T) {
		output := captureCombined(func() {
			status := run([]string{"-v", "test_files/no_data.grade"})
			assert.Equal(t, 0, status)
		})

		assert.Equal(t, expected_no_data_verbose, output)
	})

	t.Run("test file with desired with no final", func(tx *testing.T) {
		output := captureCombined(func() {
			status := run([]string{"test_files/desired_with_no_final.grade"})
			assert.Equal(t, 0, status)
		})

		assert.Equal(t, "test_files/desired_with_no_final.grade\n└── CSC 110 (88.36) (B+)\n", string(output))

		output = captureCombined(func() {
			status := run([]string{"-v", "test_files/desired_with_no_final.grade"})
			assert.Equal(t, 0, status)
		})

		assert.Equal(t, expected_desired_no_final_output, string(output))
	})

	t.Run("test run with file that has ignore=true", func(tx *testing.T) {
		output := captureCombined(func() {
			status := run([]string{"test_files/ignore.grade"})
			assert.Equal(t, 0, status)
		})

		assert.Equal(t, "test_files/ignore.grade\n└── ignore.grade (85.16) (B)\n", string(output))
	})

	// TODO: add --edit test with EDITOR set and everything valid
}

func captureCombined(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func (suite *MainTestSuite) SetupSuite() {
	suite.directory = "test_files/grades/"
	suite.inner_log_buf = new(bytes.Buffer)
	suite.inner_log = log.New(suite.inner_log_buf, "", 0)

	if _, err := os.Stat(suite.directory); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Test directory (test_files/) does not exist\n")
		os.Exit(1)
	}
}

// run after each test
func (suite *MainTestSuite) TearDownTest() {
	suite.inner_log_buf.Reset()
}

// run after end of all tests
func (suite *MainTestSuite) TearDownSuite() {
}

func TestMainTestSuite(t *testing.T) {
	suite.Run(t, new(MainTestSuite))
}
