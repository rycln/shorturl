/*
Package main provides a comprehensive static analysis tool that combines multiple Go analyzers into a single executable.

The tool integrates:
- Standard analyzers from golang.org/x/tools/go/analysis/passes
- Selected staticcheck analyzers (SA class and specific checks)
- External analyzers (errcheck and ineffassign)
- Custom analyzer (noexitanalyzer)

Included Analyzers:

Standard Analyzers:
    asmdecl      - report mismatches between assembly files and Go declarations
    assign       - check for useless assignments
    atomic       - check for common mistakes using the sync/atomic package
    bools        - check for common mistakes involving boolean operators
    buildtag     - check that +build tags are well-formed and correctly located
    cgocall      - detect some violations of the cgo pointer passing rules
    composite    - check for unkeyed composite literals
    copylock     - check for locks erroneously passed by value
    errorsas     - report passing non-pointer or non-error values to errors.As
    fieldalignment - find structs that would use less memory if their fields were sorted
    httpresponse - check for mistakes using HTTP responses
    loopclosure  - check for references to enclosing loop variables
    lostcancel   - check for failure to call a context cancellation function
    nilfunc      - check for useless comparisons between functions and nil
    printf       - check consistency of Printf format strings and arguments
    shadow       - check for possible unintended shadowing of variables
    shift        - check for shifts that equal or exceed the width of the integer
    sortslice    - check for proper usage of sort.Slice
    stdmethods   - check signature of methods of well-known interfaces
    structtag    - check that struct field tags conform to reflect.StructTag's rules
    testinggoroutine - report calls to (*testing.T).Fatal from goroutines
    tests        - check for common mistaken usages of tests and examples
    unmarshal    - report passing non-pointer or non-interface values to unmarshal
    unreachable  - check for unreachable code
    unsafeptr    - check for invalid conversions of uintptr to unsafe.Pointer
    unusedresult - check for unused results of calls to some functions

Staticcheck Analyzers:
    SA*          - all staticcheck static analysis checks (SA class)
    S1000        - style check for specific issues
    ST1000       - style check for package comments
    QF1001       - quickfix checks for automatically fixable issues

External Analyzers:
    errcheck     - check for unchecked errors
    ineffassign  - detect ineffectual assignments

Custom Analyzers:
    noexit       - detect direct calls to os.Exit in main functions

*/

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

	var staticcheckAnalyzers []*analysis.Analyzer
	for _, v := range staticcheck.Analyzers {
		if len(v.Analyzer.Name) > 2 && v.Analyzer.Name[0:2] == "SA" {
			staticcheckAnalyzers = append(staticcheckAnalyzers, v.Analyzer)
		}
	}

	checks := map[string]bool{
		"S1000":  true,
		"ST1000": true,
		"QF1001": true,
	}

	var otherAnalyzers []*analysis.Analyzer
	for _, v := range staticcheck.Analyzers {
		if checks[v.Analyzer.Name] {
			otherAnalyzers = append(otherAnalyzers, v.Analyzer)
		}
	}

	externalAnalyzers := []*analysis.Analyzer{
		errcheck.Analyzer,
	}

	var analyzers []*analysis.Analyzer
	analyzers = append(analyzers, standardAnalyzers...)
	analyzers = append(analyzers, staticcheckAnalyzers...)
	analyzers = append(analyzers, otherAnalyzers...)
	analyzers = append(analyzers, externalAnalyzers...)
	analyzers = append(analyzers, noexitanalyzer.Analyzer)

	multichecker.Main(analyzers...)
}
