/*
  Warnings:

  - Made the column `semesterId` on table `classes` required. This step will fail if there are existing NULL values in that column.

*/
-- DropForeignKey
ALTER TABLE "classes" DROP CONSTRAINT "classes_semesterId_fkey";

-- AlterTable
ALTER TABLE "classes" ALTER COLUMN "semesterId" SET NOT NULL;

-- AddForeignKey
ALTER TABLE "classes" ADD CONSTRAINT "classes_semesterId_fkey" FOREIGN KEY ("semesterId") REFERENCES "semesters"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
