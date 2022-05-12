/**
 * 20176342 Song Min Joon
 * EasyTCPClient.go
 **/

package main

import (
	// "bufio"
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

	if len(os.Args) < 2 {
		fmt.Printf("Please check your nickname argument\n")
		byebye()
		return
	}

	nickname := os.Args[1:2][0]

	fmt.Printf("%s\n", nickname)
	if len(nickname) <= 0 {
		fmt.Printf("Please check your nickname argument\n")
		byebye()
		return
	}

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

	sendPacket(conn, nickname) // send nickname to server and wait for response
	buffer := make([]byte, 1024)
	bufferSize := readPacket(conn, &buffer)
	response := string(buffer[:bufferSize]) //wait for nickname response

	if response == "duplicated" {
		fmt.Printf("your nickname %s is duplicated. please use another nickname\n")
		byebye()
		conn.Close()
		os.Exit(0)
	}

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
		fmt.Printf("handle Packet\n")
		fmt.Printf("server send : %s\n", response)

		route, _ := strconv.Atoi(strings.Split(response, "|")[0])

		msgArr := strings.Split(response, "|")

		/*
			server packet form

			(Route|Option Number|Message) client request packet
			(Route|Message) server one sided packet

			Route 0 => normal message 0|nickname|message
			Route 1 => list   just print server message 1|message
			Route 2 => dm     message 2|nickname|message
			Route 3 => disconnect  message 3|
			Route 4 => show version 4|message(Version)
			Route 5 => show rtt 5|message(rtt)
			Route 6 => show message 6|message(disconnected)

		*/

		switch route {
		case 0: // route 0 is normal message packet
			// fromNickname := msgArr[1]
			// fromMessage := msgArr[2]

			fmt.Printf("%s\n", msgArr[1])
			break
		case 2:
			//print dm
			fromNickname := msgArr[1]
			fromMessage := msgArr[2]

			fmt.Printf("from: %s>%s\n", fromNickname, fromMessage)
			break
		case 3:
			//disconnect, do nothing
			break
		case 4:
			//print version string
			fmt.Printf("%s\n", msgArr[1])
			break
		case 5:
			//print RTT
			printRTT(time.Since(lastRequestTime))
			break
		case 1: // route 1 is packet that is not related to client request (server one-sided packet)
			//print list of users< nickname, IP, port >
			fmt.Printf("%s\n", msgArr[1])
			break
		case 6:
			//server another user disconnect message
			fmt.Printf("%s\n", msgArr[1])
			break

		}
	}

}

// func processMyMessage(opt int, msg string) {
// 	//if i request server with option n, then process received packet with n

// 	switch opt {
// 	case 1:
// 		//if i request option number 1
// 		fmt.Printf("Reply from server: %s\n", msg)
// 		break
// 	case 2:
// 		//if i request option number 2
// 		fmt.Printf("Reply from server: client IP = %s, port = %s\n", string(strings.Split(msg, ":")[0]), string(strings.Split(msg, ":")[1]))
// 		break

// 	case 3:
// 		//if i request option number 3
// 		fmt.Printf("Reply from server: requests served = %s\n", msg) //print server message directly
// 		break

// 	case 4:
// 		//if i request option number 4
// 		timeD, _ := time.ParseDuration(msg)
// 		printDuration(timeD) // print server running time
// 		break

// 	}

// 	printRTT(time.Since(lastRequestTime))
// }

func handleInput(conn net.Conn) {
	for {
		time.Sleep(time.Millisecond * 100)
		// var inputstr string

		inputstr, _ := bufio.NewReader(os.Stdin).ReadString('\n')

		inputstr = inputstr[:len(inputstr)-1]

		processMyMessage(inputstr, conn)
	}

}
func processCommandOption(command string, arguments string, conn net.Conn) {
	fmt.Printf("processCommmandOption 1%s1 %s \n", command, arguments)
	requestString := "2|"
	if command == "ver" {
		fmt.Printf("command\n")

	}
	switch command {
	case "list":
		//if user command is list
		requestString += "1|"
		sendPacket(conn, requestString)

	case "dm":
		//if user command is dm
		toNickname := strings.Split(arguments, " ")[0]
		toMessage := strings.Split(arguments, " ")[1]
		requestString += "2|" + toNickname + "|" + toMessage
		sendPacket(conn, requestString)

	case "exit":
		//if user command is exit
		requestString += "3|"
		sendPacket(conn, requestString)

	case "ver":
		//if user command is ver
		fmt.Printf("ver \n")
		requestString += "4|"
		sendPacket(conn, requestString)

	case "rtt":
		//if user command is rtt
		lastRequestTime = time.Now() //startTime for print RTT
		requestString += "5|"
		sendPacket(conn, requestString)
	}
}
func processMyMessage(inputstr string, conn net.Conn) {
	/*
		if option is given, it sends packet to server and get response.
	*/

	// var temp int
	// fmt.Scanf("%s", &temp)

	lastRequestTime = time.Now() //startTime for print RTT

	if inputstr[:1] == "\\" {
		//if user input string is command
		command := ""
		arguments := ""

		if strings.Contains(inputstr, " ") == true {
			command = strings.Split(strings.Split(inputstr, " ")[0], "\\")[1]
			arguments = strings.Split(inputstr, " ")[1]

		} else {
			//if no space
			command = strings.Split(inputstr, "\\")[1]
		}

		processCommandOption(command, arguments, conn)

	} else {
		//is not command, just send normal message
		requestString := "1|" + inputstr
		sendPacket(conn, requestString)
		//send packet

	}

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
