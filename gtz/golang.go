package gtz

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/ibrt/golang-utils/errorz"
	"github.com/ibrt/golang-utils/filez"
	"github.com/ibrt/golang-utils/jsonz"
	"github.com/ibrt/golang-utils/memz"
	"github.com/ibrt/golang-utils/stringz"

	"github.com/ibrt/golang-dev/consolez"
	"github.com/ibrt/golang-dev/shellz"
)

// Known Go tools.
var (
	GoToolGoCov       = NewGoTool("github.com/axw/gocov", "gocov", "v1.2.1")
	GoToolGoCovHTML   = NewGoTool("github.com/matm/gocov-html", "cmd/gocov-html", "v1.4.0")
	GoToolGolint      = NewGoTool("golang.org/x/lint", "golint", "v0.0.0-20210508222113-6edffad5e616")
	GoToolMockGen     = NewGoTool("go.uber.org/mock", "/mockgen", "v0.5.0")
	GoToolStaticCheck = NewGoTool("honnef.co/go/tools", "cmd/staticcheck", "2024.1.1")
)

// MustLookupGoTool returns a *GoTool
func MustLookupGoTool(key string) *GoTool {
	switch key {
	case "go-cov":
		return GoToolGoCov
	case "go-cov-html":
		return GoToolGoCovHTML
	case "golint":
		return GoToolGolint
	case "mock-gen":
		return GoToolMockGen
	case "static-check":
		return GoToolStaticCheck
	default:
		panic(errorz.Errorf("unknown go tool: %v", key))
	}
}

// GoTool describes a Go tool.
type GoTool struct {
	m              *sync.Mutex
	pkg            string
	path           string
	defaultVersion string
	currentVersion string
}

// NewGoTool initializes a new Go tool.
func NewGoTool(pkg, path, defaultVersion string) *GoTool {
	if path != "" {
		path = stringz.EnsurePrefix(path, "/")
	}

	return &GoTool{
		m:              &sync.Mutex{},
		pkg:            strings.TrimSuffix(pkg, "/"),
		path:           path,
		defaultVersion: defaultVersion,
		currentVersion: "",
	}
}

// GetVersion returns the current Go tool defaultVersion (either pinned, or from go.mod).
func (t *GoTool) GetVersion() string {
	t.m.Lock()
	defer t.m.Unlock()
	return t.getVersion()
}

// GetPackage returns the Go package.
func (t *GoTool) GetPackage() string {
	t.m.Lock()
	defer t.m.Unlock()
	return t.pkg
}

// GetArgument returns an argument string suitable for Go run.
func (t *GoTool) GetArgument() string {
	t.m.Lock()
	defer t.m.Unlock()
	return fmt.Sprintf("%v%v@%v", t.pkg, t.path, t.getVersion())
}

// GetCommand returns a *shellz.Command for this Go tool.
func (t *GoTool) GetCommand() *shellz.Command {
	return shellz.NewCommand("go", "run").AddParams(t.GetArgument())
}

// MustRun runs the Go tool.
func (t *GoTool) MustRun(params ...string) {
	t.GetCommand().AddParams(params...).MustRun()
}

func (t *GoTool) getVersion() string {
	if t.currentVersion != "" {
		return t.currentVersion
	}

	func() {
		defer func() { recover() }()
		out := shellz.NewCommand("go", "list", "-m", t.pkg).SetEcho(false).MustCombinedOutputString()
		t.currentVersion = strings.TrimSpace(strings.TrimPrefix(out, t.pkg))
	}()

	if t.currentVersion == "" {
		if t.defaultVersion != "" {
			t.currentVersion = t.defaultVersion
		} else {
			t.currentVersion = "latest"
		}
	}

	return t.currentVersion
}

// GoChecksParams describes the parameters for running Go checks.
type GoChecksParams struct {
	AllPackages []string
	BuildTags   []string
}

// MustRunGoChecks runs a set of Go checks in the current working directory.
func MustRunGoChecks(params *GoChecksParams) {
	errorz.Assertf(len(params.AllPackages) > 0, "missing all packages")

	consolez.DefaultCLI.Notice("go-checks", "preparing...")

	shellz.NewCommand("go", "mod", "tidy").MustRun()
	shellz.NewCommand("go", "generate").AddParams(params.AllPackages...).MustRun()
	shellz.NewCommand("go", "fmt").AddParams(params.AllPackages...).MustRun()

	consolez.DefaultCLI.Notice("go-checks", "building...")

	shellz.NewCommand("go", "build", "-v").
		AddParamsIfTrue(len(params.BuildTags) > 0, fmt.Sprintf("-tags=%v", strings.Join(params.BuildTags, ","))).
		AddParams(params.AllPackages...).
		MustRun()

	consolez.DefaultCLI.Notice("go-checks", "linting...")

	GoToolGolint.
		GetCommand().
		AddParams("-set_exit_status").
		AddParams(params.AllPackages...).
		MustRun()

	shellz.NewCommand("go", "vet").
		AddParams(params.AllPackages...).
		MustRun()

	GoToolStaticCheck.
		GetCommand().
		AddParams(params.AllPackages...).
		MustRun()

	shellz.NewCommand("go", "mod", "tidy").
		MustRun()
}

