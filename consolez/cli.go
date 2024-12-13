package consolez

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/alecthomas/kong"
	"github.com/fatih/color"
	"github.com/ibrt/golang-utils/errorz"
	"github.com/ibrt/golang-utils/filez"
	"github.com/rodaine/table"
)

// Known icons.
var (
	IconRocket                     = "\U0001F680"
	IconHighVoltage                = "\U000026A1"
	IconBackhandIndexPointingRight = "\U0001F449"
	IconRunner                     = "\U0001F3C3"
	IconCollision                  = "\U0001F4A5"
)

// Known colors.
var (
	ColorDefault            = color.New(color.Reset)
	ColorHighlight          = color.New(color.Bold)
	ColorSecondaryHighlight = color.New(color.Bold, color.Faint)
	ColorSecondary          = color.New(color.Faint)
	ColorInfo               = color.New(color.FgCyan)
	ColorSuccess            = color.New(color.FgGreen)
	ColorWarning            = color.New(color.FgYellow)
	ColorError              = color.New(color.FgHiRed)
)

var (
	_ CLIOption = CLIOptionFunc(nil)
)

// CLIOption describes a [*CLI] option.
type CLIOption interface {
	Apply(*CLI)
}

// CLIOptionFunc describes a [*CLI] option.
type CLIOptionFunc func(*CLI)

// Apply implements the [CLIOption] interface.
func (f CLIOptionFunc) Apply(c *CLI) {
	f(c)
}

// CLIExit returns a [CLIOptionFunc] that allows to provide an alternative implementation for [os.Exit].
func CLIExit(exit func(code int)) CLIOptionFunc {
	return func(c *CLI) {
		c.exit = exit
	}
}

var (
	defaultCLI = NewCLI()
)

// DefaultCLI is a default, shared instance of [*CLI].
var (
	DefaultCLI = defaultCLI
)

// RestoreDefaultCLI restores the default value of [DefaultCLI].
func RestoreDefaultCLI() {
	DefaultCLI = defaultCLI
}

// CLI provides some utilities for printing messages in CLI tools.
type CLI struct {
	m    *sync.Mutex
	hL   int
	exit func(code int)
}

// NewCLI initializes a new [*CLI].
func NewCLI(options ...CLIOption) *CLI {
	c := &CLI{
		m:    &sync.Mutex{},
		hL:   0,
		exit: os.Exit,
	}

	for _, option := range options {
		option.Apply(c)
	}

	return c
}

// Tool introduces a command line tool.
func (c *CLI) Tool(toolName string, k *kong.Context) {
	commandParts := make([]string, 0)
	options := make([][]string, 0)

	for _, p := range k.Path {
		if p.Command != nil {
			commandParts = append(commandParts, p.Command.Name)
		} else if p.Flag != nil {
			options = append(options, []string{
				p.Flag.Summary(),
				fmt.Sprintf("%v", p.Flag.Target.Interface()),
			})
		} else if p.Positional != nil {
			options = append(options, []string{
				p.Positional.Summary(),
				fmt.Sprintf("%v", p.Positional.Target.Interface()),
			})
		}
	}

	c.Banner(toolName, strings.Join(commandParts, " "))

	if len(options) > 0 {
		fmt.Println()
		c.NewTable("Input", "Value").SetRows(options).Print()
	}
}

// Banner prints a banner.
func (c *CLI) Banner(title, tagLine string) {
	c.m.Lock()
	defer c.m.Unlock()

	fmt.Print("┌", strings.Repeat("─", len(title)+len(tagLine)+6), "┐\n")
	fmt.Print("│ ", IconRocket, " ")
	_, _ = ColorHighlight.Print(title)
	fmt.Print(" ")
	fmt.Print(tagLine)
	fmt.Print(" │\n")
	fmt.Print("└", strings.Repeat("─", len(title)+len(tagLine)+6), "┘\n")
}

// Header prints a header based on a nesting hierarchy.
// Always call the returned function to close the scope, for example by deferring it.
func (c *CLI) Header(format string, a ...any) func() {
	c.m.Lock()
	defer c.m.Unlock()

	if c.hL < 2 {
		fmt.Println()
	}

	switch c.hL {
	case 0:
		fmt.Print(IconHighVoltage)
		fmt.Print(" ")
		_, _ = ColorHighlight.Printf(format, a...)
		fmt.Println()
	case 1:
		fmt.Print(IconBackhandIndexPointingRight)
		fmt.Print(" ")
		fmt.Printf(format, a...)
		fmt.Println()
	default:
		_, _ = ColorSecondaryHighlight.Print("—— ")
		_, _ = ColorSecondaryHighlight.Printf(format, a...)
		fmt.Println()
	}

	c.hL++
	isClosed := false

	return func() {
		c.m.Lock()
		defer c.m.Unlock()

		if !isClosed {
			isClosed = true
			c.hL--
		}
	}
}

// WithHeader calls [*CLI.Header] and runs f() within its scope.
func (c *CLI) WithHeader(format string, a []any, f func()) {
	defer c.Header(format, a...)()
	f()
}

// Notice prints a notice.
func (c *CLI) Notice(scope string, highlight string, secondary ...string) {
	c.m.Lock()
	defer c.m.Unlock()

	_, _ = ColorSecondary.Printf("[%v]", alignRight(scope, 24))
	_, _ = ColorDefault.Print(" ", highlight)

	for _, v := range secondary {
		_, _ = ColorSecondary.Print(" ", v)
	}

	fmt.Println()
}

// Command prints a command.
func (c *CLI) Command(cmd string, params ...string) {
	c.m.Lock()
	defer c.m.Unlock()

	fmt.Print(IconRunner)
	fmt.Printf(" %v ", filez.MustRelForDisplay(cmd))
	_, _ = ColorSecondary.Print(strings.Join(params, " "))
	fmt.Println()
}

// NewTable creates a new table.
func (c *CLI) NewTable(columnHeaders ...any) table.Table {
	return table.
		New(columnHeaders...).
		WithHeaderFormatter(ColorHighlight.SprintfFunc()).
		WithFirstColumnFormatter(ColorWarning.SprintfFunc())
}

// Error prints an error.
func (c *CLI) Error(err error, debug bool) {
	c.m.Lock()
	defer c.m.Unlock()

	fmt.Println()
	fmt.Print(IconCollision)
	fmt.Print(" ")
	_, _ = ColorHighlight.Println("Error")
	_, _ = ColorError.Println(err.Error())

	if debug {
		fmt.Println(errorz.SDump(err))
	}
}

// Recover calls [*CLI.Error] on a recovered panic and exits.
func (c *CLI) Recover(debug bool) {
	if err := errorz.MaybeWrapRecover(recover()); err != nil {
		c.Error(err, debug)
		fmt.Println()
		c.exit(1)
	}

	fmt.Println()
}
