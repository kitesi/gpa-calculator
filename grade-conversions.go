package main

func getGradeLetter(grade float64) string {
	grade = grade * 100

	if grade >= 94 {
		return "A"
	} else if grade >= 90 {
		return "A-"
	} else if grade >= 87 {
		return "B+"
	} else if grade >= 84 {
		return "B"
	} else if grade >= 80 {
		return "B-"
	} else if grade >= 77 {
		return "C+"
	} else if grade >= 74 {
		return "C"
	} else if grade >= 70 {
		return "C-"
	} else if grade >= 67 {
		return "D+"
	} else if grade >= 64 {
		return "D"
	} else if grade >= 60 {
		return "D-"
	} else {
		return "F"
	}
}

func getGradeGPA(gradeLetter string) float64 {
	if gradeLetter == "A+" {
		return 4.3
	} else if gradeLetter == "A" {
		return 4.0
	} else if gradeLetter == "A-" {
		return 3.7
	} else if gradeLetter == "B+" {
		return 3.3
	} else if gradeLetter == "B" {
		return 3.0
	} else if gradeLetter == "B-" {
		return 2.7
	} else if gradeLetter == "C+" {
		return 2.3
	} else if gradeLetter == "C" {
		return 2.0
	} else if gradeLetter == "C-" {
		return 1.7
	} else if gradeLetter == "D+" {
		return 1.3
	} else if gradeLetter == "D" {
		return 1.0
	} else if gradeLetter == "D-" {
		return 0.7
	} else if gradeLetter == "F" {
		return 0
	} else {
		return -1
	}
}
