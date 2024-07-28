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
import { useRouter } from "next/navigation";
import NeedLogin from "@/app/ui/NeedLogin";

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
    loading: boolean;
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
        return <NeedLogin />;
    }

    function addGradeSection() {
        setGradeSections([
            ...gradeSections,
            {
                id: uuid(),
                name: "",
                weight: 0,
                data: "",
                classId: "",
            },
        ]);
    }

    function deleteGradeSection(id: string) {
        setGradeSections(gradeSections.filter((section) => section.id !== id));
    }

    function deleteClass() {
        axios
            .delete("/api/grades/" + props.className)
            .then(() => {
                toast.success("Class deleted!");
                queryClient.invalidateQueries({
                    queryKey: ["gradesData"],
                });

                setIsOpen(false);
                router.push("/");
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
            id: section.id,
            classId: "",
        }));

        if (props.editing) {
            axios
                .put("/api/grades/" + props.className, {
                    year,
                    semester,
                    className,
                    recievedGrade,
                    desiredGrade,
                    credits,
                    gradeSections: gr,
                })
                .then(() => {
                    toast.success("Class edited!");

                    if (
                        props.className !== className ||
                        props.year !== parseInt(year as string) ||
                        props.semester !== semester
                    ) {
                        queryClient.invalidateQueries({
                            queryKey: ["gradesData"],
                        });
                    }

                    router.push("/grades/" + className);
                    queryClient.invalidateQueries({
                        queryKey: ["classData"],
                    });
                })
                .catch((err) => {
                    toast.error("Failed to edit class: " + err?.response?.data);
                    console.log(err);
                });
        } else {
            axios
                .post("/api/grades/" + className, {
                    year,
                    semester,
                    className,
                    recievedGrade,
                    desiredGrade,
                    credits,
                    gradeSections: gr,
                })
                .then(() => {
                    toast.success("Class added!");
                    queryClient.invalidateQueries({
                        queryKey: ["gradesData"],
                    });
                    router.push("/grades/" + className);
                })
                .catch((err) => {
                    toast.error("Failed to add class: " + err?.response?.data);
                    console.log(err);
                });
        }
    }

    const inputClass =
        "mt-3 block w-full rounded-md border-midnight-700 bg-midnight-900 px-3 py-1.5 text-sm/6 text-white focus:outline-none data-[focus]:outline-4 data-[focus]:-outline-offset-2 data-[focus]:outline-blue-500 border-2 shadow-md disabled:text-gray-500";

    return (
        <form
            className={"flex h-full w-full overflow-auto"}
            onSubmit={(ev) => submit(ev)}
        >
            <Transition appear show={isOpen}>
                <Dialog
                    as="div"
                    className="relative z-10 focus:outline-none"
                    onClose={close}
                >
                    <div className="fixed inset-0 z-10 w-screen overflow-y-auto bg-black bg-opacity-50">
                        <div className="flex min-h-full items-center justify-center p-4">
                            <TransitionChild
                                enter="ease-out duration-300"
                                enterFrom="opacity-0 transform-[scale(95%)]"
                                enterTo="opacity-100 transform-[scale(100%)]"
                                leave="ease-in duration-200"
                                leaveFrom="opacity-100 transform-[scale(100%)]"
                                leaveTo="opacity-0 transform-[scale(95%)]"
                            >
                                <DialogPanel className="w-full max-w-md rounded-xl bg-midnight-800 p-6 drop-shadow-2xl">
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
                                            className="mr-3 inline-flex items-center gap-2 rounded-md bg-gray-700 px-3 py-1.5 text-sm/6 font-semibold text-white shadow-inner shadow-white/10 focus:outline-none data-[hover]:bg-gray-600 data-[open]:bg-gray-700 data-[focus]:outline-1 data-[focus]:outline-white"
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
            <Fieldset
                className="mb-4 max-w-lg p-5 md:m-auto"
                disabled={props.loading}
            >
                <Field className="mb-5">
                    <Label className="font-semibold after:ml-0.5 after:text-red-500 after:content-['*']">
                        Class Name{props.editing ? " (edit)" : " (create)"}
                    </Label>
                    <Input
                        className={inputClass}
                        placeholder="CS 101"
                        required
                        name="class-name"
                        defaultValue={props.className}
                        key={props.className}
                        disabled={props.loading}
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
                        className={inputClass}
                        placeholder="A, B, C, D, F, etc..."
                        name="recieved-grade"
                        defaultValue={props.recievedGrade}
                        disabled={props.loading}
                    ></Input>
                </Field>

                <Field className="mb-5">
                    <Label className="font-semibold">Desired Grade (%)</Label>
                    <Input
                        className={inputClass}
                        placeholder="87, 63, 100, etc..."
                        name="desired-grade"
                        defaultValue={props.desiredGrade}
                        disabled={props.loading}
                    ></Input>
                </Field>

                <Field className="mb-5">
                    <Label className="font-semibold after:ml-0.5 after:text-red-500 after:content-['*']">
                        Credits
                    </Label>
                    <Input
                        className={inputClass}
                        placeholder="4"
                        required
                        name="credits"
                        defaultValue={props.credits}
                        key={props.className + props.credits}
                        disabled={props.loading}
                    ></Input>
                </Field>

                <Field className="mb-5">
                    <Label className="font-semibold after:ml-0.5 after:text-red-500 after:content-['*']">
                        Year
                    </Label>
                    <Input
                        className={inputClass}
                        placeholder="2021"
                        required
                        name="year"
                        defaultValue={props.year}
                        key={props.className + props.year}
                        disabled={props.loading}
                    ></Input>
                </Field>

                <Field className="mb-5">
                    <Label className="font-semibold after:ml-0.5 after:text-red-500 after:content-['*']">
                        Semester
                    </Label>
                    <Select
                        className={clsx(
                            inputClass,
                            // Make the text of each option black on Windows
                            // "*:text-black",
                        )}
                        name="semester"
                        disabled={props.loading}
                        defaultValue={props.semester}
                    >
                        <option value="fall">Fall</option>
                        <option value="spring">Spring</option>
                        <option value="winter">Winter</option>
                        <option value="summer">Summer</option>
                    </Select>
                </Field>

                {gradeSections && (
                    <GradeSections
                        gradeSections={gradeSections}
                        inputClass={inputClass}
                        onDeleteId={deleteGradeSection}
                    />
                )}

                <Button
                    onClick={addGradeSection}
                    className={clsx(
                        "rounded-sm border-2 border-dashed border-slate-500 p-5",
                        "focus:outline-none data-[focus]:outline-2 data-[focus]:-outline-offset-2 data-[focus]:outline-blue-500",
                        gradeSections.length === 0 ? "mt-0" : "mt-5",
                    )}
                    disabled={props.loading}
                >
                    Add Grade Section (Homework, Exams, Etc...)
                </Button>

                <div className="my-5 space-x-3">
                    <Button
                        className="rounded-sm bg-my-green px-4 py-2 font-semibold text-white disabled:bg-my-neutral"
                        type="submit"
                        disabled={props.loading}
                    >
                        {props.editing ? "Edit" : "Add"}
                    </Button>

                    {props.editing && (
                        <Button
                            className="rounded-sm bg-my-red px-4 py-2 font-semibold disabled:bg-my-neutral"
                            type="button"
                            onClick={() => openDeleteConfirm()}
                            disabled={props.loading}
                        >
                            Delete
                        </Button>
                    )}
                </div>
            </Fieldset>
        </form>
    );
}
