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
        <ul className="ml-8 mt-5">
            {data.map((year) => (
                <li key={year.yearValue} className="my-1">
                    <h2 className="font-semibold">{year.yearValue}</h2>
                    <ul className="ml-5">
                        {year.semesters.map((semester) => (
                            <li key={semester.id} className="my-1">
                                <h3 className="font-medium">{semester.name}</h3>
                                <ul className="ml-5">
                                    {semester.classes.map((schoolClass) => (
                                        <li
                                            key={schoolClass.className}
                                            className="my-2"
                                        >
                                            <Link
                                                href={`/grades/${schoolClass.className}`}
                                                className="rounded-l-lg bg-gray-800 p-2"
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
