import {
    abstractFormValues,
    handleAuthorization,
    ParamsObject,
    RequestData,
} from "./utils";
import prisma from "@/prisma/client";
import { v4 as uuid } from "uuid";

export default async function POST(req: Request, { params }: ParamsObject) {
    const { error, user } = await handleAuthorization();

    if (error) {
        return error;
    }

    const input: RequestData = await req.json();
    const {
        error: formError,
        yearValue,
        desiredGrade,
        credits,
        semester,
        className,
    } = abstractFormValues(input);

    if (formError) {
        return formError;
    }

    try {
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

        let prismaSemester = await prisma.semester.findUnique({
            where: {
                name_yearId: {
                    yearId: prismaYear.id,
                    name: semester,
                },
            },
        });

        if (!prismaSemester) {
            prismaSemester = await prisma.semester.create({
                data: {
                    yearValue: yearValue,
                    yearId: prismaYear.id,
                    userId: user.id,
                    name: semester,
                },
            });
        }

        const classId = uuid();

        await prisma.class.create({
            data: {
                id: classId,
                yearId: prismaYear.id,
                userId: user.id,
                semesterId: prismaSemester.id,
                className,
                assignedGrade: input.recievedGrade,
                desiredGrade: desiredGrade,
                credits: credits,
                gradeSections: {
                    create: input.gradeSections.map((section) => ({
                        name: section.name,
                        weight: parseFloat(section.weight),
                        data: section.data,
                        id: section.id,
                        classId: classId,
                    })),
                },
            },
        });
    } catch (err) {
        if (!(err instanceof Error)) {
            return new Response("Error deleting grade", {
                status: 500,
            });
        }

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
