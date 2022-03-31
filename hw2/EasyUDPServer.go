/**
 * EasyUDPServer.go
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

var totalRequests int = 0
var startTime time.Time
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

	startTime = time.Now()
	conn, _ := net.ListenPacket("udp", ":"+serverPort)

	fmt.Printf("Server is ready to receive on port %s\n", serverPort)

	for {
		buffer := make([]byte, 1024)
		_, r_addr, _ := conn.ReadFrom(buffer)
		fmt.Printf("UDP message from %s\n", r_addr.String())

		handleMsg(conn, r_addr, string(buffer))
	}

	defer byebye()

}

func byebye() {
	fmt.Printf("Bye bye~\n")
	os.Exit(0)
}

func handleError(conn net.Conn, err error, errmsg string) {
	if conn != nil {
		conn.Close()
	}
	// fmt.Println(err)
	fmt.Println(errmsg)
}
func handleError2(conn net.Conn, errmsg string) {
	if conn != nil {
		conn.Close()
	}

	fmt.Println(errmsg)
}

func handleMsg(conn net.PacketConn, addr net.Addr, clientMsg string) {

	totalRequests++
	tempStr := clientMsg
	requestOption, _ := strconv.Atoi(strings.Split(tempStr, "|")[0])
	requestData := strings.Split(tempStr, "|")[1]
	time.Sleep(time.Millisecond * 1)

	fmt.Printf("Command %d\n\n", requestOption)
	switch requestOption {
	case 1:
		upperString := strings.ToUpper(requestData)
		sendPacket(conn, addr, upperString)
	case 2:
		// ip := conn.RemoteAddr()
		sendPacket(conn, addr, addr.String())

	case 3:
		sendPacket(conn, addr, strconv.Itoa(totalRequests))

	case 4:
		elapsed := time.Since(startTime).Truncate(time.Second).String()
		sendPacket(conn, addr, string(elapsed))

	case 5:
		conn.Close()
	default:
		conn.Close()
	}

}

func sendPacket(conn net.PacketConn, addr net.Addr, serverMsg string) {
	conn.WriteTo([]byte(serverMsg), addr)
}
