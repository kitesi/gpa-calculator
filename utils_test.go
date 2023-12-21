package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type UtilsTestSuite struct {
	suite.Suite
	directory     string
	inner_log_buf *bytes.Buffer
	inner_log     *log.Logger
}

func (suite *UtilsTestSuite) SetupSuite() {
	suite.directory = "test_files/grades/"
	suite.inner_log_buf = new(bytes.Buffer)
	suite.inner_log = log.New(suite.inner_log_buf, "", 0)

	if _, err := os.Stat(suite.directory); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Test directory (test_files/) does not exist\n")
		os.Exit(1)
	}
}

func TestUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(UtilsTestSuite))
}
