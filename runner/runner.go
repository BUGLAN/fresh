package runner

import (
	"flag"
	"io"
	"os/exec"
)

var cmdArgs []string

func init() {
	flag.Parse()
	cmdArgs = flag.Args()
}

func run() bool {
	runnerLog("Running...")


	var cmd *exec.Cmd
	if len(cmdArgs) == 0  {
		cmd = exec.Command(buildPath())
	} else {
		cmd = exec.Command(buildPath(), cmdArgs...)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fatal(err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fatal(err)
	}

	err = cmd.Start()
	if err != nil {
		fatal(err)
	}

	go io.Copy(appLogWriter{}, stderr)
	go io.Copy(appLogWriter{}, stdout)

	go func() {
		<-stopChannel
		pid := cmd.Process.Pid
		runnerLog("Killing PID %d", pid)
		cmd.Process.Kill()
	}()

	return true
}
