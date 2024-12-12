package consolez_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ibrt/golang-utils/fixturez"
	. "github.com/onsi/gomega"

	"github.com/ibrt/golang-dev/consolez"
)

type GoTestPrinterSuite struct {
	// intentionally empty
}

func TestGoTestPrinterSuite(t *testing.T) {
	fixturez.RunSuite(t, &GoTestPrinterSuite{})
}

func (*GoTestPrinterSuite) TestGoTestPrinter(g *WithT) {
	fixturez.MustBeginOutputCapture(fixturez.OutputSetupStandard, fixturez.GetOutputSetupFatihColor(false), fixturez.OutputSetupRodaineTable)
	defer fixturez.ResetOutputCapture()

	p := consolez.NewGoTestPrinter()
	p.PrintLine("other")
	p.PrintLine("ok  \t\tcoverage: 0%")
	p.PrintLine("coverage: 100%")
	p.PrintLine("--- SKIP")
	p.PrintLine("    --- SKIP")
	p.PrintLine("SKIP")
	p.PrintLine("    SKIP")
	p.PrintLine("-test.shuffle=1234")
	p.PrintLine("--- PASS")
	p.PrintLine("    --- PASS")
	p.PrintLine("PASS")
	p.PrintLine("    PASS")
	p.PrintLine("--- FAIL")
	p.PrintLine("    --- FAIL")
	p.PrintLine("FAIL")
	p.PrintLine("    FAIL")
	p.PrintLine("=== RUN")
	p.PrintLine("    === RUN")
	p.PrintLine("?   \tpkgn\t[no test files]")
	p.PrintLine("?   \tunexpected")
	p.PrintLine("ok  \tpkgn\t1s\tcoverage: 100%")
	p.PrintLine(fmt.Sprintf("FAIL\t%v\t1s", strings.Repeat("p", 1024)))
	p.PrintDone()

	outBuf, errBuf := fixturez.MustEndOutputCapture()

	g.Expect(outBuf).To(Equal(strings.Join([]string{
		"other",
		"\x1b[32mPASS    ok                                                           [no tests]\x1b[0m",
		"\x1b[2m--- SKIP\x1b[0m",
		"\x1b[2m    --- SKIP\x1b[0m",
		"\x1b[2mSKIP\x1b[0m",
		"\x1b[2m    SKIP\x1b[0m",
		"\x1b[2m-test.shuffle=1234\x1b[0m",
		"\x1b[32m--- PASS\x1b[0m",
		"\x1b[32m    --- PASS\x1b[0m",
		"\x1b[32mPASS\x1b[0m",
		"\x1b[32m    PASS\x1b[0m",
		"\x1b[91m--- FAIL\x1b[0m",
		"\x1b[91m    --- FAIL\x1b[0m",
		"\x1b[91mFAIL\x1b[0m",
		"\x1b[91m    FAIL\x1b[0m",
		"\x1b[1m=== RUN\x1b[0m",
		"\x1b[1m    === RUN\x1b[0m",
		"\x1b[2mSKIP    pkgn                                                         [no tests]\x1b[0m",
		"?   \tunexpected",
		"\x1b[32mPASS    pkgn                                                         1s        \x1b[0m",
		"\x1b[91mFAIL    ...ppppppppppppppppppppppppppppppppppppppppppppppppppppppppp 1s        \x1b[0m",
		"DONE    [SKIP: 1, PASS: 2]                                           0s        ",
		"",
	}, "\n")))

	g.Expect(errBuf).To(BeEmpty())
}
