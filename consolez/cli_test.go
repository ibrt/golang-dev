package consolez_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/alecthomas/kong"
	"github.com/ibrt/golang-utils/filez"
	"github.com/ibrt/golang-utils/fixturez"
	"github.com/ibrt/golang-utils/outz"
	. "github.com/onsi/gomega"

	"github.com/ibrt/golang-dev/consolez"
)

type CLISuite struct {
	// intentionally empty
}

func TestCLISuite(t *testing.T) {
	fixturez.RunSuite(t, &CLISuite{})
}

func (s *CLISuite) TestRestoreDefaultCLI(g *WithT) {
	defaultCLI := consolez.DefaultCLI
	consolez.DefaultCLI = consolez.NewCLI()
	consolez.RestoreDefaultCLI()
	g.Expect(consolez.DefaultCLI).To(Equal(defaultCLI))
}

func (*CLISuite) TestTool(g *WithT) {
	type CLI struct {
		Command struct {
			Arg string `arg:"" help:"Positional argument."`
		} `cmd:"" help:"Command"`
		Flag string `flag:"" help:"Flag."`
	}

	k, err := kong.New(&CLI{})
	g.Expect(err).To(Succeed())

	kCtx, err := k.Parse([]string{"--flag=flag-value", "command", "arg-value"})
	g.Expect(err).To(Succeed())

	outz.MustBeginOutputCapture(outz.OutputSetupStandard, outz.GetOutputSetupFatihColor(false), outz.OutputSetupRodaineTable)
	defer outz.ResetOutputCapture()

	consolez.DefaultCLI.Tool("Tool", kCtx)

	outBuf, errBuf := outz.MustEndOutputCapture()
	g.Expect(outBuf).To(Equal(fmt.Sprintf("┌─────────────────┐\n│ %v \x1b[1mTool\x1b[0m command │\n└─────────────────┘\n\n\x1b[1mInput          Value       \n\x1b[22m\x1b[33m--flag=STRING  \x1b[0mflag-value  \n\x1b[33m<arg>          \x1b[0marg-value   \n", consolez.IconRocket)))
	g.Expect(errBuf).To(BeEmpty())
}

func (*CLISuite) TestBanner(g *WithT) {
	outz.MustBeginOutputCapture(outz.OutputSetupStandard, outz.GetOutputSetupFatihColor(false), outz.OutputSetupRodaineTable)
	defer outz.ResetOutputCapture()

	consolez.DefaultCLI.Banner("Title", "tagline")

	outBuf, errBuf := outz.MustEndOutputCapture()
	g.Expect(outBuf).To(Equal(fmt.Sprintf("┌──────────────────┐\n│ %v \x1b[1mTitle\x1b[0m tagline │\n└──────────────────┘\n", consolez.IconRocket)))
	g.Expect(errBuf).To(BeEmpty())
}

func (*CLISuite) TestHeader(g *WithT) {
	outz.MustBeginOutputCapture(outz.OutputSetupStandard, outz.GetOutputSetupFatihColor(false), outz.OutputSetupRodaineTable)
	defer outz.ResetOutputCapture()

	f1 := consolez.DefaultCLI.Header("H1 %v", 1)
	f2 := consolez.DefaultCLI.Header("H2 %v", 1)
	f3 := consolez.DefaultCLI.Header("H3 %v", 1)
	f4 := consolez.DefaultCLI.Header("H3 %v", 1)

	f4()
	f3()
	f3()
	f2()
	f1()

	f1 = consolez.DefaultCLI.Header("H1 %v", 2)

	f1()

	outBuf, errBuf := outz.MustEndOutputCapture()

	g.Expect(outBuf).To(Equal(strings.Join(
		[]string{
			fmt.Sprintf("\n%v \x1b[1mH1 1\x1b[0m\n", consolez.IconHighVoltage),
			fmt.Sprintf("\n%v H2 1\n", consolez.IconBackhandIndexPointingRight),
			"\x1b[1;2m—— \x1b[0m\x1b[1;2mH3 1\x1b[0m\n",
			"\x1b[1;2m—— \x1b[0m\x1b[1;2mH3 1\x1b[0m\n",
			fmt.Sprintf("\n%v \x1b[1mH1 2\x1b[0m\n", consolez.IconHighVoltage),
		}, "")))

	g.Expect(errBuf).To(BeEmpty())
}

func (*CLISuite) TestWithHeader(g *WithT) {
	outz.MustBeginOutputCapture(outz.OutputSetupStandard, outz.GetOutputSetupFatihColor(false), outz.OutputSetupRodaineTable)
	defer outz.ResetOutputCapture()

	consolez.DefaultCLI.WithHeader(
		"H1 %v", []any{1},
		func() {
			consolez.DefaultCLI.WithHeader(
				"H2 %v", []any{1},
				func() {
					consolez.DefaultCLI.WithHeader(
						"H3 %v", []any{1},
						func() {
							consolez.DefaultCLI.WithHeader(
								"H3 %v", []any{1},
								func() {
									// intentionally empty
								})
						})
				})

		})

	consolez.DefaultCLI.WithHeader(
		"H1 %v", []any{2},
		func() {
			// intentionally empty
		})

	outBuf, errBuf := outz.MustEndOutputCapture()

	g.Expect(outBuf).To(Equal(strings.Join(
		[]string{
			fmt.Sprintf("\n%v \x1b[1mH1 1\x1b[0m\n", consolez.IconHighVoltage),
			fmt.Sprintf("\n%v H2 1\n", consolez.IconBackhandIndexPointingRight),
			"\x1b[1;2m—— \x1b[0m\x1b[1;2mH3 1\x1b[0m\n",
			"\x1b[1;2m—— \x1b[0m\x1b[1;2mH3 1\x1b[0m\n",
			fmt.Sprintf("\n%v \x1b[1mH1 2\x1b[0m\n", consolez.IconHighVoltage),
		}, "")))

	g.Expect(errBuf).To(BeEmpty())
}

