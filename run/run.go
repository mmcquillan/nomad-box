package run

import (
	"bufio"
	"bytes"
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
		Error("Running: " + command)
		Error(err.Error())
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		Error("Running: " + command)
		Error(err.Error())
	}
	if err := cmd.Start(); err != nil {
		Error("Running: " + command)
		Error(err.Error())
	}
	o := ReaderToString(stdout)
	if o != "" {
		Warn(o)
	}
	e := ReaderToString(stderr)
	if e != "" {
		Warn(e)
	}
	err = cmd.Wait()
	if e != "" {
		Warn(e)
	}
}

func CommandContains(command string, match string) bool {
	cmd := exec.Command("bash", "-c", command)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		Error("Running: " + command)
		Error(err.Error())
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		Error("Running: " + command)
		Error(err.Error())
	}
	if err := cmd.Start(); err != nil {
		Error("Running: " + command)
		Error(err.Error())
	}
	o := ReaderToString(stdout)
	e := ReaderToString(stderr)
	if e != "" {
		Warn(e)
	}
	err = cmd.Wait()
	if e != "" {
		Warn(e)
	}
	return strings.Contains(o, match)
}

func Process(command string, prefix string, log bool) (pid int) {
	cmd := exec.Command("bash", "-c", command)
	out, err := cmd.StdoutPipe()
	if err != nil {
		Error("Running: " + command)
		Error(err.Error())
	}
	// combine stderr + stdout (guess this wokrs)
	cmd.Stderr = cmd.Stdout
	done := make(chan struct{})
	scanner := bufio.NewScanner(out)
	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			if log {
				Out("[" + prefix + "] " + line)
			}
		}
		done <- struct{}{}
	}()
	if err := cmd.Start(); err != nil {
		Error("Running " + command)
		Error(err.Error())
	}
	go func() {
		<-done
		err = cmd.Wait()
		if err != nil {
			Warn(command)
			Warn(err.Error())
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
