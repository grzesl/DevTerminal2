package main

import (
	"devterminal2/terminal"
	"encoding/hex"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io"
	"log"
	"os"
	"slices"
	"strconv"
	"time"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/rect"
	"github.com/aarzilli/nucular/style"
	"go.bug.st/serial"
	"golang.org/x/mobile/event/key"
	"gopkg.in/ini.v1"
)

var TerminalPorts []terminal.ITerminal = []terminal.ITerminal{
	&terminal.SerialTerm{Terminal: terminal.Terminal{TypeId: terminal.TerminalSerial}},
	&terminal.TCPClientTerm{Terminal: terminal.Terminal{TypeId: terminal.TerminalTcpClient}},
	&terminal.TCPServerTerm{Terminal: terminal.Terminal{TypeId: terminal.TerminalTcpServer}},
}

var cfgIni *ini.File = nil

var SelectedTerminalPortTypeIdInt32 = 0
var SelectedTerminalPortTypeId terminal.TerminalType = 0

var SelectedPortNameId int = 0
var SelectedPortBaudRateId int = 0
var SelectedPortName string = ""

var SelectedHost string = "192.168.1.1"
var SelectedPortNumber int = 3333

var SelectedServerPortNumber int = 3333

var TextOpenDisabled bool = false
var PortDetails = ""
var TextToSend string = ""
var SendAsHex bool = false
var SendAsSymbol bool = true

var SendInLoop bool = false
var LoopInterval int32 = 1000

var led_red_RGBA *image.RGBA
var led_green_RGBA *image.RGBA

var LogBuffer = []string{"log..."}

func main() {

	var err error
	cfgIni, err = ini.Load("cfg.ini")
	if err != nil {
		prepareDefaultSettings()
	}

	f1, _ := led_red_alpha_30.Open("led_red_alpha_30.png")
	led_red, _ := png.Decode(f1)
	led_red_RGBA = image.NewRGBA(led_red.Bounds())
	draw.Draw(led_red_RGBA, led_red.Bounds(), led_red, image.Point{}, draw.Src)

	f2, _ := led_green_alpha_30.Open("led_green_alpha_30.png")
	led_green, _ := png.Decode(f2)
	led_green_RGBA = image.NewRGBA(led_green.Bounds())
	draw.Draw(led_green_RGBA, led_green.Bounds(), led_green, image.Point{}, draw.Src)

	SelectedPortName = cfgIni.Section("SerialPortRS232").Key("port_name").String()
	baudRate := cfgIni.Section("SerialPortRS232").Key("baud_rate").String()
	SelectedPortBaudRateId = int(slices.Index(getSerialPortBaudRates(), baudRate))

	SelectedHost = cfgIni.Section("TCPClient").Key("host_addr").String()
	var intval int
	intval, err = strconv.Atoi(cfgIni.Section("TCPClient").Key("host_port").String())
	SelectedPortNumber = int(intval)

	var intval3 int
	intval3, err = strconv.Atoi(cfgIni.Section("TCPServer").Key("lisent_port").String())
	SelectedServerPortNumber = int(intval3)

	var intval2 int
	intval2, err = strconv.Atoi(cfgIni.Section("Common").Key("selected_terminal").String())
	SelectedTerminalPortTypeIdInt32 = int(intval2)
	TextToSend = cfgIni.Section("Common").Key("last_send").String()

	wnd := nucular.NewMasterWindow(0, "DevTerminal 2.00", updatefn)
	wnd.SetStyle(style.FromTheme(style.DarkTheme, 1.0))
	wnd.Main()
}

func getSerialPortBaudRates() []string {
	return []string{
		"9600",
		"19200",
		"38400",
		"57600",
		"115200",
	}
}

func connectionTypesCb() []string {

	types := make([]string, len(TerminalPorts))
	for i, t := range TerminalPorts {
		types[i] = t.String()
	}
	return types
}

func sendInLoop(w *nucular.Window) {

	ticker := time.NewTicker(time.Millisecond * time.Duration(LoopInterval))

	for SendInLoop {
		sendCurrentData()
		w.Master().Changed()
		<-ticker.C
	}

}

