/**
 * 20176342 Song Min Joon
 * MultiClientTCPServer.go
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

var totalRequests int = 0		// total Request count global variable for server.
var startTime time.Time			 // for saving server start time
var serverPort string = "26342"  // for server port
var uniqueID int = 1
var clientMap map[int]net.Conn

func main() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // for exit the program gracefully
	go func() {
		for sig := range c {
			// sig is a ^C, handle it
			_ = sig
			byebye()// if this program is interrupt by Ctrl-c, print Bye bye and exits gracefully
		}
	}()

	startTime = time.Now()// records server start time for server running time
	listener, _ := net.Listen("tcp", ":"+serverPort)
	fmt.Printf("Server is ready to receive on port %s\n", serverPort)

	clientMap = make(map[int]net.Conn) //make client map for record clients 
	go printClientNum()

	for {
		//listener is waiting for tcp connection of clients.
		conn, err := listener.Accept()
		if err != nil {
			handleError(conn, err, "server accept error..")
		}

		fmt.Printf("Connection request from %s\n", conn.RemoteAddr().String())

		registerClient(conn, uniqueID)
		broadCastToAll(1, fmt.Sprintf("Client %d connected. Number of connected clients = %d", uniqueID, len(clientMap)))
		go handleMsg(conn, uniqueID) // when client is connect to server, make go-routine to communicate with client.
		uniqueID++

	}

	defer byebye() // although when client gets panic, defer should disconnect socket gracefully

}

func printClientNum() {
	for {
		time.Sleep(time.Minute * 1)
		fmt.Printf("Number of connected clients = %d\n", len(clientMap))
	}
}

func broadCastToAll(route int, msg string) {
	for _, v := range clientMap {
		sendPacket(v, strconv.Itoa(route)+"|"+msg)
	}
	fmt.Println(msg)
}

func registerClient(conn net.Conn, uniqueID int) int {
	clientMap[uniqueID] = conn
	return uniqueID
}

func unregisterClient(uniqueID int) {
	if clientMap[uniqueID] != nil {
		delete(clientMap, uniqueID)
	}
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
	// fmt.Println(err)
	// fmt.Println(errmsg)
}
func handleError2(conn net.Conn, errmsg string) {
		//handle error and print
	if conn != nil {
		conn.Close()
	}

	// fmt.Println(errmsg)
}

func handleMsg(conn net.Conn, cid int) {
	for {
		buffer := make([]byte, 1024)

		count, err := conn.Read(buffer)
		// fmt.Printf("count = %d\n", count)

		//when client sends packet

		if err != nil {
			unregisterClient(cid)
			broadCastToAll(1, fmt.Sprintf("Client %d disconnected. Number of connected clients = %d", cid, len(clientMap)))
			handleError(conn, err, "client disconnected!")
			return
		}

		if count == 0 {
			// fmt.Printf("return ! \n")
			return
		}

		_ = count

		totalRequests++

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

		tempStr := string(buffer[:count])

		// fmt.Printf("client msg %s\n", tempStr)

		requestOption, _ := strconv.Atoi(strings.Split(tempStr, "|")[0])// split client packet by '|' and takes option and convert to Integer.
		requestData := strings.Split(tempStr, "|")[1]// get message parameter from packet.
		time.Sleep(time.Millisecond * 1)// minimum delay to deliver packet to client.
		fmt.Printf("Command %d\n\n", requestOption)// print Command #

		switch requestOption {
		case 1:
			//Option 1
			upperString := strings.ToUpper(requestData)
			sendPacket(conn, "0|1|"+upperString)
		case 2:
			//Option 2

			ip := conn.RemoteAddr()
			sendPacket(conn, "0|2|"+ip.String())
		case 3:
			//Option 3

			sendPacket(conn, "0|3|"+strconv.Itoa(totalRequests))
		case 4:
			//Option 4

			elapsed := time.Since(startTime).Truncate(time.Second).String()

			sendPacket(conn, "0|4|"+string(elapsed))
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
