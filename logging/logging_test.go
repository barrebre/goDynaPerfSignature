package logging

import (
	"testing"

	"github.com/barrebre/goDynaPerfSignature/datatypes"
	"github.com/stretchr/testify/assert"
)

func TestSetLogLevel(t *testing.T) {
	type testDefs struct {
		Name           string
		Input          string
		ExpectedOutput string
	}

	tests := []testDefs{
		testDefs{
			Name:           "PASS - Nothing set. Default value",
			ExpectedOutput: "ERROR",
		},
		testDefs{
			Name:           "PASS - DEBUG set",
			Input:          "DEBUG",
			ExpectedOutput: "DEBUG",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			if test.Input != "" {
				SetLogLevel(test.Input)
			}

			assert.Equal(t, test.ExpectedOutput, getLoggingLevel())
		})
	}
}

func TestNoReturnMethods(t *testing.T) {
	perfSig := datatypes.PerformanceSignature{
		Metrics: []datatypes.Metric{
			{
				ID: "asdf",
			},
		},
	}
	LogInfo(datatypes.Logging{Message: "asdf", PerfSig: perfSig})
	LogDebug(datatypes.Logging{Message: "asdf"})
	LogError(datatypes.Logging{Message: "asdf"})
	LogSystem(datatypes.Logging{Message: "asdf"})
}

func TestShortenedFileName(t *testing.T) {
	type testDefs struct {
		Name           string
		Value          string
		ExpectedReturn string
	}

	tests := []testDefs{
		testDefs{
			Name:           "Pass - Expected Value",
			Value:          "/go/src/github.com/barrebre/goDynaPerfSignature/logging/logging.go",
			ExpectedReturn: "logging/logging.go",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			out := shortenedFileName(test.Value)
			assert.Equal(t, test.ExpectedReturn, out)
		})
	}
}

func TestShortenedMethodName(t *testing.T) {
	type testDefs struct {
		Name           string
		Value          string
		ExpectedReturn string
	}

	tests := []testDefs{
		testDefs{
			Name:           "Pass - Expected Value",
			Value:          "github.com/barrebre/goDynaPerfSignature/logging.Log",
			ExpectedReturn: "Log",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			out := shortenedMethodName(test.Value)
			assert.Equal(t, test.ExpectedReturn, out)
		})
	}
}

// Line is not checked because it will change whenever this file does and file is not checked because it shows local sensitive info
func TestTraceMethod(t *testing.T) {
	type values struct {
		File   string
		Line   int
		Method string
	}

	type testDefs struct {
		Name           string
		ExpectedMethod string
	}

	tests := []testDefs{
		testDefs{
			Name:           "Pass - Expected Value",
			ExpectedMethod: "github.com/barrebre/goDynaPerfSignature/logging.anon1",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			_, _, method := anon1()
			assert.Equal(t, test.ExpectedMethod, method)
		})
	}
}

func anon1() (string, int, string) {
	return anon2()
}

func anon2() (string, int, string) {
	return anon3()
}

func anon3() (string, int, string) {
	return traceMethod()
}
