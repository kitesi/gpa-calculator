import prisma from "@/prisma/client";
import { auth } from "@/auth";
import { Session } from "next-auth";

export async function getData(session: Session) {
    return await prisma.year.findMany({
        where: { userId: session.user?.id },
        include: {
            semesters: {
                include: {
                    classes: {
                        include: {
                            gradeSections: {
                                include: {},
                            },
                        },
                    },
                },
            },
        },
        orderBy: { yearValue: "asc" },
    });
}

export async function GET() {
    const session = await auth();

    if (!session) {
        return new Response("Unauthorized, please log in", { status: 401 });
    }

    try {
        return Response.json(await getData(session));
    } catch (err: any) {
        let errMessage = "";

        if (err?.message) {
            errMessage = err.message;
        }

        if (err?.body?.message) {
            errMessage = err.body.message;
        }

        console.error(err);

        return new Response("Error getting grades: " + errMessage, {
            status: 500,
        });
    }
}
