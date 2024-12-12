package consolez

import (
	"testing"

	"github.com/ibrt/golang-utils/fixturez"
	. "github.com/onsi/gomega"
)

type Suite struct {
	// intentionally empty
}

func TestSuite(t *testing.T) {
	fixturez.RunSuite(t, &Suite{})
}

func (*Suite) TestAlignRight(g *WithT) {
	g.Expect(alignRight("abcd", 4)).To(Equal("abcd"))
	g.Expect(alignRight("cd", 2)).To(Equal("..cd"))
	g.Expect(alignRight("abcd", 6)).To(Equal("..abcd"))
	g.Expect(alignRight("abcdef", 4)).To(Equal("...f"))
}

func (*Suite) TestTruncateLeft(g *WithT) {
	g.Expect(truncateLeft("abcd", 4)).To(Equal("abcd"))
	g.Expect(truncateLeft("cd", 2)).To(Equal("cd"))
	g.Expect(truncateLeft("abcd", 6)).To(Equal("abcd"))
	g.Expect(truncateLeft("abcdef", 4)).To(Equal("...f"))
}
