package feeder

import (
	"bufio"
	"context"
	"fmt"

	"go.bug.st/serial"
)

type Feeder struct {
	Device   string
	BaudRate int
	Cancel context.CancelFunc

	writer *bufio.Writer
	reader *bufio.Reader
}

func (f *Feeder) Open() error {
	mode := &serial.Mode{
		BaudRate: f.BaudRate,
	}
	tty, err := serial.Open(f.Device, mode)
	if err != nil {
		return err
	}
	f.reader = bufio.NewReader(tty)
	f.writer = bufio.NewWriter(tty)

	return nil
}

func (f *Feeder) Read(ctx context.Context) {
	defer f.Cancel()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			buf, _, err := f.reader.ReadLine()
			if err != nil {
				fmt.Println(err)
			}
			bufStr := string(buf)
			fmt.Println("Feeder: READING: ", bufStr)
		}
	}
}

func (f *Feeder) Write(command string) {
	f.writer.Write([]byte(fmt.Sprintf("%s\n", command)))
	f.writer.Flush()
}
