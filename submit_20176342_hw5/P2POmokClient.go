/**
 * 20176342 Song Min Joon
 * P2POmokClient.go
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

var serverName string = "nsl2.cau.ac.kr" //server host
var serverPort string = "56342"          //server port

var myTurn bool
var isFirst bool
var omok [10][10]int
var isDone bool

var addr *net.UDPAddr

type client struct {
	nickname string
	remote   string
}

var opponent client

var timer chan bool

func main() {

	var conn net.PacketConn
	var tcpconn net.Conn

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // for exit the program gracefully
	go func() {
		for sig := range c {
			// sig is a ^C, handle it

			if tcpconn != nil {
				tcpconn.Close()
			}

			if conn != nil {
				sendPacket(conn, "3|", addr)
				conn.Close()
			}

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

	tcpconn, err := net.Dial("tcp", serverName+":"+serverPort) //tcp connection

	if err != nil {
		//if server is not working, print and exit
		fmt.Printf("Please check your server is running\n")
		byebye()
		return
	}
	defer tcpconn.Close()

	localAddr := tcpconn.LocalAddr().(*net.TCPAddr) //get local port

	fmt.Printf("Client is running on port %d\n", localAddr.Port)
	fmt.Printf("Client is connected to %s\n", tcpconn.RemoteAddr().String())

	sendPacketTCP(tcpconn, nickname) // send nickname to server and wait for response
	buffer := make([]byte, 1024)
	bufferSize := readPacketTCP(tcpconn, &buffer)
	response := string(buffer[:bufferSize]) //wait for nickname response

	responseArr := strings.Split(response, "|")

	fmt.Printf("welcome %s to p2p-omok server at %s. \n", nickname, tcpconn.RemoteAddr().String())

	var targetServer string

	if responseArr[0] == "waiting" {
		//first client

		fmt.Printf("waiting for an opponent. \n")

		buffer := make([]byte, 1024)
		bufferSize := readPacketTCP(tcpconn, &buffer)
		response := string(buffer[:bufferSize]) //wait for nickname response
		responseArr := strings.Split(response, "|")

		if responseArr[0] == "matched" {
			//first client matched
			fmt.Printf("%s joined (%s). you play first. \n", responseArr[1], responseArr[2])
			targetServer = responseArr[2]
			myTurn = true
			isFirst = true

			opponent = client{responseArr[1], responseArr[2]}
		}

	} else if responseArr[0] == "success" {
		//second client
		fmt.Printf("%s is waiting for you (%s). \n", responseArr[1], responseArr[2])
		fmt.Printf("%s plays first. \n", responseArr[1])
		targetServer = responseArr[2]

		myTurn = false
		isFirst = false

		opponent = client{responseArr[1], responseArr[2]}
	}

	tcpconn.Close()

	conn, err = net.ListenPacket("udp", ":"+strconv.Itoa(localAddr.Port))

	if err != nil {
		fmt.Printf("Please check your opponent is running\n")
		return
	}
	defer conn.Close()

	udplocalAddr := conn.LocalAddr().(*net.UDPAddr) // get local addr

	fmt.Printf("Client is running on port %d\n", udplocalAddr.Port)

	server_addr, _ := net.ResolveUDPAddr("udp", targetServer) // make path for give and take packet with server
	addr = server_addr
	//region [game]

	omok = [10][10]int{{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}

	isDone = false

	printOmok()

	if myTurn {
		timer = make(chan bool)

		go func() {
			select {
			case <-timer:
				//success
				// fmt.Printf("timer cancel ddd\n")
			case <-time.After(time.Second * 10):
				sendPacket(conn, "4|", addr)

			}
		}()
	}

	//endregion

	go handlePacket(conn, server_addr)
	handleInput(conn, server_addr)

	defer conn.Close() // although when client gets panic, defer should disconnect socket gracefully
}

func byebye() {
	fmt.Printf("Bye bye~\n")
	os.Exit(0)
}

func handlePacket(conn net.PacketConn, addr *net.UDPAddr) {
	for {
		buffer := make([]byte, 1024)
		bufferSize := readPacket(conn, &buffer, addr)

		response := string(buffer[:bufferSize])

		// route, _ := strconv.Atoi(strings.Split(response, "|")[0])
		// msgArr := strings.Split(response, "|")

		/*
			server packet form

			(Route|Option Number|Message) client request packet
			(Route|Message) server one sided packet


		*/

		emit(conn, response, false)
	}

}

func handleInput(conn net.PacketConn, addr *net.UDPAddr) {
	for {
		time.Sleep(time.Millisecond * 100)
		// var inputstr string

		inputstr, _ := bufio.NewReader(os.Stdin).ReadString('\n')

		if len(inputstr) < 1 || len(strings.TrimSpace(inputstr)) == 0 {
			continue
		}

		// inputstr = inputstr[:len(inputstr)-1]

		processMyMessage(inputstr, conn, addr)
	}

}

