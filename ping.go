package tools

import (
	"fmt"
	"net"
	"os"
	"time"
)

// Ping ...
func Ping(url string, timeout time.Duration) {
	cancel := time.After(timeout)
	success := make(chan struct{})
	go func(s chan struct{}) {
		for {
			_, err := net.Dial("tcp", url)
			if err == nil {
				success <- struct{}{}
			}
			time.Sleep(time.Second)
		}
	}(success)
	select {
	case <-success:
		fmt.Println("connected")
	case <-cancel:
		fmt.Printf("Can't connected to %s\n", url)
		os.Exit(0)
	}
}
