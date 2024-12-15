package shellz_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/ibrt/golang-utils/errorz"
	"github.com/ibrt/golang-utils/filez"
	"github.com/ibrt/golang-utils/fixturez"
	"github.com/ibrt/golang-utils/outz"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"

	"github.com/ibrt/golang-dev/consolez"
	"github.com/ibrt/golang-dev/shellz"
)

type CommandSuite struct {
	// intentionally empty
}

func TestCommandSuite(t *testing.T) {
	fixturez.RunSuite(t, &CommandSuite{})
}

func (*CommandSuite) TestRun_Success(g *WithT) {
	outz.MustBeginOutputCapture(outz.OutputSetupStandard, outz.GetOutputSetupFatihColor(false), outz.OutputSetupRodaineTable)
	defer outz.ResetOutputCapture()

	g.Expect(shellz.NewCommand("cat", "-b").SetIn(strings.NewReader("input")).Run()).
		To(Succeed())

	outBuf, errBuf := outz.MustEndOutputCapture()
	g.Expect(outBuf).To(Equal(fmt.Sprintf("%v cat \x1b[2m-b\x1b[0m\n     1\tinput", consolez.IconRunner)))
	g.Expect(errBuf).To(BeEmpty())
}

func (*CommandSuite) TestRun_Error(g *WithT) {
	outz.MustBeginOutputCapture(outz.OutputSetupStandard, outz.GetOutputSetupFatihColor(false), outz.OutputSetupRodaineTable)
	defer outz.ResetOutputCapture()

	err := shellz.NewCommand("cat", "cae0e988-f55b-4803-a471-a877b686d1a8").Run()
	g.Expect(err).To(HaveOccurred())

	eErr, ok := errorz.As[*shellz.ExecutionError](err)
	g.Expect(ok).To(BeTrue())
	g.Expect(eErr.GetCommand()).To(Equal("cat"))
	g.Expect(eErr.GetParams()).To(Equal([]string{"cae0e988-f55b-4803-a471-a877b686d1a8"}))
	g.Expect(eErr.GetDir()).To(BeEmpty())
	g.Expect(eErr.GetEnv()).To(BeEmpty())
	g.Expect(eErr.GetExitCode()).To(Equal(1))
	g.Expect(eErr.Error()).To(Equal("execution error: exit status 1"))

	xErr, ok := errorz.As[*exec.ExitError](err)
	g.Expect(ok).To(BeTrue())
	g.Expect(xErr.ExitCode()).To(Equal(1))

	outBuf, errBuf := outz.MustEndOutputCapture()
	g.Expect(outBuf).To(Equal(fmt.Sprintf("%v cat \x1b[2mcae0e988-f55b-4803-a471-a877b686d1a8\x1b[0m\n", consolez.IconRunner)))
	g.Expect(errBuf).To(Equal("cat: cae0e988-f55b-4803-a471-a877b686d1a8: No such file or directory\n"))
}

func (*CommandSuite) TestMustRun_Success(g *WithT) {
	outz.MustBeginOutputCapture(outz.OutputSetupStandard, outz.GetOutputSetupFatihColor(false), outz.OutputSetupRodaineTable)
	defer outz.ResetOutputCapture()

	g.Expect(func() {
		shellz.NewCommand("cat", "-b").SetIn(strings.NewReader("input")).MustRun()
	}).ToNot(Panic())

	outBuf, errBuf := outz.MustEndOutputCapture()
	g.Expect(outBuf).To(Equal(fmt.Sprintf("%v cat \x1b[2m-b\x1b[0m\n     1\tinput", consolez.IconRunner)))
	g.Expect(errBuf).To(BeEmpty())
}

func (*CommandSuite) TestMustRun_Error(g *WithT) {
	outz.MustBeginOutputCapture(outz.OutputSetupStandard, outz.GetOutputSetupFatihColor(false), outz.OutputSetupRodaineTable)
	defer outz.ResetOutputCapture()

	g.Expect(func() {
		shellz.NewCommand("cat", "cae0e988-f55b-4803-a471-a877b686d1a8").MustRun()
	}).To(Panic())

	outBuf, errBuf := outz.MustEndOutputCapture()
	g.Expect(outBuf).To(Equal(fmt.Sprintf("%v cat \x1b[2mcae0e988-f55b-4803-a471-a877b686d1a8\x1b[0m\n", consolez.IconRunner)))
	g.Expect(errBuf).To(Equal("cat: cae0e988-f55b-4803-a471-a877b686d1a8: No such file or directory\n"))
}

