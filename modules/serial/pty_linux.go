package serial

import (
	"syscall"
	"unsafe"
	"fmt"
)

func grantpt(fd uintptr) error {
	return nil
}

func unlockpt(fd uintptr) error {
	var u uintptr
	// use TIOCSPTLCK with a zero valued arg to clear the slave pty lock
	return ioctl(fd, syscall.TIOCSPTLCK, uintptr(unsafe.Pointer(&u)))
}

func ptsname(fd int) (string, error) {
	var n uintptr
	err := ioctl(uintptr(fd), syscall.TIOCGPTN, uintptr(unsafe.Pointer(&n)))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("/dev/pts/%d", n), nil
}
