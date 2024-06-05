/*
  Warnings:

  - You are about to drop the column `gradeId` on the `assignments` table. All the data in the column will be lost.
  - You are about to drop the column `gradeId` on the `gradesections` table. All the data in the column will be lost.
  - You are about to drop the `grades` table. If the table is not empty, all the data it contains will be lost.
  - Added the required column `classId` to the `assignments` table without a default value. This is not possible if the table is not empty.
  - Added the required column `classId` to the `gradesections` table without a default value. This is not possible if the table is not empty.

*/
-- DropForeignKey
ALTER TABLE "assignments" DROP CONSTRAINT "assignments_gradeId_fkey";

-- DropForeignKey
ALTER TABLE "grades" DROP CONSTRAINT "grades_userId_fkey";

-- DropForeignKey
ALTER TABLE "grades" DROP CONSTRAINT "grades_yearValue_fkey";

-- DropForeignKey
ALTER TABLE "gradesections" DROP CONSTRAINT "gradesections_gradeId_fkey";

-- AlterTable
ALTER TABLE "assignments" DROP COLUMN "gradeId",
ADD COLUMN     "classId" TEXT NOT NULL;

-- AlterTable
ALTER TABLE "gradesections" DROP COLUMN "gradeId",
ADD COLUMN     "classId" TEXT NOT NULL;

-- DropTable
DROP TABLE "grades";

-- CreateTable
CREATE TABLE "classes" (
    "id" TEXT NOT NULL,
    "assignedGrade" TEXT NOT NULL,
    "className" TEXT NOT NULL,
    "credits" INTEGER NOT NULL,
    "desiredGrade" DOUBLE PRECISION NOT NULL,
    "yearValue" INTEGER NOT NULL,
    "semester" TEXT NOT NULL,
    "userId" TEXT NOT NULL,
    "semesterId" TEXT,

    CONSTRAINT "classes_pkey" PRIMARY KEY ("id")
);

-- AddForeignKey
ALTER TABLE "classes" ADD CONSTRAINT "classes_userId_fkey" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "classes" ADD CONSTRAINT "classes_yearValue_fkey" FOREIGN KEY ("yearValue") REFERENCES "years"("yearValue") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "classes" ADD CONSTRAINT "classes_semesterId_fkey" FOREIGN KEY ("semesterId") REFERENCES "semesters"("id") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "gradesections" ADD CONSTRAINT "gradesections_classId_fkey" FOREIGN KEY ("classId") REFERENCES "classes"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "assignments" ADD CONSTRAINT "assignments_classId_fkey" FOREIGN KEY ("classId") REFERENCES "classes"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