func (*CommandSuite) TestOutput_Success(g *WithT) {
	outz.MustBeginOutputCapture(outz.OutputSetupStandard, outz.GetOutputSetupFatihColor(false), outz.OutputSetupRodaineTable)
	defer outz.ResetOutputCapture()

	out, err := shellz.NewCommand("cat", "-b").SetIn(strings.NewReader("input")).Output(true)
	g.Expect(err).To(Succeed())
	g.Expect(out).ToNot(BeNil())
	g.Expect(string(out)).To(Equal("     1\tinput"))

	outBuf, errBuf := outz.MustEndOutputCapture()
	g.Expect(outBuf).To(BeEmpty())
	g.Expect(errBuf).To(BeEmpty())
}

func (*CommandSuite) TestOutput_Error_EchoStderrTrue(g *WithT) {
	outz.MustBeginOutputCapture(outz.OutputSetupStandard, outz.GetOutputSetupFatihColor(false), outz.OutputSetupRodaineTable)
	defer outz.ResetOutputCapture()

	out, err := shellz.NewCommand("cat", "cae0e988-f55b-4803-a471-a877b686d1a8").Output(true)
	g.Expect(out).To(BeEmpty())
	g.Expect(err).To(HaveOccurred())

	eErr, ok := errorz.As[*shellz.ExecutionError](err)
	g.Expect(ok).To(BeTrue())
	g.Expect(eErr.GetCommand()).To(Equal("cat"))
	g.Expect(eErr.GetParams()).To(Equal([]string{"cae0e988-f55b-4803-a471-a877b686d1a8"}))
	g.Expect(eErr.GetDir()).To(BeEmpty())
	g.Expect(eErr.GetEnv()).To(BeEmpty())
	g.Expect(eErr.GetExitCode()).To(Equal(1))
	g.Expect(eErr.GetCapturedStderr()).To(BeEmpty())
	g.Expect(eErr.Error()).To(Equal("execution error: exit status 1"))

	outBuf, errBuf := outz.MustEndOutputCapture()
	g.Expect(outBuf).To(BeEmpty())
	g.Expect(errBuf).To(Equal("cat: cae0e988-f55b-4803-a471-a877b686d1a8: No such file or directory\n"))
}

func (*CommandSuite) TestOutput_Error_EchoStderrFalse(g *WithT) {
	outz.MustBeginOutputCapture(outz.OutputSetupStandard, outz.GetOutputSetupFatihColor(false), outz.OutputSetupRodaineTable)
	defer outz.ResetOutputCapture()

	out, err := shellz.NewCommand("cat", "cae0e988-f55b-4803-a471-a877b686d1a8").Output(false)
	g.Expect(out).To(BeEmpty())
	g.Expect(err).To(HaveOccurred())

	eErr, ok := errorz.As[*shellz.ExecutionError](err)
	g.Expect(ok).To(BeTrue())
	g.Expect(eErr.GetCommand()).To(Equal("cat"))
	g.Expect(eErr.GetParams()).To(Equal([]string{"cae0e988-f55b-4803-a471-a877b686d1a8"}))
	g.Expect(eErr.GetDir()).To(BeEmpty())
	g.Expect(eErr.GetEnv()).To(BeEmpty())
	g.Expect(eErr.GetExitCode()).To(Equal(1))
	g.Expect(eErr.GetCapturedStderr()).To(Equal("cat: cae0e988-f55b-4803-a471-a877b686d1a8: No such file or directory\n"))
	g.Expect(eErr.Error()).To(Equal("execution error: exit status 1"))

	outBuf, errBuf := outz.MustEndOutputCapture()
	g.Expect(outBuf).To(BeEmpty())
	g.Expect(errBuf).To(BeEmpty())
}

func (*CommandSuite) TestMustOutput_Success(g *WithT) {
	outz.MustBeginOutputCapture(outz.OutputSetupStandard, outz.GetOutputSetupFatihColor(false), outz.OutputSetupRodaineTable)
	defer outz.ResetOutputCapture()

	g.Expect(func() {
		out := shellz.NewCommand("cat", "-b").SetIn(strings.NewReader("input")).MustOutput(true)
		g.Expect(out).ToNot(BeNil())
		g.Expect(string(out)).To(Equal("     1\tinput"))
	}).ToNot(Panic())

	outBuf, errBuf := outz.MustEndOutputCapture()
	g.Expect(outBuf).To(BeEmpty())
	g.Expect(errBuf).To(BeEmpty())
}

