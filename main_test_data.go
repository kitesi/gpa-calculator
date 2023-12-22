package main

var expected_gs_from_test_directory = GradeSection{
	name:    "test_files",
	credits: 64,
	gpa:     3.1875,
	gradeSubsections: []*GradeSection{
		{
			name:    "2022",
			credits: 32,
			gpa:     3.2875,
			gradeSubsections: []*GradeSection{
				{
					name:    "fall",
					credits: 16,
					gpa:     3.325,
					classes: []*SchoolClass{
						{
							name:          "cs100.grade",
							grade:         0.8516093252053647,
							credits:       4,
							totalWeight:   1.0000000223517418,
							desiredGrade:  0.8,
							explicitGrade: "A",
							gradeParts: []*GradePart{
								{
									weight:         0.20000000298023224,
									pointsRecieved: 72.19999980926514,
									pointsTotal:    80,
									name:           "Homework",
								},
								{
									weight:         0.10000000149011612,
									pointsRecieved: 68.52999973297119,
									pointsTotal:    80,
									name:           "Quizzes",
								},
								{
									weight:         0.30000001192092896,
									pointsRecieved: 36.08000183105469,
									pointsTotal:    41,
									name:           "Midterm",
								},
								{
									weight:         0.4000000059604645,
									pointsRecieved: 37.77000045776367,
									pointsTotal:    47,
									name:           "Final Exam",
								},
							},
						},
						{
							name:         "gov100.grade",
							grade:        0.8714684895819538,
							credits:      4,
							totalWeight:  1.0000000223517418,
							desiredGrade: 0.8,
							gradeParts: []*GradePart{
								{
									weight:         0.20000000298023224,
									pointsRecieved: 69.9000015258789,
									pointsTotal:    80,
									name:           "Homework",
								},
								{
									weight:         0.10000000149011612,
									pointsRecieved: 74.92000007629395,
									pointsTotal:    80,
									name:           "Quizzes",
								},
								{
									weight:         0.30000001192092896,
									pointsRecieved: 37.06999969482422,
									pointsTotal:    40,
									name:           "Midterm",
								},
								{
									weight:         0.4000000059604645,
									pointsRecieved: 37.380001068115234,
									pointsTotal:    46,
									name:           "Final Exam",
								},
							},
						},
						{
							name:         "lang100.grade",
							grade:        0.9259115123192122,
							totalWeight:  1.0000000223517418,
							credits:      4,
							desiredGrade: -1,
							gradeParts: []*GradePart{
								{
									weight:         0.20000000298023224,
									pointsRecieved: 72,
									pointsTotal:    80,
									name:           "Homework",
								},
								{
									weight:         0.10000000149011612,
									pointsRecieved: 70.66999912261963,
									pointsTotal:    80,
									name:           "Quizzes",
								},
								{
									weight:         0.30000001192092896,
									pointsRecieved: 39.45000076293945,
									pointsTotal:    43,
									name:           "Midterm",
								},
								{
									weight:         0.4000000059604645,
									pointsRecieved: 39.189998626708984,
									pointsTotal:    41,
									name:           "Final Exam",
								},
							},
						},
						{

							name:         "ma100.grade",
							grade:        0.7830559137464939,
							credits:      4,
							totalWeight:  1.0000000223517418,
							desiredGrade: -1,
							gradeParts: []*GradePart{
								{
									weight:         0.20000000298023224,
									pointsRecieved: 69.89000034332275,
									pointsTotal:    80,
									name:           "Homework",
								},
								{
									weight:         0.10000000149011612,
									pointsRecieved: 67.03999996185303,
									pointsTotal:    80,
									name:           "Quizzes",
								},
								{
									weight:         0.30000001192092896,
									pointsRecieved: 25.959999084472656,
									pointsTotal:    43,
									name:           "Midterm",
								},
								{
									weight:         0.4000000059604645,
									pointsRecieved: 35.20000076293945,
									pointsTotal:    41,
									name:           "Final Exam",
								},
							},
						},
					},
				},
				{
					name:    "spring",
					credits: 16,
					gpa:     3.25,
					classes: []*SchoolClass{
						{
							name:          "cs200.grade",
							grade:         0.9601476206513032,
							credits:       4,
							totalWeight:   1.0000000223517418,
							desiredGrade:  -1,
							explicitGrade: "A+",
							gradeParts: []*GradePart{
								{
									weight:         0.20000000298023224,
									pointsRecieved: 77.18000030517578,
									pointsTotal:    80,
									name:           "Homework",
								},
								{
									weight:         0.10000000149011612,
									pointsRecieved: 72.1200008392334,
									pointsTotal:    80,
									name:           "Quizzes",
								},
								{
									weight:         0.30000001192092896,
									pointsRecieved: 36.79999923706055,
									pointsTotal:    40,
									name:           "Midterm",
								},
								{
									weight:         0.4000000059604645,
									pointsRecieved: 42.11000061035156,
									pointsTotal:    42,
									name:           "Final Exam",
								},
							},
						},
						{
							name:         "gov200.grade",
							grade:        0.8479766759344356,
							credits:      4,
							totalWeight:  1.0000000223517418,
							desiredGrade: -1,
							gradeParts: []*GradePart{
								{
									weight:         0.20000000298023224,
									pointsRecieved: 73.67000007629395,
									pointsTotal:    80,
									name:           "Homework",
								},
								{
									weight:         0.10000000149011612,
									pointsRecieved: 71.5,
									pointsTotal:    80,
									name:           "Quizzes",
								},
								{
									weight:         0.30000001192092896,
									pointsRecieved: 34,
									pointsTotal:    45,
									name:           "Midterm",
								},
								{
									weight:         0.4000000059604645,
									pointsRecieved: 43.470001220703125,
									pointsTotal:    50,
									name:           "Final Exam",
								},
							},
						},
						{
							name:         "lang200.grade",
							grade:        0.8119374838960355,
							totalWeight:  1.0000000223517418,
							credits:      4,
							desiredGrade: -1,
							gradeParts: []*GradePart{
								{
									weight:         0.20000000298023224,
									pointsRecieved: 69.42000007629395,
									pointsTotal:    80,
									name:           "Homework",
								},
								{
									weight:         0.10000000149011612,
									pointsRecieved: 72.98999977111816,
									pointsTotal:    80,
									name:           "Quizzes",
								},
								{
									weight:         0.30000001192092896,
									pointsRecieved: 36.599998474121094,
									pointsTotal:    48,
									name:           "Midterm",
								},
								{
									weight:         0.4000000059604645,
									pointsRecieved: 39.79999923706055,
									pointsTotal:    50,
									name:           "Final Exam",
								},
							},
						},
						{
							name:         "ma200.grade",
							grade:        0.8679067898086685,
							credits:      4,
							desiredGrade: -1,
							totalWeight:  1.0000000223517418,
							gradeParts: []*GradePart{
								{
									weight:         0.20000000298023224,
									pointsRecieved: 71.89999961853027,
									pointsTotal:    80,
									name:           "Homework",
								},
								{
									weight:         0.10000000149011612,
									pointsRecieved: 73.16000175476074,
									pointsTotal:    80,
									name:           "Quizzes",
								},
								{
									weight:         0.30000001192092896,
									pointsRecieved: 46.9900016784668,
									pointsTotal:    43,
									name:           "Midterm",
								},
								{
									weight:         0.4000000059604645,
									pointsRecieved: 30.920000076293945,
									pointsTotal:    46,
									name:           "Final Exam",
								},
							},
						},
						{
							name:         "Statistics 200",
							desiredGrade: -1,
							credits:      4,
							grade:        -1,
							gradeParts: []*GradePart{
								{
									name:   "Final",
									weight: 1,
								},
							},
						},
					},
				},
			},
		},
		{
			name:    "2023",
			credits: 32,
			gpa:     3.0875000000000004,
			gradeSubsections: []*GradeSection{
				{
					name:    "fall",
					credits: 16,
					gpa:     2.825,
					classes: []*SchoolClass{
						{
							name:         "cs300.grade",
							grade:        0.8640454528473995,
							credits:      4,
							totalWeight:  1.0000000223517418,
							desiredGrade: -1,
							gradeParts: []*GradePart{
								{
									weight:         0.20000000298023224,
									pointsRecieved: 70.09999942779541,
									pointsTotal:    80,
									name:           "Homework",
								},
								{
									weight:         0.10000000149011612,
									pointsRecieved: 75.92000007629395,
									pointsTotal:    80,
									name:           "Quizzes",
								},
								{
									weight:         0.30000001192092896,
									pointsRecieved: 30.549999237060547,
									pointsTotal:    44,
									name:           "Midterm",
								},
								{
									weight:         0.4000000059604645,
									pointsRecieved: 48.20000076293945,
									pointsTotal:    50,
									name:           "Final Exam",
								},
							},
						},
						{
							name:         "gov300.grade",
							grade:        0.7933186110832144,
							credits:      4,
							totalWeight:  1.0000000223517418,
							desiredGrade: -1,
							gradeParts: []*GradePart{
								{
									weight:         0.20000000298023224,
									pointsRecieved: 70,
									pointsTotal:    80,
									name:           "Homework",
								},
								{
									weight:         0.10000000149011612,
									pointsRecieved: 74.95999908447266,
									pointsTotal:    80,
									name:           "Quizzes",
								},
								{
									weight:         0.30000001192092896,
									pointsRecieved: 20.700000762939453,
									pointsTotal:    43,
									name:           "Midterm",
								},
								{
									weight:         0.4000000059604645,
									pointsRecieved: 38.02000045776367,
									pointsTotal:    40,
									name:           "Final Exam",
								},
							},
						},
						{
							name:         "lang300.grade",
							grade:        0.8511722166753065,
							totalWeight:  1.0000000223517418,
							credits:      4,
							desiredGrade: -1,
							gradeParts: []*GradePart{
								{
									weight:         0.20000000298023224,
									pointsRecieved: 71.45999908447266,
									pointsTotal:    80,
									name:           "Homework",
								},
								{
									weight:         0.10000000149011612,
									pointsRecieved: 76.76000022888184,
									pointsTotal:    80,
									name:           "Quizzes",
								},
								{
									weight:         0.30000001192092896,
									pointsRecieved: 23.780000686645508,
									pointsTotal:    40,
									name:           "Midterm",
								},
								{
									weight:         0.4000000059604645,
									pointsRecieved: 44.79999923706055,
									pointsTotal:    45,
									name:           "Final Exam",
								},
							},
						},
						{

							name:         "ma300.grade",
							grade:        0.8510678410202177,
							credits:      4,
							desiredGrade: -1,
							totalWeight:  1.0000000223517418,
							gradeParts: []*GradePart{
								{
									weight:         0.20000000298023224,
									pointsRecieved: 80.04999923706055,
									pointsTotal:    80,
									name:           "Homework",
								},
								{
									weight:         0.10000000149011612,
									pointsRecieved: 73.60000133514404,
									pointsTotal:    80,
									name:           "Quizzes",
								},
								{
									weight:         0.30000001192092896,
									pointsRecieved: 21.15999984741211,
									pointsTotal:    42,
									name:           "Midterm",
								},
								{
									weight:         0.4000000059604645,
									pointsRecieved: 40.779998779296875,
									pointsTotal:    40,
									name:           "Final Exam",
								},
							},
						},
					},
				},
				{
					name:    "spring",
					credits: 16,
					gpa:     3.3499999999999996,
					classes: []*SchoolClass{
						{
							name:         "cs400.grade",
							grade:        0.9292773822568803,
							credits:      4,
							totalWeight:  1.0000000223517418,
							desiredGrade: -1,
							gradeParts: []*GradePart{
								{
									weight:         0.20000000298023224,
									pointsRecieved: 70.69999980926514,
									pointsTotal:    80,
									name:           "Homework",
								},
								{
									weight:         0.10000000149011612,
									pointsRecieved: 73.8400011062622,
									pointsTotal:    80,
									name:           "Quizzes",
								},
								{
									weight:         0.30000001192092896,
									pointsRecieved: 36.970001220703125,
									pointsTotal:    40,
									name:           "Midterm",
								},
								{
									weight:         0.4000000059604645,
									pointsRecieved: 40.209999084472656,
									pointsTotal:    42,
									name:           "Final Exam",
								},
							},
						},
						{
							name:         "gov400.grade",
							grade:        0.8698477377421817,
							credits:      4,
							totalWeight:  1.0000000223517418,
							desiredGrade: -1,
							gradeParts: []*GradePart{
								{
									weight:         0.20000000298023224,
									pointsRecieved: 67.98000144958496,
									pointsTotal:    80,
									name:           "Homework",
								},
								{
									weight:         0.10000000149011612,
									pointsRecieved: 79.80999946594238,
									pointsTotal:    80,
									name:           "Quizzes",
								},
								{
									weight:         0.30000001192092896,
									pointsRecieved: 28.469999313354492,
									pointsTotal:    41,
									name:           "Midterm",
								},
								{
									weight:         0.4000000059604645,
									pointsRecieved: 43.099998474121094,
									pointsTotal:    44,
									name:           "Final Exam",
								},
							},
						},
						{
							name:         "lang400.grade",
							grade:        0.9466013828702993,
							totalWeight:  1.0000000223517418,
							credits:      4,
							desiredGrade: -1,
							gradeParts: []*GradePart{
								{
									weight:         0.20000000298023224,
									pointsRecieved: 75.96999835968018,
									pointsTotal:    80,
									name:           "Homework",
								},
								{
									weight:         0.10000000149011612,
									pointsRecieved: 67.51000022888184,
									pointsTotal:    80,
									name:           "Quizzes",
								},
								{
									weight:         0.30000001192092896,
									pointsRecieved: 45.66999816894531,
									pointsTotal:    45,
									name:           "Midterm",
								},
								{
									weight:         0.4000000059604645,
									pointsRecieved: 41.380001068115234,
									pointsTotal:    45,
									name:           "Final Exam",
								},
							},
						},
						{
							name:         "ma400.grade",
							grade:        0.8166466518832215,
							credits:      4,
							desiredGrade: -1,
							totalWeight:  1.0000000223517418,
							gradeParts: []*GradePart{
								{
									weight:         0.20000000298023224,
									pointsRecieved: 68.14000034332275,
									pointsTotal:    80,
									name:           "Homework",
								},
								{
									weight:         0.10000000149011612,
									pointsRecieved: 74.93000030517578,
									pointsTotal:    80,
									name:           "Quizzes",
								},
								{
									weight:         0.30000001192092896,
									pointsRecieved: 32.58000183105469,
									pointsTotal:    41,
									name:           "Midterm",
								},
								{
									weight:         0.4000000059604645,
									pointsRecieved: 32.209999084472656,
									pointsTotal:    41,
									name:           "Final Exam",
								},
							},
						},
					},
				},
			},
		},
	},
}

