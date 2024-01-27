package terminal

import (
	"fmt"
	"time"

	"go.bug.st/serial"
)

const (
	HandhakeNone     = "HandshakeNone"
	HandhakeHardware = "HandshakeHardware"
)

type SerialTerm struct {
	Terminal
	_portName   string
	_baudRate   int
	_handshake  string
	_serialPort serial.Port
}

func OpenSerial(port string, baudRate int, handshake string) (*SerialTerm, error) {
	term := SerialTerm{
		Terminal:   Terminal{TypeId: TerminalSerial},
		_portName:  port,
		_baudRate:  baudRate,
		_handshake: handshake,
	}

	var err error
	term._serialPort, err = serial.Open(port, &serial.Mode{BaudRate: baudRate, DataBits: 8, Parity: serial.NoParity, StopBits: 0})

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	term._serialPort.SetReadTimeout(time.Millisecond * 50)

	return &term, nil
}

func (t *SerialTerm) Close() error {
	return t._serialPort.Close()
}

func (t *SerialTerm) Read(data []byte) (int, error) {
	return t._serialPort.Read(data)
}

func (t *SerialTerm) Write(data []byte) (int, error) {
	return t._serialPort.Write(data)
}
