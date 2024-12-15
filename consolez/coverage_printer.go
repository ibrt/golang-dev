package consolez

import (
	"fmt"

	"github.com/axw/gocov"
	"github.com/fatih/color"
	"github.com/ibrt/golang-utils/outz"
)

// Coverage describes collected coverage.
type Coverage struct {
	Packages []*gocov.Package
}

// CoveragePrinter implements a printer for test coverage information.
type CoveragePrinter interface {
	Print(coverage *Coverage)
}

type coveragePrinter struct {
	maxPkgLen int
	higLmt    float64
	medLmt    float64
}

// NewCoveragePrinter initializes a new [CoveragePrinter].
func NewCoveragePrinter() CoveragePrinter {
	return &coveragePrinter{
		maxPkgLen: 60,
		higLmt:    90,
		medLmt:    60,
	}
}

// Print implements the [CoveragePrinter] interface.
func (p *coveragePrinter) Print(coverage *Coverage) {
	var gtot, grch, lowPkgs, medPkgs, higPkgs int
	var gpct float64

	for _, pkg := range coverage.Packages {
		var tot, rch int
		var pct float64

		for _, fnc := range pkg.Functions {
			for _, stm := range fnc.Statements {
				tot++
				gtot++
				if stm.Reached > 0 {
					rch++
					grch++
				}
			}
		}

		if tot > 0 {
			pct = float64(rch) * 100 / float64(tot)
		} else {
			pct = 100
		}

		var pfx string
		var clr *color.Color

		switch {
		case pct >= p.higLmt:
			pfx = "HIGC"
			clr = outz.DefaultStyles.Success()
			higPkgs++
		case pct >= p.medLmt:
			pfx = "MEDC"
			clr = outz.DefaultStyles.Warning()
			medPkgs++
		default:
			pfx = "LOWC"
			clr = outz.DefaultStyles.Error()
			lowPkgs++
		}

		_, _ = clr.Printf(
			fmt.Sprintf("%%v    %%-%vv %%6v [%%v/%%v]", p.maxPkgLen),
			pfx,
			truncateLeft(pkg.Name, p.maxPkgLen),
			fmt.Sprintf("%.1f%%", pct),
			rch,
			tot)
		fmt.Print("\n")
	}

	if gtot > 0 {
		gpct = float64(grch) * 100 / float64(gtot)
	} else {
		gpct = 100
	}

	fmt.Printf(
		fmt.Sprintf("DONE    %%-%vv %%6v [%%v/%%v]\n", p.maxPkgLen),
		fmt.Sprintf("[LOWC: %v, MEDC: %v, HIGC: %v]", lowPkgs, medPkgs, higPkgs),
		fmt.Sprintf("%.1f%%", gpct),
		grch,
		gtot)
}
