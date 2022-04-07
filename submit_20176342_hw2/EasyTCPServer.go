/**
 * 20176342 Song Min Joon
 * EasyTCPServer.go
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

var totalRequests int = 0       // total Request count global variable for server.
var startTime time.Time         // for saving server start time
var serverPort string = "26342" // for server port

func main() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // for exit the program gracefully
	go func() {
		for sig := range c {
			_ = sig
			byebye() // if this program is interrupt by Ctrl-c, print Bye bye and exits gracefully
		}
	}()

	startTime = time.Now() // records server start time for server running time
	listener, _ := net.Listen("tcp", ":"+serverPort)
	fmt.Printf("Server is ready to receive on port %s\n", serverPort)

	for {
		//listener is waiting for tcp connection of clients.
		conn, err := listener.Accept()

		if err != nil {
			handleError(conn, err, "server accept error..")
		}

		fmt.Printf("Connection request from %s\n", conn.RemoteAddr().String())

		go handleMsg(conn) // when client is connect to server, make go-routine to communicate with client.
	}

	defer byebye() // although when client gets panic, defer should disconnect socket gracefully

}

func byebye() {
	fmt.Printf("Bye bye~\n")
	os.Exit(0)
}

func handleError(conn net.Conn, err error, errmsg string) {
	//handle error and print

	if conn != nil {
		conn.Close()
	}
	fmt.Println(errmsg)
}

func handleMsg(conn net.Conn) {
	for {
		buffer := make([]byte, 1024)

		count, err := conn.Read(buffer)
		//when client sends packet
		if err != nil {
			handleError(conn, err, "client disconnected!")
			return
		}

		_ = count

		totalRequests++ // number of request is added

		/*
			client packet form

			(Option Number|Message)

			1|blah blah blah...
			1|hello world!

			2|     => message is not required
			3|     => message is not required
			4|     => message is not required

			5|     => maybe not arrived??

		*/

		tempStr := string(buffer)
		requestOption, _ := strconv.Atoi(strings.Split(tempStr, "|")[0]) // split client packet by '|' and takes option and convert to Integer.
		requestData := strings.Split(tempStr, "|")[1]                    // get message parameter from packet.

		time.Sleep(time.Millisecond * 1) // minimum delay to deliver packet to client.

		fmt.Printf("Command %d\n\n", requestOption) // print Command #

		switch requestOption {
		case 1:
			//Option 1
			upperString := strings.ToUpper(requestData)
			sendPacket(conn, upperString)
		case 2:
			//Option 2
			ip := conn.RemoteAddr()
			sendPacket(conn, ip.String())
		case 3:
			//Option 3
			sendPacket(conn, strconv.Itoa(totalRequests))
		case 4:
			//Option 4
			elapsed := time.Since(startTime).Truncate(time.Second).String()
			sendPacket(conn, string(elapsed))
		case 5:
			//Option 5
			conn.Close()
		default:
			//Option default

			conn.Close()
		}
	}

}

func sendPacket(conn net.Conn, serverMsg string) {
	//send packet to client
	conn.Write([]byte(serverMsg))
}
