package terminal

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

type TCPClientTerm struct {
	Terminal
	_ip        string
	_port      int
	_tcpClient *net.TCPConn
}

func OpenTCPClient(ip string, port int) (*TCPClientTerm, error) {
	term := TCPClientTerm{
		Terminal: Terminal{TypeId: TerminalTcpClient},
		_ip:      ip,
		_port:    port,
	}

	var err error
	tcpAddr, err := net.ResolveTCPAddr("tcp", ip+":"+strconv.Itoa(port))

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	term._tcpClient, err = net.DialTCP("tcp", nil, tcpAddr)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &term, nil
}

func (t *TCPClientTerm) Close() error {
	return t._tcpClient.Close()
}

func (t *TCPClientTerm) Read(data []byte) (int, error) {
	t._tcpClient.SetReadDeadline(time.Now().Add(time.Millisecond * 50))
	n, errRead := t._tcpClient.Read(data)
	//var ok bool
	if err, ok := errRead.(net.Error); ok && err.Timeout() {
		errRead = nil
	}

	return n, errRead
}

func (t *TCPClientTerm) Write(data []byte) (int, error) {
	return t._tcpClient.Write(data)
}
