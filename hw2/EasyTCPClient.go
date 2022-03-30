/**
 * TCPClient.go
 **/

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type Packet struct {
	Option int `json:"option"`

	Data interface{} `json:"data"`
}

func main() {

	serverName := "nsl2.cau.ac.kr"
	serverPort := "26342"

	conn, err := net.Dial("tcp", serverName+":"+serverPort)

	if err != nil {
		fmt.Printf("Please check your server is running\n")
		return
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.TCPAddr)

	fmt.Printf("Client is running on port %d\n", localAddr.Port)

	for {
		handleInput(conn)
	}

	// defer conn.Close()
}

func handleInput(conn net.Conn) {
	printOption()
	fmt.Printf("Please select your option :")
	// opt, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	var opt int
	fmt.Scanf("%d", &opt)

	processOption(opt, conn)
}

func handleServerConnection(c net.Conn) {

	// we create a decoder that reads directly from the socket
	d := json.NewDecoder(c)

	var msg Packet

	err := d.Decode(&msg)
	fmt.Println(msg, err)

	c.Close()

}

func handleError(conn net.Conn, errmsg string) {
	if conn != nil {
		conn.Close()
	}
	fmt.Println(errmsg)
}

func handleSendMsg(conn net.Conn) {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Text to send : ")
		text, err := reader.ReadString('\n')
		if err != nil {
			handleError(conn, "read input failed..")
		}

		fmt.Fprintf(conn, "%s|%s", text)
	}
}

func handleRecvMsg(conn net.Conn, msgch chan string) {
	for {
		select {
		case msg := <-msgch:
			fmt.Printf("\nMessage from server : %s\n", msg)
		default:
			go recvFromServer(conn, msgch)
			time.Sleep(1000 * time.Millisecond)
		}
	}
}

func recvFromServer(conn net.Conn, msgch chan string) {
	msg, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		handleError(conn, "read msg failed..")
		os.Exit(2)
		return
	}
	msgch <- msg
}

func processOption(opt int, conn net.Conn) {

	var temp int
	fmt.Scanf("%s", &temp)

	switch opt {
	case 1:
		fmt.Printf("Input lowercase sentence: ")

		input, _ := bufio.NewReader(os.Stdin).ReadString('\n')

		requestString := strconv.Itoa(opt) + "|" + input

		startTime := time.Now()
		conn.Write([]byte(requestString))

		buffer := make([]byte, 1024)
		conn.Read(buffer)
		elapsed := time.Since(startTime).Truncate(time.Millisecond).String()
		fmt.Printf("Reply from server: %s\n", string(buffer))
		fmt.Printf("RTT = %s\n", elapsed)

	case 2:
		requestString := strconv.Itoa(opt) + "|"
		startTime := time.Now()

		conn.Write([]byte(requestString))
		buffer := make([]byte, 1024)
		conn.Read(buffer)
		elapsed := time.Since(startTime).Truncate(time.Millisecond).String()

		fmt.Printf("Reply from server: client IP = %s, port = %s\n", string(strings.Split(string(buffer), ":")[0]), string(strings.Split(string(buffer), ":")[1]))
		fmt.Printf("RTT = %s\n", elapsed)

	case 3:
		requestString := strconv.Itoa(opt) + "|"
		startTime := time.Now()

		conn.Write([]byte(requestString))

		buffer := make([]byte, 1024)
		conn.Read(buffer)

		elapsed := time.Since(startTime).Truncate(time.Millisecond).String()

		fmt.Printf("Total client request count = %s\n", string(buffer))
		fmt.Printf("RTT = %s\n", elapsed)

	case 4:
		requestString := strconv.Itoa(opt) + "|"
		startTime := time.Now()

		conn.Write([]byte(requestString))
		buffer := make([]byte, 1024)
		conn.Read(buffer)
		elapsed := time.Since(startTime).Truncate(time.Millisecond).String()

		fmt.Printf("Server started %s seconds ago\n", string(buffer))
		fmt.Printf("RTT = %s\n", elapsed)

	case 5:
		conn.Close()
		os.Exit(0)
	default:
		conn.Close()
		os.Exit(0)
	}

}

func printOption() {
	fmt.Printf("option 1) convert text to UPPER-case letters.\n")
	fmt.Printf("option 2) ask the server what the IP address and port number of the client is.\n")
	fmt.Printf("option 3) ask the server how many client requests(commands) it has served so far.\n")
	fmt.Printf("option 4) ask the server program how long it has been running for since it started.\n")
	fmt.Printf("option 5) exit client program\n")

}
