package proch

import (
	"io"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/go-errors/errors"
)

func IsBinary(pth string) bool {
	cmd := exec.Command("perl", "-E", "exit((-B $ARGV[0])?0:1);", pth)
	if err := cmd.Run(); err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			panic(err)
		}
	}
	return cmd.ProcessState.Success()
}

type Binary struct {
	Pth string
}

func NewBinary(pth string) (*Binary, error) {
	pth, err := filepath.EvalSymlinks(pth)
	if err != nil {
		return nil, err
	}
	if !IsBinary(pth) {
		return nil, &os.PathError{"binary_check", pth, errors.New("not is binary")}
	}
	return &Binary{Pth: pth}, nil
}

func (b *Binary) IsRunning(uid ...uint32) (ok bool, err error) {
	err = b.PidsS(func(pid string) error {
		return io.EOF
	}, uid...)
	if err == io.EOF {
		ok = true
		err = nil
	}
	return
}

func (b *Binary) PidsS(cb func(pid string) error, uid ...uint32) (err error) {
	var (
		buf = []byte(b.Pth)
	)
	f, err := os.Open("/proc")
	if err != nil {
		return err
	}

	if names, err := f.Readdirnames(-1); err != nil {
		return err
	} else {
		for _, name := range names {
			if isDigit(name) {
				if uid != nil {
					var ok bool
					for _, uid := range uid {
						procPth := filepath.Join("/proc", name)
						var s syscall.Stat_t
						if err := syscall.Stat(procPth, &s); err != nil {
							return err
						} else if uint32(s.Uid) == uid {
							ok = true
							break
						}
					}
					if !ok {
						continue
					}
				}

				if ok, err := IsProcB(name, b.Pth, buf); err != nil {
					return err
				} else if ok {
					if err = cb(name); err != nil {
						return err
					}
				}
			}
		}
	}
	return
}

func (b *Binary) Pids(uid ...uint32) (pids Pids, err error) {
	err = b.PidsS(func(pid string) error {
		i, _ := strconv.Atoi(pid)
		pids = append(pids, Pid(i))
		return nil
	}, uid...)
	return
}

func (b *Binary) Kill(signal ...os.Signal) error {
	u, err := user.Current()
	if err != nil {
		return err
	}

	uid64, _ := strconv.ParseInt(u.Uid, 10, 32)
	uid := uint32(uid64)

	var pids Pids
	for pids, err = b.Pids(); err == nil && len(pids) > 0; pids, err = b.Pids(uid) {
		err = pids.Kill(signal...)
	}
	return err
}
