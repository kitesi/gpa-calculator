/*
  Warnings:

  - A unique constraint covering the columns `[name,className]` on the table `gradesections` will be added. If there are existing duplicate values, this will fail.
  - Added the required column `name` to the `gradesections` table without a default value. This is not possible if the table is not empty.

*/
-- AlterTable
ALTER TABLE "gradesections" ADD COLUMN     "name" TEXT NOT NULL;

-- CreateIndex
CREATE UNIQUE INDEX "gradesections_name_className_key" ON "gradesections"("name", "className");
