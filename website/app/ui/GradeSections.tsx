import {
    Field,
    Label,
    Description,
    Input,
    Textarea,
    Button,
} from "@headlessui/react";
import { Prisma } from "@prisma/client";

type GradeSection = Prisma.GradeSectionGetPayload<{
    include: {};
}>;

export default function GradeSections({
    gradeSections,
    inputClass,
    onDeleteId,
}: {
    gradeSections: GradeSection[];
    inputClass: string;
    onDeleteId: (id: string) => void;
}) {
    return (
        <>
            {gradeSections.map((gradeSection) => (
                <>
                    <hr className="my-5 border-midnight-700" />

                    <Field key={gradeSection.id + "-name"} className="mb-5">
                        <Label className="font-semibold after:text-red-500 after:content-['*']">
                            Grade Section Name
                        </Label>
                        <Button
                            className="font-inter ml-2 rounded-lg bg-my-red px-2 py-1 text-sm font-semibold text-white"
                            type="button"
                            onClick={() => onDeleteId(gradeSection.id)}
                        >
                            Delete
                        </Button>
                        <Input
                            className={inputClass}
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
                            className={inputClass}
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
                            className={inputClass}
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
