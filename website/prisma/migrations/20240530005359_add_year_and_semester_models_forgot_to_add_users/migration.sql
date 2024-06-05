/*
  Warnings:

  - Added the required column `userId` to the `semesters` table without a default value. This is not possible if the table is not empty.
  - Added the required column `userId` to the `years` table without a default value. This is not possible if the table is not empty.

*/
-- AlterTable
ALTER TABLE "semesters" ADD COLUMN     "userId" TEXT NOT NULL;

-- AlterTable
ALTER TABLE "years" ADD COLUMN     "userId" TEXT NOT NULL;

-- AddForeignKey
ALTER TABLE "years" ADD CONSTRAINT "years_userId_fkey" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "semesters" ADD CONSTRAINT "semesters_userId_fkey" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
