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

	serverPort := "26342"
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
		// fmt.Printf("start Handle Msg ")
		buffer := make([]byte, 1024)

		count, err := conn.Read(buffer)
		if err != nil {
			handleError(conn, err, "client disconnected!")
			return
		}

		_ = count

		totalRequests++

		tempStr := string(buffer)

		time.Sleep(time.Second * 3)
		requestOption, _ := strconv.Atoi(strings.Split(tempStr, "|")[0])
		requestData := strings.Split(tempStr, "|")[1]

		// fmt.Println(requestOption)
		// fmt.Println(requestData)

		switch requestOption {
		case 1:
			// fmt.Printf("========11111======== ")

			upperString := strings.ToUpper(requestData)
			conn.Write([]byte(upperString))

		case 2:
			// fmt.Printf("========22222======== ")

			ip := conn.RemoteAddr()
			conn.Write([]byte(ip.String()))

		case 3:
			// fmt.Printf("========33333======== ")

			conn.Write([]byte(strconv.Itoa(totalRequests)))

		case 4:
			// fmt.Printf("4444 ")

			elapsed := time.Since(startTime).Truncate(time.Second).String()
			conn.Write([]byte(string(elapsed)))

		case 5:
			// fmt.Printf("5555 ")
			conn.Close()

			os.Exit(0)
		default:
			// fmt.Printf("default ")
			conn.Close()

			os.Exit(0)
		}
	}

}
