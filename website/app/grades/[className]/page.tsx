"use client";
import axios, { AxiosError } from "axios";

import type { GetClassData } from "@/app/types/data";
import { useQuery } from "@tanstack/react-query";
import AddEditClassForm from "@/app/ui/AddEditClassForm";

type Params = {
    params: {
        className: string;
    };
};

export default function SpecificGradePath({ params }: Params) {
    const { isPending, error, data } = useQuery<GetClassData>({
        queryKey: ["classData", params.className],
        queryFn: () =>
            axios
                .get("/api/grades/" + params.className)
                .then((res) => res.data),
    });

    if (isPending) return <p>Loading...</p>;
    if (error) {
        return (
            <p>
                Error: {error.message},{" "}
                {error instanceof AxiosError && error?.response?.data}
            </p>
        );
    }

    if (!data) {
        return <p>Class not found</p>;
    }

    return (
        data && (
            <AddEditClassForm
                credits={data.credits}
                year={data.year.yearValue}
                semester={data.semester.name}
                className={data.className}
                recievedGrade={data.assignedGrade || ""}
                desiredGrade={data.desiredGrade ? data.desiredGrade + "" : ""}
                gradeSections={data.gradeSections}
                editing={true}
            />
        )
    );
}