func (*CommandSuite) TestMustOutput_Error(g *WithT) {
	outz.MustBeginOutputCapture(outz.OutputSetupStandard, outz.GetOutputSetupFatihColor(false), outz.OutputSetupRodaineTable)
	defer outz.ResetOutputCapture()

	g.Expect(func() {
		shellz.NewCommand("cat", "cae0e988-f55b-4803-a471-a877b686d1a8").MustOutput(true)
	}).To(Panic())

	outBuf, errBuf := outz.MustEndOutputCapture()
	g.Expect(outBuf).To(BeEmpty())
	g.Expect(errBuf).To(Equal("cat: cae0e988-f55b-4803-a471-a877b686d1a8: No such file or directory\n"))
}

func (*CommandSuite) TestOutputString_Success(g *WithT) {
	outz.MustBeginOutputCapture(outz.OutputSetupStandard, outz.GetOutputSetupFatihColor(false), outz.OutputSetupRodaineTable)
	defer outz.ResetOutputCapture()

	out, err := shellz.NewCommand("cat", "-b").SetIn(strings.NewReader("input")).OutputString(true)
	g.Expect(err).To(Succeed())
	g.Expect(out).ToNot(BeNil())
	g.Expect(out).To(Equal("     1\tinput"))

	outBuf, errBuf := outz.MustEndOutputCapture()
	g.Expect(outBuf).To(BeEmpty())
	g.Expect(errBuf).To(BeEmpty())
}

func (*CommandSuite) TestOutputString_Error(g *WithT) {
	outz.MustBeginOutputCapture(outz.OutputSetupStandard, outz.GetOutputSetupFatihColor(false), outz.OutputSetupRodaineTable)
	defer outz.ResetOutputCapture()

	out, err := shellz.NewCommand("cat", "cae0e988-f55b-4803-a471-a877b686d1a8").OutputString(true)
	g.Expect(out).To(BeEmpty())
	g.Expect(err).To(HaveOccurred())

	eErr, ok := errorz.As[*shellz.ExecutionError](err)
	g.Expect(ok).To(BeTrue())
	g.Expect(eErr.GetCommand()).To(Equal("cat"))
	g.Expect(eErr.GetParams()).To(Equal([]string{"cae0e988-f55b-4803-a471-a877b686d1a8"}))
	g.Expect(eErr.GetDir()).To(BeEmpty())
	g.Expect(eErr.GetEnv()).To(BeEmpty())
	g.Expect(eErr.GetExitCode()).To(Equal(1))
	g.Expect(eErr.Error()).To(Equal("execution error: exit status 1"))

	outBuf, errBuf := outz.MustEndOutputCapture()
	g.Expect(outBuf).To(BeEmpty())
	g.Expect(errBuf).To(Equal("cat: cae0e988-f55b-4803-a471-a877b686d1a8: No such file or directory\n"))
}

func (*CommandSuite) TestMustOutputString_Success(g *WithT) {
	outz.MustBeginOutputCapture(outz.OutputSetupStandard, outz.GetOutputSetupFatihColor(false), outz.OutputSetupRodaineTable)
	defer outz.ResetOutputCapture()

	g.Expect(func() {
		out := shellz.NewCommand("cat", "-b").SetIn(strings.NewReader("input")).MustOutputString(true)
		g.Expect(out).ToNot(BeNil())
		g.Expect(string(out)).To(Equal("     1\tinput"))
	}).ToNot(Panic())

	outBuf, errBuf := outz.MustEndOutputCapture()
	g.Expect(outBuf).To(BeEmpty())
	g.Expect(errBuf).To(BeEmpty())
}

func (*CommandSuite) TestMustOutputString_Error(g *WithT) {
	outz.MustBeginOutputCapture(outz.OutputSetupStandard, outz.GetOutputSetupFatihColor(false), outz.OutputSetupRodaineTable)
	defer outz.ResetOutputCapture()

	g.Expect(func() {
		shellz.NewCommand("cat", "cae0e988-f55b-4803-a471-a877b686d1a8").MustOutputString(true)
	}).To(Panic())

	outBuf, errBuf := outz.MustEndOutputCapture()
	g.Expect(outBuf).To(BeEmpty())
	g.Expect(errBuf).To(Equal("cat: cae0e988-f55b-4803-a471-a877b686d1a8: No such file or directory\n"))
}

func (*CommandSuite) TestCombinedOutput_Success(g *WithT) {
	outz.MustBeginOutputCapture(outz.OutputSetupStandard, outz.GetOutputSetupFatihColor(false), outz.OutputSetupRodaineTable)
	defer outz.ResetOutputCapture()

	out, err := shellz.NewCommand("cat", "-b").SetIn(strings.NewReader("input")).CombinedOutput()
	g.Expect(err).To(Succeed())
	g.Expect(out).ToNot(BeNil())
	g.Expect(string(out)).To(Equal("     1\tinput"))

	outBuf, errBuf := outz.MustEndOutputCapture()
	g.Expect(outBuf).To(BeEmpty())
	g.Expect(errBuf).To(BeEmpty())
}

