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

var lastRequestTime time.Time

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

	go handlePacket(conn)
	handleInput(conn)

	// select {}

	defer conn.Close()
}

func byebye() {
	fmt.Printf("Bye bye~\n")
	os.Exit(0)
}

func handlePacket(conn net.Conn) {
	for {
		buffer := make([]byte, 1024)
		bufferSize := readPacket(conn, &buffer)

		response := string(buffer[:bufferSize])

		// fmt.Println(response)

		route, _ := strconv.Atoi(strings.Split(response, "|")[0])

		msgArr := strings.Split(response, "|")

		switch route {
		case 0: //내꺼
			opt, _ := strconv.Atoi(msgArr[1])

			processMyMessage(opt, msgArr[2])
			break
		case 1: //서버에서 일방적으로 보내준 패킷

			fmt.Printf("%s\n", msgArr[1])
			break

		}
	}

}

func processMyMessage(opt int, msg string) {
	switch opt {
	case 1:

		fmt.Printf("Reply from server: %s\n", msg)
		break
	case 2:

		fmt.Printf("Reply from server: client IP = %s, port = %s\n", string(strings.Split(msg, ":")[0]), string(strings.Split(msg, ":")[1]))
		break

	case 3:

		fmt.Printf("Reply from server: requests served = %s\n", msg)
		break

	case 4:
		timeD, _ := time.ParseDuration(msg)
		printDuration(timeD)
		break

	}

	printRTT(time.Since(lastRequestTime))
}

func handleInput(conn net.Conn) {
	for {
		time.Sleep(time.Millisecond * 100)
		printOption()
		fmt.Printf("Please select your option :")
		var opt int
		fmt.Scanf("%d", &opt)
		processOption(opt, conn)
	}

}

func processOption(opt int, conn net.Conn) {

	lastRequestTime = time.Now()

	var temp int
	fmt.Scanf("%s", &temp)

	switch opt {
	case 1:
		fmt.Printf("Input lowercase sentence: ")
		input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		lastRequestTime = time.Now()
		requestString := strconv.Itoa(opt) + "|" + input
		sendPacket(conn, requestString)

	case 2:
		requestString := strconv.Itoa(opt) + "|"
		sendPacket(conn, requestString)

	case 3:
		requestString := strconv.Itoa(opt) + "|"
		sendPacket(conn, requestString)

	case 4:
		requestString := strconv.Itoa(opt) + "|"
		sendPacket(conn, requestString)

	case 5:
		byebye()
		conn.Close()
		os.Exit(0)
	default:
		byebye()
		conn.Close()
		os.Exit(0)
	}
	// printRTT(time.Since(lastRequestTime))

}

func printRTT(d time.Duration) {
	d = d.Round(time.Millisecond)

	s := d / time.Millisecond

	fmt.Printf("RTT = %dms \n\n\n", s)

}
func printDuration(d time.Duration) {
	d = d.Round(time.Second)

	h := d / time.Hour
	d -= h * time.Hour

	m := d / time.Minute
	d -= m * time.Minute

	s := d / time.Second
	d -= s * time.Second

	fmt.Printf("Reply from server: run time = %02d:%02d:%02d\n", h, m, s)
}

func sendPacket(conn net.Conn, requestString string) {
	conn.Write([]byte(requestString))
}

func readPacket(conn net.Conn, buffer *[]byte) int {
	count, err := conn.Read(*buffer)
	if err != nil {
		fmt.Println("connection is closed by server")
		byebye()
		conn.Close()
		os.Exit(0)
	}
	return count
}

func printOption() {
	fmt.Printf("<Menu>\n")
	fmt.Printf("option 1) convert text to UPPER-case letters.\n")
	fmt.Printf("option 2) ask the server what the IP address and port number of the client is.\n")
	fmt.Printf("option 3) ask the server how many client requests(commands) it has served so far.\n")
	fmt.Printf("option 4) ask the server program how long it has been running for since it started.\n")
	fmt.Printf("option 5) exit client program\n")

}
