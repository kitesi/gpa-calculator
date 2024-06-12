import { GradeSection } from "@prisma/client";
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

    try {
        const input: RequestData = await req.json();

        const oldClass = await prisma.class.findUnique({
            where: {
                className_userId: {
                    className: params.className,
                    userId: user.id,
                },
            },
            include: {
                semester: true,
                year: true,
                gradeSections: true,
            },
        });

        if (!oldClass) {
            return new Response("Class not found", { status: 404 });
        }

        const oldClassName = oldClass.className;
        const oldSemester = oldClass.semester.name;
        const oldYearValue = oldClass.year.yearValue;

        const {
            error: formError,
            yearValue,
            assignedGrade,
            className,
            credits,
            desiredGrade,
            semester,
            gradeSections,
        } = abstractFormValues(input, oldClass.id);

        if (formError) {
            return formError;
        }

        let prismaYearId = oldClass.yearId;
        let prismaSemesterId = oldClass.semesterId;

        // // if the year is different and we don't have it made yet we need to create a new year
        if (yearValue !== oldYearValue) {
            prismaYearId = (
                await prisma.year.upsert({
                    where: {
                        yearValue_userId: {
                            yearValue: yearValue,
                            userId: user.id,
                        },
                    },

                    update: {},
                    create: {
                        yearValue: yearValue,
                        userId: user.id,
                    },
                })
            ).id;
        }

        if (yearValue !== oldYearValue || semester !== oldSemester) {
            prismaSemesterId = (
                await prisma.semester.upsert({
                    where: {
                        name_yearId: {
                            name: semester,
                            yearId: prismaYearId,
                        },
                    },
                    update: {},
                    create: {
                        name: semester,
                        userId: user.id,
                        yearId: prismaYearId,
                        yearValue: yearValue,
                    },
                })
            ).id;
        }

        const toDelete: string[] = [];
        const toCreate: GradeSection[] = [];
        const toUpdate: ReturnType<typeof prisma.gradeSection.update>[] = [];

        for (const gs of oldClass.gradeSections) {
            if (!gradeSections.some((s) => s.id === gs.id)) {
                toDelete.push(gs.id);
            }
        }

        for (const section of gradeSections) {
            const oldSection = oldClass.gradeSections.find(
                (s) => s.id === section.id,
            );

            if (!oldSection) {
                toCreate.push(section);
            } else if (
                oldSection.name !== section.name ||
                oldSection.weight !== section.weight ||
                oldSection.data !== section.data
            ) {
                toUpdate.push(
                    prisma.gradeSection.update({
                        where: { id: oldSection.id },
                        data: section,
                    }),
                );
            }
        }

        await prisma.$transaction([
            prisma.gradeSection.createMany({
                data: toCreate,
            }),
            prisma.gradeSection.deleteMany({
                where: {
                    id: {
                        in: toDelete,
                    },
                },
            }),
            ...toUpdate,
            prisma.class.update({
                where: {
                    className_userId: {
                        className: oldClassName,
                        userId: user.id,
                    },
                },
                data: {
                    yearId: prismaYearId,
                    semesterId: prismaSemesterId,
                    assignedGrade,
                    desiredGrade,
                    className,
                    credits,
                },
            }),
        ]);

        return new Response("OK", {
            status: 200,
        });
    } catch (err) {
        if (!(err instanceof Error)) {
            return new Response("Error updating grade", {
                status: 500,
            });
        }

        return new Response("Error updating grade: " + err?.message, {
            status: 500,
        });
    }
}