func sendCurrentData() error {
	var outputData []byte
	var err error

	if SendAsHex {
		outputData, err = hex.DecodeString(TextToSend)
		if err != nil {
			return err
		}
	} else if SendAsSymbol {
		outputData = terminal.ConvertStrToChars(TextToSend)
	} else {
		outputData = []byte(TextToSend)
	}

	logData(fmt.Sprintf("%s -> %s",
		time.Now().Local().Format("[15:04:05.000]"),
		terminal.ConvertCharsToStr(outputData)))

	TerminalPorts[SelectedTerminalPortTypeId].Write(outputData)

	return nil
}

func readCurretData() {
	data := make([]byte, 64)

	n, err := TerminalPorts[SelectedTerminalPortTypeId].Read(data)

	if err == io.EOF {
		terminalClose()
		logData(err.Error())
		return
	}

	if err != nil {
		return
	}

	if n > 0 {
		logData(fmt.Sprintf("%s <- %s",
			time.Now().Local().Format("[15:04:05.000]"),
			terminal.ConvertCharsToStr(data[0:n])))
	}
}

func constantReadTerminal() {
	for TextOpenDisabled {
		readCurretData()
		time.Sleep(time.Millisecond * 10)
	}
}

func terminalOpen() error {
	var arg1, arg2, arg3 string
	var err error
	var term terminal.ITerminal

	if SelectedTerminalPortTypeId == terminal.TerminalSerial {
		arg1 = SelectedPortName
		arg2 = getSerialPortBaudRates()[SelectedPortBaudRateId]
		arg3 = ""
	} else if SelectedTerminalPortTypeId == terminal.TerminalTcpClient {
		arg1 = SelectedHost
		arg2 = strconv.Itoa(int(SelectedPortNumber))
		arg3 = ""
	} else if SelectedTerminalPortTypeId == terminal.TerminalTcpServer {
		arg1 = strconv.Itoa(int(SelectedServerPortNumber))
		arg2 = ""
		arg3 = ""
	}

	term, err = TerminalPorts[SelectedTerminalPortTypeId].Open(
		SelectedTerminalPortTypeId,
		arg1,
		arg2,
		arg3)

	if err != nil {
		return err
	} else {
		TerminalPorts[SelectedTerminalPortTypeId] = term
		go constantReadTerminal()
	}

	return nil
}

var editorPortName nucular.TextEditor
var editorHostName nucular.TextEditor

func showSettingsPopup(w *nucular.Window, terminalType terminal.TerminalType) {
	w.Master().PopupOpen("Settings", nucular.WindowMovable|nucular.WindowTitle|nucular.WindowDynamic|nucular.WindowNoScrollbar, rect.Rect{20, 100, 300, 150}, true, func(w2 *nucular.Window) {

		if terminalType == terminal.TerminalSerial {

			portNames, err := serial.GetPortsList()

			portNames = append([]string{"(none)"}, portNames...)

			baudRates := getSerialPortBaudRates()
			if err != nil {
				portNames = append(portNames, err.Error()) //show error in combo
			}
			w2.Row(20).Ratio(0.4, 0.1, 0.4)
			editorPortName.Flags = nucular.EditField
			editorPortName.Maxlen = 255
			editorPortName.Buffer = []rune(SelectedPortName)
			editorPortName.Filter = nucular.FilterDefault
			editorPortName.Edit(w2)

			if w2.ButtonText("X") {
				SelectedPortName = ""
				SelectedPortNameId = 0
				editorPortName.Buffer = []rune(SelectedPortName)
			}
			SelectedPortName = string(editorPortName.Buffer)
			SelectedPortNameId = w2.ComboSimple(portNames, SelectedPortNameId, 20)

			if SelectedPortNameId > 0 {
				SelectedPortName = portNames[SelectedPortNameId]
			}

			w2.Row(20).Ratio(0.4)

			SelectedPortBaudRateId = w2.ComboSimple(baudRates, SelectedPortBaudRateId, 20)

		} else if terminalType == terminal.TerminalTcpClient {
			w2.Row(20).Ratio(0.5, 0.5)
			editorHostName.Flags = nucular.EditField
			editorHostName.Maxlen = 255
			editorHostName.Buffer = []rune(SelectedHost)
			editorHostName.Filter = nucular.FilterDefault
			editorHostName.Edit(w2)
			SelectedHost = string(editorHostName.Buffer)
			w2.PropertyInt("Port", 0, &SelectedPortNumber, 65535, 1, 1)

		} else if terminalType == terminal.TerminalTcpServer {
			w2.Row(20).Ratio(0.5, 0.5)
			w2.PropertyInt("Port", 0, &SelectedServerPortNumber, 65535, 1, 1)

		}

		w2.Row(20).Ratio(0.3, 0.4, 0.3)

		if w2.ButtonText("Save") {
			applyTerminalPortSettings(SelectedTerminalPortTypeId)
		}

		w2.Spacing(1)

		if w2.ButtonText("Close") {
			w2.Close()
		}
	})
}