func (*CommandSuite) TestMustCombinedOutput_Success(g *WithT) {
	outz.MustBeginOutputCapture(outz.OutputSetupStandard, outz.GetOutputSetupFatihColor(false), outz.OutputSetupRodaineTable)
	defer outz.ResetOutputCapture()

	g.Expect(func() {
		out := shellz.NewCommand("cat", "-b").SetIn(strings.NewReader("input")).MustCombinedOutput()
		g.Expect(out).ToNot(BeNil())
		g.Expect(string(out)).To(Equal("     1\tinput"))
	}).ToNot(Panic())

	outBuf, errBuf := outz.MustEndOutputCapture()
	g.Expect(outBuf).To(BeEmpty())
	g.Expect(errBuf).To(BeEmpty())
}

func (*CommandSuite) TestCombinedOutput_Error(g *WithT) {
	outz.MustBeginOutputCapture(outz.OutputSetupStandard, outz.GetOutputSetupFatihColor(false), outz.OutputSetupRodaineTable)
	defer outz.ResetOutputCapture()

	out, err := shellz.NewCommand("cat", "cae0e988-f55b-4803-a471-a877b686d1a8").CombinedOutput()
	g.Expect(out).To(BeEmpty())
	g.Expect(err).To(HaveOccurred())

	eErr, ok := errorz.As[*shellz.ExecutionError](err)
	g.Expect(ok).To(BeTrue())
	g.Expect(eErr.GetCommand()).To(Equal("cat"))
	g.Expect(eErr.GetParams()).To(Equal([]string{"cae0e988-f55b-4803-a471-a877b686d1a8"}))
	g.Expect(eErr.GetDir()).To(BeEmpty())
	g.Expect(eErr.GetEnv()).To(BeEmpty())
	g.Expect(eErr.GetExitCode()).To(Equal(1))
	g.Expect(eErr.Error()).To(Equal("execution error: exit status 1"))

	outBuf, errBuf := outz.MustEndOutputCapture()
	g.Expect(outBuf).To(BeEmpty())
	g.Expect(errBuf).To(BeEmpty())
}

func (*CommandSuite) TestMustCombinedOutput_Error(g *WithT) {
	outz.MustBeginOutputCapture(outz.OutputSetupStandard, outz.GetOutputSetupFatihColor(false), outz.OutputSetupRodaineTable)
	defer outz.ResetOutputCapture()

	g.Expect(func() {
		shellz.NewCommand("cat", "cae0e988-f55b-4803-a471-a877b686d1a8").MustCombinedOutput()
	}).To(Panic())

	outBuf, errBuf := outz.MustEndOutputCapture()
	g.Expect(outBuf).To(BeEmpty())
	g.Expect(errBuf).To(BeEmpty())
}

func (*CommandSuite) TestCombinedOutputString_Success(g *WithT) {
	outz.MustBeginOutputCapture(outz.OutputSetupStandard, outz.GetOutputSetupFatihColor(false), outz.OutputSetupRodaineTable)
	defer outz.ResetOutputCapture()

	out, err := shellz.NewCommand("cat", "-b").
		SetIn(strings.NewReader("input")).
		CombinedOutputString()
	g.Expect(err).To(Succeed())
	g.Expect(out).ToNot(BeNil())
	g.Expect(out).To(Equal("     1\tinput"))

	outBuf, errBuf := outz.MustEndOutputCapture()
	g.Expect(outBuf).To(BeEmpty())
	g.Expect(errBuf).To(BeEmpty())
}

func (*CommandSuite) TestCombinedOutputString_Error(g *WithT) {
	outz.MustBeginOutputCapture(outz.OutputSetupStandard, outz.GetOutputSetupFatihColor(false), outz.OutputSetupRodaineTable)
	defer outz.ResetOutputCapture()

	out, err := shellz.NewCommand("cat", "cae0e988-f55b-4803-a471-a877b686d1a8").CombinedOutputString()
	g.Expect(out).To(BeEmpty())
	g.Expect(err).To(HaveOccurred())

	eErr, ok := errorz.As[*shellz.ExecutionError](err)
	g.Expect(ok).To(BeTrue())
	g.Expect(eErr.GetCommand()).To(Equal("cat"))
	g.Expect(eErr.GetParams()).To(Equal([]string{"cae0e988-f55b-4803-a471-a877b686d1a8"}))
	g.Expect(eErr.GetDir()).To(BeEmpty())
	g.Expect(eErr.GetEnv()).To(BeEmpty())
	g.Expect(eErr.GetExitCode()).To(Equal(1))
	g.Expect(eErr.Error()).To(Equal("execution error: exit status 1"))

	outBuf, errBuf := outz.MustEndOutputCapture()
	g.Expect(outBuf).To(BeEmpty())
	g.Expect(errBuf).To(BeEmpty())
}

