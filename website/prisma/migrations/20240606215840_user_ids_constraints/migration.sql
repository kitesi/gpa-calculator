/*
  Warnings:

  - The primary key for the `classes` table will be changed. If it partially fails, the table could be left without primary key constraint.
  - The primary key for the `years` table will be changed. If it partially fails, the table could be left without primary key constraint.
  - A unique constraint covering the columns `[className,userId]` on the table `classes` will be added. If there are existing duplicate values, this will fail.
  - A unique constraint covering the columns `[name,classId]` on the table `gradesections` will be added. If there are existing duplicate values, this will fail.
  - A unique constraint covering the columns `[name,yearId]` on the table `semesters` will be added. If there are existing duplicate values, this will fail.
  - A unique constraint covering the columns `[yearValue,userId]` on the table `years` will be added. If there are existing duplicate values, this will fail.
  - The required column `id` was added to the `classes` table with a prisma-level default value. This is not possible if the table is not empty. Please add this column as optional, then populate it before making it required.
  - Added the required column `semesterName` to the `classes` table without a default value. This is not possible if the table is not empty.
  - Added the required column `yearId` to the `classes` table without a default value. This is not possible if the table is not empty.
  - Added the required column `classId` to the `gradesections` table without a default value. This is not possible if the table is not empty.
  - Added the required column `yearId` to the `semesters` table without a default value. This is not possible if the table is not empty.
  - The required column `id` was added to the `years` table with a prisma-level default value. This is not possible if the table is not empty. Please add this column as optional, then populate it before making it required.

*/
-- DropForeignKey
ALTER TABLE "classes" DROP CONSTRAINT "classes_yearValue_fkey";

-- DropForeignKey
ALTER TABLE "gradesections" DROP CONSTRAINT "gradesections_className_fkey";

-- DropForeignKey
ALTER TABLE "semesters" DROP CONSTRAINT "semesters_yearValue_fkey";

-- DropIndex
DROP INDEX "gradesections_name_className_key";

-- DropIndex
DROP INDEX "semesters_name_yearValue_key";

-- AlterTable
ALTER TABLE "classes" DROP CONSTRAINT "classes_pkey",
ADD COLUMN     "id" TEXT NOT NULL,
ADD COLUMN     "semesterName" TEXT NOT NULL,
ADD COLUMN     "yearId" TEXT NOT NULL,
ADD CONSTRAINT "classes_pkey" PRIMARY KEY ("id");

-- AlterTable
ALTER TABLE "gradesections" ADD COLUMN     "classId" TEXT NOT NULL;

-- AlterTable
ALTER TABLE "semesters" ADD COLUMN     "yearId" TEXT NOT NULL;

-- AlterTable
ALTER TABLE "years" DROP CONSTRAINT "years_pkey",
ADD COLUMN     "id" TEXT NOT NULL,
ADD CONSTRAINT "years_pkey" PRIMARY KEY ("id");

-- CreateIndex
CREATE UNIQUE INDEX "classes_className_userId_key" ON "classes"("className", "userId");

-- CreateIndex
CREATE UNIQUE INDEX "gradesections_name_classId_key" ON "gradesections"("name", "classId");

-- CreateIndex
CREATE UNIQUE INDEX "semesters_name_yearId_key" ON "semesters"("name", "yearId");

-- CreateIndex
CREATE UNIQUE INDEX "years_yearValue_userId_key" ON "years"("yearValue", "userId");

-- AddForeignKey
ALTER TABLE "semesters" ADD CONSTRAINT "semesters_yearId_fkey" FOREIGN KEY ("yearId") REFERENCES "years"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "classes" ADD CONSTRAINT "classes_yearId_fkey" FOREIGN KEY ("yearId") REFERENCES "years"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "gradesections" ADD CONSTRAINT "gradesections_className_fkey" FOREIGN KEY ("className") REFERENCES "classes"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
