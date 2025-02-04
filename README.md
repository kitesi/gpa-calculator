# GPA Calculator

This is a simple grade calculator written in go. It gathers the data from local
files written in a special but simple format.

## Calculation of Grades & GPA

-   GPA ranges from 0.0 to 4.0
-   A is the highest grades at 94%
-   Read `./grade-conversions.go` for full details

There is no config to change the calculation currently, so I would recommend
modifying the source code yourself. However, you can also just assign your
received grade to a class in the grade file itself.

## Structure of Folder

I recommend having a folder with all your grades, with the structure as such:

```
grades/
├── 2023
│   ├── fall
│   │   ├── cs100.grade
│   │   ├── gov100.grade
│   │   ├── lang100.grade
│   │   └── ma100.grade
│   └── spring
│       ├── cs200.grade
│       ├── gov200.grade
│       ├── lang200.grade
│       └── ma200.grade
└── 2024
    ├── fall
    │   ├── cs300.grade
    │   ├── gov300.grade
    │   ├── lang300.grade
    │   └── ma300.grade
    └── spring
        ├── cs400.grade
        ├── gov400.grade
        ├── lang400.grade
        └── ma400.grade
```

The files contain the data for each class, and they should have the extension "grade."
However, you do not need to have this structure. You can structure your folders however you
like, or you can just run the program on a single file. Note: 2023 and 2024 refers to the starting
academic year, but obviously spring occurs in the following year. A bit confusing, maybe you would
want to rename them "freshman" and "sophomore."

### Grade File

The grade file is just a simple text file with specific syntax. Grade
parts/categories and denoted with a ">" prefix and then their name, and the
children elements are the `weight`, `data`, and `drop_lowest`.

```
> Homework
weight = 0.2
data =
    20/20, # Assignment 1: What is a Derivative
    17/20, # Assignment 2: Exploring Change
    19/20  # ...
drop_lowest = 2 # drop the lowest 2 grades

> Quizes
weight = 0.4
data = 40/50, 10/10, 9/10, 8/10
drop_lowest # drop the lowest grade

> Final Exam
weight = 0.4
data = 58/60
```

Comments start with a pound symbol (#) and the program does not consider
anything past that character for that line. Indentation and extra white space
do not matter.

If you wish to specify options for the entire grade file, use the line "~ Meta"
at the start of your file. For example, each class is assumed to be 4 credits,
if you wish to change this use the credits option (there are 3 other recognized
options, but you can add whatever you like).

```
~ Meta
# Amount of credits as an int.
credits = 3

# Nice name for when printing.
name = "MA160 Multivariable Calculus & Series"

# Your desired grade as a float from 0-100. If you have this set to a value,
# when you run the program with verbosity (-v), it will tell you what you need to
# get on your final to get this grade.
desired_grade = 94

# Describe your actual grade with a letter. This is useful for curves and other
# unique systems that this program doesn't handle. For example, this program
# doesn't assign A+ automatically because in my school, as a professor has to
# manually assign an A+.

grade = "A+"

# Ignore this file when the program is run through a directory.
# Useful for when you want theoretical grades for other classes.
ignore = true
# or
# ignore

# unrecognized option but it's ok
location = Grand Hall 202

# Start a new grade part: homework given every thursday
> Homework
# ...snip...
```

## Output

This is the structure of the output with no verbosity:

```shell
$ gpa-calculator ~/grades
/home/seven/grades (3.08)
├── 2023 (3.12)
│   ├── fall (3.08)
│   │   ├── cs100.grade (85.16) (B)
│   │   ├── "GOV100: Introduction to Race & Politics"(87.15) (B+)
│   │   ├── lang100.grade (92.59) (A-)
│   │   └── math100.grade (78.31) (C+)
│   └── spring (3.17)
│       ├── gov200.grade (84.80) (B)
│       ├── lang200.grade (81.19) (B-)
│       ├── ma200.grade (86.79) (B)
│       └── cs200.grade (96.01) (A)
└── 2024 (3.06)
    ├── fall (2.91)
    │   ├── cs300.grade (86.40) (B)
    │   ├── gov300.grade (79.33) (C+)
    │   ├── lang300.grade (85.12) (B)
    │   └── ma300.grade (85.11) (B)
    └── spring (3.35)
        ├── cs400.grade (92.93) (A-)
        ├── gov400.grade (86.98) (B)
        ├── lang400.grade (94.66) (A)
        └── ma400.grade (81.66) (B-)
```

This is the structure with verbosity:

```shell
$ gpa-calculator ~/grades-1 - v
/home/seven/grades-1 (3.08)
├── 2023 (3.12)
│   ├── fall (3.08)
│   │   ├── cs100.grade (85.16) (B)
│   │   │    ├── Homework (90.25) (A-)
│   │   │    ├── Quizes (85.66) (B)
│   │   │    ├── Mid Term (88.00) (B+)
│   │   │    ├── Final Eaxm (80.36) (B-)
│   │   │    └── to get a 80.00% you need at least a 72.26% on the final
# ...snip...
```

## Installation

Install with brew:

```shell
brew tap kitesi/gpa-calculator https://github.com/kitesi/gpa-calculator
brew install gpa-calculator
```

Install with Go:

```shell
go install github.com/kitesi/gpa-calculator@latest
```

Otherwise, go to the releases page and install from there.

## Usage

```shell
$ gpa-calculator [file] [-e|--edit] [-h|--help] [-v|--verbose] [--version] [-u|--unweighted]
```

This program only takes one positional argument, the file/folder to examine. If
not given, it will default to the environment variable `$GRADES_DIR`. If a file
is given, but it is not found, the program will try its best to find it if
`$GRADES_DIR` is set.

```shell
$ echo $GRADES_DIR
/home/seven/grades/

$ gpa-calculator cs100 -v
found file: /home/seven/grades/2023/fall/cs100.grade
└── cs100.grade (98.81) (A)
# ...
```

The edit option will open up a specified file in your editor of choice
(`$EDITOR`). If verbose is on it will display the subsections for a class as
well, like the homework, quizzes, etc. It will also show what you need to get
on the final to match your `desired_grade`. I personally have `gpa` aliased to
`gpa-calculator -v`.

### Syntax Highlighting

Syntax highlighting on vim with tokyonight-night theme:

![syntax highlighting](./syntax-highlighting.png)

If you would like (light) syntax highlighting for the grade files and you use
vim, you can copy `./grade.vim` into `$VIMRUNTIME/syntax/grade.vim`. Then
associate the file extension with the filetype by adding this to your vim
config:

```vim
autocmd BufNewFile,BufRead *.grade setf grade
```

Syntax highlighting for vscode might be implemented in the future.

## Future

-   add configuration
    -   default folder
    -   editor
    -   GPA correlation (A+ = 4.3, A = 4.0)
    -   Grade correlation (A = 93, A-=90, B+=87)
-   consider more grading systems (AP, IP, 1-100 scale)
-   add syntax highlighting to grade file on vscode
-   make easier install
-   colored output
