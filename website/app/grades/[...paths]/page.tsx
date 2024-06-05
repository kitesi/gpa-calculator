"use client";
import axios, { AxiosError } from "axios";

import type { GetYearsData } from "@/app/types/data";
import { useQuery } from "@tanstack/react-query";

type Params = {
    params: {
        paths: string[];
    };
};

export default function SpecificGradePath({ params }: Params) {
    const { isPending, error, data } = useQuery<GetYearsData>({
        queryKey: ["classData"],
        queryFn: () =>
            axios
                .get("/api/grades/" + params.paths.join("/"))
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

    return <p>Success!</p>;
}