func isInteger(v string) bool {
	if _, err := strconv.Atoi(v); err == nil {
		return true
	}
	return false
}
func processMyMessage(inputstr string, conn net.PacketConn, addr *net.UDPAddr) {
	/*
		if option is given, it sends packet to server and get response.
	*/

	/*

		define route

		0|message    normal chat
		1|x y		 place stone
		2|gg
		3|exit

	*/
	if inputstr[:2] == "\\\\" {
		//place stone
		if isDone {
			fmt.Printf("This game is already done. \n")
			return
		}

		if checkWin() != 0 {
			fmt.Printf("This game is already done. \n")
			return
		}

		inputstr = inputstr[:len(inputstr)-1]
		if !myTurn {
			fmt.Printf("Not your turn \n")
			return
		}

		temp := strings.Split(inputstr, "\\\\")[1]
		tempArr := strings.Split(temp, " ")

		var x, y int
		x = -1
		y = -1

		xyCount := 0
		for _, v := range tempArr {
			if isInteger(v) {
				if xyCount == 0 {
					x, _ = strconv.Atoi(v)
					xyCount++
				} else if xyCount == 1 {
					y, _ = strconv.Atoi(v)
					xyCount++
					break
				}
			}

		}

		if x == -1 || y == -1 {
			fmt.Printf("Invalid Command \n")
			return
		}

		if canPlace(x, y) {
			sendPacket(conn, "1|"+strconv.Itoa(x)+" "+strconv.Itoa(y), addr)
		} else {
			fmt.Printf("Invalid move! \n")
		}

	} else if inputstr[:1] == "\\" {
		//command

		inputstr = inputstr[:len(inputstr)-1]

		command := strings.Split(inputstr, "\\")[1]

		switch command {
		case "gg":
			if isDone {
				fmt.Printf("This game is already done. \n")
				return
			}

			if checkWin() != 0 {
				fmt.Printf("This game is already done. \n")
				return
			}

			sendPacket(conn, "2|", addr)

			break
		case "exit":
			sendPacket(conn, "3|", addr)

			break

		default:
			fmt.Printf("Invalid Command \n")
			return
		}
	} else {
		//normal chat

		sendPacket(conn, "0|"+inputstr, addr)

	}

}

func sendPacketTCP(conn net.Conn, requestString string) {
	//send Packet to server
	conn.Write([]byte(requestString))
}

