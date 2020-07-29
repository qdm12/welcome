# Welcome

```
| |  | |    | |
| |  | | ___| | ___ ___  _ __ ___   ___
| |/\| |/ _ \ |/ __/ _ \| '_ ` _ \ / _ \
\  /\  /  __/ | (_| (_) | | | | | |  __/
 \/  \/ \___|_|\___\___/|_| |_| |_|\___|
```

This is a Golang static binary I use on my servers when I login.

## Features

It checks and displays several things

1. Your hostname as *ASCII art* (random font by default)
1. The date and time
1. The server uptime
1. The total RAM and CPU % usage
1. For all ZFS volumes (if any)
    - Capacity left in %
    - Health status
    - Data data errors
1. Usage % of other partitions
1. Docker (if installed)
    - Docker version
    - Docker compose version (only with `--compose` as it takes one second)
    - Number of containers running
    - Unhealthy and restarting containers as warnings
1. Network information
    - hostname
    - Main LAN IP address
    - Public IP address (only with `--network` flag as it takes 100-300 milliseconds)
1. Checks multiple websites are up (only with `--network` flag as it takes 100-300 milliseconds)

## Setup

```sh
wget -qO ~/welcome https://github.com/qdm12/welcome/releases/download/v0.1.0/welcome_0.1.0_linux_amd64
chmod +x ~/welcome
# And then, depending on your shell
echo "~/welcome" >> ~/.zshrc
```

## Usage

```sh
Usage of welcome:
  -compose
        show docker-compose version (slow)
  -network
        verify network connectivity
  -requiredContainers string
        comma separated list of required running container names to check for (default "dns,ddns")
  -websitesToCheck string
        comma separated list of websites to check, only enabled if --network is specified (default "https://qqq.ninja,https://1.1.1.1")
```

## TODOs

- Check run as root
- Fix RAM usage
- Read partitions using advanced `df --help`
- Check for non-imported encrypted zpool and prompt to import them `zpool import poolname -l`
- If Docker is not running, launch it `systemctl start docker`
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