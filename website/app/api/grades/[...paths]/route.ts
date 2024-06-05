import prisma from "@/prisma/client";
import { auth } from "@/auth";

export async function getData(
    year: string,
    semester: string,
    className: string,
) {
    return await prisma.class.findUnique({
        where: {
            yearValue: parseInt(year),
            className: className,
            semesterId: semester,
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

export async function GET(
    req: Request,
    { params }: { params: { paths: string[] } },
) {
    const session = await auth();
    if (!session) {
        return new Response("Unauthorized, please log in", { status: 401 });
    }

    const year = params.paths[0];
    const semester = params.paths[1];
    const className = params.paths[2];

    if (!year || !semester || !className) {
        return new Response(
            "Invalid request, missing one of: year, semester, className",
            { status: 400 },
        );
    }

    try {
        return Response.json(await getData(year, semester, className));
    } catch (err: any) {
        let errMessage = "";

        if (err?.message) {
            errMessage = err.message;
        }

        if (err?.body?.message) {
            errMessage = err.body.message;
        }

        return new Response("Error retrieving grade: " + errMessage, {
            status: 500,
        });
    }
}

interface RequestData {
    recievedGrade: string;
    desiredGrade: string;
    credits: string;
    gradeSections: {
        name: string;
        weight: string;
        data: string;
        className: string;
    }[];
}

export async function POST(
    req: Request,
    { params }: { params: { paths: string[] } },
) {
    const session = await auth();

    if (!session || !session.user || !session.user.email) {
        return new Response("Unauthorized, please log in", { status: 401 });
    }

    const input: RequestData = await req.json();

    if (!input) {
        return new Response("Invalid request, missing body", { status: 400 });
    }

    if (!input.credits) {
        return new Response(
            "Invalid request, missing one of: recievedGrade, desiredGrade, credits",
            { status: 400 },
        );
    }

    const [year, semester, className] = params.paths;

    if (!year || !semester || !className) {
        return new Response(
            "Invalid request, missing one of: year, semester, className",
            { status: 400 },
        );
    }

    const yearValue = parseInt(year);
    const desiredGrade = input.desiredGrade
        ? parseFloat(input.desiredGrade)
        : 0;
    const credits = parseInt(input.credits);

    if (isNaN(yearValue)) {
        return new Response("Invalid year value", { status: 400 });
    }

    if (input.desiredGrade != "" && isNaN(desiredGrade)) {
        return new Response("Invalid desired grade value", { status: 400 });
    }

    if (isNaN(credits)) {
        return new Response("Invalid credits value", { status: 400 });
    }

    // if assigned grade is not A-F with +/-
    if (
        (input.recievedGrade !== "" &&
            !/^[A-F][+-]?$/.test(input.recievedGrade)) ||
        input.recievedGrade === "F-"
    ) {
        return new Response(
            "Invalid recieved grade value should be A-F with possibly +/- but not F-",
            { status: 400 },
        );
    }

    const prismaUser = await prisma.user.findUnique({
        where: { email: session.user.email },
    });

    if (!prismaUser) {
        return new Response("User not found", { status: 400 });
    }

    try {
        let prismaYear = await prisma.year.findUnique({
            where: { yearValue: yearValue, userId: prismaUser.id },
        });

        if (!prismaYear) {
            await prisma.year.create({
                data: {
                    yearValue: yearValue,
                    userId: prismaUser.id,
                },
            });
        }

        const prismaSemester = await prisma.semester.findUnique({
            where: { id: yearValue + "-" + semester },
        });

        if (!prismaSemester) {
            await prisma.semester.create({
                data: {
                    id: yearValue + "-" + semester,
                    yearValue: yearValue,
                    userId: prismaUser.id,
                    name: semester,
                },
            });
        }

        await prisma.class.create({
            data: {
                yearValue,
                semesterId: yearValue + "-" + semester,
                className,
                assignedGrade: input.recievedGrade,
                desiredGrade: desiredGrade,
                credits: credits,
                userId: prismaUser.id,
                gradeSections: {
                    create: input.gradeSections.map((section) => ({
                        name: section.name,
                        weight: parseFloat(section.weight),
                        data: section.data,
                    })),
                },
            },
        });
    } catch (err) {
        return new Response("Error creating grade: " + err?.message, {
            status: 500,
        });
    }

    return new Response(
        "POST request to /api/grades/" + params.paths.join("/"),
        {
            status: 200,
        },
    );
}
