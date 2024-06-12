import { handleAuthorization, ParamsObject } from "./utils";
import prisma from "@/prisma/client";

export default async function DELETE(req: Request, { params }: ParamsObject) {
    const { error, user } = await handleAuthorization();

    if (error) {
        return error;
    }

    const className = params.className;

    try {
        const prismaClass = await prisma.class.findUnique({
            where: {
                className_userId: {
                    className,
                    userId: user.id,
                },
            },
        });

        if (!prismaClass) {
            return new Response("Class not found", { status: 404 });
        }

        await prisma.gradeSection.deleteMany({
            where: {
                classId: prismaClass.id,
            },
        });

        await prisma.class.delete({
            where: {
                id: prismaClass.id,
            },
        });

        return new Response("OK", {
            status: 200,
        });
    } catch (err) {
        if (!(err instanceof Error)) {
            return new Response("Error deleting grade", {
                status: 500,
            });
        }

        return new Response("Error deleting grade: " + err?.message, {
            status: 500,
        });
    }
}
