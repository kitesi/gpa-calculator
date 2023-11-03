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
	} else if grade >= 0 {
		return "F"
	} else {
		return "?"
	}
}

func getGradeGPA(grade float64) float64 {
	grade = grade * 100

	if grade >= 94 {
		return 4.0
	} else if grade >= 90 {
		return 3.7
	} else if grade >= 87 {
		return 3.3
	} else if grade >= 84 {
		return 3.0
	} else if grade >= 80 {
		return 2.7
	} else if grade >= 77 {
		return 2.3
	} else if grade >= 74 {
		return 2.0
	} else if grade >= 70 {
		return 1.7
	} else if grade >= 67 {
		return 1.3
	} else if grade >= 64 {
		return 1.0
	} else if grade >= 60 {
		return 0.7
	} else {
		return 0
	}
}
