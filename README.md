# 🛳️ Container in Go

This project is a minimal container implemented in Go, inspired by Liz Rice's talk at GOTO Conference 2018. It demonstrates fundamental container concepts like namespaces, control groups (cgroups), and chroot.

## 📜 Table of Contents

- [Overview](#-overview)
- [Features](#-features)
- [File Structure](#-file-structure)
- [Setup](#%EF%B8%8F-setup)
- [Usage](#%EF%B8%8F-usage)
- [Acknowledgements](#-acknowledgements)

## 📖 Overview

This container program uses Linux namespaces and cgroups to isolate and limit processes, mimicking some of the core features of containers like Docker. It leverages the `CLONE_NEWUTS`, `CLONE_NEWPID`, and `CLONE_NEWNS` flags for process isolation and mounts a basic Ubuntu filesystem.

## ✨ Features

- Process isolation using namespaces
- Resource control using cgroups (supports both v1 and v2)
- Custom chroot environment
- Minimal footprint and lightweight implementation

## 🗂 File Structure

```plaintext
.
├── go.mod          # Go module file
├── main.go         # Main program
├── README.md       # Project documentation
└── ubuntufs.tar    # Ubuntu filesystem archive
```

## ⚙️ Setup

### Prerequisites

- Go 1.18 or later
- Linux environment with support for namespaces and cgroups
- Root privileges for chroot and cgroup setup

### Installation

1. Clone this repository:

   ```bash
   git clone https://github.com/iamfitsum/containers-go
   cd containers-go
   ```

2. Extract the Ubuntu filesystem archive:

   ```bash
   mkdir -p /home/$(whoami)/ubuntu-fs
   tar -xvf ubuntufs.tar -C /home/$(whoami)/ubuntu-fs
   ```

3. Update the `main.go` file to reflect the correct path:

   Open `main.go` and replace `/home/vibe/ubuntu-fs` with `/home/[your-username]/ubuntu-fs`.

4. Ensure the necessary cgroup directories are writable:

   ```bash
   sudo chmod -R 755 /sys/fs/cgroup
   ```

## ▶️ Usage

Run the program using the `go run` command:

### Example Commands

1. To run a command inside the container:

   ```bash
   sudo go run main.go run <command> <args>
   ```

   Example:

   ```bash
   sudo go run main.go run /bin/bash
   ```

2. The program supports both cgroup v1 and v2 automatically. Ensure your system has the correct cgroup configuration.

### Expected Behavior

- Process isolation: The process runs in a separate PID and UTS namespace.
- Resource limitation: The cgroup restricts the process to a maximum of 20 PIDs.
- Root filesystem isolation: The process runs in the chroot environment specified in `/home/[your-username]/ubuntu-fs`.

## 🙏 Acknowledgements

Special thanks to [Liz Rice](https://github.com/lizrice) for the inspiration from her GOTO Conference talk on creating containers from scratch. Watch her full presentation [here](https://www.youtube.com/watch?v=8fi7uSYlOdc).

## 📜 License

This project is licensed under the [MIT License](LICENSE).
