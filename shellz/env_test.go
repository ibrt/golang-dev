package shellz_test

import (
	"os"
	"testing"

	"github.com/ibrt/golang-utils/fixturez"
	. "github.com/onsi/gomega"

	"github.com/ibrt/golang-dev/shellz"
)

type EnvironSuite struct {
	// intentionally empty
}

func TestEnvironSuite(t *testing.T) {
	fixturez.RunSuite(t, &EnvironSuite{})
}

func (*EnvironSuite) TestEnviron(g *WithT) {
	g.Expect(shellz.UnmarshalEnviron(
		[]string{
			"PREFIX_K1=v1",
			"PREFIX_K2=v2=x",
			"K3=v3",
		}, "PREFIX_")).
		To(Equal(
			map[string]string{
				"PREFIX_K1": "v1",
				"PREFIX_K2": "v2=x",
			}))

	g.Expect(shellz.UnmarshalEnviron(
		[]string{
			"PREFIX_K1=v1",
			"PREFIX_K2=v2=x",
			"K3=v3",
		}, "")).
		To(Equal(
			map[string]string{
				"PREFIX_K1": "v1",
				"PREFIX_K2": "v2=x",
				"K3":        "v3",
			}))

	g.Expect(shellz.MarshalEnviron(
		map[string]string{
			"PREFIX_K1": "v1",
			"PREFIX_K2": "v2=x",
			"K3":        "v3",
		})).
		To(ConsistOf(
			"PREFIX_K1=v1",
			"PREFIX_K2=v2=x",
			"K3=v3"))
}

func (*EnvironSuite) TestWithEnv(g *WithT) {
	g.Expect(os.Setenv("k1", "v11")).To(Succeed())

	shellz.WithEnv(map[string]string{"k1": "v12", "k2": "v2"}, func() {
		g.Expect(os.Getenv("k1")).To(Equal("v12"))
		g.Expect(os.Getenv("k2")).To(Equal("v2"))
	})

	g.Expect(os.Getenv("k1")).To(Equal("v11"))
	_, ok := os.LookupEnv("k2")
	g.Expect(ok).To(BeFalse())
}
