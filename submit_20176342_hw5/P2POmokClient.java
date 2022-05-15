/**
 * 20176342 Song Min Joon
 * P2POmokClient.java
 **/

// package submit_20176342_hw5;

import java.util.*;
import java.io.DataOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.DatagramPacket;
import java.net.DatagramSocket;
import java.net.InetAddress;
import java.net.Socket;
import java.net.SocketException;
import java.net.UnknownHostException;

class client {
	String nickname;
	String remote;

	client(){

	}
	client(String nickname, String remote){
		this.nickname = nickname;
		this.remote = remote;
	}
}
public class P2POmokClient {

    final static String serverName = "localhost"; // server host
    final static int serverPort = 56342;// server port
	static client opponent;
	static Timer timer;
	static TimerTask task;
	static Socket tcpconn;
	static DatagramSocket conn;

	static String targetServer;

	static boolean myTurn;
	static boolean isFirst;
	static Omok omok;
    public static class ByeByeThread extends Thread {
        // ByeBye Thread for graceful exit program.
        Socket conn;

        ByeByeThread(Socket conn) {
            this.conn = conn;
        }

        public void run() {
            try {
                Thread.sleep(200);
                byebye();

                try {
                    this.conn.close();
                } catch (IOException e) {

                }

            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
                e.printStackTrace();

            }
        }
    }
    public static void main(String[] args) {
		if(args.length < 2){
			System.out.println("Please check your nickname argunemt");
			byebye();
			// return;
		}

		String nickname = "args[1]";
		System.out.println(nickname);
		if(nickname.length() <= 0){
			System.out.println("Please check your nickname argunemt");
			byebye();
			// return;

		}

        try {
            tcpconn = new Socket(serverName, serverPort); // create TCP Socket
            Runtime.getRuntime().addShutdownHook(new ByeByeThread(tcpconn)); // shutdown hook for graceful exit

            

        } catch (UnknownHostException e) {
            //server not connected
            System.out.println("Please check your server is running");
            System.exit(0);
        } catch (IOException e) {
            //server not connected
            System.out.println("Please check your server is running");
            System.exit(0);

        }


		OutputStream os = tcpconn.getOutputStream(); // output stream
		InputStream is = tcpconn.getInputStream(); // input stream


		InetAddress localAddr = tcpconn.getLocalAddress();
		int localPort = tcpconn.getLocalPort();
		String localIP = tcpconn.getLocalAddress().toString();

		System.out.println("Client is running on "+tcpconn.getLocalAddress()+":"+tcpconn.getLocalPort());

		sendPacketTCP(os, nickname);

		String res = readPacketTCP(is);
		String[] responseArr = res.split("|");

		System.out.printf("welcome %s to p2p-omok server at %s. \n",nickname, tcpconn.getRemoteSocketAddress().toString());
		
		if(responseArr[0] == "waiting"){
			System.out.printf("waiting for an opponent. \n");

			res = readPacketTCP(is);

			responseArr = res.split("|");

			if(responseArr[0] == "matched"){
				System.out.printf("%s joined (%s). you play first. \n", responseArr[1], responseArr[2]);
				targetServer = responseArr[2];
				myTurn = true;
				isFirst = true;
				opponent = new client(responseArr[1], responseArr[1]);

			}
		}
		else if(responseArr[0] == "success"){
			//second client
			System.out.printf("%s is waiting for you (%s). \n", responseArr[1], responseArr[2]);
			System.out.printf("%s plays first. \n", responseArr[1]);
			targetServer = responseArr[2];

			myTurn = false;
			isFirst = false;

			opponent = new client(responseArr[1], responseArr[2]);
		}

		tcpconn.close();

		try{
			conn = new DatagramSocket(localPort);
			
		} catch(SocketException e){
			System.out.println("please check your udp socket : "+localPort);
			byebye();
			return;
		}

		System.out.printf("Client is running on %s \n",conn.getLocalAddress());

		omok = new Omok();
		omok.initOmok();
		omok.printOmok();

		if(myTurn){

		}


		

    }

	public static void emit(DatagramSocket conn, String str, boolean isMine){
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

    public static void byebye() {
        //print bye bye
        System.out.println("Bye bye~");
    }

    public static void sendPacketTCP(OutputStream os, String requestString) {
        //use outputStream to send packet to server
        try {
            os.write(requestString.getBytes());
            os.flush();
        } catch (IOException e) {
            // do nothing
        }

    }

    public static String readPacketTCP(InputStream is) {
	    //read Packet from server using InputStream.

        byte[] data = new byte[1024];

        try {
            int n = is.read(data);
            String res = new String(data, 0, n);
            return res;

        } catch (IOException e) {

            System.out.println("disconnected by server");
            System.exit(0);

            return "";
        }

    }

    public static void sendPacket(DatagramSocket conn, String requestString) {
        // use datagram Packet to send packet to server

        try {
            InetAddress serverIP = InetAddress.getByName(serverName);

            DatagramPacket dp = new DatagramPacket(requestString.getBytes(), requestString.getBytes().length, serverIP,
                    serverPort);

            conn.send(dp);

        } catch (IOException e) {
            // do nothing
        }

    }

    public static String readPacket(DatagramSocket conn) {
        // read Packet from server using Datagram.

        try {
            byte[] buffer = new byte[1024];
            DatagramPacket response = new DatagramPacket(buffer, buffer.length);
            conn.receive(response); // receive
            String res = new String(buffer, 0, response.getLength());

            return res;

        } catch (IOException e) {

            return null;
        }

    }


	static class HandleInput extends Thread {
        //HandleInput thread
        Socket conn;
        OutputStream os;
        InputStream is;

        HandleInput(Socket conn, OutputStream os, InputStream is) {
            this.conn = conn;
            this.os = os;
            this.is = is;
        }

        public void run() {
            while (true) {
                try {
                    //infinite loop
                    //print 5 options and take option input from user
                    System.out.printf("Please select your option :");
                    Scanner sc = new Scanner(System.in);
                    int opt = sc.nextInt();
                } catch (NoSuchElementException e) {
                    System.exit(0);
                }

            }

        }
    }

}

class Omok {
	int[][] omok;
	boolean isDone;
	Omok(){
		this.omok = new int[10][10];
		this.initOmok();
		this.isDone = false;
	}

