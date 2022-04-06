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

var serverName string = "nsl2.cau.ac.kr"
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

func processOption(opt int, conn net.Conn) {

	startTime := time.Now()

	var temp int
	fmt.Scanf("%s", &temp)

	switch opt {
	case 1:
		fmt.Printf("Input lowercase sentence: ")
		input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		startTime = time.Now()
		requestString := strconv.Itoa(opt) + "|" + input
		sendPacket(conn, requestString)
		buffer := make([]byte, 1024)
		bufferSize := readPacket(conn, &buffer)
		fmt.Printf("Reply from server: %s\n", string(buffer[:bufferSize]))

	case 2:
		requestString := strconv.Itoa(opt) + "|"
		sendPacket(conn, requestString)
		buffer := make([]byte, 1024)
		bufferSize := readPacket(conn, &buffer)
		fmt.Printf("Reply from server: client IP = %s, port = %s\n", string(strings.Split(string(buffer[:bufferSize]), ":")[0]), string(strings.Split(string(buffer[:bufferSize]), ":")[1]))

	case 3:
		requestString := strconv.Itoa(opt) + "|"
		sendPacket(conn, requestString)
		buffer := make([]byte, 1024)
		bufferSize := readPacket(conn, &buffer)
		fmt.Printf("Reply from server: requests served = %s\n", string(buffer[:bufferSize]))

	case 4:
		requestString := strconv.Itoa(opt) + "|"
		sendPacket(conn, requestString)
		buffer := make([]byte, 1024)
		bufferSize := readPacket(conn, &buffer)
		timeD, _ := time.ParseDuration(string(buffer[:bufferSize]))

		printDuration(timeD)

	case 5:
		byebye()
		conn.Close()
		os.Exit(0)
	default:
		byebye()
		conn.Close()
		os.Exit(0)
	}
	printRTT(time.Since(startTime))

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
