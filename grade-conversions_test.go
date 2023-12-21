package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GradeConversionsTestSuite struct {
	suite.Suite
}

func (suite *GradeConversionsTestSuite) TestGetGradeGPA() {
	t := suite.T()

	assert.Equal(t, getGradeGPA("A+"), 4.3)
	assert.Equal(t, getGradeGPA("A"), 4.0)
	assert.Equal(t, getGradeGPA("A-"), 3.7)
	assert.Equal(t, getGradeGPA("B+"), 3.3)
	assert.Equal(t, getGradeGPA("B"), 3.0)
	assert.Equal(t, getGradeGPA("B-"), 2.7)
	assert.Equal(t, getGradeGPA("C+"), 2.3)
	assert.Equal(t, getGradeGPA("C"), 2.0)
	assert.Equal(t, getGradeGPA("C-"), 1.7)
	assert.Equal(t, getGradeGPA("D+"), 1.3)
	assert.Equal(t, getGradeGPA("D"), 1.0)
	assert.Equal(t, getGradeGPA("D-"), 0.7)
	assert.Equal(t, getGradeGPA("F"), 0.0)
}

func (suite *GradeConversionsTestSuite) TestGetGradeLetter() {
	t := suite.T()

	assert.Equal(t, getGradeLetter(0.99), "A")
	assert.Equal(t, getGradeLetter(0.94), "A")
	assert.Equal(t, getGradeLetter(0.93), "A-")
	assert.Equal(t, getGradeLetter(0.89), "B+")
	assert.Equal(t, getGradeLetter(0.84), "B")
	assert.Equal(t, getGradeLetter(0.83), "B-")
	assert.Equal(t, getGradeLetter(0.79), "C+")
	assert.Equal(t, getGradeLetter(0.74), "C")
	assert.Equal(t, getGradeLetter(0.73), "C-")
	assert.Equal(t, getGradeLetter(0.69), "D+")
	assert.Equal(t, getGradeLetter(0.64), "D")
	assert.Equal(t, getGradeLetter(0.63), "D-")
	assert.Equal(t, getGradeLetter(0.59), "F")
}

func TestGradeConversionsTestSuite(t *testing.T) {
	suite.Run(t, new(GradeConversionsTestSuite))
}
