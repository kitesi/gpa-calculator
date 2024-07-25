"use client";
import axios from "axios";

import type { GetYearsData } from "@/app/types/data";
import { useQuery } from "@tanstack/react-query";
import { useSession } from "next-auth/react";
import { FolderIcon } from "@heroicons/react/20/solid";
import Link from "next/link";
import { useRouter } from "next/router";
import clsx from "clsx";
import { usePathname } from "next/navigation";

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
        <ul className="w-full">
            {data.map((year) =>
                year.semesters.map((semester) =>
                    semester.classes.map((schoolClass) => (
                        <li key={schoolClass.className} className="w-full">
                            <Link
                                href={`/grades/${schoolClass.className}`}
                                className={clsx(
                                    "block w-full border-b-[1px] border-r-2 border-b-midnight-700 border-r-midnight-800 py-4 text-center",
                                    currentPath ==
                                        `/grades/${schoolClass.className}` &&
                                        "border-r-orange-800",
                                )}
                            >
                                {schoolClass.className} - {year.yearValue} -{" "}
                                {semester.name}
                            </Link>
                        </li>
                    )),
                ),
            )}
        </ul>
    );
}
