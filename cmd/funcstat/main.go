package main

import (
	"github.com/gostaticanalysis/funcstat"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(funcstat.Analyzer) }
