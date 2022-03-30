/**
 * TCPClient.go
 **/

package main

import ("bufio"; "fmt"; "net"; "os")

func main() {

    serverName := "localhost"
    serverPort := "26342"

    conn, _:= net.Dial("tcp", serverName+":"+serverPort)

    localAddr := conn.LocalAddr().(*net.TCPAddr)
    fmt.Printf("Client is running on port %d\n", localAddr.Port)

    fmt.Printf("Input lowercase sentence: ")
    input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
    conn.Write([]byte(input))

    buffer := make([]byte, 1024)
    conn.Read(buffer)
    fmt.Printf("Reply from server: %s", string(buffer))

    conn.Close()
}


// option 1) convert text to UPPER-case letters. // a feature that SimpleEcho programs already have.
// option 2) ask the server what the IP address and port number of the client is.
// option 3) ask the server how many client requests(commands) it has served so far.
// option 4) ask the server program how long it has been running for since it started.
// option 5) exit client program