	void initOmok(){
		for(var i = 0; i<10; i++){
			for( var j = 0; j<10; j++){
				this.omok[i][j]=0;
			}
		}
	}

	void printOmok(){
		System.out.printf("  y 0 1 2 3 4 5 6 7 8 9  \n");
		System.out.printf("x -----------------------\n");
		for (var i = 0; i<10; i++) {
			System.out.printf("%d |", i);
	
			for (var j = 0; j<10; j++) {
				int val = omok[i][j];
				if (val == 0) {
					System.out.printf(" +");
				} else if (val == 1) {
					System.out.printf(" O");
				} else if (val == 2) {
					System.out.printf(" @");
				}
			}
	
			System.out.printf(" |\n");
		}
		System.out.printf("  -----------------------\n");
	}
	void placeStone(int x, int y, int num){
		if (omok[x][y] != 0) {
			System.out.printf("Unknown err (%d, %d) is already %d\n", x, y, omok[x][y]);
			return;
		}
		omok[x][y] = num;
	}
	int checkWin(){
		int lastChecked = 0;
		int checkSum = 0;
	
		for(var i = 0; i<10; i++){
			lastChecked = 0;
			checkSum = 0;

			for( var j = 0; j<10; j++){
				int val = omok[i][j];

				if (val == lastChecked) {
					if (val != 0) {
						checkSum++;
						if (checkSum >= 5) {
							return lastChecked;
						}
					}
				} else {
					checkSum = 1;
				}
				lastChecked = val;
			}
		}
		
		lastChecked = 0;
		checkSum = 0;
	
		for(var i = 0; i<10; i++){
			lastChecked = 0;
			checkSum = 0;

			for( var j = 0; j<10; j++){
				int val = omok[j][i];

				if (val == lastChecked) {
					if (val != 0) {
						checkSum++;
						if (checkSum >= 5) {
							return lastChecked;
						}
					}
				} else {
					checkSum = 1;
				}
				lastChecked = val;
			}
		}
	
		int startX = 9;
		int startY = 0;
	
		int currX = startX;
		int currY = startY;
	
		lastChecked = 0;
		checkSum = 0;

		while(true) {
			if (startX < 0 || startX > 9 || startY < 0 || startY > 9) {
				break;
			}
			int val = omok[currX][currY];
	
			if (val == lastChecked) {
				if (val != 0) {
					checkSum++;
					if (checkSum >= 5) {
						return lastChecked;
					}
				}
			} else {
				checkSum = 1;
			}
			lastChecked = val;
	
			currX += 1;
			currY += 1;
	
			if (currX >= 10 || currY >= 10) {
				startX -= 1;
	
				currX = startX;
				currY = startY;
	
				lastChecked = 0;
				checkSum = 0;
			}
		}
	
		startX = 0;
		startY = 1;
	
		currX = startX;
		currY = startY;
	
		lastChecked = 0;
		checkSum = 0;
		while(true) {
			if (startX < 0 || startX > 9 || startY < 0 || startY > 9) {
				break;
			}
			int val = omok[currX][currY];
	
			if (val == lastChecked) {
				if (val != 0) {
					checkSum++;
					if (checkSum >= 5) {
						return lastChecked;
					}
				}
			} else {
				checkSum = 1;
			}
			lastChecked = val;
	
			currX += 1;
			currY += 1;
	
			if (currX >= 10 || currY >= 10) {
				startY += 1;
	
				currX = startX;
				currY = startY;
	
				lastChecked = 0;
				checkSum = 0;
			}
		}
	
		startX = 0;
		startY = 0;
	
		currX = startX;
		currY = startY;
	
		lastChecked = 0;
		checkSum = 0;
		while(true) {
			if (startX < 0 || startX > 9 || startY < 0 || startY > 9) {
				break;
			}
			int val = omok[currX][currY];
	
			if (val == lastChecked) {
				if (val != 0) {
					checkSum++;
					if (checkSum >= 5) {
						return lastChecked;
					}
				}
			} else {
				checkSum = 1;
			}
			lastChecked = val;
	
			currX += 1;
			currY -= 1;
	
			if (currX >= 10 || currY < 0) {
				startY += 1;
	
				currX = startX;
				currY = startY;
	
				lastChecked = 0;
				checkSum = 0;
			}
		}
	
		startX = 1;
		startY = 9;
	
		currX = startX;
		currY = startY;
	
		lastChecked = 0;
		checkSum = 0;
		while(true) {
			if (startX < 0 || startX > 9 || startY < 0 || startY > 9) {
				break;
			}
			int val = omok[currX][currY];
	
			if (val == lastChecked) {
				if (val != 0) {
					checkSum++;
					if (checkSum >= 5) {
						return lastChecked;
					}
				}
			} else {
				checkSum = 1;
			}
			lastChecked = val;
	
			currX += 1;
			currY -= 1;
	
			if (currX >= 10 || currY < 0) {
				startX += 1;
	
				currX = startX;
				currY = startY;
	
				lastChecked = 0;
				checkSum = 0;
			}
		}
	
		return 0;
	}
	boolean canPlace(int x, int y){
		if (x < 0 || x >= 10 || y < 0 || y >= 10) {
			return false;
		}
		return omok[x][y] == 0;
	}

}