package terminal

import (
	"devterminal2/utils"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

type TCPServerTerm struct {
	Terminal
	_port        int
	_tcpServer   *net.TCPListener
	_tcpInClient *net.TCPConn
	_active      bool
	_lastError   error
	_netOutput   *utils.QueueEasy
}

func OpenTCPServer(port int) (*TCPServerTerm, error) {
	term := TCPServerTerm{
		Terminal:   Terminal{TypeId: TerminalTcpServer},
		_netOutput: &utils.QueueEasy{},
		_port:      port,
		_active:    true,
	}

	var err error
	var addr *net.TCPAddr

	addr, err = net.ResolveTCPAddr("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	term._tcpServer, err = net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	go lisentServer(&term)

	return &term, nil
}

func lisentServer(t *TCPServerTerm) {

	if t._tcpInClient != nil {
		t._tcpInClient.Close()
		t._tcpInClient = nil
	}

	for t._active {
		t._tcpInClient, t._lastError = t._tcpServer.AcceptTCP()

		if t._lastError != nil {
			break
		}

		go handleConn(t, t._tcpInClient)

	}
}

func handleConn(t *TCPServerTerm, conn *net.TCPConn) {
	for {
		var n int
		var arr []byte
		var err error

		arr, err = t._netOutput.PickArray()

		n, err = conn.Write(arr)
		if err != nil {
			break
		}
		if n > 0 {
			t._netOutput.PopArray(len(arr))
		}
	}
}

func (t *TCPServerTerm) Close() error {
	t._active = false
	return t._tcpServer.Close()
}

func (t *TCPServerTerm) Read(data []byte) (int, error) {

	if t._tcpInClient == nil {
		return 0, nil
	}

	t._tcpInClient.SetReadDeadline(time.Now().Add(time.Millisecond * 50))
	n, errRead := t._tcpInClient.Read(data)
	//var ok bool
	if err, ok := errRead.(net.Error); ok && err.Timeout() {
		errRead = nil

	}

	if errRead == io.EOF {
		t._tcpInClient.Close()
		t._tcpInClient = nil
		errRead = nil
	}

	return n, errRead
}

func (t *TCPServerTerm) Write(data []byte) (int, error) {
	t._netOutput.PushArray(data)
	return len(data), nil
}
