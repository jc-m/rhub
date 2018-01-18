package serial

import (
	"syscall"
	"log"
	"unsafe"
	"errors"
	"io"
	"github.com/jc-m/rhub/modules"
	"fmt"
)

type pty struct {
	channels []modules.Channels
	ptmx     int
	slave int
	portName string
	state byte
}

func ioctl(fd uintptr, flag, data uintptr) error {
	if _, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, flag, data); err != 0 {
		return err
	}

	return nil
}
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


func (p *pty) receiveloop() {
	buffer := make([]byte, 1024)
	for {
		n, err := syscall.Read(p.ptmx, buffer)

		if n > 0 {
			log.Printf("[DEBUG] Pty: Received %d bytes", n)

			b := make([]byte, n)
			copy(b, buffer[:n])
			log.Printf("[DEBUG] Pty: Sending %+v", b)

			p.channels[0].Out <- modules.Message{Id:p.portName, Body:b}
		}
		if n <= 0 {
			if err != nil {
				if err == io.EOF || err == io.ErrUnexpectedEOF {
					log.Print("[DEBUG] Pty: EOF")
				}
				if err != nil {
					log.Printf("[ERROR] Pty: Read error : %s", err)
					break
				}
			}

			if err == nil {
				log.Print("[DEBUG] Null Read")
				p.state = STATE_CLOSED
			}
		}
		if p.state == STATE_CLOSED {
			log.Printf("[DEBUG] Pty: Port is closed : %s", p.portName)
			break
		}
	}
	log.Print("[DEBUG] Pty: Terminating Receiving loop")
	p.channels[0].Ctl <- true
}


func (p *pty) sendloop() {
	for {
		select {
		case r := <-p.channels[0].In:
			n, err := syscall.Write(p.ptmx, r.Body)
			if err != nil {
				panic(err)
			}
			log.Printf("[DEBUG] Pty: Sent %d bytes", n)

		case <-p.channels[0].Ctl:
			log.Print("[DEBUG] Pty: Terminating Sending loop")
			return
		}
	}
}

func (p *pty) ptyOpen() error {
	master, err := syscall.Open("/dev/ptmx", syscall.O_RDWR|syscall.O_NOCTTY|syscall.O_CLOEXEC, 0620)
	if err != nil {
		log.Printf("[ERROR] Pty: Cannot open master %s",err)
		return err
	}
	fd := uintptr(master)


	defer func() {
		if err != nil {
			_ = syscall.Close(master) // Best effort.
		}
	}()


	// Grant/unlock slave.
	if err := grantpt(fd); err != nil {
		log.Printf("[ERROR] Pty: Cannot grant slave %s",err)
		panic(err)
	}
	if err := unlockpt(fd); err != nil {
		log.Printf("[ERROR] Pty: Cannot unlock slave %s",err)

		panic(err)
	}

	sname, err := ptsname(master)
	if err != nil {
		log.Printf("[ERROR] Pty: Cannot get slave %s",err)
		return err
	}

	log.Printf("[DEBUG] Pty: slave name: %s",sname)

	// Keep the pty open so that the other end can close/open at will without causing an EOF error
	x, err := syscall.Open(sname, syscall.O_RDWR|syscall.O_NOCTTY, 0620 )
	if err != nil {
		log.Printf("[ERROR] Pty: Cannot open slave %s",err)
		return err
	}

	p.portName = sname
	p.ptmx = master
	p.slave = x
	p.state = STATE_OPEN
	return nil
}

func (p *pty) Close() {
	log.Print("[DEBUG] Pty: Closing")
	close(p.channels[0].Ctl)
	syscall.Close(p.slave)
	syscall.Close(p.ptmx)

}

func (p *pty) GetType() int {
	return modules.DRIVER
}

func (p *pty) GetName() string {
	return p.portName
}

func (m *pty) AddChannels(channels modules.Channels) ([]modules.Channels, error)  {
	if len(m.channels) > 0 {
		return nil, fmt.Errorf("Module supports only on pair of channels")
	}
	m.channels = append(m.channels, channels)
	return m.channels, nil
}


// Reads from a pair of channels
// and write to a serial Port
// TODO provide the notion of command with a separator
func (p *pty) Open()  error {

	if err := p.ptyOpen(); err != nil {
		log.Printf("[ERROR] Pty: Cannot open port: %s", err)
		return err
	}

	go p.receiveloop()
	go p.sendloop()

	return nil
}

func NewPty() modules.Module {

	return &pty {
		channels: make([]modules.Channels,0),
	}
}
