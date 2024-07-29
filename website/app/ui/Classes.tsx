"use client";
import axios from "axios";
import type { GetYearsData } from "@/app/types/data";
import { useQuery } from "@tanstack/react-query";
import { useSession } from "next-auth/react";
import { usePathname } from "next/navigation";
import Link from "next/link";
import clsx from "clsx";

export default function SchoolClasses() {
    const { data: session } = useSession();
    const { isPending, error, data } = useQuery<GetYearsData>({
        queryKey: ["gradesData"],
        queryFn: () => axios.get("/api/grades").then((res) => res.data),
    });

    if (!session) {
        return;
    }

    if (isPending) return <p className="p-5 font-semibold">Loading...</p>;
    if (error)
        return <p className="p-5 font-semibold">Error: {error.message}</p>;

    const currentPath = usePathname();

    return (
        <div className="flex w-full flex-col p-4 pl-6">
            {data.map((year, yearIndex) => (
                <div key={year.yearValue}>
                    <div className="mb-1 font-bold text-white">
                        {year.yearValue}
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
                                {semester.name}
                            </div>
                            {semester.classes.map((schoolClass, classIndex) => (
                                <div
                                    key={schoolClass.className}
                                    className="ml-4 leading-none"
                                >
                                    <div className="mb-1 text-base font-medium text-white">
                                        {classIndex !==
                                        semester.classes.length - 1
                                            ? "├──"
                                            : "└──"}{" "}
                                        <Link
                                            href={`/grades/${schoolClass.className}`}
                                            className={clsx(
                                                "w-full rounded-sm bg-midnight-800 px-3 py-0",
                                                currentPath ==
                                                    `/grades/${schoolClass.className}`
                                                    ? "font-semibold text-white"
                                                    : "text-blue-400",
                                            )}
                                        >
                                            {schoolClass.className}
                                        </Link>
                                    </div>
                                </div>
                            ))}
                        </div>
                    ))}
                </div>
            ))}
        </div>
    );
}
