package main

import (
	"github.com/kisielk/errcheck/errcheck"
	"github.com/rycln/shorturl/cmd/staticlint/noexitanalyzer"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"honnef.co/go/tools/staticcheck"
)

func main() {
	// Стандартные анализаторы из golang.org/x/tools/go/analysis/passes
	standardAnalyzers := []*analysis.Analyzer{
		asmdecl.Analyzer,
		assign.Analyzer,
		atomic.Analyzer,
		bools.Analyzer,
		buildtag.Analyzer,
		cgocall.Analyzer,
		composite.Analyzer,
		copylock.Analyzer,
		errorsas.Analyzer,
		fieldalignment.Analyzer,
		httpresponse.Analyzer,
		loopclosure.Analyzer,
		lostcancel.Analyzer,
		nilfunc.Analyzer,
		printf.Analyzer,
		shadow.Analyzer,
		shift.Analyzer,
		sortslice.Analyzer,
		stdmethods.Analyzer,
		structtag.Analyzer,
		testinggoroutine.Analyzer,
		tests.Analyzer,
		unmarshal.Analyzer,
		unreachable.Analyzer,
		unsafeptr.Analyzer,
		unusedresult.Analyzer,
	}

	// Анализаторы staticcheck класса SA (static analysis)
	var staticcheckAnalyzers []*analysis.Analyzer
	for _, v := range staticcheck.Analyzers {
		// Берем только анализаторы класса SA
		if len(v.Analyzer.Name) > 2 && v.Analyzer.Name[0:2] == "SA" {
			staticcheckAnalyzers = append(staticcheckAnalyzers, v.Analyzer)
		}
	}

	// По одному анализатору из других классов staticcheck
	checks := map[string]bool{
		"S1000":  true,
		"ST1000": true,
		"QF1001": true,
	}

	var otherAnalyzers []*analysis.Analyzer
	for _, v := range staticcheck.Analyzers {
		// добавляем в массив нужные проверки
		if checks[v.Analyzer.Name] {
			otherAnalyzers = append(otherAnalyzers, v.Analyzer)
		}
	}

	// Два публичных анализатора
	externalAnalyzers := []*analysis.Analyzer{
		errcheck.Analyzer,
	}

	// Собираем все анализаторы вместе
	var analyzers []*analysis.Analyzer
	analyzers = append(analyzers, standardAnalyzers...)
	analyzers = append(analyzers, staticcheckAnalyzers...)
	analyzers = append(analyzers, otherAnalyzers...)
	analyzers = append(analyzers, externalAnalyzers...)
	analyzers = append(analyzers, noexitanalyzer.Analyzer)

	multichecker.Main(analyzers...)
}
