/*
  Warnings:

  - You are about to drop the `Comment` table. If the table is not empty, all the data it contains will be lost.
  - You are about to drop the `Post` table. If the table is not empty, all the data it contains will be lost.

*/
-- DropForeignKey
ALTER TABLE "Comment" DROP CONSTRAINT "Comment_postId_fkey";

-- DropForeignKey
ALTER TABLE "Comment" DROP CONSTRAINT "Comment_userId_fkey";

-- DropForeignKey
ALTER TABLE "Post" DROP CONSTRAINT "Post_userId_fkey";

-- DropTable
DROP TABLE "Comment";

-- DropTable
DROP TABLE "Post";

-- CreateTable
CREATE TABLE "grades" (
    "id" TEXT NOT NULL,
    "assignedGrade" TEXT NOT NULL,
    "className" TEXT NOT NULL,
    "credits" INTEGER NOT NULL,
    "desiredGrade" DOUBLE PRECISION NOT NULL,
    "year" INTEGER NOT NULL,
    "semester" TEXT NOT NULL,
    "userId" TEXT,

    CONSTRAINT "grades_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "gradesections" (
    "id" TEXT NOT NULL,
    "weight" DOUBLE PRECISION NOT NULL,
    "gradeId" TEXT NOT NULL,

    CONSTRAINT "gradesections_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "assignments" (
    "id" TEXT NOT NULL,
    "gradeId" TEXT NOT NULL,
    "gradeSectionId" TEXT NOT NULL,
    "totalPoints" DOUBLE PRECISION NOT NULL,
    "receivedPoints" DOUBLE PRECISION NOT NULL,

    CONSTRAINT "assignments_pkey" PRIMARY KEY ("id")
);

-- AddForeignKey
ALTER TABLE "grades" ADD CONSTRAINT "grades_userId_fkey" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "gradesections" ADD CONSTRAINT "gradesections_gradeId_fkey" FOREIGN KEY ("gradeId") REFERENCES "grades"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "assignments" ADD CONSTRAINT "assignments_gradeSectionId_fkey" FOREIGN KEY ("gradeSectionId") REFERENCES "gradesections"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "assignments" ADD CONSTRAINT "assignments_gradeId_fkey" FOREIGN KEY ("gradeId") REFERENCES "grades"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
