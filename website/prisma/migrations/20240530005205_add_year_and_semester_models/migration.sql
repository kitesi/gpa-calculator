/*
  Warnings:

  - You are about to drop the column `year` on the `grades` table. All the data in the column will be lost.
  - Added the required column `yearValue` to the `grades` table without a default value. This is not possible if the table is not empty.
  - Made the column `userId` on table `grades` required. This step will fail if there are existing NULL values in that column.

*/
-- DropForeignKey
ALTER TABLE "grades" DROP CONSTRAINT "grades_userId_fkey";

-- AlterTable
ALTER TABLE "grades" DROP COLUMN "year",
ADD COLUMN     "yearValue" INTEGER NOT NULL,
ALTER COLUMN "userId" SET NOT NULL;

-- CreateTable
CREATE TABLE "years" (
    "yearValue" INTEGER NOT NULL,

    CONSTRAINT "years_pkey" PRIMARY KEY ("yearValue")
);

-- CreateTable
CREATE TABLE "semesters" (
    "id" TEXT NOT NULL,
    "name" TEXT NOT NULL,
    "yearValue" INTEGER NOT NULL,

    CONSTRAINT "semesters_pkey" PRIMARY KEY ("id")
);

-- AddForeignKey
ALTER TABLE "semesters" ADD CONSTRAINT "semesters_yearValue_fkey" FOREIGN KEY ("yearValue") REFERENCES "years"("yearValue") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "grades" ADD CONSTRAINT "grades_userId_fkey" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "grades" ADD CONSTRAINT "grades_yearValue_fkey" FOREIGN KEY ("yearValue") REFERENCES "years"("yearValue") ON DELETE RESTRICT ON UPDATE CASCADE;
