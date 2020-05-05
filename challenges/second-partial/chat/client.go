// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 227.

// Netcat is a simple read/write client for TCP servers.
package main

import (
	"io"
	"log"
	"net"
	"os"
	"flag"
	"fmt"
)

//!+
func main() {
	username := flag.String("user", "____", "a string")
	server := flag.String("server", "localhost:8000", "a string")
	flag.Parse()

	if *username == "____" {
		log.Fatal("Must provide a username\n")
	}

	conn, err := net.Dial("tcp", *server)
	// conn.Write([]byte(user))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(conn, *username+"\n")
	done := make(chan struct{})
	go func() {
		io.Copy(os.Stdout, conn) // NOTE: ignoring errors
		log.Println("done")
		done <- struct{}{} // signal the main goroutine
	}()
	mustCopy(conn, os.Stdin)
	conn.Close()
	<-done // wait for background goroutine to finish
}

//!-

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
