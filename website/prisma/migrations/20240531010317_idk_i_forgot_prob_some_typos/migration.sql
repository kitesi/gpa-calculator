/*
  Warnings:

  - You are about to drop the column `classId` on the `assignments` table. All the data in the column will be lost.
  - The primary key for the `classes` table will be changed. If it partially fails, the table could be left without primary key constraint.
  - You are about to drop the column `id` on the `classes` table. All the data in the column will be lost.
  - You are about to drop the column `semester` on the `classes` table. All the data in the column will be lost.
  - You are about to drop the column `classId` on the `gradesections` table. All the data in the column will be lost.
  - Added the required column `className` to the `assignments` table without a default value. This is not possible if the table is not empty.
  - Added the required column `className` to the `gradesections` table without a default value. This is not possible if the table is not empty.

*/
-- DropForeignKey
ALTER TABLE "assignments" DROP CONSTRAINT "assignments_classId_fkey";

-- DropForeignKey
ALTER TABLE "gradesections" DROP CONSTRAINT "gradesections_classId_fkey";

-- AlterTable
ALTER TABLE "assignments" DROP COLUMN "classId",
ADD COLUMN     "className" TEXT NOT NULL;

-- AlterTable
ALTER TABLE "classes" DROP CONSTRAINT "classes_pkey",
DROP COLUMN "id",
DROP COLUMN "semester",
ADD CONSTRAINT "classes_pkey" PRIMARY KEY ("className");

-- AlterTable
ALTER TABLE "gradesections" DROP COLUMN "classId",
ADD COLUMN     "className" TEXT NOT NULL;

-- AddForeignKey
ALTER TABLE "gradesections" ADD CONSTRAINT "gradesections_className_fkey" FOREIGN KEY ("className") REFERENCES "classes"("className") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "assignments" ADD CONSTRAINT "assignments_className_fkey" FOREIGN KEY ("className") REFERENCES "classes"("className") ON DELETE RESTRICT ON UPDATE CASCADE;
