package shellz

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"

	"github.com/ibrt/golang-utils/errorz"
	"github.com/ibrt/golang-utils/memz"

	"github.com/ibrt/golang-dev/consolez"
)

var (
	defaultExecutor = &RealExecutor{}
)

var (
	// DefaultExecutor is the default Executor for commands.
	DefaultExecutor Executor = defaultExecutor
)

// RestoreDefaultExecutor restores the default executor.
func RestoreDefaultExecutor() {
	DefaultExecutor = defaultExecutor
}

// Executor implements the OS-level operations related to a command.
type Executor interface {
	ExecCmdCombinedOutput(c *Command, cmd *exec.Cmd) ([]byte, error)
	ExecCmdOutput(c *Command, cmd *exec.Cmd) ([]byte, error)
	ExecCmdRun(c *Command, cmd *exec.Cmd) error
	ExecCmdStart(c *Command, cmd *exec.Cmd) error
	ExecCmdWait(c *Command, cmd *exec.Cmd) error
	ExecLookPath(c *Command, file string) (string, error)
	OSChdir(c *Command, dir string) error
	SyscallExec(c *Command, argv0 string, argv []string, envv []string) error
}

// RealExecutor implements the Executor interface and actually runs commands on the host.
type RealExecutor struct {
	// intentionally empty
}

// ExecCmdCombinedOutput implements the Executor interface.
func (e *RealExecutor) ExecCmdCombinedOutput(_ *Command, cmd *exec.Cmd) ([]byte, error) {
	return cmd.CombinedOutput()
}

// ExecCmdOutput implements the Executor interface.
func (e *RealExecutor) ExecCmdOutput(_ *Command, cmd *exec.Cmd) ([]byte, error) {
	return cmd.Output()
}

// ExecCmdRun implements the Executor interface.
func (e *RealExecutor) ExecCmdRun(_ *Command, cmd *exec.Cmd) error {
	return cmd.Run()
}

// ExecCmdStart implements the Executor interface.
func (e *RealExecutor) ExecCmdStart(_ *Command, cmd *exec.Cmd) error {
	return cmd.Start()
}

// ExecCmdWait implements the Executor interface.
func (e *RealExecutor) ExecCmdWait(_ *Command, cmd *exec.Cmd) error {
	return cmd.Wait()
}

// ExecLookPath implements the Executor interface.
func (e *RealExecutor) ExecLookPath(_ *Command, file string) (string, error) {
	return exec.LookPath(file)
}

// OSChdir implements the Executor interface.
func (e *RealExecutor) OSChdir(_ *Command, dir string) error {
	return os.Chdir(dir)
}

// SyscallExec implements the Executor interface.
func (e *RealExecutor) SyscallExec(_ *Command, argv0 string, argv []string, envv []string) error {
	return syscall.Exec(argv0, argv, envv)
}

var (
	_ error               = (*ExecutionError)(nil)
	_ errorz.UnwrapSingle = (*ExecutionError)(nil)
)

// ExecutionError describes an error.
type ExecutionError struct {
	cmd            string
	params         []string
	dir            string
	env            map[string]string
	exitCode       int
	capturedStderr string
	err            error
}

// NewExecutionError initializes a new execution error.
func NewExecutionError(err error, c *Command) *ExecutionError {
	e := &ExecutionError{
		cmd:            c.cmd,
		params:         memz.ShallowCopySlice(c.params),
		dir:            c.dir,
		env:            memz.ShallowCopyMap(c.env),
		exitCode:       -1,
		capturedStderr: "",
		err:            err,
	}

	if eErr, ok := errorz.As[*exec.ExitError](err); ok {
		e.exitCode = eErr.ExitCode()

		if len(eErr.Stderr) > 0 {
			e.capturedStderr = string(eErr.Stderr)
		}
	}

	return e
}

// GetCommand returns the originating command.
func (e *ExecutionError) GetCommand() string {
	return e.cmd
}

// GetParams returns the originating params.
func (e *ExecutionError) GetParams() []string {
	return e.params
}

// GetDir returns the originating dir.
func (e *ExecutionError) GetDir() string {
	return e.dir
}

// GetEnv returns the originating env.
func (e *ExecutionError) GetEnv() map[string]string {
	return e.env
}

// GetExitCode returns the originating exit code.
func (e *ExecutionError) GetExitCode() int {
	return e.exitCode
}

// GetCapturedStderr returns the originating captured standard error (if available).
func (e *ExecutionError) GetCapturedStderr() string {
	return e.capturedStderr
}