func prepareDefaultSettings() {
	f, err := os.Create("cfg.ini")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	cfgIni, err = ini.Load("cfg.ini")
	if err != nil {
		log.Fatal(err)
	}
	cfgIni.Section("SerialPortRS232").Key("port_name").SetValue("COM1")
	cfgIni.Section("SerialPortRS232").Key("baud_rate").SetValue(getSerialPortBaudRates()[0])
	cfgIni.Section("TCPClient").Key("host_addr").SetValue("192.168.1.1")
	cfgIni.Section("TCPClient").Key("host_port").SetValue(strconv.Itoa(int(3333)))
	cfgIni.Section("TCPServer").Key("lisent_port").SetValue(strconv.Itoa(int(3333)))
	cfgIni.Section("Common").Key("selected_terminal").SetValue(strconv.Itoa(int(0)))
	cfgIni.Section("Common").Key("last_send").SetValue("TEST")
	cfgIni.SaveTo("cfg.ini")

	fmt.Println("prepareDefaultSettings Done")
}

func applyTerminalPortSettings(tt terminal.TerminalType) {

	if tt == terminal.TerminalSerial {
		cfgIni.Section("SerialPortRS232").Key("port_name").SetValue(SelectedPortName)
		cfgIni.Section("SerialPortRS232").Key("baud_rate").SetValue(getSerialPortBaudRates()[SelectedPortBaudRateId])
	} else if tt == terminal.TerminalTcpClient {
		cfgIni.Section("TCPClient").Key("host_addr").SetValue(SelectedHost)
		cfgIni.Section("TCPClient").Key("host_port").SetValue(strconv.Itoa(int(SelectedPortNumber)))
	} else if tt == terminal.TerminalTcpServer {
		cfgIni.Section("TCPServer").Key("lisent_port").SetValue(strconv.Itoa(int(SelectedServerPortNumber)))
	}
	cfgIni.SaveTo("cfg.ini")
}

func applyCommonSettings() {
	cfgIni.Section("Common").Key("selected_terminal").SetValue(strconv.Itoa(int(SelectedTerminalPortTypeIdInt32)))
	cfgIni.Section("Common").Key("last_send").SetValue(TextToSend)
	cfgIni.SaveTo("cfg.ini")
}

func terminalClose() error {
	err := TerminalPorts[SelectedTerminalPortTypeId].Close()
	if err != nil {
		return err
	}
	SendInLoop = false
	return nil
}

func showMessage(w *nucular.Window, hdr string, msg string) {

	w.Master().PopupOpen(hdr, nucular.WindowMovable|nucular.WindowTitle|nucular.WindowDynamic|nucular.WindowNoScrollbar, rect.Rect{20, 100, 400, 150}, true, func(w *nucular.Window) {
		w.Row(25).Dynamic(1)
		w.Label(msg, "LC")
		w.Row(25).Static(0, 100)
		w.Spacing(1)
		if w.ButtonText("OK") {
			w.Close()
		}
	})
}

var listDemoSelected = -1
var listDemoCnt = 0

