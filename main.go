package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"os/user"
	"strings"
	"syscall"
	"time"

	"github.com/qdm12/golibs/command"
	"github.com/qdm12/welcome/pkg/display"
	"github.com/qdm12/welcome/pkg/docker"
	"github.com/qdm12/welcome/pkg/hardware"
	"github.com/qdm12/welcome/pkg/network"
)

func main() {
	os.Exit(_main(context.Background()))
}

type (
	Hardware interface {
		PartitionsUsage(ctx context.Context) (partitionsUsage map[string]int, warnings []string, err error)
		IsZpoolInstalled(ctx context.Context) bool
		GetPools(ctx context.Context) (poolNames []string, err error)
		GetPoolHealth(ctx context.Context, poolName string) (health string, err error)
		GetPoolErrors(ctx context.Context, poolName string) (errors string, err error)
		GetPoolCapacity(ctx context.Context, poolName string) (capacity int, err error)
	}

	Docker interface {
		IsRunning(ctx context.Context) (running bool)
		Version(ctx context.Context) string
		ComposeVersion(ctx context.Context) string
		CountContainers(ctx context.Context) (count int, err error)
		AreContainerRunning(ctx context.Context, requiredContainerNames []string) (containersNotRunning []string, err error)
		BadContainers(ctx context.Context) (containerNameToState map[string]string, err error)
	}

	Display interface {
		Error(arg ...interface{})
		Warning(arg ...interface{})
		FormatRandomASCIIArt(s string) string
	}
)

func _main(ctx context.Context) int {
	networkFlag := flag.Bool("network", false, "verify network connectivity")
	composeFlag := flag.Bool("compose", false, "show docker-compose version (slow)")
	requiredContainerNamesFlag := flag.String("requiredContainers", "dns,ddns",
		"comma separated list of required running container names to check for")
	websitesToCheckFlag := flag.String("websitesToCheck",
		"https://qqq.ninja,https://1.1.1.1",
		"comma separated list of websites to check, only enabled if --network is specified")
	flag.Parse()

	requiredContainerNames := strings.Split(*requiredContainerNamesFlag, ",")
	websitesToCheck := strings.Split(*websitesToCheckFlag, ",")

	display := display.New()
	cmd := command.NewCmder()
	docker, err := docker.New(cmd)
	if err != nil {
		display.Error(err)
		return 1
	}
	hardware := hardware.New(cmd, "/var/lib/docker") // TODO docker root path auto-detection
	network, err := network.New(net.DefaultResolver)
	if err != nil {
		display.Error(err)
		return 1
	}

	currentUser, err := user.Current()
	if err != nil {
		display.Error(err)
		return 1
	}
	runningAsRoot := currentUser.Uid == "0"

	signalsCh := make(chan os.Signal, 1)
	signal.Notify(signalsCh,
		syscall.SIGINT,
		syscall.SIGTERM,
		os.Interrupt,
	)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		select {
		case <-ctx.Done():
		case signal := <-signalsCh:
			display.Warning("Caught OS signal %s, exiting...", signal)
			cancel()
		}
	}()

	hostname(display)

	// Time and uptime
	t := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("%s | Up %s\n", t, hardware.Uptime())

	// RAM and CPU usage
	UsageRAM := hardware.RAMPercentUsage()
	UsageCPU, err := hardware.CPUPercentUsage()
	if err != nil {
		display.Error(err)
		UsageCPU = -1
	}
	fmt.Printf("RAM %d%% | CPU %d%%\n", UsageRAM, UsageCPU)

	zfs(ctx, hardware, display)

	if runningAsRoot {
		partitions(ctx, hardware, display)
	} else {
		display.Warning("ignoring partitions because you are not running as root")
	}

	// OS information
	processesCount := hardware.ProcessesCount()
	fmt.Printf("%d processes running\n", processesCount)

	doDocker(ctx, docker, display, *composeFlag, requiredContainerNames)

	if *networkFlag {
		doNetwork(ctx, display, network, websitesToCheck)
	}

	if ctx.Err() != nil {
		return 1
	}
	return 0
}

