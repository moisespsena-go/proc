package proch

import (
	"bytes"
	"os"
	"unicode"
)

func isDigit(str string) bool {
	for _, c := range str {
		if c == '-' || c == '+' || !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

func IsProcB(pid string, binPath string, buf []byte) (bool, error) {
	pth := "/proc/" + pid + "/cmdline"
	f, err := os.Open(pth)
	if err != nil {
		return false, &os.PathError{"open", pth, err}
	}
	defer f.Close()
	n, err := f.Read(buf)
	if n != len(binPath) {
		return false, nil
	}
	if bytes.Compare(buf, []byte(binPath)) == 0 {
		return true, nil
	}
	return false, nil
}

func IsProc(pid string, binPath string) (bool, error) {
	return IsProcB(pid, binPath, []byte(binPath))
}
