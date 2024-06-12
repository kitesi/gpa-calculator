/*
  Warnings:

  - You are about to drop the column `className` on the `gradesections` table. All the data in the column will be lost.

*/
-- DropForeignKey
ALTER TABLE "gradesections" DROP CONSTRAINT "gradesections_className_fkey";

-- AlterTable
ALTER TABLE "gradesections" DROP COLUMN "className";

-- AddForeignKey
ALTER TABLE "gradesections" ADD CONSTRAINT "gradesections_classId_fkey" FOREIGN KEY ("classId") REFERENCES "classes"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
