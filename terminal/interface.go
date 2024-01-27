package terminal

import (
	"errors"
	"strconv"
	"strings"
)

type TerminalType int

const (
	TerminalSerial TerminalType = iota
	TerminalTcpClient
	TerminalTcpServer
)

type ITerminal interface {
	Open(tt TerminalType, args ...string) (ITerminal, error)
	Close() error
	Read(data []byte) (int, error)
	Write(data []byte) (int, error)
	String() string
}

type Terminal struct {
	TypeId TerminalType
}

func (t *Terminal) Open(tt TerminalType, args ...string) (ITerminal, error) {

	var term ITerminal
	var openErr error

	switch tt {
	case TerminalSerial:
		{
			intVal, err := strconv.Atoi(args[1])
			if err != nil {
				openErr = err
				break
			}

			serialTerm, err := OpenSerial(args[0],
				intVal, args[2])

			term = serialTerm
			openErr = err
		}
	case TerminalTcpClient:
		{
			intVal, err := strconv.Atoi(args[1])
			if err != nil {
				openErr = err
				break
			}

			rcpclientTerm, err := OpenTCPClient(args[0],
				intVal)

			term = rcpclientTerm
			openErr = err
		}
	case TerminalTcpServer:
		{
			intVal, err := strconv.Atoi(args[0])
			if err != nil {
				openErr = err
				break
			}

			rcpclientTerm, err := OpenTCPServer(intVal)

			term = rcpclientTerm
			openErr = err
		}
	default:
		{
			term = nil
			openErr = errors.New("Unknow Terminal Type")
		}
	}

	return term, openErr
}

func (t *Terminal) String() string {

	switch t.TypeId {
	case TerminalSerial:
		return "Serial"
	case TerminalTcpClient:
		return "TCPClient"
	case TerminalTcpServer:
		return "TCPServer"
	}

	return ""
}

func ConvertStrToChars(str string) []byte {
	newString := strings.ReplaceAll(str, "<NUL>", "\x00")
	newString = strings.ReplaceAll(newString, "<SOH>", "\x01")
	newString = strings.ReplaceAll(newString, "<STX>", "\x02")
	newString = strings.ReplaceAll(newString, "<ETX>", "\x03")
	newString = strings.ReplaceAll(newString, "<EOT>", "\x04")
	newString = strings.ReplaceAll(newString, "<ENQ>", "\x05")
	newString = strings.ReplaceAll(newString, "<ACK>", "\x06")
	newString = strings.ReplaceAll(newString, "<BEL>", "\x07")
	newString = strings.ReplaceAll(newString, "<BS>", "\x08")
	newString = strings.ReplaceAll(newString, "<TAB>", "\x09")
	newString = strings.ReplaceAll(newString, "<LF>", "\x0A")
	newString = strings.ReplaceAll(newString, "<VT>", "\x0B")
	newString = strings.ReplaceAll(newString, "<FF>", "\x0C")
	newString = strings.ReplaceAll(newString, "<CR>", "\x0D")
	newString = strings.ReplaceAll(newString, "<SO>", "\x0E")
	newString = strings.ReplaceAll(newString, "<SI>", "\x0F")
	newString = strings.ReplaceAll(newString, "<DLE>", "\x10")
	newString = strings.ReplaceAll(newString, "<DC1>", "\x11")
	newString = strings.ReplaceAll(newString, "<DC2>", "\x12")
	newString = strings.ReplaceAll(newString, "<DC3>", "\x13")
	newString = strings.ReplaceAll(newString, "<DC4>", "\x14")
	newString = strings.ReplaceAll(newString, "<NAK>", "\x15")
	newString = strings.ReplaceAll(newString, "<SYN>", "\x16")
	newString = strings.ReplaceAll(newString, "<ETB>", "\x17")
	newString = strings.ReplaceAll(newString, "<CAN>", "\x18")
	newString = strings.ReplaceAll(newString, "<EM>", "\x19")
	newString = strings.ReplaceAll(newString, "<SUB>", "\x1A")
	newString = strings.ReplaceAll(newString, "<ESC>", "\x1B")
	newString = strings.ReplaceAll(newString, "<FS>", "\x1C")
	newString = strings.ReplaceAll(newString, "<GS>", "\x1D")
	newString = strings.ReplaceAll(newString, "<RS>", "\x1E")
	newString = strings.ReplaceAll(newString, "<US>", "\x1F")

	return []byte(newString)
}

func ConvertCharsToStr(data []byte) string {
	var str string = ""

	for _, v := range data {
		var charStr string = ""

		switch v {
		case 0x0:
			charStr = "<NUL>"
		case 0x1:
			charStr = "<SOH>"
		case 0x2:
			charStr = "<STX>"
		case 0x3:
			charStr = "<ETX>"
		case 0x4:
			charStr = "<EOT>"
		case 0x5:
			charStr = "<ENQ>"
		case 0x6:
			charStr = "<ACK>"
		case 0x7:
			charStr = "<BEL>"
		case 0x8:
			charStr = "<BS>"
		case 0x9:
			charStr = "<TAB>"
		case 0xA:
			charStr = "<LF>"
		case 0xB:
			charStr = "<VT>"
		case 0xC:
			charStr = "<FF>"
		case 0xD:
			charStr = "<CR>"
		case 0xE:
			charStr = "<SO>"
		case 0xF:
			charStr = "<SI>"
		case 0x10:
			charStr = "<DLE>"
		case 0x11:
			charStr = "<DC1>"
		case 0x12:
			charStr = "<DC2>"
		case 0x13:
			charStr = "<DC3>"
		case 0x14:
			charStr = "<DC4>"
		case 0x15:
			charStr = "<NAK>"
		case 0x16:
			charStr = "<SYN>"
		case 0x17:
			charStr = "<ETB>"
		case 0x18:
			charStr = "<CAN>"
		case 0x19:
			charStr = "<EM>"
		case 0x1A:
			charStr = "<SUB>"
		case 0x1B:
			charStr = "<ESC>"
		case 0x1C:
			charStr = "<FS>"
		case 0x1D:
			charStr = "<GS>"
		case 0x1E:
			charStr = "<RS>"
		case 0x1F:
			charStr = "<US>"

		default:
			charStr = string(v)
		}

		str += charStr
	}

	return str
}
