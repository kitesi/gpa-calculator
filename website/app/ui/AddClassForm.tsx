"use client";
import {
    Button,
    Description,
    Field,
    Fieldset,
    Input,
    Label,
    Legend,
    Select,
    Textarea,
} from "@headlessui/react";
import GradeSections from "./GradeSections";
import clsx from "clsx";
import { FormEvent, useState } from "react";
import { v4 as uuid } from "uuid";
import { Prisma } from "@prisma/client";
import axios from "axios";
import { useSession } from "next-auth/react";
import toast from "react-hot-toast";

type GradeSection = Prisma.GradeSectionGetPayload<{
    include: {};
}>;

interface Props {
    credits: number;
    year: number;
}

export default function AddClassForm(props: Props) {
    const [gradeSections, setGradeSections] = useState<GradeSection[]>([]);
    const { data: session } = useSession();

    if (!session) {
        return (
            <p className="p-4 font-semibold">
                You must be logged in to add a class.
            </p>
        );
    }

    function add() {
        setGradeSections([
            ...gradeSections,
            { id: uuid(), className: "", name: "", weight: 0, data: "" },
        ]);
    }

    function submit(ev: FormEvent<HTMLFormElement>) {
        ev.preventDefault();

        const form = new FormData(ev.currentTarget);
        const className = form.get("class-name");
        const recievedGrade = form.get("recieved-grade");
        const desiredGrade = form.get("desired-grade");

        const credits = form.get("credits");
        const year = form.get("year");

        const semester = form.get("semester");
        const gr = gradeSections.map((section) => ({
            name: form.get(section.id + "-section-name"),
            weight: form.get(section.id + "-section-weight"),
            data: form.get(section.id + "-section-data"),
            className: className,
        }));

        axios
            .post("/api/grades/" + year + "/" + semester + "/" + className, {
                recievedGrade,
                desiredGrade,
                credits,
                gradeSections: gr,
            })
            .then(() => toast.success("Class added!"))
            .catch((err) => {
                toast.error("Failed to add class: " + err?.response?.data);
                console.log(err);
            });
    }

    return (
        <form
            className={"h-full w-full overflow-scroll p-10"}
            onSubmit={(ev) => submit(ev)}
        >
            <Fieldset className="max-w-md bg-slate-800 p-5">
                <Field className="mb-5">
                    <Label className="font-semibold after:ml-0.5 after:text-red-500 after:content-['*']">
                        Class Name
                    </Label>
                    <Input
                        className={clsx(
                            "mt-3 block w-full rounded-lg border-none bg-white px-3 py-1.5 text-sm/6 text-black",
                            "focus:outline-none data-[focus]:outline-4 data-[focus]:-outline-offset-2 data-[focus]:outline-blue-500",
                        )}
                        placeholder="CS 101"
                        required
                        name="class-name"
                    ></Input>
                </Field>

                <Field className="mb-5">
                    <Label className="font-semibold">
                        Recieved Grade (Letter)
                    </Label>
                    <Description className="mt-2 text-sm text-slate-300">
                        If you have already recieved a grade for this class
                        enter it here, no grade sections will be considered if
                        this is filled out.
                    </Description>
                    <Input
                        className={clsx(
                            "mt-3 block w-full rounded-lg border-none bg-white px-3 py-1.5 text-sm/6 text-black",
                            "focus:outline-none data-[focus]:outline-4 data-[focus]:-outline-offset-2 data-[focus]:outline-blue-500",
                        )}
                        placeholder="A, B, C, D, F, etc..."
                        name="recieved-grade"
                    ></Input>
                </Field>

                <Field className="mb-5">
                    <Label className="font-semibold">Desired Grade (%)</Label>
                    <Input
                        className={clsx(
                            "mt-3 block w-full rounded-lg border-none bg-white px-3 py-1.5 text-sm/6 text-black",
                            "focus:outline-none data-[focus]:outline-4 data-[focus]:-outline-offset-2 data-[focus]:outline-blue-500",
                        )}
                        placeholder="87, 63, 100, etc..."
                        name="desired-grade"
                    ></Input>
                </Field>

                <Field className="mb-5">
                    <Label className="font-semibold after:ml-0.5 after:text-red-500 after:content-['*']">
                        Credits
                    </Label>
                    <Input
                        className={clsx(
                            "mt-3 block w-full rounded-lg border-none bg-white px-3 py-1.5 text-sm/6 text-black",
                            "focus:outline-none data-[focus]:outline-4 data-[focus]:-outline-offset-2 data-[focus]:outline-blue-500",
                        )}
                        placeholder="4"
                        required
                        name="credits"
                        value={props.credits}
                    ></Input>
                </Field>

                <Field className="mb-5">
                    <Label className="font-semibold after:ml-0.5 after:text-red-500 after:content-['*']">
                        Year
                    </Label>
                    <Input
                        className={clsx(
                            "mt-3 block w-full rounded-lg border-none bg-white px-3 py-1.5 text-sm/6 text-black",
                            "focus:outline-none data-[focus]:outline-4 data-[focus]:-outline-offset-2 data-[focus]:outline-blue-500",
                        )}
                        placeholder="2021"
                        required
                        name="year"
                        value={props.year}
                    ></Input>
                </Field>

                <Field className="mb-5">
                    <Label className="font-semibold after:ml-0.5 after:text-red-500 after:content-['*']">
                        Semester
                    </Label>
                    <Select
                        className={clsx(
                            "mt-3 block w-full rounded-lg border-none bg-white px-3 py-1.5 text-sm/6 text-black",
                            "focus:outline-none data-[focus]:outline-4 data-[focus]:-outline-offset-2 data-[focus]:outline-blue-500",
                            // Make the text of each option black on Windows
                            "*:text-black",
                        )}
                        name="semester"
                    >
                        <option value="fall">Fall</option>
                        <option value="spring">Spring</option>
                        <option value="winter">Winter</option>
                        <option value="summer">Summer</option>
                    </Select>
                </Field>

                <GradeSections gradeSections={gradeSections} />

                <Button
                    onClick={add}
                    className={clsx(
                        "rounded-md border-2 border-dashed border-slate-500 p-5",
                        "focus:outline-none data-[focus]:outline-2 data-[focus]:-outline-offset-2 data-[focus]:outline-blue-500",
                        gradeSections.length === 0 ? "mt-0" : "mt-5",
                    )}
                >
                    Add Grade Section (Homework, Exams, Etc...)
                </Button>

                <div className="mt-5 space-x-3">
                    <Button
                        className="rounded-lg bg-green-500 px-4 py-2 font-semibold text-black"
                        type="submit"
                    >
                        Submit
                    </Button>
                    <Button
                        className="rounded-lg bg-red-500 px-4 py-2 font-semibold"
                        type="submit"
                    >
                        Delete
                    </Button>
                </div>
            </Fieldset>
        </form>
    );
}
