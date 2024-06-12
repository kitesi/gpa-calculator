import { ParamsObject, handleAuthorization, getClassData } from "./utils";

export default async function GET(req: Request, { params }: ParamsObject) {
    const { error, session, user } = await handleAuthorization();

    if (error) {
        return error;
    }

    try {
        return Response.json(await getClassData(params.className, user.id));
    } catch (err: any) {
        let errMessage = "";

        if (err?.message) {
            errMessage = err.message;
        }

        if (err?.body?.message) {
            errMessage = err.body.message;
        }

        return new Response("Error retrieving grade: " + errMessage, {
            status: 500,
        });
    }
}
