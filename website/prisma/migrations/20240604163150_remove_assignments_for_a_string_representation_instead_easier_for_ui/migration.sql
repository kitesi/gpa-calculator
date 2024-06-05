/*
  Warnings:

  - You are about to drop the `assignments` table. If the table is not empty, all the data it contains will be lost.
  - Added the required column `data` to the `gradesections` table without a default value. This is not possible if the table is not empty.

*/
-- DropForeignKey
ALTER TABLE "assignments" DROP CONSTRAINT "assignments_className_fkey";

-- DropForeignKey
ALTER TABLE "assignments" DROP CONSTRAINT "assignments_gradeSectionId_fkey";

-- AlterTable
ALTER TABLE "gradesections" ADD COLUMN     "data" TEXT NOT NULL;

-- DropTable
DROP TABLE "assignments";
