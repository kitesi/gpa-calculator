import {
    Button,
    Field,
    Label,
    Description,
    Input,
    Textarea,
} from "@headlessui/react";
import { Prisma } from "@prisma/client";
import clsx from "clsx";

type GradeSection = Prisma.GradeSectionGetPayload<{
    include: {};
}>;

export default function GradeSections({
    gradeSections,
}: {
    gradeSections: GradeSection[];
}) {
    return (
        <>
            {gradeSections.map((gradeSection) => (
                <>
                    <hr className="my-5 border-slate-500" />

                    <Field key={gradeSection.id + "-name"} className="mb-5">
                        <Label className="font-semibold after:text-red-500 after:content-['*']">
                            Grade Section Name
                        </Label>
                        <Input
                            className={clsx(
                                "mt-3 block w-full rounded-lg border-none bg-white px-3 py-1.5 text-sm/6 text-black",
                                "focus:outline-none data-[focus]:outline-2 data-[focus]:-outline-offset-2 data-[focus]:outline-white/25",
                            )}
                            placeholder="Homework, Exams, etc..."
                            required
                            name={gradeSection.id + "-section-name"}
                            defaultValue={gradeSection.name}
                        ></Input>
                    </Field>

                    <Field key={gradeSection.id + "-weight"} className="mb-5">
                        <Label className="font-semibold after:text-red-500 after:content-['*']">
                            Weight (%)
                        </Label>
                        <Input
                            className={clsx(
                                "mt-3 block w-full rounded-lg border-none bg-white px-3 py-1.5 text-sm/6 text-black",
                                "focus:outline-none data-[focus]:outline-2 data-[focus]:-outline-offset-2 data-[focus]:outline-white/25",
                            )}
                            placeholder="20"
                            required
                            name={gradeSection.id + "-section-weight"}
                            defaultValue={gradeSection.weight.toString()}
                        ></Input>
                    </Field>

                    <Field key={gradeSection.id + "-add-data"}>
                        <Label className="font-semibold">Add Data (x/y)</Label>

                        <Description className="mb-3 mt-2 text-sm/6 text-slate-300">
                            Add data for this grade section in the format of
                            x/y, where x is the grade recieved and y is the
                            total possible points. You can input multiple at a
                            time with a comma separating them.
                        </Description>

                        <Textarea
                            className={clsx(
                                "block w-full flex-1 rounded-lg border-none bg-white px-3 py-1.5 text-sm/6 text-black",
                                "focus:outline-none data-[focus]:outline-2 data-[focus]:-outline-offset-2 data-[focus]:outline-white/25",
                            )}
                            rows={4}
                            placeholder={`10/10,
5/6,10/10,10/10,
7/10,
100/100`}
                            name={gradeSection.id + "-section-data"}
                            defaultValue={gradeSection.data}
                        />
                    </Field>
                </>
            ))}
        </>
    );
}
