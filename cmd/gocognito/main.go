package main

import (
	"github.com/nrnrk/gocognito"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(gocognito.Analyzer)
}