var expected_ignore_class = SchoolClass{
	name:         "ignore.grade",
	credits:      4,
	grade:        0.8516093252053647,
	totalWeight:  1.0000000223517418,
	desiredGrade: 1,
	gradeParts: []*GradePart{
		{
			name:           "Homework",
			weight:         0.20000000298023224,
			pointsRecieved: 72.19999980926514,
			pointsTotal:    80,
		},
		{
			name:           "Quizzes",
			weight:         0.10000000149011612,
			pointsRecieved: 68.52999973297119,
			pointsTotal:    80,
		},
		{
			name:           "Midterm",
			weight:         0.30000001192092896,
			pointsRecieved: 36.08000183105469,
			pointsTotal:    41,
		},
		{
			name:           "Final Exam",
			weight:         0.4000000059604645,
			pointsRecieved: 37.77000045776367,
			pointsTotal:    47,
		},
	},
}

var expected_printed_output_non_verbose = `├── 2022 (3.29)
│   ├── fall (3.33)
│   │   ├── cs100.grade (85.16) (A)
│   │   ├── gov100.grade (87.15) (B+)
│   │   ├── lang100.grade (92.59) (A-)
│   │   └── ma100.grade (78.31) (C+)
│   └── spring (3.25)
│       ├── Statistics 200 (unset)
│       ├── cs200.grade (96.01) (A+)
│       ├── gov200.grade (84.80) (B)
│       ├── lang200.grade (81.19) (B-)
│       └── ma200.grade (86.79) (B)
└── 2023 (3.09)
    ├── fall (2.83)
    │   ├── cs300.grade (86.40) (B)
    │   ├── gov300.grade (79.33) (C+)
    │   ├── lang300.grade (85.12) (B)
    │   └── ma300.grade (85.11) (B)
    └── spring (3.35)
        ├── cs400.grade (92.93) (A-)
        ├── gov400.grade (86.98) (B)
        ├── lang400.grade (94.66) (A)
        └── ma400.grade (81.66) (B-)
`

