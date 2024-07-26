import { Prisma } from "@prisma/client";

export default function (err: Prisma.PrismaClientKnownRequestError) {
    switch (err.code) {
        case "P2002":
            return "Class with that name already exists.";
        case "P2025":
            return "Class not found.";
        default:
            return err.message;
    }
}
