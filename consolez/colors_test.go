package consolez_test

import (
	"testing"

	"github.com/ibrt/golang-utils/fixturez"
	. "github.com/onsi/gomega"

	"github.com/ibrt/golang-dev/consolez"
)

type ColorsSuite struct {
	// intentionally empty
}

func TestColorsSuite(t *testing.T) {
	fixturez.RunSuite(t, &ColorsSuite{})
}

func (s *ColorsSuite) TestColors(g *WithT) {
	fixturez.MustBeginOutputCapture(fixturez.GetOutputSetupFatihColor(false))
	defer fixturez.ResetOutputCapture()

	g.Expect(consolez.GetColorDefault().Print("default")).Error().To(Succeed())
	g.Expect(consolez.GetColorHighlight().Print("highlight")).Error().To(Succeed())
	g.Expect(consolez.GetColorSecondaryHighlight().Print("secondaryHighlight")).Error().To(Succeed())
	g.Expect(consolez.GetColorSecondary().Print("secondary")).Error().To(Succeed())
	g.Expect(consolez.GetColorInfo().Print("info")).Error().To(Succeed())
	g.Expect(consolez.GetColorSuccess().Print("success")).Error().To(Succeed())
	g.Expect(consolez.GetColorWarning().Print("warning")).Error().To(Succeed())
	g.Expect(consolez.GetColorError().Print("error")).Error().To(Succeed())

	outBuf, errBuf := fixturez.MustEndOutputCapture()
	g.Expect(outBuf).To(Equal("\x1b[0mdefault\x1b[0m\x1b[1mhighlight\x1b[0m\x1b[1;2msecondaryHighlight\x1b[0m\x1b[2msecondary\x1b[0m\x1b[36minfo\x1b[0m\x1b[32msuccess\x1b[0m\x1b[33mwarning\x1b[0m\x1b[91merror\x1b[0m"))
	g.Expect(errBuf).To(BeEmpty())
}
