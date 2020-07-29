package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
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

func _main(ctx context.Context) int {
	networkFlag := flag.Bool("network", false, "verify network connectivity")
	composeFlag := flag.Bool("compose", false, "show docker-compose version (slow)")
	requiredContainerNamesFlag := flag.String("requiredContainers", "dns,ddns", "comma separated list of required running container names to check for")
	websitesToCheckFlag := flag.String("websitesToCheck", "https://qqq.ninja,https://1.1.1.1", "comma separated list of websites to check, only enabled if --network is specified")
	flag.Parse()

	requiredContainerNames := strings.Split(*requiredContainerNamesFlag, ",")
	websitesToCheck := strings.Split(*websitesToCheckFlag, ",")

	display := display.New()
	cmd := command.NewCommander()
	docker := docker.New(cmd)
	hardware := hardware.New(cmd, "/var/lib/docker") // TODO docker root path auto-detection
	network := network.New()

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

	partitions(ctx, hardware, display)

	// OS information
	processesCount := hardware.ProcessesCount()
	fmt.Printf("%d processes running\n", processesCount)

	doDocker(ctx, docker, display, *composeFlag, requiredContainerNames)

	if *networkFlag {
		doNetwork(display, network, websitesToCheck)
	}

	if ctx.Err() != nil {
		return 1
	}
	return 0
}

func hostname(display display.Display) {
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

func zfs(ctx context.Context, hardware hardware.Hardware, display display.Display) {
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

func partitions(ctx context.Context, hardware hardware.Hardware, display display.Display) {
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

func doDocker(ctx context.Context, docker docker.Docker, display display.Display, composeCheck bool, requiredContainerNames []string) {
	if !docker.IsInstalled(ctx) {
		display.Error("Docker is not installed")
		return
	} else if !docker.IsRunning(ctx) {
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
		if !docker.IsComposeInstalled(ctx) {
			display.Warning("Docker-Compose is not installed")
		} else {
			dockerComposeVersion := docker.ComposeVersion(ctx)
			if err != nil {
				display.Error("Cannot get docker-compose version: %s", err)
			} else {
				dockerData = append(dockerData, "Compose "+dockerComposeVersion)
			}
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
	badStatus, err := docker.BadContainers(ctx)
	if err != nil {
		display.Error(err)
	}
	for _, status := range badStatus {
		display.Warning(status)
	}
}

func doNetwork(display display.Display, network network.Network, websitesToCheck []string) {
	var netData []string
	privateIP, err := network.GetOutboundIP()
	if err != nil {
		display.Error(err)
	} else {
		netData = append(netData, privateIP)
	}
	publicIP, err := network.GetPublicIP()
	if err != nil {
		display.Error(err)
	} else {
		netData = append(netData, publicIP)
	}
	fmt.Println(strings.Join(netData, " | "))

	errors := network.CheckMultipleHTTPConnections(websitesToCheck)
	for _, err := range errors {
		display.Warning(err)
	}
}
