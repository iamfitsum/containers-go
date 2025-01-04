package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

// docker 		  run image <cmd> <params>
// go run main.go run       <cmd> <params>

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go run <cmd> <params>")
		os.Exit(1)
	}
	
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("bad command")
	}
}

func run() {
	fmt.Printf("Running %v as %d\n", os.Args[2:], os.Getpid())
	
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		Unshareflags: syscall.CLONE_NEWNS,
	}
	must(cmd.Run())
}

func child() {
	fmt.Printf("Running %v as %d\n", os.Args[2:], os.Getpid())
	
	if isCgroupV2() {
		cgV2()
	} else {
		cgV1()
	}
	
	must(syscall.Sethostname([]byte("container")))
	must(syscall.Chroot("/home/vibe/ubuntu-fs"))
	must(syscall.Chdir("/"))
	must(syscall.Mount("proc", "proc", "proc", 0, ""))

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	must(cmd.Run())

	must(syscall.Unmount("/proc", 0))
}

// Detect if the system is using cgroup v2
func isCgroupV2() bool {
	fi, err := os.Stat("/sys/fs/cgroup/cgroup.controllers")
	return err == nil && !fi.IsDir()
}

// cgroup v2
func cgV2() {
	cgroups := "/sys/fs/cgroup/"
	fit := filepath.Join(cgroups, "fit")

	// Create the cgroup directory
	if err := os.MkdirAll(fit, 0755); err != nil && !os.IsExist(err) {
		must(err)
	}

	// Write the pids.max file
	if err := os.WriteFile(filepath.Join(fit, "pids.max"), []byte("20"), 0700); err != nil {
		must(err)
	}

	// Write the cgroup.procs file
	if err := os.WriteFile(filepath.Join(fit, "cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700); err != nil {
		must(err)
	}
}

// cgroup v1
func cgV1() {
	cgroups := "/sys/fs/cgroup/"
	pids := filepath.Join(cgroups, "pids")
	fit := filepath.Join(pids, "fit")

	// Create the cgroup directory
	if err := os.MkdirAll(fit, 0755); err != nil && !os.IsExist(err) {
		must(err)
	}

	// Write the pids.max file
	if err := os.WriteFile(filepath.Join(fit, "pids.max"), []byte("20"), 0700); err != nil {
		must(err)
	}

	// Write the notify_on_release file
	if err := os.WriteFile(filepath.Join(fit, "notify_on_release"), []byte("1"), 0700); err != nil {
		must(err)
	}

	// Write the cgroup.procs file
	if err := os.WriteFile(filepath.Join(fit, "cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700); err != nil {
		must(err)
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
