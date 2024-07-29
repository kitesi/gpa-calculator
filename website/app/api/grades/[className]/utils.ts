import prisma from "@/prisma/client";
import { auth } from "@/auth";
import { GradeSection } from "@prisma/client";
import { Prisma } from "@prisma/client";

export interface RequestData {
    year: string;
    semester: string;
    className: string;
    recievedGrade: string;
    desiredGrade: string;
    credits: string;
    gradeSections: {
        classId: string;
        name: string;
        weight: string;
        data: string;
        id: string;
    }[];
}

export interface ParamsObject {
    params: {
        className: string;
    };
}

export async function getClassData(className: string, userId: string) {
    return await prisma.class.findUnique({
        where: {
            className_userId: {
                className: className,
                userId,
            },
        },
        include: {
            gradeSections: {
                include: {},
            },
            year: true,
            semester: true,
        },
    });
}

export async function handleAuthorization() {
    const session = await auth();

    if (!session || !session.user || !session.user.email) {
        return {
            error: new Response("Unauthorized, please log in", { status: 401 }),
            session,
        };
    }

    const user = await prisma.user.findUnique({
        where: { email: session.user.email },
    });

    if (!user) {
        return {
            error: new Response("User not found", { status: 400 }),
            session,
        };
    }

    return { error: null, session, user };
}

export function abstractYearSemesterClass(
    year: string,
    semester: string,
    className: string,
) {
    const yearValue = parseInt(year);

    if (isNaN(yearValue)) {
        return {
            error: new Response("Invalid year value", { status: 400 }),
        };
    }

    if (
        semester !== "spring" &&
        semester !== "fall" &&
        semester !== "summer" &&
        semester !== "winter"
    ) {
        return {
            error: new Response("Invalid semester value", { status: 400 }),
        };
    }

    if (className === "") {
        return {
            error: new Response("Missing class name", { status: 400 }),
        };
    }

    return {
        yearValue,
        semester,
        className,
    };
}

export function abstractFormValues(input: RequestData, classId: string) {
    if (!input) {
        return {
            error: new Response("Invalid request, missing body", {
                status: 400,
            }),
        };
    }

    if (!input.credits) {
        return {
            error: new Response(
                "Invalid request, missing one of: recievedGrade, desiredGrade, credits",
                { status: 400 },
            ),
        };
    }

    const desiredGrade = input.desiredGrade
        ? parseFloat(input.desiredGrade)
        : 0;
    const credits = parseInt(input.credits);
    const { error, className, yearValue, semester } = abstractYearSemesterClass(
        input.year,
        input.semester,
        input.className,
    );

    if (error) {
        return { error };
    }

    if (
        input.desiredGrade != "" &&
        (isNaN(desiredGrade) || desiredGrade < 0 || desiredGrade > 100)
    ) {
        return {
            error: new Response("Invalid desired grade value", { status: 400 }),
        };
    }

    if (isNaN(credits) || credits < 0) {
        return {
            error: new Response("Invalid credits value", { status: 400 }),
        };
    }

    // if assigned grade is not A-F with +/-
    if (
        (input.recievedGrade !== "" &&
            !/^[A-F][+-]?$/.test(input.recievedGrade)) ||
        input.recievedGrade === "F-" ||
        input.recievedGrade === "F+"
    ) {
        return {
            error: new Response(
                "Invalid recieved grade value should be A-F with possibly +/- but not F-",
                { status: 400 },
            ),
        };
    }

    const gradeSections: Omit<GradeSection, "classId">[] = [];

    for (const section of input.gradeSections) {
        const { error } = parseGradeSectionData(section.data);

        if (error) {
            return {
                error: new Response(error, { status: 400 }),
            };
        }

        const weight = parseFloat(section.weight);

        if (isNaN(weight) || weight < 0 || weight > 100) {
            return {
                error: new Response("Invalid weight value", { status: 400 }),
            };
        }

        gradeSections.push({
            name: section.name,
            weight,
            data: section.data,
            id: section.id,
        });
    }

    return {
        error: null,
        yearValue,
        desiredGrade,
        assignedGrade: input.recievedGrade,
        credits,
        semester,
        className,
        gradeSections,
    };
}

export function getFriendlyErrorMessage(
    err: Prisma.PrismaClientKnownRequestError,
) {
    switch (err.code) {
        case "P2002":
            return "Class with that name already exists.";
        case "P2025":
            return "Class not found.";
        default:
            return err.message;
    }
}

export function parseGradeSectionData(data: string): {
    error?: string;
    value?: {
        recieved: number;
        total: number;
    }[];
} {
    const lines = data.split("\n");
    const value: {
        recieved: number;
        total: number;
    }[] = [];

    for (let line of lines) {
        let commentIndex = line.indexOf("#");

        if (commentIndex !== -1) {
            line = line.slice(0, commentIndex);
        }

        const scores = line.trim().split(",");

        for (const score of scores) {
            if (score === "") {
                continue;
            }

            const scoreFractions = score.split("/");

            if (scoreFractions.length != 2) {
                return {
                    error: `Invalid score value: ${score}`,
                };
            }

            const numerator = parseFloat(scoreFractions[0]);
            const denominator = parseFloat(scoreFractions[1]);

            if (
                !denominator ||
                isNaN(denominator) ||
                !numerator ||
                isNaN(numerator)
            ) {
                return {
                    error: `Invalid score value: ${score}`,
                };
            }

            value.push({
                recieved: numerator,
                total: denominator,
            });
        }
    }

    return { value: value };
}

export function getGradeGPA(grade: string) {
    switch (grade) {
        case "A+":
            return 4.3;
        case "A":
            return 4;
        case "A-":
            return 3.7;
        case "B+":
            return 3.3;
        case "B":
            return 3;
        case "B-":
            return 2.7;
        case "C+":
            return 2.3;
        case "C":
            return 2.0;
        case "C-":
            return 1.7;
        case "D+":
            return 1.3;
        case "D":
            return 1.0;
        case "D-":
            return 0.7;
        case "F":
            return 0;
    }

    return -1;
}

export function getGradeLetter(grade: number): string {
    grade = grade * 100;

    if (grade >= 94) {
        return "A";
    } else if (grade >= 90) {
        return "A-";
    } else if (grade >= 87) {
        return "B+";
    } else if (grade >= 84) {
        return "B";
    } else if (grade >= 80) {
        return "B-";
    } else if (grade >= 77) {
        return "C+";
    } else if (grade >= 74) {
        return "C";
    } else if (grade >= 70) {
        return "C-";
    } else if (grade >= 67) {
        return "D+";
    } else if (grade >= 64) {
        return "D";
    } else if (grade >= 60) {
        return "D-";
    } else {
        return "F";
    }
}
