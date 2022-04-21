package gocognito_test

import (
	"testing"

	"github.com/nrnrk/gocognito"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	gocognito.Analyzer.Flags.Set("over", "0")
	analysistest.Run(t, testdata, gocognito.Analyzer, "a")
}

func TestAnalyzerOver3(t *testing.T) {
	testdata := analysistest.TestData()
	gocognito.Analyzer.Flags.Set("over", "3")
	analysistest.Run(t, testdata, gocognito.Analyzer, "b")
}
