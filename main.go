package main

import (
	"github.com/mlen/ussd/pack"
	"github.com/tarm/serial"

	"bufio"
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"strings"
)

var (
	port = flag.String("port", "/dev/ttyUSB0", "serial port to open")
	baud = flag.Int("baud", 9600, "baud rate")
)

type Mode int

const (
	_ Mode = iota
	Send
	Cancel
)

func mustOpenPort(port string, baud int) io.ReadWriteCloser {
	c := &serial.Config{Name: port, Baud: baud}
	p, err := serial.OpenPort(c)
	if err != nil {
		panic(err)
	}
	return p
}

func sendCusd(port io.Writer, mode Mode, msg string) error {
	var command string
	if mode == Cancel {
		command = fmt.Sprintf("AT+CUSD=%d", mode)
	} else {
		command = fmt.Sprintf("AT+CUSD=%d,\"%02X\",15\r", mode, pack.Pack7Bit([]byte(msg)))
	}

	_, err := port.Write([]byte(command))
	return err
}

func parseCusd(msg string) (string, error) {
	encoded := strings.Split(msg, "\"")[1]
	data, err := hex.DecodeString(encoded)

	if err != nil {
		return "", err
	}

	return string(pack.Unpack7Bit(data)), nil
}

func reader(ctx context.Context, port io.ReadCloser, lines chan<- string) {
	rd := bufio.NewReader(port)

	for {
		select {
		case <-ctx.Done():
			port.Close()
			return
		default:
			s, err := rd.ReadString('\r')
			if err != nil && err != io.EOF {
				panic(err)
			}

			s = strings.Trim(s, "\r\n")
			if strings.HasPrefix(s, "+CUSD: ") {
				lines <- strings.Trim(s, "\r\n")
			}
		}
	}
}

func printResponse(line string) {
	response, err := parseCusd(line)
	if err != nil {
		panic(err)
	}

	fmt.Println(response)
}

func main() {
	flag.Parse()

	r := mustOpenPort(*port, *baud)

	ctx, cancel := context.WithCancel(context.Background())
	lines := make(chan string)
	go reader(ctx, r, lines)

	r.Write([]byte("\rAT\r"))
	sendCusd(r, Send, flag.Arg(0))

	line := <-lines
	switch line[7] {
	case '0':
		printResponse(line)

	case '1':
		printResponse(line)
		sendCusd(r, Cancel, "")

	default:
		fmt.Println("Error:", line[7:])
	}

	cancel()
}
