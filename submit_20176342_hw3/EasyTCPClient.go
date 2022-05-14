/**
 * 20176342 Song Min Joon
 * EasyTCPClient.go
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

var serverName string = "localhost" //server host
var serverPort string = "26342"     //server port

var lastRequestTime time.Time

func main() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // for exit the program gracefully
	go func() {
		for sig := range c {
			// sig is a ^C, handle it
			_ = sig
			byebye() // print byebye func
		}
	}()

	conn, err := net.Dial("tcp", serverName+":"+serverPort) //tcp connection

	if err != nil {
		//if server is not working, print and exit
		fmt.Printf("Please check your server is running\n")
		byebye()
		return
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.TCPAddr) //get local port

	fmt.Printf("Client is running on port %d\n", localAddr.Port)

	go handlePacket(conn)
	handleInput(conn)

	// select {}

	defer conn.Close() // although when client gets panic, defer should disconnect socket gracefully
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

		fmt.Printf("response : %s \n", response)

		route, _ := strconv.Atoi(strings.Split(response, "|")[0])

		msgArr := strings.Split(response, "|")


		/*
			server packet form

			(Route|Option Number|Message) client request packet
			(Route|Message) server one sided packet

			0|1|BLAH BLAH BLAH...
			0|1|HELLO WORLD!

			0|2|127.0.0.1:6342    //my ip and port
			0|3|30 				  //count of server serves client requests
			0|4|01h30m12s		  //server running time

			1|client #3 is connected   //server broadcast message (server one-sided)
			1|client #2 is disconnected   //server broadcast message (server one-sided)

		*/

		switch route {
		case 0: // route 0 is packet that related to client request
			opt, _ := strconv.Atoi(msgArr[1])

			processMyMessage(opt, msgArr[2])
			break
		case 1: // route 1 is packet that is not related to client request (server one-sided packet)

			fmt.Printf("%s\n", msgArr[1])
			break

		}
	}

}

func processMyMessage(opt int, msg string) {
	//if i request server with option n, then process received packet with n

	switch opt {
	case 1:
		//if i request option number 1
		fmt.Printf("Reply from server: %s\n", msg)
		break
	case 2:
		//if i request option number 2
		fmt.Printf("Reply from server: client IP = %s, port = %s\n", string(strings.Split(msg, ":")[0]), string(strings.Split(msg, ":")[1]))
		break

	case 3:
		//if i request option number 3
		fmt.Printf("Reply from server: requests served = %s\n", msg) //print server message directly
		break

	case 4:
		//if i request option number 4
		timeD, _ := time.ParseDuration(msg)
		printDuration(timeD) // print server running time
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
	/*
		if option is given, it sends packet to server and get response.
	*/
	lastRequestTime = time.Now() //startTime for print RTT

	var temp int
	fmt.Scanf("%s", &temp)

	switch opt {
	case 1:
		//if my option is 1
		fmt.Printf("Input lowercase sentence: ")
		input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		lastRequestTime = time.Now()
		requestString := strconv.Itoa(opt) + "|" + input
		sendPacket(conn, requestString)
		//send packet
	case 2:
		// Option 2
		requestString := strconv.Itoa(opt) + "|"
		sendPacket(conn, requestString)

	case 3:
		// Option 3
		requestString := strconv.Itoa(opt) + "|"
		sendPacket(conn, requestString)

	case 4:
		// Option 4
		requestString := strconv.Itoa(opt) + "|"
		sendPacket(conn, requestString)

	case 5:
		// Option 5
		byebye()
		conn.Close()
		os.Exit(0)
	default:
		// not Option 1~5 (default)
		byebye()
		conn.Close()
		os.Exit(0)
	}
	// printRTT(time.Since(lastRequestTime))

}

func printRTT(d time.Duration) {
	d = d.Round(time.Millisecond)

	s := d / time.Millisecond

	fmt.Printf("RTT = %dms \n\n\n", s) // print RTT Since startTime

}
func printDuration(d time.Duration) {
	//print server running time in proper form(HH:MM:ss)

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
	//send Packet to server
	conn.Write([]byte(requestString))
}

func readPacket(conn net.Conn, buffer *[]byte) int {
	//read Packet from server and saves to buffer and return buffer size.
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
	//print Menu and 5 Options.
	fmt.Printf("<Menu>\n")
	fmt.Printf("option 1) convert text to UPPER-case letters.\n")
	fmt.Printf("option 2) ask the server what the IP address and port number of the client is.\n")
	fmt.Printf("option 3) ask the server how many client requests(commands) it has served so far.\n")
	fmt.Printf("option 4) ask the server program how long it has been running for since it started.\n")
	fmt.Printf("option 5) exit client program\n")

}
