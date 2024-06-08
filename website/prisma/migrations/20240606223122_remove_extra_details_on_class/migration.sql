/*
  Warnings:

  - You are about to drop the column `semesterName` on the `classes` table. All the data in the column will be lost.
  - You are about to drop the column `yearValue` on the `classes` table. All the data in the column will be lost.

*/
-- AlterTable
ALTER TABLE "classes" DROP COLUMN "semesterName",
DROP COLUMN "yearValue";