func hostname(display Display) {
	bytes, err := ioutil.ReadFile("/proc/sys/kernel/hostname")
	if err != nil {
		display.Error("Cannot get hostname: %s", err)
		return
	}
	hostname := string(bytes)
	hostname = strings.TrimSpace(hostname)
	hostname = strings.TrimSuffix(hostname, "\n")
	fmt.Println(display.FormatRandomASCIIArt(hostname))
}

func zfs(ctx context.Context, hardware Hardware, display Display) {
	if !hardware.IsZpoolInstalled(ctx) {
		return
	}
	capacities := []string{}
	pools, err := hardware.GetPools(ctx)
	if err != nil {
		display.Error(err)
	}
	for _, pool := range pools {
		health, err := hardware.GetPoolHealth(ctx, pool)
		if err != nil {
			display.Error(err)
		} else if health != "" {
			display.Warning("ZFS pool %s is %s", pool, health)
		}
		errors, err := hardware.GetPoolErrors(ctx, pool)
		if err != nil {
			display.Error(err)
		} else if errors != "" {
			display.Warning("ZFS pool %s: %s", pool, errors)
		}
		capacity, err := hardware.GetPoolCapacity(ctx, pool)
		if err != nil {
			display.Error(err)
		}
		capacities = append(capacities, fmt.Sprintf("%s %d%%", pool, capacity))
	}
	fmt.Println(strings.Join(capacities, " | "))
}

func partitions(ctx context.Context, hardware Hardware, display Display) {
	partitionsUsage, warnings, err := hardware.PartitionsUsage(ctx)
	for _, warning := range warnings {
		display.Warning(warning)
	}
	if err != nil {
		display.Error(err)
	} else {
		ss := []string{}
		for fs, use := range partitionsUsage {
			ss = append(ss, fmt.Sprintf("%s %d%%", fs, use))
		}
		s := strings.Join(ss, " | ")
		fmt.Println(s)
	}
}

func doDocker(ctx context.Context, docker Docker, display Display,
	composeCheck bool, requiredContainerNames []string) {
	if !docker.IsRunning(ctx) {
		display.Error("Docker is not running")
		return
	}

	// TODO Get RAM and CPU usage of Docker containers using Docker Go SDK

	dockerData := []string{}
	dockerVersion := docker.Version(ctx)
	dockerData = append(dockerData, "Docker "+dockerVersion)
	containersCount, err := docker.CountContainers(ctx)
	if err != nil {
		display.Error("Cannot count containers: %s", err)
	} else {
		dockerData = append(dockerData, fmt.Sprintf("%d containers", containersCount))
	}
	if composeCheck {
		dockerComposeVersion := docker.ComposeVersion(ctx)
		if dockerComposeVersion == "" {
			display.Warning("Docker Compose is not installed")
		} else {
			dockerData = append(dockerData, "Compose "+dockerComposeVersion)
		}
	}
	fmt.Println(strings.Join(dockerData, " | "))

	containersNotRunning, err := docker.AreContainerRunning(ctx, requiredContainerNames)
	if err != nil {
		display.Error(err)
	}
	for _, container := range containersNotRunning {
		display.Warning("Container %s is not running", container)
	}
	badContainers, err := docker.BadContainers(ctx)
	if err != nil {
		display.Error(err)
	}
	for name, status := range badContainers {
		display.Warning("Container %s: %s", name, status)
	}
}

func doNetwork(ctx context.Context, display Display,
	network network.Interface, websitesToCheck []string) {
	var netData []string
	privateIP, err := network.OutboundIP(ctx)
	if err != nil {
		display.Error(err)
	} else {
		netData = append(netData, privateIP)
	}
	publicIP, err := network.PublicIP(ctx)
	if err != nil {
		display.Error(err)
	} else {
		netData = append(netData, publicIP)
	}
	fmt.Println(strings.Join(netData, " | "))

	errors := network.Check(ctx, websitesToCheck)
	for _, err := range errors {
		if err != nil {
			display.Warning(err)
		}
	}
}