// Error implements the error interface.
func (e *ExecutionError) Error() string {
	return "execution error: " + e.err.Error()
}

// Unwrap implements the errorz.UnwrapSingle interface.
func (e *ExecutionError) Unwrap() error {
	return e.err
}

// Command describes a command to be spawned in a shell.
type Command struct {
	cmd    string
	params []string

	dir      string
	env      map[string]string
	in       io.Reader
	echo     *bool
	executor Executor
}

// NewCommand creates a new Command.
func NewCommand(cmd string, params ...string) *Command {
	return &Command{
		cmd:      cmd,
		params:   memz.ShallowCopySlice(params),
		env:      make(map[string]string),
		executor: DefaultExecutor,
	}
}

// AddParams adds the given params to the command.
func (c *Command) AddParams(params ...string) *Command {
	cc := c.clone()
	cc.params = append(cc.params, params...)
	return cc
}

// AddParamsIfTrue adds the given params to the command if the condition is true.
func (c *Command) AddParamsIfTrue(cond bool, params ...string) *Command {
	if cond {
		return c.AddParams(params...)
	}

	return c
}

// GetParams returns the current params.
func (c *Command) GetParams() []string {
	return memz.ShallowCopySlice(c.params)
}

// SetDir sets the working directory on the command.
func (c *Command) SetDir(dir string) *Command {
	cc := c.clone()
	cc.dir = dir
	return cc
}

// GetDir returns the current dir.
func (c *Command) GetDir() string {
	return c.dir
}

// SetEnv sets an environment variable on the command.
func (c *Command) SetEnv(k, v string) *Command {
	cc := c.clone()
	cc.env[k] = v
	return cc
}

// MergeEnv sets all the environment variables on the command.
func (c *Command) MergeEnv(env map[string]string) *Command {
	cc := c.clone()
	cc.env = memz.MergeMaps(cc.env, env)
	return cc
}

// GetEnv returns the current env.
func (c *Command) GetEnv() map[string]string {
	return memz.ShallowCopyMap(c.env)
}

// SetIn sets the input to the command.
func (c *Command) SetIn(in io.Reader) *Command {
	cc := c.clone()
	cc.in = in
	return cc
}

// GetIn returns the current input.
func (c *Command) GetIn() io.Reader {
	return c.in
}

// SetEcho configures echo.
func (c *Command) SetEcho(echo bool) *Command {
	cc := c.clone()
	cc.echo = memz.Ptr(echo)
	return cc
}

// GetEcho returns the current echo configuration.
func (c *Command) GetEcho() *bool {
	if c.echo == nil {
		return nil
	}
	return memz.Ptr(*c.echo)
}

// SetExecutor sets the Executor for the command.
func (c *Command) SetExecutor(executor Executor) *Command {
	cc := c.clone()
	cc.executor = executor
	return cc
}

// Run runs the command.
func (c *Command) Run() error {
	c.maybeEcho(true)
	cmd := c.newCmd()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := c.executor.ExecCmdRun(c, cmd); err != nil {
		return NewExecutionError(err, c)
	}

	return nil
}

// MustRun is like run but panics on error.
func (c *Command) MustRun() {
	errorz.MaybeMustWrap(c.Run())
}

// Output runs the command and returns a buffer containing the resulting standard output.
// Standard error is not redirected.
func (c *Command) Output(echoStderr bool) ([]byte, error) {
	c.maybeEcho(false)
	cmd := c.newCmd()

	if echoStderr {
		cmd.Stderr = os.Stderr
	} else {
		cmd.Stderr = nil
	}

	out, err := c.executor.ExecCmdOutput(c, cmd)
	if err != nil {
		return nil, NewExecutionError(err, c)
	}

	return out, nil
}

// MustOutput is like Output but panics on error.
func (c *Command) MustOutput(echoStderr bool) []byte {
	out, err := c.Output(echoStderr)
	errorz.MaybeMustWrap(err)
	return out
}

// OutputString is like Output but returns a string.
func (c *Command) OutputString(echoStderr bool) (string, error) {
	buf, err := c.Output(echoStderr)
	if err != nil {
		return "", errorz.Wrap(err)
	}

	return string(buf), nil
}

// MustOutputString is like OutputString but panics on error.
func (c *Command) MustOutputString(echoStderr bool) string {
	buf, err := c.OutputString(echoStderr)
	errorz.MaybeMustWrap(err)
	return buf
}

