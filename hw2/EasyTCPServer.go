/**
 * TCPServer.go
 **/

package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var totalRequests int = 0
var startTime time.Time
var serverPort string = "26342"

func main() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for sig := range c {
			// sig is a ^C, handle it
			_ = sig
			byebye()
		}
	}()

	startTime = time.Now()
	listener, _ := net.Listen("tcp", ":"+serverPort)
	fmt.Printf("Server is ready to receive on port %s\n", serverPort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			handleError(conn, err, "server accept error..")
		}

		fmt.Printf("Connection request from %s\n", conn.RemoteAddr().String())

		go handleMsg(conn)
	}

	defer byebye()

}

func byebye() {
	fmt.Printf("Bye bye~\n")
	os.Exit(0)
}

func handleError(conn net.Conn, err error, errmsg string) {
	if conn != nil {
		conn.Close()
	}
	// fmt.Println(err)
	fmt.Println(errmsg)
}
func handleError2(conn net.Conn, errmsg string) {
	if conn != nil {
		conn.Close()
	}

	fmt.Println(errmsg)
}

func handleMsg(conn net.Conn) {
	for {
		buffer := make([]byte, 1024)

		count, err := conn.Read(buffer)
		if err != nil {
			handleError(conn, err, "client disconnected!")
			return
		}

		_ = count

		totalRequests++
		tempStr := string(buffer)
		requestOption, _ := strconv.Atoi(strings.Split(tempStr, "|")[0])
		requestData := strings.Split(tempStr, "|")[1]
		time.Sleep(time.Millisecond * 1)
		fmt.Printf("Command %d\n\n", requestOption)

		switch requestOption {
		case 1:
			upperString := strings.ToUpper(requestData)
			sendPacket(conn, upperString)
		case 2:
			ip := conn.RemoteAddr()
			sendPacket(conn, ip.String())
		case 3:
			sendPacket(conn, strconv.Itoa(totalRequests))
		case 4:
			elapsed := time.Since(startTime).Truncate(time.Second).String()

			sendPacket(conn, string(elapsed))
		case 5:
			conn.Close()
			os.Exit(0)
		default:
			conn.Close()
			os.Exit(0)
		}
	}

}

func sendPacket(conn net.Conn, serverMsg string) {
	conn.Write([]byte(serverMsg))
}