var expected_printed_output_verbose = `├── 2022 (3.29)
│   ├── fall (3.33)
│   │   ├── cs100.grade (85.16) (A)
│   │   │    ├── Homework (90.25) (A-)
│   │   │    ├── Quizzes (85.66) (B)
│   │   │    ├── Midterm (88.00) (B+)
│   │   │    ├── Final Exam (80.36) (B-)
│   │   │    └── to get a 80.00% you need at least a 67.46% on the final
│   │   ├── gov100.grade (87.15) (B+)
│   │   │    ├── Homework (87.38) (B+)
│   │   │    ├── Quizzes (93.65) (A-)
│   │   │    ├── Midterm (92.67) (A-)
│   │   │    ├── Final Exam (81.26) (B-)
│   │   │    └── to get a 80.00% you need at least a 63.39% on the final
│   │   ├── lang100.grade (92.59) (A-)
│   │   │    ├── Homework (90.00) (A-)
│   │   │    ├── Quizzes (88.34) (B+)
│   │   │    ├── Midterm (91.74) (A-)
│   │   │    └── Final Exam (95.59) (A)
│   │   └── ma100.grade (78.31) (C+)
│   │        ├── Homework (87.36) (B+)
│   │        ├── Quizzes (83.80) (B-)
│   │        ├── Midterm (60.37) (D-)
│   │        └── Final Exam (85.85) (B)
│   └── spring (3.25)
│       ├── Statistics 200 (unset)
│       │    └── Final (unset)
│       ├── cs200.grade (96.01) (A+)
│       │    ├── Homework (96.48) (A)
│       │    ├── Quizzes (90.15) (A-)
│       │    ├── Midterm (92.00) (A-)
│       │    └── Final Exam (100.26) (A)
│       ├── gov200.grade (84.80) (B)
│       │    ├── Homework (92.09) (A-)
│       │    ├── Quizzes (89.38) (B+)
│       │    ├── Midterm (75.56) (C)
│       │    └── Final Exam (86.94) (B)
│       ├── lang200.grade (81.19) (B-)
│       │    ├── Homework (86.78) (B)
│       │    ├── Quizzes (91.24) (A-)
│       │    ├── Midterm (76.25) (C)
│       │    └── Final Exam (79.60) (C+)
│       └── ma200.grade (86.79) (B)
│            ├── Homework (89.87) (B+)
│            ├── Quizzes (91.45) (A-)
│            ├── Midterm (109.28) (A)
│            └── Final Exam (67.22) (D+)
└── 2023 (3.09)
    ├── fall (2.83)
    │   ├── cs300.grade (86.40) (B)
    │   │    ├── Homework (87.62) (B+)
    │   │    ├── Quizzes (94.90) (A)
    │   │    ├── Midterm (69.43) (D+)
    │   │    └── Final Exam (96.40) (A)
    │   ├── gov300.grade (79.33) (C+)
    │   │    ├── Homework (87.50) (B+)
    │   │    ├── Quizzes (93.70) (A-)
    │   │    ├── Midterm (48.14) (F)
    │   │    └── Final Exam (95.05) (A)
    │   ├── lang300.grade (85.12) (B)
    │   │    ├── Homework (89.32) (B+)
    │   │    ├── Quizzes (95.95) (A)
    │   │    ├── Midterm (59.45) (F)
    │   │    └── Final Exam (99.56) (A)
    │   └── ma300.grade (85.11) (B)
    │        ├── Homework (100.06) (A)
    │        ├── Quizzes (92.00) (A-)
    │        ├── Midterm (50.38) (F)
    │        └── Final Exam (101.95) (A)
    └── spring (3.35)
        ├── cs400.grade (92.93) (A-)
        │    ├── Homework (88.37) (B+)
        │    ├── Quizzes (92.30) (A-)
        │    ├── Midterm (92.43) (A-)
        │    └── Final Exam (95.74) (A)
        ├── gov400.grade (86.98) (B)
        │    ├── Homework (84.98) (B)
        │    ├── Quizzes (99.76) (A)
        │    ├── Midterm (69.44) (D+)
        │    └── Final Exam (97.95) (A)
        ├── lang400.grade (94.66) (A)
        │    ├── Homework (94.96) (A)
        │    ├── Quizzes (84.39) (B)
        │    ├── Midterm (101.49) (A)
        │    └── Final Exam (91.96) (A-)
        └── ma400.grade (81.66) (B-)
             ├── Homework (85.18) (B)
             ├── Quizzes (93.66) (A-)
             ├── Midterm (79.46) (C+)
             └── Final Exam (78.56) (C+)
`

