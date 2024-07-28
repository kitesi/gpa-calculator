"use client";
import { useQuery } from "@tanstack/react-query";
import { useSession } from "next-auth/react";
import {
    Disclosure,
    DisclosureButton,
    DisclosurePanel,
} from "@headlessui/react";
import axios from "axios";
import {
    parseGradeSectionData,
    getGradeGPA,
    getGradeLetter,
} from "@/app/api/grades/[className]/utils";

import type { GetYearsData } from "@/app/types/data";
import { ChevronDownIcon } from "@heroicons/react/20/solid";

function formatGpa(gpa: number) {
    return gpa.toFixed(2);
}

function handleYear(
    gpaMap: { [key: string]: number },
    creditsMap: { [key: string]: number },
    year: GetYearsData[0],
) {
    let totalCreditsAdded = 0;

    year.semesters.map((semester) => {
        let semesterTotalCreditsAdded = 0;
        semester.classes.map((gradeClass) => {
            const gpa = gpaMap[gradeClass.id];

            if (gpa === -1) {
                return;
            }

            if (gpaMap[semester.id] === undefined) {
                gpaMap[semester.id] = 0;
            }

            semesterTotalCreditsAdded += gpa * gradeClass.credits;
            gpaMap[semester.id] +=
                (gpa * gradeClass.credits) / creditsMap[semester.id];
        });

        if (gpaMap[year.yearValue] === undefined) {
            gpaMap[year.yearValue] = 0;
        }

        gpaMap[year.yearValue] +=
            semesterTotalCreditsAdded / creditsMap[year.yearValue];
        totalCreditsAdded += semesterTotalCreditsAdded;
    });

    return totalCreditsAdded;
}

function handleClass(
    gradeClass: GetYearsData[0]["semesters"][0]["classes"][0],
    gradesMap: { [key: string]: string },
) {
    if (gradeClass.assignedGrade !== null && gradeClass.assignedGrade !== "") {
        gradesMap[gradeClass.id] = gradeClass.assignedGrade;
        return getGradeGPA(gradeClass.assignedGrade);
    }

    let totalWeight = 0;
    let totalGrades = 0;

    for (const section of gradeClass.gradeSections) {
        const { value: scores } = parseGradeSectionData(section.data);

        let pointsTotal = 0;
        let pointsRecieved = 0;

        for (const score of scores!) {
            pointsTotal += score.total;
            pointsRecieved += score.recieved;
        }

        if (pointsTotal == 0) {
            continue;
        }

        totalWeight += section.weight;
        totalGrades += (pointsRecieved / pointsTotal) * section.weight;
    }

    if (totalWeight == 0) {
        gradesMap[gradeClass.id] = "?";
        return -1;
    } else {
        gradesMap[gradeClass.id] = getGradeLetter(totalGrades / totalWeight);
        return getGradeGPA(getGradeLetter(totalGrades / totalWeight));
    }
}

export default function GradesDisplayed() {
    const { data: session } = useSession();

    const { isPending, error, data } = useQuery<GetYearsData>({
        queryKey: ["gradesData"],
        queryFn: () => axios.get("/api/grades").then((res) => res.data),
    });

    if (!session) {
        return;
    }

    if (!data) return;

    const creditsMap: { [k: string]: number } = {};
    const gpaMap: { [k: string]: number } = {};
    const gradesMap: { [k: string]: string } = {};

    // set up creditsMap and gradesMap
    data.map((year) => {
        year.semesters.map((semester) => {
            semester.classes.map((gradeClass) => {
                const grade = handleClass(gradeClass, gradesMap);
                gpaMap[gradeClass.id] = grade;

                if (grade !== -1) {
                    if (creditsMap[year.yearValue] === undefined) {
                        creditsMap[year.yearValue] = 0;
                    }

                    if (creditsMap[semester.id] === undefined) {
                        creditsMap[semester.id] = 0;
                    }

                    creditsMap[year.yearValue] += gradeClass.credits;
                    creditsMap[semester.id] += gradeClass.credits;
                }
            });
        });
    });

    // setup gpaMap
    let totalCreditsAdded = 0;
    let totalCredits = 0;
    data.forEach((year) => {
        totalCreditsAdded += handleYear(gpaMap, creditsMap, year);
        totalCredits += creditsMap[year.yearValue];
    });

    return (
        <div className="flex h-full w-full flex-col p-4">
            <div className="mx-auto my-20 w-full max-w-lg rounded-md border-2 border-midnight-700 bg-midnight-900 shadow-lg lg:my-auto">
                <h2 className="border-b-2 border-b-midnight-700 p-2 pl-4 text-xl font-bold text-white">
                    {" "}
                    Grades (GPA: {formatGpa(totalCreditsAdded / totalCredits)})
                </h2>
                <div className="p-4">
                    {data.map((year, yearIndex) => (
                        <div key={year.yearValue} className={``}>
                            <div className="mb-1 font-bold text-white">
                                {year.yearValue}{" "}
                                <span className="font-normal text-gray-400">
                                    ({formatGpa(gpaMap[year.yearValue])})
                                </span>
                            </div>
                            {year.semesters.map((semester, semIndex) => (
                                <div
                                    key={semester.id}
                                    className={`ml-4 leading-none ${semIndex !== year.semesters.length - 1 || yearIndex !== data.length - 1 ? "border-l-[1.5px] pl-4" : ""}`}
                                >
                                    <div className="mb-1 -translate-x-1 transform font-semibold text-white">
                                        {semIndex !== year.semesters.length - 1
                                            ? "├──"
                                            : "└──"}{" "}
                                        {semester.name}{" "}
                                        <span className="font-normal text-gray-400">
                                            ({formatGpa(gpaMap[semester.id])})
                                        </span>
                                    </div>
                                    {semester.classes.map(
                                        (gradeClass, classIndex) => (
                                            <div
                                                key={gradeClass.id}
                                                className="ml-4 leading-none"
                                            >
                                                <div className="mb- text-base font-medium text-white">
                                                    {classIndex !==
                                                    semester.classes.length - 1
                                                        ? "├──"
                                                        : "└──"}{" "}
                                                    {gradeClass.className}{" "}
                                                    <span className="text-base font-normal text-gray-400">
                                                        {gradesMap[
                                                            gradeClass.id
                                                        ] === "?"
                                                            ? ""
                                                            : gradesMap[
                                                                  gradeClass.id
                                                              ]}
                                                    </span>
                                                </div>
                                            </div>
                                        ),
                                    )}
                                </div>
                            ))}
                        </div>
                    ))}
                </div>
            </div>
        </div>

        // <>

        //     <Disclosure as="div" className="p-6" defaultOpen={true}>
        //         <DisclosureButton className="group flex w-full items-center justify-between">
        //             <span className="text-sm/6 font-medium text-white group-data-[hover]:text-white/80">
        //                 {year.yearValue} ({gpaMap[year.yearValue]})
        //             </span>
        //             <ChevronDownIcon className="size-5 fill-white/60 group-data-[open]:rotate-180 group-data-[hover]:fill-white/50" />
        //         </DisclosureButton>
        //         <DisclosurePanel className="mt-2 text-sm/5 text-white/50">
        //             If you're unhappy with your purchase, we'll
        //             refund you in full.
        //         </DisclosurePanel>
        //     </Disclosure>
        // </>
    );
}
