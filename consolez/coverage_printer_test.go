package consolez_test

import (
	"strings"
	"testing"

	"github.com/axw/gocov"
	"github.com/ibrt/golang-utils/fixturez"
	. "github.com/onsi/gomega"

	"github.com/ibrt/golang-dev/consolez"
)

type CoveragePrinterSuite struct {
	// intentionally empty
}

func TestCoveragePrinterSuite(t *testing.T) {
	fixturez.RunSuite(t, &CoveragePrinterSuite{})
}

func (*CoveragePrinterSuite) TestCoveragePrinter(g *WithT) {
	fixturez.MustBeginOutputCapture(fixturez.OutputSetupStandard, fixturez.GetOutputSetupFatihColor(false), fixturez.OutputSetupRodaineTable)
	defer fixturez.ResetOutputCapture()

	consolez.NewCoveragePrinter().Print(&consolez.Coverage{
		Packages: []*gocov.Package{
			{
				Name: "lowp",
				Functions: []*gocov.Function{
					{
						Statements: []*gocov.Statement{
							{Reached: 1},
							{Reached: 0},
							{Reached: 0},
							{Reached: 0},
						},
					},
				},
			},
			{
				Name: "medp",
				Functions: []*gocov.Function{
					{
						Statements: []*gocov.Statement{
							{Reached: 1},
							{Reached: 1},
							{Reached: 1},
							{Reached: 0},
						},
					},
				},
			},
			{
				Name: "higp",
				Functions: []*gocov.Function{
					{
						Statements: []*gocov.Statement{
							{Reached: 1},
							{Reached: 1},
							{Reached: 1},
							{Reached: 1},
						},
					},
				},
			},
			{
				Name: "nostm",
				Functions: []*gocov.Function{
					{
						Statements: []*gocov.Statement{},
					},
				},
			},
			{
				Name: strings.Repeat("p", 1024),
				Functions: []*gocov.Function{
					{
						Statements: []*gocov.Statement{},
					},
				},
			},
		},
	})

	outBuf, errBuf := fixturez.MustEndOutputCapture()

	g.Expect(outBuf).To(Equal(strings.Join([]string{
		"\x1b[91mLOWC    lowp                                                          25.0% [1/4]\x1b[0m",
		"\x1b[33mMEDC    medp                                                          75.0% [3/4]\x1b[0m",
		"\x1b[32mHIGC    higp                                                         100.0% [4/4]\x1b[0m",
		"\x1b[32mHIGC    nostm                                                        100.0% [0/0]\x1b[0m",
		"\x1b[32mHIGC    ...ppppppppppppppppppppppppppppppppppppppppppppppppppppppppp 100.0% [0/0]\x1b[0m",
		/*   */ "DONE    [LOWC: 1, MEDC: 1, HIGC: 3]                                   66.7% [8/12]",
		"",
	}, "\n")))

	g.Expect(errBuf).To(BeEmpty())
}

func (*CoveragePrinterSuite) TestCoveragePrinterNoStatements(g *WithT) {
	fixturez.MustBeginOutputCapture(fixturez.OutputSetupStandard, fixturez.GetOutputSetupFatihColor(false), fixturez.OutputSetupRodaineTable)
	defer fixturez.ResetOutputCapture()

	consolez.NewCoveragePrinter().Print(&consolez.Coverage{
		Packages: []*gocov.Package{
			{
				Name: "nostm",
				Functions: []*gocov.Function{
					{
						Statements: []*gocov.Statement{},
					},
				},
			},
		},
	})

	outBuf, errBuf := fixturez.MustEndOutputCapture()

	g.Expect(outBuf).To(Equal(strings.Join([]string{
		"\x1b[32mHIGC    nostm                                                        100.0% [0/0]\x1b[0m",
		/*   */ "DONE    [LOWC: 0, MEDC: 0, HIGC: 1]                                  100.0% [0/0]",
		"",
	}, "\n")))

	g.Expect(errBuf).To(BeEmpty())
}
