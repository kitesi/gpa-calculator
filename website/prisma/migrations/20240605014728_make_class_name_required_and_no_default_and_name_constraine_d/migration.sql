/*
  Warnings:

  - A unique constraint covering the columns `[name,yearValue]` on the table `semesters` will be added. If there are existing duplicate values, this will fail.

*/
-- CreateIndex
CREATE UNIQUE INDEX "semesters_name_yearValue_key" ON "semesters"("name", "yearValue");