var expected_fuzzy_directory_output = `test_files/grades/2022/fall (3.33)
├── cs100.grade (85.16) (A)
├── gov100.grade (87.15) (B+)
├── lang100.grade (92.59) (A-)
└── ma100.grade (78.31) (C+)
`

var expected_desired_no_final_output = `test_files/desired_with_no_final.grade
└── CSC 110 (88.36) (B+)
     ├── Homework (90.25) (A-)
     ├── Quizzes (85.66) (B)
     ├── Midterm (88.00) (B+)
error [CSC 110]: could not find a grade part that starts with 'final'
`

var expected_no_data_verbose = `test_files/no_data.grade
└── no_data.grade (unset)
     ├── Homework (unset)
     ├── Quizzes (unset)
     ├── Midterm (unset)
     ├── Final Exam (unset)
     └── to get a 80.00% you need at least a 80.00% on the final
`

var expected_files_with_dirs_output = `├── grades (3.19)
│   ├── 2022 (3.29)
│   │   ├── fall (3.33)
│   │   │   ├── cs100.grade (85.16) (A)
│   │   │   ├── gov100.grade (87.15) (B+)
│   │   │   ├── lang100.grade (92.59) (A-)
│   │   │   └── ma100.grade (78.31) (C+)
│   │   └── spring (3.25)
│   │       ├── Statistics 200 (unset)
│   │       ├── cs200.grade (96.01) (A+)
│   │       ├── gov200.grade (84.80) (B)
│   │       ├── lang200.grade (81.19) (B-)
│   │       └── ma200.grade (86.79) (B)
│   └── 2023 (3.09)
│       ├── fall (2.83)
│       │   ├── cs300.grade (86.40) (B)
│       │   ├── gov300.grade (79.33) (C+)
│       │   ├── lang300.grade (85.12) (B)
│       │   └── ma300.grade (85.11) (B)
│       └── spring (3.35)
│           ├── cs400.grade (92.93) (A-)
│           ├── gov400.grade (86.98) (B)
│           ├── lang400.grade (94.66) (A)
│           └── ma400.grade (81.66) (B-)
├── CSC 110 (85.16) (A+)
├── CSC 110 (88.36) (B+)
├── CSC 110 (85.16) (B)
├── CSC 110 (85.16) (B)
├── comments.grade (83.94) (B-)
├── credits.grade (85.16) (B)
├── multiline.grade (83.94) (B-)
├── no_data.grade (unset)
├── non_multiline.grade (83.94) (B-)
└── set_grade.grade (85.16) (A)
`
