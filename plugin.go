package main

import (
	"bufio"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
)

// Plugin defines the plugin parameters.
type Plugin struct {
	Flags  []string
	Action string
	Source string
	Target string
}

// Exec runs the plugin
func (p *Plugin) Exec() error {
	var cmd *exec.Cmd

	if p.Action == "" {
		logrus.Fatal("you must provide command to run")
	}

	cmd = exec.Command("rclone", p.Action)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "RCLONE_STATS=5s")
	//cmd.Env = append(cmd.Env, "RCLONE_LOG_LEVEL=INFO")
	//cmd.Env = append(cmd.Env, "VERBOSE=2")
	//	p.Flags = append(p.Flags, []string{"--stats=5s", "-v=1"}...)

	switch p.Action {
	case "sync", "move":
		if p.Source == "" || p.Target == "" {
			logrus.Fatal("you must provide source and target")
		}
		cmd.Args = append(cmd.Args, p.Source)
		cmd.Args = append(cmd.Args, p.Target)
		if len(p.Flags) > 0 {
			cmd.Args = append(cmd.Args, p.Flags...)
		}
	default:
		logrus.Fatalf("unsupported command %s", p.Action)
	}

	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		logrus.Fatalf("Error creating StdoutPipe: %s", err)
	}
	defer stdOut.Close()
	stdOutScan := bufio.NewScanner(stdOut)
	go func() {
		for stdOutScan.Scan() {
			logrus.Printf("%s\n", stdOutScan.Text())
		}
	}()

	stdErr, err := cmd.StderrPipe()
	if err != nil {
		logrus.Fatalf("Error creating StderrPipe: %s", err)
	}
	defer stdErr.Close()
	stdErrScan := bufio.NewScanner(stdErr)
	go func() {
		for stdErrScan.Scan() {
			logrus.Printf("%s\n", stdErrScan.Text())
		}
	}()

	logrus.Infof("%s", strings.Join(cmd.Args, " "))
	err = cmd.Start()
	if err != nil {
		logrus.Fatalf("Error: %s", err)
	}

	logrus.Info("waiting for finish")
	err = cmd.Wait()
	if err != nil {
		logrus.Fatalf("Error: %s", err)
	}

	// for broken docker
	time.Sleep(5 * time.Second)
	return nil
}
