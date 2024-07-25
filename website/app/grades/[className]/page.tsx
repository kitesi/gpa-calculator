"use client";
import axios, { AxiosError } from "axios";

import type { GetClassData } from "@/app/types/data";
import { useQuery } from "@tanstack/react-query";
import AddEditClassForm from "@/app/ui/AddEditClassForm";
import { useSession } from "next-auth/react";
import NeedLogin from "@/app/ui/NeedLogin";
import Error from "@/app/ui/Error";

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

    const { data: session } = useSession();

    if (!session) {
        return <NeedLogin />;
    }

    if (isPending)
        return (
            <AddEditClassForm
                loading={true}
                credits={0}
                year={0}
                semester={""}
                className={"Loading"}
                recievedGrade=""
                desiredGrade=""
                gradeSections={[]}
                editing={true}
            />
        );

    if (error) {
        let message = error.message;

        if (error instanceof AxiosError) {
            message += " - " + error?.response?.data;
        }

        return <Error message={error.message} />;
    }

    if (!data) {
        return <Error message={"Class not found."} />;
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
                loading={false}
            />
        )
    );
}
