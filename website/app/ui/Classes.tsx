"use client";
import axios from "axios";

import type { GetYearsData } from "@/app/types/data";
import { useQuery } from "@tanstack/react-query";
import { useSession } from "next-auth/react";
import { FolderIcon } from "@heroicons/react/20/solid";
import Link from "next/link";

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

    return (
        <ul className="ml-5 mt-5">
            {data.map((year) => (
                <li key={year.yearValue}>
                    <h2>{year.yearValue}</h2>
                    <ul className="ml-5">
                        {year.semesters.map((semester) => (
                            <li key={semester.id}>
                                <h3>{semester.name}</h3>
                                <ul className="ml-5">
                                    {semester.classes.map((schoolClass) => (
                                        <li
                                            key={schoolClass.className}
                                            className="text-blue-300 underline"
                                        >
                                            <Link
                                                href={`/grades/${year.yearValue}/${semester.name}/${schoolClass.className}`}
                                            >
                                                {schoolClass.className}
                                            </Link>
                                        </li>
                                    ))}
                                </ul>
                            </li>
                        ))}
                    </ul>
                </li>
            ))}
        </ul>
    );
}