// GoTestsParams describes the parameters for running Go tests.
type GoTestsParams struct {
	AllPackages      []string
	SelectedPackages []string
	BuildTags        []string
	TestRegexp       string
	IgnoreCache      bool
	Verbose          *bool
	CoverageDirPath  string
	OpenCoverage     bool
}

// MustRunGoTests runs a set of Go tests in the current working directory.
func MustRunGoTests(params *GoTestsParams) {
	errorz.Assertf(len(params.AllPackages) > 0, "missing all packages")
	errorz.Assertf(params.CoverageDirPath != "", "missing coverage dir path")

	consolez.DefaultCLI.Notice("go-tests", "preparing coverage directory...")

	filez.MustPrepareDir(params.CoverageDirPath, 0777)

	consolez.DefaultCLI.Notice("go-tests", "generating Go code...")

	shellz.NewCommand("go", "generate").
		AddParams(params.AllPackages...).
		MustRun()

	cmd := shellz.NewCommand("go", "test").
		AddParams("-trimpath", "-race", "-failfast", "-shuffle=on").
		AddParams("-covermode=atomic", fmt.Sprintf("-coverprofile=%v", filepath.Join(params.CoverageDirPath, "coverage.out")))

	if params.IgnoreCache {
		cmd = cmd.AddParams("-count=1")
	}

	if params.TestRegexp != "" {
		cmd = cmd.AddParams(fmt.Sprintf("-run=%v", params.TestRegexp))
	}

	if memz.ValNilToZero(params.Verbose) ||
		(params.Verbose == nil && len(params.SelectedPackages) == 1 && !strings.HasSuffix(params.SelectedPackages[0], "...")) {
		cmd = cmd.AddParams("-v")
	}

	consolez.DefaultCLI.Notice("go-tests", "running tests...")

	if len(params.SelectedPackages) > 0 {
		cmd = cmd.AddParams(params.SelectedPackages...)
	} else {
		cmd = cmd.AddParams(params.AllPackages...)
	}

	p := consolez.NewGoTestPrinter()
	cmd.MustLines(p.PrintLine)
	p.PrintDone()

	coverageJSON := processGoCoverage(params)
	consolez.NewCoveragePrinter().Print(jsonz.MustUnmarshal[*consolez.Coverage](coverageJSON))

	if params.OpenCoverage {
		openGoCoverage(params, coverageJSON)
	}
}

func processGoCoverage(params *GoTestsParams) []byte {
	consolez.DefaultCLI.Notice("go-tests", "processing coverage...")

	coverageOutLines := memz.FilterSlice(
		strings.Split(filez.MustReadFileString(filepath.Join(params.CoverageDirPath, "coverage.out")), "\n"),
		func(l string) bool {
			return !strings.Contains(l, ".gen.go:") &&
				!strings.Contains(l, "_generated.go:") &&
				!strings.Contains(l, ".nocov.go:")
		})

	filez.MustWriteFileString(
		filepath.Join(params.CoverageDirPath, "coverage.out"),
		0777, 0666,
		strings.Join(coverageOutLines, "\n"))

	coverageJSON := GoToolGoCov.
		GetCommand().
		AddParams("convert", filepath.Join(params.CoverageDirPath, "coverage.out")).
		MustOutput(false)

	filez.MustWriteFile(
		filepath.Join(params.CoverageDirPath, "coverage.json"),
		0777, 0666,
		coverageJSON)

	return coverageJSON
}

func openGoCoverage(params *GoTestsParams, coverageJSON []byte) {
	consolez.DefaultCLI.Notice("go-tests", "opening coverage...")

	coverageHTML := GoToolGoCovHTML.
		GetCommand().
		AddParams("-t", "golang").
		SetIn(bytes.NewReader(coverageJSON)).
		MustOutput(false)

	filez.MustWriteFile(
		filepath.Join(params.CoverageDirPath, "coverage.html"),
		0777, 0666,
		coverageHTML)

	shellz.NewCommand("open").
		AddParams(filepath.Join(params.CoverageDirPath, "coverage.html")).
		MustRun()
}
