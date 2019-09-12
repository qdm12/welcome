# Welcome

This is a Golang static binary I use on my servers when I login.

## Purpose

It checks and displays several things
1. Your hostname as *ASCII art* (random font by default)
1. The date and time
1. The server uptime
1. The total RAM and CPU usage
1. For all ZFS volumes
    - Capacity left in %
    - Health status
    - Data data errors
1. Docker
    - Docker and Docker-Compose versions
    - Number of containers running
    - Unhealthy and restarting containers as warnings
1. The hostname, LAN ip address, and public IP address (using duckduckgo.com)
1. Checks multiple websites are up

## Building

1. Open VS code, install the recommended extensions
1. Open the folder in the dev container
1. In the dev container, open a bash terminal (you can use **Ctrl**+**Shift**+**`**)
1. Build for Linux on an amd64 CPU with

    ```sh
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w"
    ```

1. You can then copy `welcome` to your server at `~` and add `./welcome` to your *.zshrc* or *.bashrc* in example.

Note that it cannot really run in a Docker container as it needs info from the host machine.

## TODOs

- Change colors depending on % used
- Add emojis
- Go routines
    - Bandwidth download, then upload
    - Ping
    - Docker Stats
- Use Docker API, see [this](https://docs.docker.com/develop/sdk/examples/)
- ZFS
    - Use `zpool events`
    - Check scrub time
    - Check error numbers with `zpool status POOLNAME`