package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/leoleovich/serialchiller/feeder"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	f := &feeder.Feeder{
		Cancel: cancel,
	}

	flag.StringVar(&f.Device, "device", "", "Serial device to interact with")
	flag.IntVar(&f.BaudRate, "baudrate", 115200, "Serial device speed")
	flag.Parse()

	if f.Device == "" {
		panic("Must set a valid device name")
	}

	if err := f.Open(); err != nil {
		panic(err)
	}


	go f.Read(ctx)

	fmt.Println("Enter your commands:")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return
		default:
			cmd := scanner.Text()
			f.Write(cmd)
		}
	}
}
