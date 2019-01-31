package proch

import (
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

type Pid int64

func (p Pid) Kill(signal ...os.Signal) error {
	var args []string
	for _, sig := range signal {
		args = append(args, "-"+strconv.Itoa(int(sig.(syscall.Signal))))
	}
	args = append(args, strconv.Itoa(int(p)))
	cmd := exec.Command("kill", args...)
	return cmd.Run()
}

type Pids []Pid

func (p Pids) Kill(signal ...os.Signal) error {
	if len(p) == 0 {
		return nil
	}

	var args []string
	for _, sig := range signal {
		args = append(args, "-"+strconv.Itoa(int(sig.(syscall.Signal))))
	}
	for _, pid := range p {
		args = append(args, strconv.Itoa(int(pid)))
	}
	cmd := exec.Command("kill", args...)
	err := cmd.Run()
	return err
}
