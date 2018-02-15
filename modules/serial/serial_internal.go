
package serial

import (
	"io"
	"syscall"
	"log"
	"errors"
)

const (
	TCSANOW   = 0
	TCSADRAIN = 1
	TCSAFLUSH = 2
)

type termStatus syscall.Termios

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
	fd := uintptr(p.fd)
	Tcflush(fd, syscall.TCIFLUSH)
}

func openInternal(options SerialConfig) (io.ReadWriteCloser, error) {
	var status termStatus

	log.Printf("[DEBUG] Serial: Openning %s", options.Port)

	// Open will block without the O_NONBLOCK
	port, err := syscall.Open(options.Port, syscall.O_RDWR|syscall.O_NOCTTY|syscall.O_CLOEXEC|syscall.O_NONBLOCK, 0620)
	if err != nil {
		return nil, err
	}


	fd := uintptr(port)

	Tcflush(fd, syscall.TCIFLUSH)


	err = Tcgetattr(fd, &status)
	if err != nil {
		syscall.Close(port)
		return nil, err
	}

	status.Cflag = syscall.CS8 | syscall.CREAD | syscall.CLOCAL
	status.Cflag &^= syscall.CSIZE | syscall.PARENB
	status.Cflag |= syscall.CLOCAL | syscall.CREAD | syscall.CS8



	switch options.StopBits {
	case 1:
		status.Cflag &^= syscall.CSTOPB
	case 2:
		status.Cflag |= syscall.CSTOPB
	default:
		return nil, errors.New("Unknown StopBits value")
	}

	// raw input
	status.Cflag &^= syscall.ICANON | syscall.ECHO | syscall.ECHOE | syscall.ISIG

	// raw output
	status.Oflag &^= syscall.OPOST
	// software flow control disabled
	status.Iflag &^= syscall.IXON
	// do not translate CR to NL
	status.Iflag &^= syscall.ICRNL

	if options.RtsCts {
		status.Cflag &= CRTSCTS
	} else {
		status.Cflag &^= CRTSCTS
	}

	status.setSpeed(options.Baud)

	err = Tcsetattr(fd, TCSANOW, &status)
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

	// clear the DTR bit
	flag := syscall.TIOCM_DTR

	Tiocmbic(fd, &flag)


	Tcflush(fd, syscall.TCIFLUSH)

	p := &Port{
		fd: port,
	}
	return p,nil
}