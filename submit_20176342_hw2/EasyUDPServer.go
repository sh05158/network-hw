/**
 * 20176342 Song Min Joon
 * EasyUDPServer.go
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
	conn, _ := net.ListenPacket("udp", ":"+serverPort)

	fmt.Printf("Server is ready to receive on port %s\n", serverPort)

	for {
		buffer := make([]byte, 1024)
		_, r_addr, _ := conn.ReadFrom(buffer)
		fmt.Printf("UDP message from %s\n", r_addr.String())

		//if server receives packet from path which connects client and server
		//pass buffer(packet) to handleMsg func to process message
		handleMsg(conn, r_addr, string(buffer))
	}

	defer byebye()

}

func byebye() {
	fmt.Printf("Bye bye~\n")
	os.Exit(0)
}

func handleMsg(conn net.PacketConn, addr net.Addr, clientMsg string) {

	totalRequests++
	tempStr := clientMsg

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

	requestOption, _ := strconv.Atoi(strings.Split(tempStr, "|")[0]) // split client packet by '|' and takes option and convert to Integer.
	requestData := strings.Split(tempStr, "|")[1]                    // get message parameter from packet.
	time.Sleep(time.Millisecond * 1)                                 // minimum delay to deliver packet to client.

	fmt.Printf("Command %d\n\n", requestOption) // print Command #
	switch requestOption {
	case 1:
		//Option 1
		upperString := strings.ToUpper(requestData)
		sendPacket(conn, addr, upperString)
	case 2:
		//Option 2
		sendPacket(conn, addr, addr.String())
	case 3:
		//Option 3
		sendPacket(conn, addr, strconv.Itoa(totalRequests))
	case 4:
		//Option 4
		elapsed := time.Since(startTime).Truncate(time.Second).String()
		sendPacket(conn, addr, string(elapsed))
	case 5:
		//Option 5
		conn.Close()
	default:
		//Option default
		conn.Close()
	}

}

func sendPacket(conn net.PacketConn, addr net.Addr, serverMsg string) {
	//send packet to client
	conn.WriteTo([]byte(serverMsg), addr)
}
