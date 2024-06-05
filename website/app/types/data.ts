import { getData as getYearsData } from "@/app/api/grades/route";
import { getData as getClassData } from "@/app/api/grades/[...paths]/route";

export type GetYearsData = Awaited<ReturnType<typeof getYearsData>>;
export type GetClassData = Awaited<ReturnType<typeof getClassData>>;
