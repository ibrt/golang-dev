package dcz

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"time"

	ct "github.com/compose-spec/compose-go/types"
	"github.com/fatih/color"
	"github.com/ibrt/golang-utils/errorz"
	"github.com/ibrt/golang-utils/memz"
	"gopkg.in/yaml.v3"

	"github.com/ibrt/golang-dev/shellz"
)

var (
	// DefaultRuntimeGOOS allows to inject different values of "runtime.GOOS" for tests.
	DefaultRuntimeGOOS = runtime.GOOS
)

// RestoreDefaultRuntimeGOOS restores the default value of DefaultRuntimeGOOS.
func RestoreDefaultRuntimeGOOS() {
	DefaultRuntimeGOOS = runtime.GOOS
}

// DockerCompose helps operate a Docker Compose config.
type DockerCompose struct {
	projectName string
	profiles    []string
	cfg         *ct.Config
}

// NewDockerCompose initializes a new Docker Compose.
func NewDockerCompose(cfg *ct.Config) *DockerCompose {
	errorz.Assertf(cfg != nil, "missing config")

	return &DockerCompose{
		projectName: "",
		profiles:    nil,
		cfg:         cfg,
	}
}

// WithProjectName returns a clone of the Docker Compose with the given project name set.
func (d *DockerCompose) WithProjectName(projectName string) *DockerCompose {
	return &DockerCompose{
		projectName: projectName,
		profiles:    memz.ShallowCopySlice(d.profiles),
		cfg:         d.cfg,
	}
}

// WithProfiles returns a clone of the Docker Compose with the given profiles set.
func (d *DockerCompose) WithProfiles(profiles ...string) *DockerCompose {
	return &DockerCompose{
		projectName: d.projectName,
		profiles:    memz.ShallowCopySlice(profiles),
		cfg:         d.cfg,
	}
}

// GetProjectName returns the Docker Compose project name if set.
func (d *DockerCompose) GetProjectName() string {
	return d.projectName
}

// GetProfiles returns the Docker Compose profiles if set.
func (d *DockerCompose) GetProfiles() []string {
	return d.profiles
}

// GetConfig returns the Docker Compose config.
func (d *DockerCompose) GetConfig() *ct.Config {
	return d.cfg
}

// GetMarshaledConfig returns the Docker Compose config marshaled to YAML.
func (d *DockerCompose) GetMarshaledConfig() []byte {
	buf, err := yaml.Marshal(d.cfg)
	errorz.MaybeMustWrap(err)
	return buf
}

// GetUpCommand returns a pre-configured "docker compose up" command.
func (d *DockerCompose) GetUpCommand() *shellz.Command {
	return d.GetCommand().
		AddParams("up", "--detach", "--build", "--pull", "always", "--force-recreate", "--remove-orphans", "--wait")
}

// GetDownCommand returns a pre-configured "docker compose down" command.
func (d *DockerCompose) GetDownCommand() *shellz.Command {
	return d.GetCommand().
		AddParams("down", "--volumes", "--rmi", "local", "--remove-orphans", "--timeout", "5")
}

// GetPSCommand returns a pre-configured "docker compose ps" command.
func (d *DockerCompose) GetPSCommand() *shellz.Command {
	return d.GetCommand().
		AddParams("ps", "--all", "--orphans")
}

// GetCommand returns a pre-configured "docker compose" command.
func (d *DockerCompose) GetCommand() *shellz.Command {
	cmd := shellz.NewCommand("docker", "compose")

	if d.projectName != "" {
		cmd = cmd.AddParams("--project-name", d.projectName)
	}

	for _, profile := range d.profiles {
		cmd = cmd.AddParams("--profile", profile)
	}

	if color.NoColor || os.Getenv("CI") != "" {
		cmd = cmd.AddParams("--ansi", "never", "--progress", "plain")
	}

	return cmd.
		AddParams("-f", "-").
		SetIn(bytes.NewReader(d.GetMarshaledConfig()))
}

// NewDockerComposeConfigDeploy is a Docker Compose config helper.
func NewDockerComposeConfigDeploy(memoryLimitMB int64, replicas uint64) *ct.DeployConfig {
	dc := &ct.DeployConfig{
		RestartPolicy: &ct.RestartPolicy{
			Condition:   "on-failure",
			Delay:       memz.Ptr(ct.Duration(10 * time.Second)),
			MaxAttempts: memz.Ptr[uint64](3),
		},
	}

	if memoryLimitMB > 0 {
		dc.Resources = ct.Resources{
			Limits: &ct.Resource{
				MemoryBytes: ct.UnitBytes(memoryLimitMB * 1024 * 1024),
			},
		}
	}

	if replicas > 1 {
		dc.Replicas = memz.Ptr(replicas)
		dc.EndpointMode = "dnsrr" // DNS Round-Robing
	}

	return dc
}

// NewDockerComposeConfigExtraHosts is a Docker Compose config helper.
func NewDockerComposeConfigExtraHosts(extraHostsMaps ...map[string]string) ct.HostsList {
	hostsList := ct.HostsList{}

	if DefaultRuntimeGOOS == "linux" {
		hostsList["host.docker.internal"] = "host-gateway"
	}

	for _, extraHostsMap := range extraHostsMaps {
		for k, v := range extraHostsMap {
			hostsList[k] = v
		}
	}

	return hostsList
}

// NewDockerComposeConfigHealthCheckShell is a Docker Compose config helper.
func NewDockerComposeConfigHealthCheckShell(format string, a ...any) *ct.HealthCheckConfig {
	return &ct.HealthCheckConfig{
		StartPeriod: memz.Ptr(ct.Duration(30 * time.Second)),
		Interval:    memz.Ptr(ct.Duration(5 * time.Second)),
		Timeout:     memz.Ptr(ct.Duration(3 * time.Second)),
		Retries:     memz.Ptr(uint64(3)),
		Test: ct.HealthCheckTest{
			"CMD-SHELL",
			fmt.Sprintf(format, a...),
		},
	}
}