func (*CLISuite) TestNotice(g *WithT) {
	outz.MustBeginOutputCapture(outz.OutputSetupStandard, outz.GetOutputSetupFatihColor(false), outz.OutputSetupRodaineTable)
	defer outz.ResetOutputCapture()

	consolez.DefaultCLI.Notice("scope", "p1", "p2", "p3")

	outBuf, errBuf := outz.MustEndOutputCapture()
	g.Expect(outBuf).To(Equal("\x1b[2m[...................scope]\x1b[0m\x1b[0m p1\x1b[0m\x1b[2m p2\x1b[0m\x1b[2m p3\x1b[0m\n"))
	g.Expect(errBuf).To(BeEmpty())
}

func (*CLISuite) TestCommand_Rel(g *WithT) {
	outz.MustBeginOutputCapture(outz.OutputSetupStandard, outz.GetOutputSetupFatihColor(false), outz.OutputSetupRodaineTable)
	defer outz.ResetOutputCapture()

	consolez.DefaultCLI.Command("cmd", "p1", "p2")

	outBuf, errBuf := outz.MustEndOutputCapture()
	g.Expect(outBuf).To(Equal(fmt.Sprintf("%v cmd \x1b[2mp1 p2\x1b[0m\n", consolez.IconRunner)))
	g.Expect(errBuf).To(BeEmpty())
}

func (*CLISuite) TestCommand_Abs(g *WithT) {
	outz.MustBeginOutputCapture(outz.OutputSetupStandard, outz.GetOutputSetupFatihColor(false), outz.OutputSetupRodaineTable)
	defer outz.ResetOutputCapture()

	consolez.DefaultCLI.Command(filez.MustAbs("cmd"), "p1", "p2")

	outBuf, errBuf := outz.MustEndOutputCapture()
	g.Expect(outBuf).To(Equal(fmt.Sprintf("%v cmd \x1b[2mp1 p2\x1b[0m\n", consolez.IconRunner)))
	g.Expect(errBuf).To(BeEmpty())
}

func (*CLISuite) TestNewTable(g *WithT) {
	outz.MustBeginOutputCapture(outz.OutputSetupStandard, outz.GetOutputSetupFatihColor(false), outz.OutputSetupRodaineTable)
	defer outz.ResetOutputCapture()

	consolez.DefaultCLI.NewTable("A", "B").AddRow("a", "b").Print()

	outBuf, errBuf := outz.MustEndOutputCapture()
	g.Expect(outBuf).To(Equal("\x1b[1mA  B  \n\x1b[22m\x1b[33ma  \x1b[0mb  \n"))
	g.Expect(errBuf).To(BeEmpty())
}

func (*CLISuite) TestError_DebugFalse(g *WithT) {
	outz.MustBeginOutputCapture(outz.OutputSetupStandard, outz.GetOutputSetupFatihColor(false), outz.OutputSetupRodaineTable)
	defer outz.ResetOutputCapture()

	consolez.DefaultCLI.Error(fmt.Errorf("test error"), false)

	outBuf, errBuf := outz.MustEndOutputCapture()
	g.Expect(outBuf).To(Equal(fmt.Sprintf("\n%v \x1b[1mError\x1b[22m\n\x1b[91mtest error\x1b[0m\n", consolez.IconCollision)))
	g.Expect(errBuf).To(BeEmpty())
}

func (*CLISuite) TestError_DebugTrue(g *WithT) {
	outz.MustBeginOutputCapture(outz.OutputSetupStandard, outz.GetOutputSetupFatihColor(false), outz.OutputSetupRodaineTable)
	defer outz.ResetOutputCapture()

	consolez.DefaultCLI.Error(fmt.Errorf("test error"), true)

	outBuf, errBuf := outz.MustEndOutputCapture()
	g.Expect(outBuf).To(HavePrefix(fmt.Sprintf("\n%v \x1b[1mError\x1b[22m\n\x1b[91mtest error\x1b[0m\n(errorz.dump)", consolez.IconCollision)))
	g.Expect(errBuf).To(BeEmpty())
}

func (*CLISuite) TestRecover(g *WithT) {
	c := consolez.NewCLI().
		SetExit(func(code int) {
			g.Expect(code).To(Equal(1))
		})

	outz.MustBeginOutputCapture(outz.OutputSetupStandard, outz.GetOutputSetupFatihColor(false), outz.OutputSetupRodaineTable)
	defer outz.ResetOutputCapture()

	defer func() {
		outBuf, errBuf := outz.MustEndOutputCapture()
		g.Expect(outBuf).To(HavePrefix(fmt.Sprintf("\n%v \x1b[1mError\x1b[22m\n\x1b[91mtest panic\x1b[0m\n(errorz.dump)", consolez.IconCollision)))
		g.Expect(errBuf).To(BeEmpty())
	}()

	defer c.Recover(true)
	panic("test panic")
}
