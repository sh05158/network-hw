/**
 * 20176342 Song Min Joon
 * EasyTCPClient.java
 **/

package submit_20176342_hw2;

import java.util.*;
import java.io.DataOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.Socket;
import java.net.UnknownHostException;

public class EasyTCPClient {

    final static String serverName = "nsl2.cau.ac.kr"; // server host
    final static int serverPort = 26342;// server port

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
        Socket conn = null;

        try {
            conn = new Socket(serverName, serverPort); // create TCP Socket
            Runtime.getRuntime().addShutdownHook(new ByeByeThread(conn)); // shutdown hook for graceful exit

            OutputStream os = conn.getOutputStream(); // output stream
            InputStream is = conn.getInputStream(); // input stream

            HandleInput H = new HandleInput(conn, os, is); // sub thread for input
            H.start(); //thread start

        } catch (UnknownHostException e) {
            //server not connected
            System.out.println("Please check your server is running");
            System.exit(0);
        } catch (IOException e) {
            //server not connected
            System.out.println("Please check your server is running");
            System.exit(0);

        }

    }

    public static void byebye() {
        //print bye bye
        System.out.println("Bye bye~");
    }

    public static void processOption(int opt, Socket conn, OutputStream os, InputStream is) {
        long startTime = new Date().getTime(); //startTime for print RTT

        try {
            String requestString;
            String reply;

            switch (opt) {
                case 1:
		            // Option 1

                    System.out.printf("Input lowercase sentence: ");
                    Scanner sc = new Scanner(System.in);
                    String input = sc.nextLine();
                    requestString = opt + "|" + input;

                    startTime = new Date().getTime();

                    sendPacket(os, requestString);

                    reply = readPacket(is);

                    System.out.printf("Reply from server: %s\n", reply);
                    break;
                case 2:
		            // Option 2

                    requestString = opt + "|";
                    sendPacket(os, requestString);
                    reply = readPacket(is);
                    System.out.printf("Reply from server: client IP = %s, port = %s\n", reply.split(":")[0],
                            reply.split(":")[1]);
                    break;
                case 3:
		            // Option 3

                    requestString = opt + "|";
                    sendPacket(os, requestString);
                    reply = readPacket(is);
                    System.out.printf("Reply from server: requests served = %s\n", reply);
                    break;
                case 4:
		            // Option 4

                    requestString = opt + "|";
                    sendPacket(os, requestString);
                    reply = readPacket(is);
                    printDuration(reply);
                    break;
                case 5:
		            // Option 5

                    conn.close();
                    System.exit(0);
                    break;
                default:
                    // not Option 1~5 (default)
                    conn.close();
                    System.exit(0);
                    break;

            }

        } catch (UnknownHostException e) {
            e.printStackTrace();
        } catch (IOException e) {
            e.printStackTrace();
        }

        printRTT(startTime);
    }

    public static void printRTT(long startTime) {
        //print RTT function

        System.out.printf("RTT = %dms\n\n\n", (new Date().getTime() - startTime));
    }

    public static void printDuration(String t) {
        //print server running time in proper form(HH:MM:ss)

        int h = 0;
        int m = 0;
        int s = 0;

        if (t.contains("h")) {
            h = Integer.parseInt(t.split("h")[0]);
            t = t.split("h")[1];
        }

        if (t.contains("m")) {
            m = Integer.parseInt(t.split("m")[0]);
            t = t.split("m")[1];
        }

        if (t.contains("s")) {
            s = Integer.parseInt(t.split("s")[0]);
        }

        System.out.printf("Reply from server: run time = %02d:%02d:%02d\n", h, m, s);
    }

    public static void sendPacket(OutputStream os, String requestString) {
        //use outputStream to send packet to server
        try {
            os.write(requestString.getBytes());
            os.flush();
        } catch (IOException e) {
            // do nothing
        }

    }

    public static String readPacket(InputStream is) {
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

    public static void printOption() {
    	//print Menu and 5 Options.

        System.out.printf("<Menu>\n");
        System.out.printf("option 1) convert text to UPPER-case letters.\n");
        System.out.printf("option 2) ask the server what the IP address and port number of the client is.\n");
        System.out.printf("option 3) ask the server how many client requests(commands) it has served so far.\n");
        System.out.printf("option 4) ask the server program how long it has been running for since it started.\n");
        System.out.printf("option 5) exit client program\n");
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
                    printOption();
                    System.out.printf("Please select your option :");
                    Scanner sc = new Scanner(System.in);
                    int opt = sc.nextInt();
                    processOption(opt, this.conn, this.os, this.is);
                } catch (NoSuchElementException e) {
                    System.exit(0);
                }

            }

        }
    }
}