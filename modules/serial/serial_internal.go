
package serial

import (
	"io"
	"syscall"
	"log"
	"errors"
	"golang.org/x/sys/unix"
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


	status.Cflag = unix.CSIZE | unix.PARENB | unix.CLOCAL | unix.CREAD | unix.CS8


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


	// clear the DTR bit
	flag := syscall.TIOCM_DTR

	Tiocmbic(fd, &flag)


	Tcflush(fd, syscall.TCIFLUSH)

	status, _ = unix.IoctlGetTermios(port, ioctlReadTermios)

	log.Printf("[DEBUG] Serial: Read back Status %+v", status)

	p := &Port{
		fd: port,
	}
	return p,nil
}