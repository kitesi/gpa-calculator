"use client";
import {
    Button,
    Description,
    Field,
    Fieldset,
    Input,
    Label,
    Select,
    Dialog,
    DialogPanel,
    DialogTitle,
    Transition,
    TransitionChild,
} from "@headlessui/react";
import GradeSections from "./GradeSections";
import clsx from "clsx";
import { FormEvent, useState } from "react";
import { v4 as uuid } from "uuid";
import { Prisma } from "@prisma/client";
import axios from "axios";
import { useSession } from "next-auth/react";
import toast from "react-hot-toast";
import { useQueryClient } from "@tanstack/react-query";
import { redirect, useRouter } from "next/navigation";

type GradeSection = Prisma.GradeSectionGetPayload<{
    include: {};
}>;

interface Props {
    credits: number;
    year: number;
    semester: string;
    gradeSections: GradeSection[];
    className: string;
    recievedGrade: string;
    desiredGrade: string;
    editing: boolean;
}

export default function AddEditClassForm(props: Props) {
    const [gradeSections, setGradeSections] = useState<GradeSection[]>(
        props.gradeSections,
    );
    const [isOpen, setIsOpen] = useState(false);
    const router = useRouter();

    function openDeleteConfirm() {
        setIsOpen(true);
    }

    function close() {
        setIsOpen(false);
    }

    const { data: session } = useSession();
    const queryClient = useQueryClient();

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
            {
                id: uuid(),
                className: "",
                name: "",
                weight: 0,
                data: "",
                classId: "",
            },
        ]);
    }

    function deleteClass() {
        axios
            .delete(
                "/api/grades/" +
                    props.year +
                    "/" +
                    props.semester +
                    "/" +
                    props.className,
            )
            .then(() => {
                toast.success("Class deleted!");
                queryClient.invalidateQueries({
                    queryKey: ["gradesData"],
                });

                setIsOpen(false);
                router.push("/grades");
            })
            .catch((err) => {
                toast.error("Failed to delete class: " + err?.response?.data);
                console.log(err);
            });
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

        if (props.editing) {
            axios
                .put(
                    "/api/grades/" +
                        props.year +
                        "/" +
                        props.semester +
                        "/" +
                        props.className,
                    {
                        year,
                        semester,
                        className,
                        recievedGrade,
                        desiredGrade,
                        credits,
                        gradeSections: gr,
                    },
                )
                .then(() => {
                    toast.success("Class edited!");
                    queryClient.invalidateQueries({
                        queryKey: ["gradesData"],
                    });
                    router.push(
                        "/grades/" + year + "/" + semester + "/" + className,
                    );
                })
                .catch((err) => {
                    toast.error("Failed to edit class: " + err?.response?.data);
                    console.log(err);
                });
        } else {
            axios
                .post(
                    "/api/grades/" + year + "/" + semester + "/" + className,
                    {
                        year,
                        semester,
                        className,
                        recievedGrade,
                        desiredGrade,
                        credits,
                        gradeSections: gr,
                    },
                )
                .then(() => {
                    toast.success("Class added!");
                    queryClient.invalidateQueries({
                        queryKey: ["gradesData"],
                    });
                })
                .catch((err) => {
                    toast.error("Failed to add class: " + err?.response?.data);
                    console.log(err);
                });
        }
    }

    return (
        <form
            className={"h-full w-full overflow-scroll p-10"}
            onSubmit={(ev) => submit(ev)}
        >
            <Transition appear show={isOpen}>
                <Dialog
                    as="div"
                    className="relative z-10 focus:outline-none"
                    onClose={close}
                >
                    <div className="fixed inset-0 z-10 w-screen overflow-y-auto">
                        <div className="flex min-h-full items-center justify-center p-4">
                            <TransitionChild
                                enter="ease-out duration-300"
                                enterFrom="opacity-0 transform-[scale(95%)]"
                                enterTo="opacity-100 transform-[scale(100%)]"
                                leave="ease-in duration-200"
                                leaveFrom="opacity-100 transform-[scale(100%)]"
                                leaveTo="opacity-0 transform-[scale(95%)]"
                            >
                                <DialogPanel className="w-full max-w-md rounded-xl bg-black p-6">
                                    <DialogTitle
                                        as="h3"
                                        className="text-base/7 font-medium text-white"
                                    >
                                        Confirm Deletion
                                    </DialogTitle>
                                    <p className="mt-2 text-sm/6 text-white/50">
                                        Are you sure you want to delete this
                                        class? This action cannot be undone.
                                        Please double check all of the details
                                    </p>
                                    <div className="mt-4">
                                        <Button
                                            className="inline-flex items-center gap-2 rounded-md bg-gray-700 px-3 py-1.5 text-sm/6 font-semibold text-white shadow-inner shadow-white/10 focus:outline-none data-[hover]:bg-gray-600 data-[open]:bg-gray-700 data-[focus]:outline-1 data-[focus]:outline-white"
                                            onClick={deleteClass}
                                        >
                                            Delete
                                        </Button>
                                        <Button
                                            className="inline-flex items-center gap-2 rounded-md bg-gray-700 px-3 py-1.5 text-sm/6 font-semibold text-white shadow-inner shadow-white/10 focus:outline-none data-[hover]:bg-gray-600 data-[open]:bg-gray-700 data-[focus]:outline-1 data-[focus]:outline-white"
                                            onClick={() => setIsOpen(false)}
                                        >
                                            Close
                                        </Button>
                                    </div>
                                </DialogPanel>
                            </TransitionChild>
                        </div>
                    </div>
                </Dialog>
            </Transition>
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
                        defaultValue={props.className}
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
                        defaultValue={props.recievedGrade}
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
                        defaultValue={props.desiredGrade}
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
                        defaultValue={props.credits}
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
                        defaultValue={props.year}
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
                        <option
                            value="fall"
                            selected={props.semester === "fall"}
                        >
                            Fall
                        </option>
                        <option
                            value="spring"
                            selected={props.semester === "spring"}
                        >
                            Spring
                        </option>
                        <option
                            value="winter"
                            selected={props.semester === "winter"}
                        >
                            Winter
                        </option>
                        <option
                            value="summer"
                            selected={props.semester === "summer"}
                        >
                            Summer
                        </option>
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
                        {props.editing ? "Edit" : "Add"}
                    </Button>

                    {props.editing && (
                        <Button
                            className="rounded-lg bg-red-500 px-4 py-2 font-semibold"
                            type="button"
                            onClick={() => openDeleteConfirm()}
                        >
                            Delete
                        </Button>
                    )}
                </div>
            </Fieldset>
        </form>
    );
}
