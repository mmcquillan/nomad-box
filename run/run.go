package run

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/process"
)

func Command(command string) {
	cmd := exec.Command("bash", "-c", command)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("   ERROR: Running " + command)
		fmt.Printf("   %+v\n", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println("   ERROR: Running " + command)
		fmt.Printf("   %+v\n", err)
	}
	if err := cmd.Start(); err != nil {
		fmt.Println("   ERROR: Running " + command)
		fmt.Printf("   %+v\n", err)
	}
	o := ReaderToString(stdout)
	if o != "" {
		fmt.Println("   WARN: " + o)
	}
	e := ReaderToString(stderr)
	if e != "" {
		fmt.Println("   WARN: " + e)
	}
	err = cmd.Wait()
	if e != "" {
		fmt.Println("   WARN: " + e)
	}
}

func CommandContains(command string, match string) bool {
	cmd := exec.Command("bash", "-c", command)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("   ERROR: Running " + command)
		fmt.Printf("   %+v\n", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println("   ERROR: Running " + command)
		fmt.Printf("   %+v\n", err)
	}
	if err := cmd.Start(); err != nil {
		fmt.Println("   ERROR: Running " + command)
		fmt.Printf("   %+v\n", err)
	}
	o := ReaderToString(stdout)
	e := ReaderToString(stderr)
	if e != "" {
		fmt.Println("   WARN: " + e)
	}
	err = cmd.Wait()
	if e != "" {
		fmt.Println("   WARN: " + e)
	}
	return strings.Contains(o, match)
}

func Process(command string, prefix string, log bool) (pid int) {
	cmd := exec.Command("bash", "-c", command)
	out, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("   ERROR: Running " + command)
		fmt.Printf("   %+v\n", err)
	}
	// combine stderr + stdout (guess this wokrs)
	cmd.Stderr = cmd.Stdout
	done := make(chan struct{})
	scanner := bufio.NewScanner(out)
	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			if log {
				fmt.Println("   [" + prefix + "]  " + line)
			}
		}
		done <- struct{}{}
	}()
	if err := cmd.Start(); err != nil {
		fmt.Println("   ERROR: Running " + command)
		fmt.Printf("   %+v\n", err)
	}
	go func() {
		<-done
		err = cmd.Wait()
		if err != nil {
			fmt.Println("   WARN: " + command)
			fmt.Printf("   %+v\n", err)
		}
	}()
	for cmd.Process.Pid == 0 {
		time.Sleep(2 * time.Second)
	}
	return cmd.Process.Pid
}

func CheckProcess(pid int) bool {
	exists, err := process.PidExists(int32(pid))
	if err != nil {
		return true
	}
	return exists
}

func ReaderToString(reader io.ReadCloser) (out string) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	out = buf.String()
	return out
}
