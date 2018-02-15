// +build darwin

package serial

import (
	"io"
	"syscall"
	"log"
	"errors"
	"fmt"
	"unsafe"
)
// #include <termios.h>
// #include <unistd.h>
import "C"

type Port struct{
	fd int
}

func (p Port) Close() error {
	return syscall.Close(p.fd)
}

func (p Port) Read(b []byte) (n int, err error) {
	return syscall.Read(p.fd, b)
}

func (p Port) Write(b []byte) (n int, err error) {
	return syscall.Write(p.fd,b)
}

func (p Port) Flush() {
	fd := C.int(p.fd)
	C.tcflush(fd, C.TCIFLUSH)
}

func openInternal(options SerialConfig) (io.ReadWriteCloser, error) {
	log.Printf("[DEBUG] Serial: Openning %s", options.Port)

	// Open will block without the O_NONBLOCK
	port, err := syscall.Open(options.Port, syscall.O_RDWR|syscall.O_NOCTTY|syscall.O_CLOEXEC|syscall.O_NONBLOCK, 0620)
	if err != nil {
		return nil, err
	}

	fd := C.int(port)

	if C.isatty(fd) != 1 {
		syscall.Close(port)
		return nil, errors.New("File is not a tty")
	}

	C.tcflush(fd, C.TCIFLUSH)

	var st C.struct_termios

	_, err = C.tcgetattr(fd, &st)
	if err != nil {
		syscall.Close(port)
		return nil, err
	}

	st.c_cflag &= ^C.tcflag_t(C.CSIZE | C.PARENB)

	st.c_cflag |= (C.CLOCAL | C.CREAD | C.CS8)


	switch options.StopBits {
	case 1:
		st.c_cflag &= ^C.tcflag_t(C.CSTOPB)
	case 2:
		st.c_cflag |= C.CSTOPB
	default:
		return nil, errors.New("Unknown StopBits value")
	}

	// raw input
	st.c_cflag &= ^C.tcflag_t(C.ICANON | C.ECHO | C.ECHOE | C.ISIG)

	// raw output
	st.c_oflag &= ^C.tcflag_t(C.OPOST)
	// software flow control disabled
	st.c_iflag &= ^C.tcflag_t(C.IXON)
	// do not translate CR to NL
	st.c_iflag &= ^C.tcflag_t(C.ICRNL)

	if options.RtsCts {
		st.c_cflag |= C.CRTSCTS
	} else {
		st.c_cflag &= ^C.tcflag_t(C.CRTSCTS)
	}


	var speed C.speed_t

	switch options.Baud {
	case 115200:
		speed = C.B115200
	case 57600:
		speed = C.B57600
	case 38400:
		speed = C.B38400
	case 19200:
		speed = C.B19200
	case 9600:
		speed = C.B9600
	case 4800:
		speed = C.B4800
	case 2400:
		speed = C.B2400
	case 1200:
		speed = C.B1200
	case 600:
		speed = C.B600
	case 300:
		speed = C.B300
	case 200:
		speed = C.B200
	case 150:
		speed = C.B150
	case 134:
		speed = C.B134
	case 110:
		speed = C.B110
	case 75:
		speed = C.B75
	case 50:
		speed = C.B50
	default:
		syscall.Close(port)
		return nil, fmt.Errorf("Unknown Baud value %v", options.Baud)
	}

	_, err = C.cfsetispeed(&st, speed)
	if err != nil {
		syscall.Close(port)
		return nil, err
	}

	_, err = C.cfsetospeed(&st, speed)
	if err != nil {
		syscall.Close(port)
		return nil, err
	}

	_, err = C.tcsetattr(fd, C.TCSANOW, &st)
	if err != nil {
		syscall.Close(port)
		return nil, err
	}

	// Need to change back the port to blocking.

	nonblockErr := syscall.SetNonblock(port, false)
	if nonblockErr != nil {
		syscall.Close(port)
		return nil, nonblockErr
	}

	var oldflags C.int
	var e error

	e = ioctl(uintptr(port), syscall.TIOCMGET, uintptr(unsafe.Pointer(&oldflags)))
	if errno, ok := e.(syscall.Errno); ok {
		if errno != syscall.ENOTTY {
			syscall.Close(port)
			return nil, e
		}
	}

	// clear the DTR bit
	oldflags &= ^C.int(C.TIOCM_DTR);

	e = ioctl(uintptr(port), syscall.TIOCMSET, uintptr(unsafe.Pointer(&oldflags)))

	if errno, ok := e.(syscall.Errno); ok {
		if errno != syscall.ENOTTY {
			syscall.Close(port)
			return nil, e
		}
	}

	C.tcflush(fd, C.TCIFLUSH)

	p := &Port{
		fd: port,
	}
	return p,nil
}