func (*CommandSuite) TestMustCombinedOutputString_Success(g *WithT) {
	outz.MustBeginOutputCapture(outz.OutputSetupStandard, outz.GetOutputSetupFatihColor(false), outz.OutputSetupRodaineTable)
	defer outz.ResetOutputCapture()

	g.Expect(func() {
		g.Expect(shellz.NewCommand("cat", "-b").SetIn(strings.NewReader("input")).MustCombinedOutputString()).
			To(Equal("     1\tinput"))
	}).ToNot(Panic())

	outBuf, errBuf := outz.MustEndOutputCapture()
	g.Expect(outBuf).To(BeEmpty())
	g.Expect(errBuf).To(BeEmpty())
}

func (*CommandSuite) TestMustCombinedOutputString_Error(g *WithT) {
	outz.MustBeginOutputCapture(outz.OutputSetupStandard, outz.GetOutputSetupFatihColor(false), outz.OutputSetupRodaineTable)
	defer outz.ResetOutputCapture()

	g.Expect(func() {
		shellz.NewCommand("cat", "cae0e988-f55b-4803-a471-a877b686d1a8").MustCombinedOutputString()
	}).To(Panic())

	outBuf, errBuf := outz.MustEndOutputCapture()
	g.Expect(outBuf).To(BeEmpty())
	g.Expect(errBuf).To(BeEmpty())
}

func (*CommandSuite) TestLines_Success(g *WithT) {
	outz.MustBeginOutputCapture(outz.OutputSetupStandard, outz.GetOutputSetupFatihColor(false), outz.OutputSetupRodaineTable)
	defer outz.ResetOutputCapture()

	m := &sync.Mutex{}
	receivedLines := make([]string, 0)

	g.Expect(
		shellz.NewCommand("cat").
			SetIn(strings.NewReader("1\n2\n3\n")).
			Lines(func(line string) {
				m.Lock()
				defer m.Unlock()
				receivedLines = append(receivedLines, line)
			})).
		To(Succeed())

	g.Expect(receivedLines).To(Equal([]string{
		"1",
		"2",
		"3",
	}))

	outBuf, errBuf := outz.MustEndOutputCapture()
	g.Expect(outBuf).To(Equal(fmt.Sprintf("%v cat \x1b[2m\x1b[0m\n", consolez.IconRunner)))
	g.Expect(errBuf).To(BeEmpty())
}

func (*CommandSuite) TestLines_Success_LongLine(g *WithT) {
	outz.MustBeginOutputCapture(outz.OutputSetupStandard, outz.GetOutputSetupFatihColor(false), outz.OutputSetupRodaineTable)
	defer outz.ResetOutputCapture()

	longLine := strings.Repeat("x", 8*1024)
	m := &sync.Mutex{}
	receivedLines := make([]string, 0)

	g.Expect(
		shellz.NewCommand("cat").
			SetIn(strings.NewReader(longLine)).
			Lines(func(line string) {
				m.Lock()
				defer m.Unlock()
				receivedLines = append(receivedLines, line)
			})).
		To(Succeed())

	func() {
		m.Lock()
		defer m.Unlock()

		g.Expect(receivedLines).To(Equal([]string{
			longLine,
		}))
	}()

	outBuf, errBuf := outz.MustEndOutputCapture()
	g.Expect(outBuf).To(Equal(fmt.Sprintf("%v cat \x1b[2m\x1b[0m\n", consolez.IconRunner)))
	g.Expect(errBuf).To(BeEmpty())
}

