# GPA Calculator

This is a simple grade calculator written in go.
It gathers the data from local files written in
a special but simple format.

## Calculation of Grades & GPA

-   GPA ranges from 0.0 to 4.0
-   A is the highest grades at 94%
-   Read `./grade-conversions.go` for full details

There is no config to change the calculation currently, so I would recommend
modifying the source code yourself.

## Structure of Folder

I recommend having a folder with all your grades, with the structure as such:

```
grades/
├── 2022
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
└── 2023
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

The files, which contain the data for each class, all have the extension "grade", but it doesn't matter what extension you use.

### Grade File

The grade file is just a simple text file with specific syntax.
Grade categories and denoted with a ">" prefix and then their name,
and the children elements are the `weight` and `data`

```
> Homework
weight = 0.2
data = 20/20, 17/20, 19/20

> Midterm
weight = 0.4
data = 40/50

> Final Exam
weight = 0.4
data = 58/60
```

Comments can only be declared with a line that starts with a pound symbol (#).
Indentation and white space do not matter.

If you wish to specify options for the entire grade file, use the line "~ Meta"
at the start of your file. For example, each class is assumed to be 4 credits,
if you wish to change this use the credits option (there are 3 other recognized
options, but you can add whatever you like).

```
~ Meta
# amount of credits as an int
credits = 3

# nice name for when printing
name = "MA160 Multivariable Calculus & Series"

# your desired grade as a float from 0-100
desired_grade = 94

# describe your actual grade with a letter. this is useful if you recieve an A+,
# because the program doesn't assign A+ automatically (in my school and many, the
# professor has to assign an A+)

grade = "A+"

# unrecognized option but is ok
location = Grand Hall 202

# homework given every thursday
> Homework
# ...snip...
```

## Output

This is the general output for now:

```
├── 2022
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
└── 2023
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

## Output

This is the general and only output for now assuming you do `gpa grades/`:

```shell
$ gpa-calculator ~/grades-1
grades-1 (3.11)
├── 2022 (3.12)
│   ├── fall (3.08)
│   │   ├── cs100.grade (85.16) [B]
│   │   ├── gov100.grade (87.15) [B+]
│   │   ├── lang100.grade (92.59) [A-]
│   │   └── ma100.grade (78.31) [C+]
│   └── spring (3.17)
│       ├── ma200.grade (86.79) [B]
│       ├── cs200.grade (96.01) [A]
│       ├── gov200.grade (84.80) [B]
│       └── lang200.grade (81.19) [B-]
└── 2023 (3.09)
    ├── fall (2.83)
    │   ├── gov300.grade (79.33) [C+]
    │   ├── lang300.grade (85.12) [B]
    │   ├── ma300.grade (85.11) [B]
    │   └── cs300.grade (86.40) [B]
    └── spring (3.35)
        ├── ma400.grade (81.66) [B-]
        ├── cs400.grade (92.93) [A-]
        ├── gov400.grade (86.98) [B]
        └── lang400.grade (94.66) [A]
```

## Installation

TODO

## Usage

```shell
$ gpa-calculator <folder> [-h|--help] [-v|--verbose]
```

If verbose is on it will display the subsections for a class as well like the homework, quizzes, etc.
I personally have `gpa` aliased to `gpa-calculator $my_folder`

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

-   add root project config file like ~/grades/config.json
-   add tests
-   add syntax highlighting to grade file on vscode
