package serial


import (
	"syscall"
	"unsafe"
	"errors"
)

func _IOC_PARM_LEN(ioctl uintptr) uintptr {
	const 	(
		_IOC_PARAM_SHIFT = 13

		_IOC_PARAM_MASK  = (1 << _IOC_PARAM_SHIFT) - 1
	)
	return (ioctl >> 16) & _IOC_PARAM_MASK
}

func grantpt(fd uintptr) error {
	return ioctl(fd, syscall.TIOCPTYGRANT, 0)
}

func unlockpt(fd uintptr) error {
	return ioctl(fd, syscall.TIOCPTYUNLK, 0)
}


func ptsname(f int) (string, error) {
	n := make([]byte, _IOC_PARM_LEN(syscall.TIOCPTYGNAME))

	err := ioctl(uintptr(f), syscall.TIOCPTYGNAME, uintptr(unsafe.Pointer(&n[0])))
	if err != nil {
		return "", err
	}

	for i, c := range n {
		if c == 0 {
			return string(n[:i]), nil
		}
	}
	return "", errors.New("TIOCPTYGNAME string not NUL-terminated")
}