func (*CommandSuite) TestLines_Error(g *WithT) {
	outz.MustBeginOutputCapture(outz.OutputSetupStandard, outz.GetOutputSetupFatihColor(false), outz.OutputSetupRodaineTable)
	defer outz.ResetOutputCapture()

	m := &sync.Mutex{}
	receivedLines := make([]string, 0)

	err := shellz.NewCommand("cat", "cae0e988-f55b-4803-a471-a877b686d1a8").Lines(func(line string) {
		m.Lock()
		defer m.Unlock()
		receivedLines = append(receivedLines, line)
	})
	g.Expect(err).To(HaveOccurred())

	eErr, ok := errorz.As[*shellz.ExecutionError](err)
	g.Expect(ok).To(BeTrue())
	g.Expect(eErr.GetCommand()).To(Equal("cat"))
	g.Expect(eErr.GetParams()).To(Equal([]string{"cae0e988-f55b-4803-a471-a877b686d1a8"}))
	g.Expect(eErr.GetDir()).To(BeEmpty())
	g.Expect(eErr.GetEnv()).To(BeEmpty())
	g.Expect(eErr.GetExitCode()).To(Equal(1))
	g.Expect(eErr.Error()).To(Equal("execution error: exit status 1"))

	g.Expect(receivedLines).To(Equal([]string{
		"cat: cae0e988-f55b-4803-a471-a877b686d1a8: No such file or directory",
	}))

	outBuf, errBuf := outz.MustEndOutputCapture()
	g.Expect(outBuf).To(Equal(fmt.Sprintf("%v cat \x1b[2mcae0e988-f55b-4803-a471-a877b686d1a8\x1b[0m\n", consolez.IconRunner)))
	g.Expect(errBuf).To(BeEmpty())
}

func (*CommandSuite) TestLines_Error_Start(g *WithT) {
	outz.MustBeginOutputCapture(outz.OutputSetupStandard, outz.GetOutputSetupFatihColor(false), outz.OutputSetupRodaineTable)
	defer outz.ResetOutputCapture()

	m := &sync.Mutex{}
	receivedLines := make([]string, 0)

	err := shellz.NewCommand("cae0e988-f55b-4803-a471-a877b686d1a8").Lines(func(line string) {
		m.Lock()
		defer m.Unlock()
		receivedLines = append(receivedLines, line)
	})
	g.Expect(err).To(HaveOccurred())

	eErr, ok := errorz.As[*shellz.ExecutionError](err)
	g.Expect(ok).To(BeTrue())
	g.Expect(eErr.GetCommand()).To(Equal("cae0e988-f55b-4803-a471-a877b686d1a8"))
	g.Expect(eErr.GetParams()).To(BeEmpty())
	g.Expect(eErr.GetDir()).To(BeEmpty())
	g.Expect(eErr.GetEnv()).To(BeEmpty())
	g.Expect(eErr.GetExitCode()).To(Equal(-1))
	g.Expect(eErr.Error()).To(Equal("execution error: exec: \"cae0e988-f55b-4803-a471-a877b686d1a8\": executable file not found in $PATH"))

	g.Expect(receivedLines).To(Equal([]string{}))

	outBuf, errBuf := outz.MustEndOutputCapture()
	g.Expect(outBuf).To(Equal(fmt.Sprintf("%v cae0e988-f55b-4803-a471-a877b686d1a8 \x1b[2m\x1b[0m\n", consolez.IconRunner)))
	g.Expect(errBuf).To(BeEmpty())
}

func (*CommandSuite) TestLines_Error_LongLine(g *WithT) {
	outz.MustBeginOutputCapture(outz.OutputSetupStandard, outz.GetOutputSetupFatihColor(false), outz.OutputSetupRodaineTable)
	defer outz.ResetOutputCapture()

	longLine := strings.Repeat("x", 8*1024)
	m := &sync.Mutex{}
	receivedLines := make([]string, 0)

	err := shellz.NewCommand("cat", longLine).Lines(func(line string) {
		m.Lock()
		defer m.Unlock()
		receivedLines = append(receivedLines, line)
	})
	g.Expect(err).To(HaveOccurred())

	eErr, ok := errorz.As[*shellz.ExecutionError](err)
	g.Expect(ok).To(BeTrue())
	g.Expect(eErr.GetCommand()).To(Equal("cat"))
	g.Expect(eErr.GetParams()).To(Equal([]string{longLine}))
	g.Expect(eErr.GetDir()).To(BeEmpty())
	g.Expect(eErr.GetEnv()).To(BeEmpty())
	g.Expect(eErr.GetExitCode()).To(Equal(1))
	g.Expect(eErr.Error()).To(Equal("execution error: exit status 1"))

	g.Expect(receivedLines).To(Equal([]string{
		fmt.Sprintf("cat: %v: File name too long", longLine),
	}))

	outBuf, errBuf := outz.MustEndOutputCapture()
	g.Expect(outBuf).To(Equal(fmt.Sprintf("%v cat \x1b[2m%v\x1b[0m\n", consolez.IconRunner, longLine)))
	g.Expect(errBuf).To(BeEmpty())
}

