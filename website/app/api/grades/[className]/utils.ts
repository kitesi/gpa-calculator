import prisma from "@/prisma/client";
import { auth } from "@/auth";

export interface RequestData {
    year: string;
    semester: string;
    className: string;
    recievedGrade: string;
    desiredGrade: string;
    credits: string;
    gradeSections: {
        className: string;
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
        input.recievedGrade === "F-"
    ) {
        return {
            error: new Response(
                "Invalid recieved grade value should be A-F with possibly +/- but not F-",
                { status: 400 },
            ),
        };
    }

    return {
        error: null,
        yearValue,
        desiredGrade,
        assignedGrade: input.recievedGrade,
        credits,
        semester,
        className,
        gradeSections: input.gradeSections.map((section) => ({
            name: section.name,
            weight: parseFloat(section.weight),
            data: section.data,
            id: section.id,
            classId: classId,
        })),
    };
}
