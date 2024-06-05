import { PrismaClient } from "@prisma/client";

/**
 * @type {import('@prisma/client').PrismaClient}
 */
const client = globalThis.prisma || new PrismaClient();

if (process.env.NODE_ENV === "development") {
    globalThis.prisma = client;
}

export default client;