func (*CommandSuite) TestMustLines_Success(g *WithT) {
	outz.MustBeginOutputCapture(outz.OutputSetupStandard, outz.GetOutputSetupFatihColor(false), outz.OutputSetupRodaineTable)
	defer outz.ResetOutputCapture()

	m := &sync.Mutex{}
	receivedLines := make([]string, 0)

	g.Expect(
		func() {
			shellz.NewCommand("cat").
				SetIn(strings.NewReader("1\n2\n3\n")).
				MustLines(func(line string) {
					m.Lock()
					defer m.Unlock()
					receivedLines = append(receivedLines, line)
				})
		}).
		ToNot(Panic())

	g.Expect(receivedLines).To(Equal([]string{
		"1",
		"2",
		"3",
	}))

	outBuf, errBuf := outz.MustEndOutputCapture()
	g.Expect(outBuf).To(Equal(fmt.Sprintf("%v cat \x1b[2m\x1b[0m\n", consolez.IconRunner)))
	g.Expect(errBuf).To(BeEmpty())
}

func (*CommandSuite) TestMustLines_Error(g *WithT) {
	outz.MustBeginOutputCapture(outz.OutputSetupStandard, outz.GetOutputSetupFatihColor(false), outz.OutputSetupRodaineTable)
	defer outz.ResetOutputCapture()

	m := &sync.Mutex{}
	receivedLines := make([]string, 0)

	g.Expect(
		func() {
			shellz.NewCommand("cat", "cae0e988-f55b-4803-a471-a877b686d1a8").
				MustLines(func(line string) {
					m.Lock()
					defer m.Unlock()
					receivedLines = append(receivedLines, line)
				})
		}).
		To(Panic())

	g.Expect(receivedLines).To(Equal([]string{
		"cat: cae0e988-f55b-4803-a471-a877b686d1a8: No such file or directory",
	}))

	outBuf, errBuf := outz.MustEndOutputCapture()
	g.Expect(outBuf).To(Equal(fmt.Sprintf("%v cat \x1b[2mcae0e988-f55b-4803-a471-a877b686d1a8\x1b[0m\n", consolez.IconRunner)))
	g.Expect(errBuf).To(BeEmpty())
}

// TestExecExecutor is a mock shellz.Executor used by TestExec.
type TestExecExecutor struct {
	*shellz.RealExecutor
	g    *WithT
	fail bool
}

// SyscallExec implements the shellz.Executor interface.
func (m *TestExecExecutor) SyscallExec(c *shellz.Command, argv0 string, argv []string, envv []string) error {
	m.g.Expect(argv0).To(HaveSuffix("ls"))
	m.g.Expect(argv).To(Equal([]string{"ls", "."}))
	m.g.Expect(envv).To(ContainElement("K=V"))

	if m.fail {
		return errorz.Errorf("test error")
	}

	return nil
}

func (s *CommandSuite) TestExec_Success(g *WithT) {
	shellz.DefaultExecutor = &TestExecExecutor{RealExecutor: &shellz.RealExecutor{}, g: g, fail: false}
	defer shellz.RestoreDefaultExecutor()

	g.Expect(shellz.NewCommand("ls", ".").SetEnv("K", "V").Exec()).To(BeNil())
}

func (*CommandSuite) TestExec_PreparationErrors(g *WithT) {
	g.Expect(shellz.NewCommand("cae0e988-f55b-4803-a471-a877b686d1a8").Exec()).
		To(MatchError(`execution error: exec: "cae0e988-f55b-4803-a471-a877b686d1a8": executable file not found in $PATH`))

	g.Expect(shellz.NewCommand("cat").SetDir("cae0e988-f55b-4803-a471-a877b686d1a8").Exec()).
		To(MatchError(`execution error: chdir cae0e988-f55b-4803-a471-a877b686d1a8: no such file or directory`))
}

func (s *CommandSuite) TestExec_ExecutionError(g *WithT) {
	g.Expect(
		shellz.NewCommand("ls", ".").
			SetExecutor(&TestExecExecutor{
				RealExecutor: &shellz.RealExecutor{},
				g:            g,
				fail:         true,
			}).
			SetEnv("K", "V").
			Exec()).
		To(MatchError("execution error: test error"))
}

func (*CommandSuite) TestMustExec_Error(g *WithT) {
	g.Expect(
		func() {
			shellz.NewCommand("cae0e988-f55b-4803-a471-a877b686d1a8").MustExec()
		}).
		To(PanicWith(MatchError(`execution error: exec: "cae0e988-f55b-4803-a471-a877b686d1a8": executable file not found in $PATH`)))
}

