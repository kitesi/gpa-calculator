import { getData as getYearsData } from "@/app/api/grades/route";
import { getClassData as getClassData } from "@/app/api/grades/[...paths]/utils";

export type GetYearsData = Awaited<ReturnType<typeof getYearsData>>;
export type GetClassData = Awaited<ReturnType<typeof getClassData>>;
