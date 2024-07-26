import {
    abstractFormValues,
    handleAuthorization,
    ParamsObject,
    RequestData,
} from "./utils";
import prisma from "@/prisma/client";
import { Prisma } from "@prisma/client";
import { v4 as uuid } from "uuid";
import friendlyErrorMessage from "./friendlyErrorMessage";

export default async function POST(req: Request, { params }: ParamsObject) {
    const { error, user } = await handleAuthorization();

    if (error) {
        return error;
    }

    const input: RequestData = await req.json();
    const classId = uuid();

    const {
        error: formError,
        yearValue,
        desiredGrade,
        credits,
        semester,
        className,
        gradeSections,
    } = abstractFormValues(input, classId);

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
                    create: gradeSections,
                },
            },
        });
    } catch (err) {
        if (!(err instanceof Error)) {
            return new Response("Error creating grade", {
                status: 500,
            });
        }

        if (err instanceof Prisma.PrismaClientKnownRequestError) {
            return new Response(friendlyErrorMessage(err), { status: 400 });
        }

        return new Response("Error creating grade: " + err?.message, {
            status: 500,
        });
    }

    return new Response("OK", {
        status: 200,
    });
}
