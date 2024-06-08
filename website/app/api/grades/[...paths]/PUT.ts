import {
    handleAuthorization,
    ParamsObject,
    RequestData,
    abstractFormValues,
    abstractYearSemesterClass,
} from "./utils";
import prisma from "@/prisma/client";

export default async function PUT(req: Request, { params }: ParamsObject) {
    const { error, user } = await handleAuthorization();

    if (error) {
        return error;
    }

    const input: RequestData = await req.json();
    const [y, oldSemester, oldClassName] = params.paths;
    const { error: errorFromPath, yearValue: oldYear } =
        abstractYearSemesterClass(y, oldSemester, oldClassName);

    if (errorFromPath) {
        return errorFromPath;
    }

    const oldClass = await prisma.class.findUnique({
        where: {
            className_userId: {
                className: params.paths[2],
                userId: user.id,
            },
        },
        include: {
            semester: true,
            year: true,
        },
    });

    if (!oldClass) {
        return new Response("Class not found", { status: 404 });
    }

    if (
        oldClass.semester.name !== oldSemester ||
        oldClass.year.yearValue !== oldYear
    ) {
        return new Response("Year or Semester doesn't match", { status: 401 });
    }

    const {
        error: formError,
        yearValue,
        assignedGrade,
        className,
        credits,
        desiredGrade,
        semester,
    } = abstractFormValues(input);

    if (formError) {
        return formError;
    }

    let prismaYearId = oldClass.yearId;

    // if the year is different and we don't have it made yet we need to create a new year
    if (yearValue !== oldYear) {
        let prismaYear = await prisma.year.findUnique({
            where: {
                yearValue_userId: {
                    yearValue: yearValue,
                    userId: user.id,
                },
            },
        });

        if (!prismaYear) {
            prismaYear = await prisma.year.create({
                data: {
                    yearValue: yearValue,
                    userId: user.id,
                },
            });
        }

        prismaYearId = prismaYear.id;
    }

    let prismaSemesterId = oldClass.semesterId;

    if (semester !== oldSemester) {
        let prismaSemester = await prisma.semester.findUnique({
            where: {
                name_yearId: {
                    name: input.semester,
                    yearId: prismaYearId,
                },
            },
        });

        if (!prismaSemester) {
            prismaSemester = await prisma.semester.create({
                data: {
                    name: input.semester,
                    userId: user.id,
                    yearId: prismaYearId,
                    yearValue: yearValue,
                },
            });
        }

        prismaSemesterId = prismaSemester.id;
    }

    prisma.class.update({
        where: {
            className_userId: {
                className: params.paths[2],
                userId: user.id,
            },
        },
        data: {
            yearId: prismaYearId,
            semesterId: prismaSemesterId,
            assignedGrade,
            className,
            credits,
            desiredGrade,
            gradeSections: {},
        },
    });

    return new Response(
        "PUT request to /api/grades/" + params.paths.join("/"),
        {
            status: 200,
        },
    );
}