func (*CommandSuite) TestAddParams(g *WithT) {
	outz.MustBeginOutputCapture(outz.OutputSetupStandard, outz.GetOutputSetupFatihColor(false), outz.OutputSetupRodaineTable)
	defer outz.ResetOutputCapture()

	cmd := shellz.NewCommand("cat").SetIn(strings.NewReader("input")).AddParams("-b")
	g.Expect(cmd.GetParams()).To(Equal([]string{"-b"}))
	g.Expect(cmd.Run()).To(Succeed())

	outBuf, errBuf := outz.MustEndOutputCapture()
	g.Expect(outBuf).To(Equal(fmt.Sprintf("%v cat \x1b[2m-b\x1b[0m\n     1\tinput", consolez.IconRunner)))
	g.Expect(errBuf).To(BeEmpty())
}

func (*CommandSuite) TestIfTrueAddParams(g *WithT) {
	g.Expect(
		shellz.NewCommand("cmd").
			AddParamsIfTrue(true, "yes").
			AddParamsIfTrue(false, "no").GetParams()).
		To(Equal([]string{"yes"}))
}

func (*CommandSuite) TestSetDir(g *WithT) {
	outz.MustBeginOutputCapture(outz.OutputSetupStandard, outz.GetOutputSetupFatihColor(false), outz.OutputSetupRodaineTable)
	defer outz.ResetOutputCapture()

	filePath := filez.MustCreateTempFileString("content")
	defer func() { errorz.MaybeMustWrap(os.Remove(filePath)) }()

	dirPath := filepath.Dir(filePath)

	cmd := shellz.NewCommand("cat", filepath.Base(filePath)).SetDir(dirPath).SetEcho(false)
	g.Expect(cmd.GetDir()).To(Equal(dirPath))
	g.Expect(cmd.Run()).To(Succeed())

	outBuf, errBuf := outz.MustEndOutputCapture()
	g.Expect(outBuf).To(Equal("content"))
	g.Expect(errBuf).To(BeEmpty())
}

func (*CommandSuite) TestSetEnv(g *WithT) {
	outz.MustBeginOutputCapture(outz.OutputSetupStandard, outz.GetOutputSetupFatihColor(false), outz.OutputSetupRodaineTable)
	defer outz.ResetOutputCapture()

	cmd := shellz.NewCommand("env").SetEnv("k-cae0e988-f55b-4803-a471-a877b686d1a8", "v-cae0e988-f55b-4803-a471-a877b686d1a8")
	g.Expect(cmd.GetEnv()).To(Equal(map[string]string{"k-cae0e988-f55b-4803-a471-a877b686d1a8": "v-cae0e988-f55b-4803-a471-a877b686d1a8"}))
	g.Expect(cmd.Run()).To(Succeed())

	outBuf, errBuf := outz.MustEndOutputCapture()
	g.Expect(outBuf).To(ContainSubstring("k-cae0e988-f55b-4803-a471-a877b686d1a8=v-cae0e988-f55b-4803-a471-a877b686d1a8"))
	g.Expect(errBuf).To(BeEmpty())
}

func (*CommandSuite) TestMergeEnv(g *WithT) {
	outz.MustBeginOutputCapture(outz.OutputSetupStandard, outz.GetOutputSetupFatihColor(false), outz.OutputSetupRodaineTable)
	defer outz.ResetOutputCapture()

	cmd := shellz.NewCommand("env").MergeEnv(map[string]string{"k-cae0e988-f55b-4803-a471-a877b686d1a8": "v-cae0e988-f55b-4803-a471-a877b686d1a8"})
	g.Expect(cmd.GetEnv()).To(Equal(map[string]string{"k-cae0e988-f55b-4803-a471-a877b686d1a8": "v-cae0e988-f55b-4803-a471-a877b686d1a8"}))
	g.Expect(cmd.Run()).To(Succeed())

	outBuf, errBuf := outz.MustEndOutputCapture()
	g.Expect(outBuf).To(ContainSubstring("k-cae0e988-f55b-4803-a471-a877b686d1a8=v-cae0e988-f55b-4803-a471-a877b686d1a8"))
	g.Expect(errBuf).To(BeEmpty())
}

func (*CommandSuite) TestSetIn(g *WithT) {
	r := strings.NewReader("")
	cmd := shellz.NewCommand("cmd").SetIn(r)
	g.Expect(cmd.GetIn()).To(Equal(r))
}

func (*CommandSuite) TestSetEcho(g *WithT) {
	cmd := shellz.NewCommand("cmd")
	g.Expect(cmd.GetEcho()).To(BeNil())
	cmd = cmd.SetEcho(true)
	g.Expect(cmd.GetEcho()).To(PointTo(BeTrue()))
	cmd = cmd.SetEcho(false)
	g.Expect(cmd.GetEcho()).To(PointTo(BeFalse()))
}
