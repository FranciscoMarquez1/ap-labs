// Clock2 is a concurrent TCP server that periodically writes the time.
package main

import (
	"io"
	"log"
	"net"
	"time"
	"os"
	"fmt"
	"flag"
)

func handleConn(c net.Conn, timezone string) {
	defer c.Close()
	for {
		// t, err := TimeIn(time.Now(), "Asia/Shanghai")
		location, err := time.LoadLocation(timezone)
		if err == nil {
			_, err := io.WriteString(c, timezone + " " + time.Now().In(location).Format("15:04:05\n"))
			if err != nil {
				return // e.g., client disconnected
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func main() {
	port := flag.String("port", "9090", "a string")
	flag.Parse()
	fmt.Println("port:", *port)
	tz := os.Getenv("TZ")
	listener, err := net.Listen("tcp", "localhost:"+ *port)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		go handleConn(conn, tz) // handle connections concurrently
	}
}
