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
var clientMap map[int]client

type client struct{
	nickname string
	uniqueID int
	conn net.Conn
	ip string
	port string
}
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

	clientMap = make(map[int]client) //make client map for record clients 

	for {
		//listener is waiting for tcp connection of clients.
		conn, err := listener.Accept()
		if err != nil {
			handleError(conn, err, "server accept error..")
		}

		fmt.Printf("Connection request from %s\n", conn.RemoteAddr().String())

		nickBuffer := make([]byte, 1024)

		count, err := conn.Read(nickBuffer)

		if err != nil || count == 0 {
			continue;
			return
		}

		targetNick := string(nickBuffer)

		_, isDuplicate := getClientByNickname(targetNick)

		if isDuplicate {
			sendPacket(conn, "duplicated")
			conn.Close()
			continue
		}

		sendPacket(conn, "duplicated")


		newIP := strings.Split(conn.RemoteAddr().String(), ":")[0]
		newPort := strings.Split(conn.RemoteAddr().String(), ":")[1]

		newClient := client{targetNick, uniqueID, conn, newIP, newPort }

		registerClient(newClient, uniqueID)
		
		broadCastToAll(1, fmt.Sprintf("Client %d connected. Number of connected clients = %d", uniqueID, len(clientMap)))
		go handleMsg(newClient, uniqueID) // when client is connect to server, make go-routine to communicate with client.
		uniqueID++

	}

	defer byebye() // although when client gets panic, defer should disconnect socket gracefully

}

func getClientByNickname(nick string)( client, bool) {
	for _, v := range clientMap {
		if v.nickname == nick{
			return v, true
		}
	}
	return client{}, false
}

func printClientNum() {
	for {
		time.Sleep(time.Minute * 1)
		fmt.Printf("Number of connected clients = %d\n", len(clientMap))
	}
}

func broadCastToAll(route int, msg string) {
	for _, v := range clientMap {
		sendPacket(v.conn, strconv.Itoa(route)+"|"+msg)
	}
	fmt.Println(msg)
}

func broadCastExceptMe(route int, msg string, client client) {
	for _, v := range clientMap {
		if client.uniqueID != v.uniqueID{
			//do not send to myself
			sendPacket(v.conn, strconv.Itoa(route)+"|"+msg)
		}
	}
	fmt.Println(msg)
}


func registerClient(client client, uniqueID int) int {
	clientMap[uniqueID] = client
	return uniqueID
}

func unregisterClient(uniqueID int) {
	if _, ok := clientMap[uniqueID];ok {
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

func handleMsg(client client, cid int) {
	for {
		buffer := make([]byte, 1024)

		count, err := client.conn.Read(buffer)
		// fmt.Printf("count = %d\n", count)

		//when client sends packet

		if err != nil {
			unregisterClient(cid)
			broadCastToAll(6, fmt.Sprintf("Client %d disconnected. Number of connected clients = %d", cid, len(clientMap)))
			handleError(client.conn, err, "client disconnected!")
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

			(isCommand | Command Option)

			isCommand => 1   not command( normal chat )
			isCommand => 2    is command


			Command Option => 1 (list) show the nickname IP Port of all connected users
			Command Option => 2 (dm)  dm destination message
			Command Option => 3 (exit) disconnect
			Command Option => 4 (ver) show version
			Command Option => 5 (rtt) show rtt

			2|2|hulk hello there



		*/

		tempStr := string(buffer)

		if strings.Contains(strings.ToUpper(tempStr), "I HATE PROFESSOR") {
			//if client message includes 'i hate professor' disconnect socket 
			client.conn.Close()
			unregisterClient(client.uniqueID)
			broadCastToAll(6, fmt.Sprintf("Client %d disconnected. Number of connected clients = %d", cid, len(clientMap)))
			sendPacket(client.conn, "5|")
		}

		// fmt.Printf("client msg %s\n", tempStr)

		requestOption, _ := strconv.Atoi(strings.Split(tempStr, "|")[0])// split client packet by '|' and takes option and convert to Integer.

		// requestData := strings.Split(tempStr, "|")[1]// get message parameter from packet.
		time.Sleep(time.Millisecond * 1)// minimum delay to deliver packet to client.
		// fmt.Printf("Command %d\n\n", requestOption)// print Command #

		switch requestOption {
		case 0:
			//  is not command (normal message)
			message := strings.Split(tempStr, "|")[1]

			broadCastExceptMe(0, client.nickname+"> "+message, client)
		case 1:
			//Option 2
			commandOption, _ := strconv.Atoi(strings.Split(tempStr, "|")[1])

			switch commandOption {
			case 1:
				sendPacket(client.conn, "1|"+getClientListString())
			case 2:
				dmOption := strings.Split(tempStr, "|")[2]
				dmTarget := strings.Split(dmOption, " ")[0]
				dmMessage := strings.Split(dmOption, " ")[1]

				targetClient, success := getClientByNickname(dmTarget)

				if success {
					sendPacket(targetClient.conn, "2|"+client.nickname+"|"+dmMessage)
				

				}

			case 3:
				sendPacket(client.conn, "3|")
			case 4:
				sendPacket(client.conn, "4|Chat TCP Version 0.1")
			case 5:
				sendPacket(client.conn, "5|")
			}
		default:
			//Option default
			// do nothing
			// conn.Close()
		}
	}

}

func getClientListString() string{
	returnStr := ""

	for _, v := range clientMap {
		returnStr += "\n<"+v.nickname+", "+v.ip+", "+v.port+">"
	}

	return returnStr
}

func sendPacket(conn net.Conn, serverMsg string) {
	//send packet to client

	conn.Write([]byte(serverMsg))
}
