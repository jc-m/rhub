
package serial

import (
	"io"
	"syscall"
	"log"
	"errors"
	"golang.org/x/sys/unix"
	"unsafe"
)


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

func setRts(fd uintptr,b  bool) error {
	state := syscall.TIOCM_RTS
	if err := ioctl(fd, syscall.TIOCMGET, uintptr(unsafe.Pointer(&state))); err != nil {
		return err
	}
	if b {
		state |= syscall.TIOCM_RTS
	} else {
		state &^= syscall.TIOCM_RTS
	}
	return ioctl(fd, syscall.TIOCMSET, uintptr(unsafe.Pointer(&state)))
}

func setDTR(fd uintptr,b  bool) error {
	state := syscall.TIOCM_DTR
	if err := ioctl(fd, syscall.TIOCMGET, uintptr(unsafe.Pointer(&state))); err != nil {
		return err
	}
	if b {
		state |= syscall.TIOCM_DTR
	} else {
		state &^= syscall.TIOCM_DTR
	}
	return ioctl(fd, syscall.TIOCMSET, uintptr(unsafe.Pointer(&state)))
}

func openInternal(options SerialConfig) (io.ReadWriteCloser, error) {

	log.Printf("[DEBUG] Serial: Openning %s", options.Port)

	// Open will block without the O_NONBLOCK
	port, err := syscall.Open(options.Port, syscall.O_RDWR|syscall.O_NOCTTY|syscall.O_NONBLOCK, 0600)
	if err != nil {
		return nil, err
	}


	fd := uintptr(port)

	Tcflush(fd, syscall.TCIFLUSH)


	status, geterr := unix.IoctlGetTermios(port, ioctlReadTermios)
	if geterr != nil {
		syscall.Close(port)
		return nil, err
	}

	log.Printf("[DEBUG] Serial: Status %+v", status)


	status.Cflag = unix.CSIZE  | unix.CLOCAL | unix.CREAD | unix.CS8
	status.Cflag &^= unix.PARENB

	switch options.StopBits {
	case 1:
		status.Cflag &^= unix.CSTOPB
	case 2:
		status.Cflag |= unix.CSTOPB
	default:
		return nil, errors.New("Unknown StopBits value")
	}

	// raw input
	status.Lflag &^= unix.ICANON | unix.ECHO | unix.ECHOE | unix.ISIG | unix.IEXTEN

	// raw output
	status.Oflag &^= unix.OPOST
	// software flow control disabled
	status.Iflag &^= unix.IXON
	// do not translate CR to NL
	status.Iflag &^= unix.ICRNL


	if options.RtsCts {
		status.Cflag |= unix.CRTSCTS
	} else {
		status.Cflag &^= unix.CRTSCTS
	}
	log.Printf("[DEBUG] Serial: New Status %+v", status)

	setSpeed(status, options.Baud)

	status.Cc[unix.VMIN] = 1
	status.Cc[unix.VTIME] = 0

	log.Printf("[DEBUG] Serial: Applying new status %+v", status)

	unix.IoctlSetTermios(port, ioctlWriteTermios, status)

	// Need to change back the port to blocking.

	nonblockErr := syscall.SetNonblock(port, false)
	if nonblockErr != nil {
		syscall.Close(port)
		return nil, nonblockErr
	}

	setDTR(fd, false)

	Tcflush(fd, syscall.TCIFLUSH)

	status, _ = unix.IoctlGetTermios(port, ioctlReadTermios)

	log.Printf("[DEBUG] Serial: Read back Status %+v", status)

	p := &Port{
		fd: port,
	}
	return p,nil
}