// CombinedOutput runs the command and returns a buffer containing the resulting combined standard output and error.
func (c *Command) CombinedOutput() ([]byte, error) {
	c.maybeEcho(false)

	out, err := c.executor.ExecCmdCombinedOutput(c, c.newCmd())
	if err != nil {
		return nil, NewExecutionError(err, c)
	}

	return out, nil
}

// MustCombinedOutput is like CombinedOutput but panics on error.
func (c *Command) MustCombinedOutput() []byte {
	out, err := c.CombinedOutput()
	errorz.MaybeMustWrap(err)
	return out
}

// CombinedOutputString is like CombinedOutput but returns a string.
func (c *Command) CombinedOutputString() (string, error) {
	buf, err := c.CombinedOutput()
	if err != nil {
		return "", errorz.Wrap(err)
	}

	return string(buf), nil
}

// MustCombinedOutputString is like CombinedOutputString but panics on error.
func (c *Command) MustCombinedOutputString() string {
	buf, err := c.CombinedOutputString()
	errorz.MaybeMustWrap(err)
	return buf
}

// Lines runs the command and calls "lineFunc" with each line of output.
func (c *Command) Lines(lineFunc func(string)) error {
	c.maybeEcho(true)
	cmd := c.newCmd()

	outR, err := cmd.StdoutPipe()
	errorz.MaybeMustWrap(err)

	errR, err := cmd.StderrPipe()
	errorz.MaybeMustWrap(err)

	m := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	wg.Add(2)

	callLineFunc := func(line string) {
		m.Lock()
		defer m.Unlock()
		lineFunc(line)
	}

	go c.handleLines(wg, outR, callLineFunc)
	go c.handleLines(wg, errR, callLineFunc)

	if err := c.executor.ExecCmdStart(c, cmd); err != nil {
		return NewExecutionError(err, c)
	}

	wg.Wait()

	if err := c.executor.ExecCmdWait(c, cmd); err != nil {
		return NewExecutionError(err, c)
	}

	return nil
}

func (c *Command) handleLines(wg *sync.WaitGroup, r io.Reader, lineFunc func(string)) {
	defer wg.Done()
	defer func() { recover() }()

	br := bufio.NewReader(r)
	s := &strings.Builder{}

	for {
		if buf, isPrefix, err := br.ReadLine(); err == nil {
			_, _ = s.Write(buf)

			if !isPrefix {
				lineFunc(s.String())
				s.Reset()
			}
		} else {
			break
		}
	}

	if s.Len() > 0 {
		lineFunc(s.String())
		s.Reset()
	}
}

// MustLines is like lines but panics on error.
func (c *Command) MustLines(lineFunc func(string)) {
	errorz.MaybeMustWrap(c.Lines(lineFunc))
}

// Exec execs the command (i.e. replaces the current process).
func (c *Command) Exec() error {
	c.maybeEcho(true)

	binFilePath, err := c.executor.ExecLookPath(c, c.cmd)
	if err != nil {
		return NewExecutionError(err, c)
	}

	if c.dir != "" {
		if err := c.executor.OSChdir(c, c.dir); err != nil {
			return NewExecutionError(err, c)
		}
	}

	if err := c.executor.SyscallExec(c, binFilePath, append([]string{c.cmd}, c.params...), c.newEnviron()); err != nil {
		return NewExecutionError(err, c)
	}

	return nil // Note: unreachable with default implementation.
}

// MustExec is like Exec but panics on error.
func (c *Command) MustExec() {
	errorz.MaybeMustWrap(c.Exec())
}

func (c *Command) maybeEcho(defaultEcho bool) {
	if (c.echo == nil && !defaultEcho) || (c.echo != nil && !*c.echo) {
		return
	}

	consolez.DefaultCLI.Command(c.cmd, c.params...)
}

func (c *Command) newCmd() *exec.Cmd {
	cmd := exec.Command(c.cmd, c.params...)
	cmd.Dir = c.dir
	cmd.Env = c.newEnviron()
	cmd.Stdin = c.in
	return cmd
}

func (c *Command) newEnviron() []string {
	env := os.Environ()

	for k, v := range c.env {
		env = append(env, fmt.Sprintf("%v=%v", k, v))
	}

	return env
}

func (c *Command) clone() *Command {
	cc := &Command{
		cmd:      c.cmd,
		params:   memz.ShallowCopySlice(c.params),
		dir:      c.dir,
		env:      memz.ShallowCopyMap(c.env),
		in:       c.in,
		echo:     nil,
		executor: c.executor,
	}

	if c.echo != nil {
		cc.echo = memz.Ptr(*c.echo)
	}

	return cc
}
