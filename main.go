package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"welcome/pkg/display"
	"welcome/pkg/docker"
	"welcome/pkg/hardware"
	"welcome/pkg/network"
	"welcome/pkg/terminal"
)

var websitesToCheck = []string{
	"https://1.1.1.1",
	"https://qqq.ninja",
	"https://lux.qqq.ninja",
}

func main() {
	CHECKNETWORK := false
	CHECKCOMPOSE := false
	for _, arg := range os.Args[1:len(os.Args)] {
		if arg == "--network" || arg == "-n" {
			CHECKNETWORK = true
		} else if arg == "--compose" {
			CHECKCOMPOSE = true
		}
	}
	hostname, err := terminal.RunCommand("hostname")
	if err != nil {
		display.Error("Cannot get hostname: %s", err)
	}
	fmt.Println(display.GetRandomAsciiArt(hostname))

	// Time and uptime
	t := time.Now().Format("2006-01-02 15:04:05")
	upTime := hardware.GetUptime()
	fmt.Printf("%s | Up %s\n", t, upTime)

	// RAM and CPU usage
	UsageRam := hardware.GetRAMUsage()
	UsageCpu, err := hardware.GetCPUUsage()
	if err != nil {
		display.Error("%s", err)
		UsageCpu = -1
	}
	fmt.Printf("RAM %d%% | CPU %d%%\n", UsageRam, UsageCpu)

	// ZFS
	ZpoolIsInstalled := hardware.IsZpoolInstalled()
	if ZpoolIsInstalled {
		capacities := []string{}
		pools, err := hardware.GetPools()
		if err != nil {
			display.Error("%s", err)
		}
		for _, pool := range pools {
			health, err := hardware.GetPoolHealth(pool)
			if err != nil {
				display.Error("%s", err)
			} else if health != "" {
				display.Warning("ZFS pool %s is %s", pool, health)
			}
			errors, err := hardware.GetPoolErrors(pool)
			if err != nil {
				display.Error("%s", err)
			} else if errors != "" {
				display.Warning("ZFS pool %s: %s", pool, errors)
			}
			capacity, err := hardware.GetPoolCapacity(pool)
			if err != nil {
				display.Error("%s", err)
			}
			capacities = append(capacities, fmt.Sprintf("%s %d%%", pool, capacity))
		}
		fmt.Println(strings.Join(capacities, " | "))
	}

	// TODO disk usage without ZFS

	// OS information
	processes := hardware.GetProcesses()
	fmt.Printf("%d processes running\n", processes)

	// Docker
	dockerInstalled := docker.IsDockerInstalled()
	if !dockerInstalled {
		display.Error("Docker is not installed")
	}
	dockerRunning := false
	if dockerInstalled {
		dockerRunning = docker.IsDockerRunning()
		if !dockerRunning {
			display.Error("Docker is not running")
		}
	}

	// TODO Get RAM and CPU usage of Docker containers (SLOW, do in goroutine)
	// dockerStats, _ := terminal.RunCommand("docker", "stats", "--no-stream", "--format", "'{{.MemUsage}}'")

	dockerData := []string{}
	dockerVersion := docker.GetDockerVersion()
	dockerData = append(dockerData, "Docker "+dockerVersion)
	if dockerRunning {
		containersCount, err := docker.CountContainers()
		if err != nil {
			display.Error("Cannot count containers: %s", err)
		} else {
			dockerData = append(dockerData, fmt.Sprintf("%d containers", containersCount))
		}
	}
	if CHECKCOMPOSE {
		if !docker.IsDockerComposeInstalled() {
			display.Warning("Docker-Compose is not installed")
		} else {
			dockerComposeVersion := docker.GetDockerComposeVersion()
			if err != nil {
				display.Error("%s", err)
			} else {
				dockerData = append(dockerData, "Compose "+dockerComposeVersion)
			}
		}
	}
	fmt.Println(strings.Join(dockerData, " | "))
	if dockerRunning {
		containersNotRunning, err := docker.IsContainerRunning("dns", "ddns", "sftp", "samba")
		if err != nil {
			display.Error("Cannot check for running containers: %s", err)
		}
		for _, container := range containersNotRunning {
			display.Warning("Container %s is not running", container)
		}
		badStatus, err := docker.GetBadContainers()
		if err != nil {
			display.Error("Cannot get bad containers: %s", err)
		}
		for _, status := range badStatus {
			display.Warning(status)
		}
	}

	// Networking
	netData := []string{hostname}
	privateIP, err := network.GetOutboundIP()
	if err != nil {
		display.Error("Cannot get private IP address: %s", err)
	}
	netData = append(netData, privateIP)
	if CHECKNETWORK {
		publicIP, err := network.GetPublicIP()
		if err != nil {
			display.Error("Cannot get public IP address: %s", err)
		}
		netData = append(netData, publicIP)
	}
	fmt.Println(strings.Join(netData, " | "))

	if CHECKNETWORK {
		errors := network.CheckMultipleHTTPConnections(websitesToCheck)
		for _, err := range errors {
			display.Warning("%s", err)
		}
	}
}