func readPacketTCP(conn net.Conn, buffer *[]byte) int {
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

func sendPacket(conn net.PacketConn, requestString string, addr *net.UDPAddr) {
	//send Packet to server
	conn.WriteTo([]byte(requestString), addr)

	emit(conn, requestString, true)
}

func readPacket(conn net.PacketConn, buffer *[]byte, addr *net.UDPAddr) int {
	//read Packet from server and saves to buffer and return buffer size.
	count, _, err := conn.ReadFrom(*buffer)
	if err != nil {
		fmt.Println("connection is closed by server")
		byebye()
		conn.Close()
		os.Exit(0)
	}
	return count
}

func emit(conn net.PacketConn, str string, isMine bool) {
	/*

		define route

		0|message    normal chat
		1|x y		 place stone
		2|gg
		3|exit
		4|timeout

	*/

	route, _ := strconv.Atoi(strings.Split(str, "|")[0])
	msgArr := strings.Split(str, "|")

	switch route {
	case 0:
		if isMine {
			//my chatting does not need to display!
			return
		} else {
			fmt.Printf("%s> %s\n", opponent.nickname, msgArr[1])
		}
		break
	case 1:
		x, _ := strconv.Atoi(strings.Split(msgArr[1], " ")[0])
		y, _ := strconv.Atoi(strings.Split(msgArr[1], " ")[1])
		if isMine {
			timer <- true
			if isFirst {
				placeStone(x, y, 1)
			} else {
				placeStone(x, y, 2)
			}
			myTurn = false

		} else {
			timer = make(chan bool)

			go func() {
				select {
				case <-timer:
					// fmt.Printf("timer cancel\n")
					//success
				case <-time.After(time.Second * 10):
					sendPacket(conn, "4|", addr)

				}
			}()

			if isFirst {
				placeStone(x, y, 2)
			} else {
				placeStone(x, y, 1)
			}
			myTurn = true
		}

		printOmok()

		checkwin := checkWin()

		if checkwin == 1 {
			if isFirst {
				fmt.Printf("you win\n")
			} else {
				fmt.Printf("you lose\n")
			}
			isDone = true
		} else if checkwin == 2 {
			if isFirst {
				fmt.Printf("you lose\n")
			} else {
				fmt.Printf("you win\n")
			}
			isDone = true
		}

		break
	case 2:
		if isMine {
			//if i gave up the game
			fmt.Printf("you lose\n")
		} else {
			//if opponent give up the game
			fmt.Printf("%s give up this game ! \n", opponent.nickname)
			fmt.Printf("you win\n")
		}
		isDone = true
		break
	case 3:
		if !isDone {
			if isMine {
				//if i gave up the game
				fmt.Printf("you lose\n")
				byebye()
				os.Exit(0)
			} else {
				//if opponent give up the game
				fmt.Printf("%s left the game ! \n", opponent.nickname)
				fmt.Printf("you win\n")
			}
			isDone = true
		} else {
			if isMine {
				//if i gave up the game
				byebye()
				os.Exit(0)
			} else {
				//if opponent give up the game
				fmt.Printf("%s left the game ! \n", opponent.nickname)
			}
		}

		break

	case 4:
		if !isDone {
			if isMine {
				//if i gave up the game
				fmt.Printf("time out..(10s)\n")
				fmt.Printf("you lose\n")

			} else {
				//if opponent give up the game
				fmt.Printf("%s time out..(10s) \n", opponent.nickname)
				fmt.Printf("you win\n")
			}
			isDone = true
		}

		break

	}
}

func printOmok() {
	fmt.Printf("  y 0 1 2 3 4 5 6 7 8 9  \n")
	fmt.Printf("x -----------------------\n")
	for c, row := range omok {
		fmt.Printf("%d |", c)

		for _, val := range row {
			if val == 0 {
				fmt.Printf(" +")
			} else if val == 1 {
				fmt.Printf(" O")
			} else if val == 2 {
				fmt.Printf(" @")
			}
		}

		fmt.Printf(" |\n")
	}
	fmt.Printf("  -----------------------\n")

}

func placeStone(x int, y int, num int) {
	if omok[x][y] != 0 {
		fmt.Printf("Unknown err (%d, %d) is already %d\n", x, y, omok[x][y])
		return
	}
	omok[x][y] = num
}

func checkWin() int {
	lastChecked := 0
	checkSum := 0

	for _, row := range omok {
		lastChecked = 0
		checkSum = 0
		for _, val := range row {
			if val == lastChecked {
				if val != 0 {
					checkSum++
					if checkSum >= 5 {
						return lastChecked
					}
				}
			} else {
				checkSum = 1
			}
			lastChecked = val

		}
	}

	lastChecked = 0
	checkSum = 0

	for i := 0; i <= 9; i++ {
		lastChecked = 0
		checkSum = 0
		for j := 0; j <= 9; j++ {
			val := omok[j][i]

			if val == lastChecked {
				if val != 0 {
					checkSum++
					if checkSum >= 5 {
						return lastChecked
					}
				}
			} else {
				checkSum = 1
			}
			lastChecked = val

		}
	}

	startX := 9
	startY := 0

	currX := startX
	currY := startY

	lastChecked = 0
	checkSum = 0
	for {
		if startX < 0 || startX > 9 || startY < 0 || startY > 9 {
			break
		}
		val := omok[currX][currY]

		if val == lastChecked {
			if val != 0 {
				checkSum++
				if checkSum >= 5 {
					return lastChecked
				}
			}
		} else {
			checkSum = 1
		}
		lastChecked = val

		currX += 1
		currY += 1

		if currX >= 10 || currY >= 10 {
			startX -= 1

			currX = startX
			currY = startY

			lastChecked = 0
			checkSum = 0
		}
	}

	startX = 0
	startY = 1

	currX = startX
	currY = startY

	lastChecked = 0
	checkSum = 0
	for {
		if startX < 0 || startX > 9 || startY < 0 || startY > 9 {
			break
		}
		val := omok[currX][currY]

		if val == lastChecked {
			if val != 0 {
				checkSum++
				if checkSum >= 5 {
					return lastChecked
				}
			}
		} else {
			checkSum = 1
		}
		lastChecked = val

		currX += 1
		currY += 1

		if currX >= 10 || currY >= 10 {
			startY += 1

			currX = startX
			currY = startY

			lastChecked = 0
			checkSum = 0
		}
	}

	startX = 0
	startY = 0

	currX = startX
	currY = startY

	lastChecked = 0
	checkSum = 0
	for {
		if startX < 0 || startX > 9 || startY < 0 || startY > 9 {
			break
		}
		val := omok[currX][currY]

		if val == lastChecked {
			if val != 0 {
				checkSum++
				if checkSum >= 5 {
					return lastChecked
				}
			}
		} else {
			checkSum = 1
		}
		lastChecked = val

		currX += 1
		currY -= 1

		if currX >= 10 || currY < 0 {
			startY += 1

			currX = startX
			currY = startY

			lastChecked = 0
			checkSum = 0
		}
	}

	startX = 1
	startY = 9

	currX = startX
	currY = startY

	lastChecked = 0
	checkSum = 0
	for {
		if startX < 0 || startX > 9 || startY < 0 || startY > 9 {
			break
		}
		val := omok[currX][currY]

		if val == lastChecked {
			if val != 0 {
				checkSum++
				if checkSum >= 5 {
					return lastChecked
				}
			}
		} else {
			checkSum = 1
		}
		lastChecked = val

		currX += 1
		currY -= 1

		if currX >= 10 || currY < 0 {
			startX += 1

			currX = startX
			currY = startY

			lastChecked = 0
			checkSum = 0
		}
	}

	return 0
}

func canPlace(x int, y int) bool {
	if x < 0 || x >= 10 || y < 0 || y >= 10 {
		return false
	}
	return omok[x][y] == 0
}
