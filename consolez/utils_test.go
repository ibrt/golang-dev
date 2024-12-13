package consolez

import (
	"testing"

	"github.com/ibrt/golang-utils/fixturez"
	. "github.com/onsi/gomega"
)

type UtilsSuite struct {
	// intentionally empty
}

func TestUtilsSuite(t *testing.T) {
	fixturez.RunSuite(t, &UtilsSuite{})
}

func (*UtilsSuite) TestAlignRight(g *WithT) {
	g.Expect(alignRight("abcd", 4)).To(Equal("abcd"))
	g.Expect(alignRight("cd", 2)).To(Equal("..cd"))
	g.Expect(alignRight("abcd", 6)).To(Equal("..abcd"))
	g.Expect(alignRight("abcdef", 4)).To(Equal("...f"))
}

func (*UtilsSuite) TestTruncateLeft(g *WithT) {
	g.Expect(truncateLeft("abcd", 4)).To(Equal("abcd"))
	g.Expect(truncateLeft("cd", 2)).To(Equal("cd"))
	g.Expect(truncateLeft("abcd", 6)).To(Equal("abcd"))
	g.Expect(truncateLeft("abcdef", 4)).To(Equal("...f"))
}
