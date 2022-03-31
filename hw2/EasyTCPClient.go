/**
 * TCPClient.go
 **/

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var serverName string = "localhost"
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

	conn, err := net.Dial("tcp", serverName+":"+serverPort)

	if err != nil {
		fmt.Printf("Please check your server is running\n")
		byebye()
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

func byebye() {
	fmt.Printf("Bye bye~\n")
	os.Exit(0)
}

func handleInput(conn net.Conn) {
	printOption()
	fmt.Printf("Please select your option :")
	var opt int
	fmt.Scanf("%d", &opt)
	processOption(opt, conn)
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

	startTime := time.Now()

	// var temp int
	// fmt.Scanf("%s", &temp)

	switch opt {
	case 1:
		fmt.Printf("Input lowercase sentence: ")
		input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		startTime = time.Now()
		requestString := strconv.Itoa(opt) + "|" + input
		sendPacket(conn, requestString)
		buffer := make([]byte, 1024)
		readPacket(conn, &buffer)
		fmt.Printf("Reply from server: %s\n", string(buffer))

	case 2:
		requestString := strconv.Itoa(opt) + "|"
		sendPacket(conn, requestString)
		buffer := make([]byte, 1024)
		readPacket(conn, &buffer)
		fmt.Printf("Reply from server: client IP = %s, port = %s\n", string(strings.Split(string(buffer), ":")[0]), string(strings.Split(string(buffer), ":")[1]))

	case 3:
		requestString := strconv.Itoa(opt) + "|"
		sendPacket(conn, requestString)
		buffer := make([]byte, 1024)
		readPacket(conn, &buffer)
		fmt.Printf("Total client request count = %s\n", string(buffer))

	case 4:
		requestString := strconv.Itoa(opt) + "|"
		sendPacket(conn, requestString)
		buffer := make([]byte, 1024)
		readPacket(conn, &buffer)
		fmt.Printf("Server started %s seconds ago\n", string(buffer))

	case 5:
		byebye()
		conn.Close()
		os.Exit(0)
	default:
		byebye()
		conn.Close()
		os.Exit(0)
	}
	elapsed := time.Since(startTime).Truncate(time.Millisecond).String()
	fmt.Printf("RTT = %s \n", elapsed)

}

func sendPacket(conn net.Conn, requestString string) {
	conn.Write([]byte(requestString))
}

func readPacket(conn net.Conn, buffer *[]byte) {
	_, err := conn.Read(*buffer)
	if err != nil {
		fmt.Println("connection is closed by server")
		byebye()
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