func logList(w *nucular.Window) {
	var N = len(LogBuffer)
	recenter := false
	for _, e := range w.Input().Keyboard.Keys {
		switch e.Code {
		case key.CodeDownArrow:
			listDemoSelected++
			if listDemoSelected >= N {
				listDemoSelected = N - 1
			}
			recenter = true
		case key.CodeUpArrow:
			listDemoSelected--
			if listDemoSelected < -1 {
				listDemoSelected = -1
			}
			recenter = true
		}
	}
	w.Row(w.Bounds.H - 145).Dynamic(1)
	if gl, w := nucular.GroupListStart(w, N, "list", nucular.WindowNoHScrollbar); w != nil {

		if TextOpenDisabled {
			recenter = true
			listDemoSelected = N - 1
		} else {
			recenter = false
		}

		if !recenter {
			gl.SkipToVisible(20)
		}
		w.Row(20).Dynamic(1)
		cnt := 0
		for gl.Next() {
			cnt++
			i := gl.Index()
			selected := i == listDemoSelected
			w.SelectableLabel(LogBuffer[i], "LC", &selected)
			if selected {
				listDemoSelected = i
				if recenter {
					gl.Center()
				}
			}
		}
	}
}

func logData(msg string) {
	LogBuffer = append(LogBuffer, msg)
}

var editorSendData nucular.TextEditor

func updatefn(w *nucular.Window) {
	w.Row(30).Static(30, 120, 120, 80)

	if TextOpenDisabled {
		w.Image(led_green_RGBA)
	} else {
		w.Image(led_red_RGBA)
	}

	SelectedTerminalPortTypeId = terminal.TerminalType(SelectedTerminalPortTypeIdInt32)
	if SelectedTerminalPortTypeId == terminal.TerminalSerial {
		PortDetails = SelectedPortName + " :" + getSerialPortBaudRates()[SelectedPortBaudRateId]
	} else if SelectedTerminalPortTypeId == terminal.TerminalTcpClient {
		PortDetails = SelectedHost + " :" + strconv.Itoa(int(SelectedPortNumber))
	} else if SelectedTerminalPortTypeId == terminal.TerminalTcpServer {
		PortDetails = strconv.Itoa(int(SelectedServerPortNumber))
	}

	var openCloseText = "Open"
	if TextOpenDisabled {
		openCloseText = "Close"
		w.Label(connectionTypesCb()[SelectedTerminalPortTypeIdInt32], "LC")
		w.Label(PortDetails, "LC")
	} else {
		SelectedTerminalPortTypeIdInt32 = w.ComboSimple(connectionTypesCb(), SelectedTerminalPortTypeIdInt32, 20)
		SelectedTerminalPortTypeId = terminal.TerminalType(SelectedTerminalPortTypeIdInt32)
		if w.ButtonText("Settings") {
			showSettingsPopup(w, SelectedTerminalPortTypeId)
		}
	}

	if w.ButtonText(openCloseText) {
		if TextOpenDisabled { //opened
			err := terminalClose()
			if err != nil {
				showMessage(w, "Error", err.Error())
			} else {
				logData("terminal closed...")
			}
			TextOpenDisabled = !TextOpenDisabled
		} else {
			err := terminalOpen()
			if err != nil {
				showMessage(w, "Error", err.Error())
			} else {
				logData("terminal opened...")
				TextOpenDisabled = !TextOpenDisabled
			}
		}
	}

	logList(w)

	w.Row(40).Ratio(0.9, 0.1)
	editorSendData.Flags = nucular.EditField
	editorSendData.Maxlen = 255
	editorSendData.Buffer = []rune(TextToSend)
	editorSendData.Filter = nucular.FilterDefault
	editorSendData.Edit(w)
	TextToSend = string(editorSendData.Buffer)

	if TextOpenDisabled {
		if w.ButtonText("SEND ->") {
			sendCurrentData()
		}
	} else {
		w.Label("CLOSED", "LC")
	}

	w.Row(40).Ratio(0.1, 0.1, 0.2)

	if w.CheckboxText("Loop", &SendInLoop) {
		if !TextOpenDisabled {
			SendInLoop = false
			return
		}

		if SendInLoop {
			go sendInLoop(w)
		}
	}
	if w.CheckboxText("Hex", &SendAsHex) {
		SendAsSymbol = false
	}
	if w.CheckboxText("Symbol", &SendAsSymbol) {
		SendAsHex = false
	}

}
