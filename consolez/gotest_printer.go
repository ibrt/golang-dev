package consolez

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
)

// GoTestPrinter implements a printer for "go test" output.
type GoTestPrinter interface {
	PrintLine(line string)
	PrintDone()
}

type goTestPrinter struct {
	startTime time.Time
	passPkgs  int
	skipPkgs  int
	maxPkgLen int
}

// NewGoTestPrinter initializes a new [GoTestPrinter].
func NewGoTestPrinter() GoTestPrinter {
	return &goTestPrinter{
		startTime: time.Now(),
		passPkgs:  0,
		maxPkgLen: 60,
	}
}

// PrintLine implements the [GoTestPrinter] interface.
func (p *goTestPrinter) PrintLine(line string) {
	trimmedLine := strings.TrimSpace(line)

	switch {
	case p.maybeHandleSummaryLine(line), strings.HasPrefix(line, "coverage:"):
		// do nothing
	case p.isSecondary(trimmedLine):
		_, _ = ColorSecondary.Print(line)
		fmt.Print("\n")
	case p.isSuccess(trimmedLine):
		_, _ = ColorSuccess.Print(line)
		fmt.Print("\n")
	case p.isError(trimmedLine):
		_, _ = ColorError.Print(line)
		fmt.Print("\n")
	case p.isHighlight(trimmedLine):
		_, _ = ColorHighlight.Print(line)
		fmt.Print("\n")
	default:
		_, _ = fmt.Println(line)
	}
}

// PrintDone prints a final line.
func (p *goTestPrinter) PrintDone() {
	fmt.Printf(
		fmt.Sprintf("DONE    %%-%vv %%-10s\n", p.maxPkgLen),
		fmt.Sprintf("[SKIP: %v, PASS: %v]", p.skipPkgs, p.passPkgs),
		time.Since(p.startTime).Truncate(time.Millisecond*10))
}

func (p *goTestPrinter) maybeHandleSummaryLine(line string) bool {
	var pfx string
	var clr *color.Color

	switch {
	case strings.HasPrefix(line, "?   \t"), strings.HasPrefix(line, "\t"):
		pfx = "SKIP"
		clr = ColorSecondary
	case strings.HasPrefix(line, "ok  \t"):
		pfx = "PASS"
		clr = ColorSuccess
	case strings.HasPrefix(line, "FAIL\t"):
		pfx = "FAIL"
		clr = ColorError
	default:
		return false
	}

	rawLineParts := strings.Split(line, "\t")
	lineParts := make([]string, 0, len(rawLineParts))
	hasCoverage := false

	for _, rawLinePart := range rawLineParts {
		linePart := strings.TrimSpace(rawLinePart)

		if linePart == "" {
			continue
		} else if strings.HasPrefix(linePart, "coverage:") {
			hasCoverage = true
			continue
		} else if d, err := time.ParseDuration(linePart); err == nil {
			linePart = fmt.Sprintf("%-10s", d.Truncate(time.Millisecond*10))
		} else if linePart == "[no test files]" {
			linePart = "[no tests]"
		}

		lineParts = append(lineParts, linePart)
	}

	if len(lineParts) == 1 && hasCoverage {
		lineParts = []string{
			"",
			lineParts[0],
			"[no tests]",
		}
	}

	if len(lineParts) < 3 {
		return false
	}

	if pfx == "SKIP" {
		p.skipPkgs++
	}

	if pfx == "PASS" {
		p.passPkgs++
	}

	_, _ = clr.Printf(
		fmt.Sprintf("%%v    %%-%vv %%v", p.maxPkgLen),
		pfx,
		truncateLeft(lineParts[1], p.maxPkgLen),
		strings.Join(lineParts[2:], " "))
	fmt.Print("\n")

	return true
}

func (p *goTestPrinter) isSecondary(trimmedLine string) bool {
	switch {
	case strings.HasPrefix(trimmedLine, "--- SKIP"),
		trimmedLine == "SKIP",
		strings.HasPrefix(trimmedLine, "-test.shuffle"):
		return true
	default:
		return false
	}
}

func (p *goTestPrinter) isSuccess(trimmedLine string) bool {
	switch {
	case strings.HasPrefix(trimmedLine, "--- PASS"),
		trimmedLine == "PASS":
		return true
	default:
		return false
	}
}

func (p *goTestPrinter) isError(trimmedLine string) bool {
	switch {
	case strings.HasPrefix(trimmedLine, "--- FAIL"),
		trimmedLine == "FAIL",
		strings.Contains(trimmedLine, "no tests to run"):
		return true
	default:
		return false
	}
}

func (p *goTestPrinter) isHighlight(trimmedLine string) bool {
	switch {
	case strings.HasPrefix(trimmedLine, "=== RUN"),
		strings.Contains(trimmedLine, "[BeforeSuite]"),
		strings.Contains(trimmedLine, "[AfterSuite]"),
		strings.Contains(trimmedLine, "[BeforeTest]"),
		strings.Contains(trimmedLine, "[AfterTest]"),
		strings.Contains(trimmedLine, "[TestMethod]"):
		return true
	default:
		return false
	}
}
