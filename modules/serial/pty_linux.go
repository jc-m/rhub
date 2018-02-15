package serial

import (
	"strconv"
	"syscall"
	"unsafe"
)

func grantpt(fd uintptr) error {
	return nil
}

func unlockpt(fd uintptr) error {
	var u _C_int
	// use TIOCSPTLCK with a zero valued arg to clear the slave pty lock
	return ioctl(fd, syscall.TIOCSPTLCK, uintptr(unsafe.Pointer(&u)))
}

func ptsname(fd int) (string, error) {
	var n _C_uint
	err := ioctl(uintptr(fd), syscall.TIOCGPTN, uintptr(unsafe.Pointer(&n)))
	if err != nil {
		return "", err
	}
	return "/dev/pts/" + strconv.Itoa(int(n)), nil
}
