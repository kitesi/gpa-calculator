import { getData } from "@/app/api/grades/route";
export type GetYearsData = Awaited<ReturnType<typeof getData>>;
