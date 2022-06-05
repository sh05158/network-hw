/**
 * 20176342 Song Min Joon
 * P2POmokServer.go
 **/
package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var serverPort string = "56342" // for server port
var uniqueID int = 1
var clientMap map[int]client

type client struct {
	nickname string
	uniqueID int
	conn     net.Conn
	ip       string
	port     string
	remote   string
}

func main() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // for exit the program gracefully
	go func() {
		for sig := range c {
			// sig is a ^C, handle it
			_ = sig
			byebye() // if this program is interrupt by Ctrl-c, print Bye bye and exits gracefully
		}
	}()

	listener, _ := net.Listen("tcp", ":"+serverPort)
	fmt.Printf("Server is ready to receive on port %s\n", serverPort)

	clientMap = make(map[int]client) //make client map for record clients

	for {
		//listener is waiting for tcp connection of clients.
		conn, err := listener.Accept()
		if err != nil {
			handleError(conn, err, "server accept error..")
		}


		nickBuffer := make([]byte, 1024)

		count, err := conn.Read(nickBuffer)

		if err != nil || count == 0 {
			continue
			return
		}

		targetNick := string(nickBuffer[:count])



		remoteAddr := conn.RemoteAddr().String()
		lastIdx := strings.LastIndex(remoteAddr, ":")

		newIP := remoteAddr[0:lastIdx]
		newPort := remoteAddr[lastIdx+1:]

		fmt.Printf("%s joined from from %s. UDP port %s\n", targetNick,conn.RemoteAddr().String(),newPort)

		newClient := client{targetNick, uniqueID, conn, newIP, newPort, remoteAddr}

		registerClient(newClient, uniqueID)

		go handleMsg(newClient, uniqueID) // when client is connect to server, make go-routine to communicate with client.
		uniqueID++

		if len(clientMap) == 1 {
			sendPacket(newClient, "waiting||")
			fmt.Printf("1 user connected, waiting for another\n")
		} else if len(clientMap) == 2 {
			var otherRemote client
			for _, v := range clientMap {
				if v.uniqueID != newClient.uniqueID {
					otherRemote = v
					sendPacket(v, "matched|"+newClient.nickname+"|"+newClient.remote)
				}
			}
			sendPacket(newClient, "success|"+otherRemote.nickname+"|"+otherRemote.remote)

			fmt.Printf("2 users connected, notifying %s and %s\n",otherRemote.nickname,newClient.nickname )


			for _, v := range clientMap {
				unregisterClient(v.uniqueID)
			}
			fmt.Printf("%s and %s disconnected.\n",otherRemote.nickname,newClient.nickname )

		}

	}

	defer byebye() // although when client gets panic, defer should disconnect socket gracefully

}

func handleMsg(client client, cid int) {
	for {
		buffer := make([]byte, 1024)

		count, err := client.conn.Read(buffer)

		//when client sends packet

		if err != nil {
			unregisterClient(cid)
			// broadCastToAll(6, fmt.Sprintf("%s is disconnected. There are %d users in the chat room", client.nickname, len(clientMap)))
			if len(clientMap) == 1{
				fmt.Printf("%s disconnected. 1 User left in server.\n",client.nickname)
			} else if len(clientMap) == 0{
				fmt.Printf("%s disconnected. No User left in server.\n",client.nickname)
			}
			handleError(client.conn, err, "%s disconnected!")
			return
		}

		if count == 0 {
			// fmt.Printf("return ! \n")
			return
		}
	}
}
func registerClient(client client, uniqueID int) int {
	clientMap[uniqueID] = client
	return uniqueID
}

func unregisterClient(uniqueID int) {
	if _, ok := clientMap[uniqueID]; ok {
		delete(clientMap, uniqueID)
	}
}

func byebye() {
	fmt.Printf("Bye bye~\n")
	os.Exit(0)
}

func sendPacket(client client, serverMsg string) {
	//send packet to client
	time.Sleep(time.Millisecond * 1) // minimum delay to deliver packet to client.
	client.conn.Write([]byte(serverMsg))
	fmt.Printf("send packet %s => %s \n", client.nickname, serverMsg)
}

func handleError(conn net.Conn, err error, errmsg string) {
	//handle error and print
	if conn != nil {
		conn.Close()
	}
}
