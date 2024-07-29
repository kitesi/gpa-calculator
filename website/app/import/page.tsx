"use client";
import {
    Button,
    Description,
    Field,
    Fieldset,
    Input,
    Label,
    Select,
} from "@headlessui/react";
import axios from "axios";
import { useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import toast from "react-hot-toast";
import { v4 as uuid } from "uuid";

interface GradeSection {
    name: string;
    weight: number;
    data: string;
    id: string;
    classId: string;
}

interface SchoolClass {
    gradeSections: GradeSection[];
    credits: number;
    explicitGrade: string;
    desiredGrade: number;
}

export default function ImportPage() {
    const queryClient = useQueryClient();
    const inputClass =
        "mt-3 block w-full rounded-md border-midnight-700 bg-midnight-900 px-3 py-1.5 text-sm/6 text-white focus:outline-none data-[focus]:outline-4 data-[focus]:-outline-offset-2 data-[focus]:outline-blue-500 border-2 shadow-md disabled:text-gray-500";

    const [selectedFile, setSelectedFile] = useState<string | null>(null);
    const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        const file = event.target.files?.[0];
        if (file) {
            setSelectedFile(file.name);
        }
    };

    const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();

        const formData = new FormData(event.currentTarget);
        const className = formData.get("class-name") as string;
        const year = formData.get("year") as string;
        const semester = formData.get("semester") as string;

        if (!selectedFile) {
            toast.error("Please select a file to import.");
            return;
        }

        const file =
            document.querySelector<HTMLInputElement>("#grade-file")?.files?.[0];
        if (!file) {
            toast.error("Failed to read the file.");
            return;
        }

        const reader = new FileReader();

        reader.onload = (e: ProgressEvent<FileReader>) => {
            if (e.target && typeof e.target.result === "string") {
                const fileContents = e.target.result;
                try {
                    const parsedClass = parseFile(fileContents);
                    console.log(parsedClass);
                    axios
                        .post("/api/grades/" + className, {
                            year,
                            semester,
                            className,
                            recievedGrade: parsedClass.explicitGrade,
                            desiredGrade: parsedClass.desiredGrade,
                            credits: parsedClass.credits,
                            gradeSections: parsedClass.gradeSections,
                        })
                        .then(() => {
                            toast.success("Class imported!");
                            queryClient.invalidateQueries({
                                queryKey: ["gradesData"],
                            });
                        })
                        .catch((err) => {
                            toast.error(
                                "Failed to add class: " + err?.response?.data,
                            );
                            console.log(err);
                        });
                } catch (err) {
                    toast.error("Error: " + ((err as Error)?.message || ""));
                }
            }
        };
        reader.readAsText(file);
    };

    const parseFile = (fileContent: string): SchoolClass => {
        const lines = fileContent.split(/\r?\n/);
        const gradeSections: GradeSection[] = [];

        let inMetaOptions = false;
        let userExplicitGrade = "";
        let currentGradePartIndex = -1;
        let desiredGrade = -1.0;
        let credits = 4;

        for (let lineIndex = 0; lineIndex < lines.length; lineIndex++) {
            let line = lines[lineIndex].trim();

            if (line === "" || line.startsWith("#")) {
                continue;
            }

            if (line.startsWith("~ Meta")) {
                if (inMetaOptions) {
                    throw new Error(
                        `Meta options already started on line ${lineIndex + 1}`,
                    );
                }

                inMetaOptions = true;
                continue;
            }

            if (inMetaOptions && !line.startsWith(">")) {
                const [fieldName, fieldValue] = parseOptionLine(line);
                switch (fieldName) {
                    case "credits":
                        credits = parseInt(fieldValue, 10);
                        if (isNaN(credits))
                            throw new Error(
                                `Invalid credits value on line ${lineIndex + 1}`,
                            );
                        break;
                    case "grade":
                        userExplicitGrade = fieldValue;
                        break;
                    case "desired_grade":
                        desiredGrade = parseFloat(fieldValue);
                        if (isNaN(desiredGrade))
                            throw new Error(
                                `Invalid desired_grade value on line ${lineIndex + 1}`,
                            );
                        break;
                    // ignore the ignore field here
                    case "ignore":
                        break;
                    default:
                        throw new Error(
                            `Unknown meta field on line ${lineIndex + 1}`,
                        );
                }
                continue;
            }

            if (line.startsWith(">")) {
                inMetaOptions = false;
                const gradePartName = line.slice(1).trim();
                if (gradePartName === "")
                    throw new Error(
                        `Empty grade part name on line ${lineIndex + 1}`,
                    );

                for (const gradePart of gradeSections) {
                    if (gradePart.name === gradePartName) {
                        throw new Error(
                            `Duplicate grade part name on line ${lineIndex + 1}`,
                        );
                    }
                }

                currentGradePartIndex++;
                gradeSections.push({
                    name: gradePartName,
                    weight: 0,
                    data: "",
                    id: uuid(),
                    classId: "",
                });
                continue;
            } else if (currentGradePartIndex === -1) {
                throw new Error(
                    `Data outside grade part on line ${lineIndex + 1}`,
                );
            }

            let optionsString = line.trim();
            let nextLineIndex = lineIndex + 1;

            while (nextLineIndex < lines.length) {
                let nextLine = lines[nextLineIndex].trim();

                if (
                    nextLine.startsWith(">") ||
                    nextLine.startsWith("~") ||
                    nextLine.includes("=")
                ) {
                    break;
                }

                optionsString += "\n" + nextLine;
                nextLineIndex++;
            }

            const [fieldName, fieldValue] = parseOptionLine(optionsString);
            if (fieldName === "weight") {
                const weight = parseFloat(fieldValue) * 100;
                if (isNaN(weight))
                    throw new Error(
                        `Invalid weight value on line ${lineIndex + 1}`,
                    );
                gradeSections[currentGradePartIndex].weight = weight;
            } else if (fieldName === "data") {
                gradeSections[currentGradePartIndex].data = optionsString
                    .trimEnd()
                    .replace(/data\s?\=\s?/, "");
            } else {
                throw new Error(`Unknown field name on line ${lineIndex + 1}`);
            }

            lineIndex = nextLineIndex - 1;
        }

        return {
            gradeSections,
            credits,
            explicitGrade: userExplicitGrade,
            desiredGrade,
        };
    };

    const parseOptionLine = (line: string): [string, string] => {
        const parts = line
            .split("=")
            .map((part) => part.trim().replace(/^\"|\"$/g, ""));

        if (parts.length !== 2) {
            throw new Error(`Invalid line format: ${line}`);
        }

        const commentIndex = parts[1].indexOf("#");

        if (commentIndex !== -1) {
            parts[1] = parts[1].slice(0, commentIndex).trim();
        }

        return [parts[0], parts[1]];
    };

    return (
        <form
            className={"flex h-full w-full overflow-auto"}
            onSubmit={handleSubmit}
        >
            <Fieldset className="mb-4 w-full max-w-lg space-y-5 p-5 md:m-auto">
                <div>
                    <h1 className="border-b-4 border-b-slate-700 text-lg font-bold lg:text-2xl">
                        Import a Grade File!
                    </h1>
                    <h2 className="mt-3 text-sm/6 leading-7 text-slate-300">
                        From here you can import a grade file you have locally
                        into the website. The offline/cli version prefers
                        numbers in decimal form when specifying the weight of a
                        grade section, while the website is the opposite. When
                        importing a file, it will be converted to the website's
                        format.
                    </h2>
                </div>
                <Field>
                    <Label className="font-semibold after:ml-0.5 after:text-red-500 after:content-['*']">
                        Class Name
                    </Label>

                    <Description className="mb-3 mt-2 text-sm/6 text-slate-300">
                        This can also be retrieved from the meta section with
                        the field "name".
                    </Description>

                    <Input
                        className={inputClass}
                        placeholder="CS 101"
                        required
                        name="class-name"
                    ></Input>
                </Field>
                <Field>
                    <Label className="font-semibold after:ml-0.5 after:text-red-500 after:content-['*']">
                        Year
                    </Label>
                    <Description className="mb-3 mt-2 text-sm/6 text-slate-300">
                        This can also be retrieved from the meta section with
                        the field "year".
                    </Description>
                    <Input
                        className={inputClass}
                        placeholder="2021"
                        required
                        name="year"
                    ></Input>
                </Field>
                <Field>
                    <Label className="font-semibold after:ml-0.5 after:text-red-500 after:content-['*']">
                        Semester
                    </Label>
                    <Description className="mb-3 mt-2 text-sm/6 text-slate-300">
                        This can also be retrieved from the meta section with
                        the field "semester".
                    </Description>
                    <Select className={inputClass} name="semester">
                        <option value="fall">Fall</option>
                        <option value="spring">Spring</option>
                        <option value="winter">Winter</option>
                        <option value="summer">Summer</option>
                    </Select>
                </Field>
                <Field>
                    <Label className="mb-3 block text-sm font-semibold text-gray-900 dark:text-white">
                        Upload file{" "}
                        <span className="ml-0.5 text-red-500">*</span>
                    </Label>
                    <div className="relative flex items-center rounded-md border border-midnight-700 bg-midnight-900">
                        <Input
                            type="file"
                            name="grade-file"
                            id="grade-file"
                            className="absolute inset-0 h-full w-full cursor-pointer opacity-0"
                            onChange={handleFileChange}
                        />
                        <label
                            htmlFor="grade-file"
                            className="cursor-pointer rounded-[0.32rem] rounded-br-none rounded-tr-none bg-gray-700 px-3 py-2 text-center text-sm font-semibold text-gray-200 focus:outline-none"
                        >
                            Browse...
                        </label>
                        {selectedFile && (
                            <span className="ml-4 block text-sm">
                                Selected file: {selectedFile}
                            </span>
                        )}
                    </div>
                </Field>
                <Button
                    className="w-20 rounded-md bg-my-green px-4 py-2 text-sm font-semibold text-white disabled:bg-my-neutral"
                    type="submit"
                >
                    Add
                </Button>
            </Fieldset>
        </form>
    );
}
