/**
 * 20176342 Song Min Joon
 * EasyUDPServer.java
 **/

package submit_20176342_hw2;

import java.util.*;
import java.io.DataInputStream;
import java.io.DataOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.ServerSocket;
import java.net.Socket;
import java.net.SocketAddress;
import java.net.SocketException;
import java.net.UnknownHostException;
import java.sql.Time;
import java.net.DatagramPacket;
import java.net.DatagramSocket;
import java.net.InetAddress;

public class EasyUDPServer {

    final public static int serverPort = 26342;// for server port
    public static int totalRequest = 0;// total Request count global variable for server.
    public static long startTime;// for saving server start time

    public static void main(String args[]) {

        startTime = new Date().getTime(); // save server start time(curr time)
        DatagramSocket conn = null;

        try {
            conn = new DatagramSocket(serverPort); // create listener for udp
        } catch (SocketException e) {
            System.out.println("maybe port is already used.." + serverPort);
            System.exit(0);
        }

        System.out.printf("Server is ready to receive on port %s\n", serverPort);
        Runtime.getRuntime().addShutdownHook(new ByeByeThread(conn));// shutdown hook for graceful exit

        ServerReceiver th = new ServerReceiver(conn);
        // cause there is no 'connection' in udp, server open udp path and waiting for
        // any packet
        // by start ServerReceiver thread
        th.start();

    }

    public static void byebye() {
        // print bye bye~
        System.out.println("Bye bye~");
    }

    public static String milliToTimeFormat(long ms) {
        // convert Date.now(millisecond) format to HH:MM:ss
        int ss = (int) Math.floor(ms / 1000);

        String res = "";

        int h = ss / 60 / 60;
        ss = ss - (h * 60 * 60);

        int m = ss / 60;
        ss = ss - (m * 60);

        int s = ss;

        if (h >= 1) {
            if (h < 10) {
                res += "0" + h;
            } else {
                res += h;
            }
            res += "h";
        }

        if (m >= 1) {
            if (m < 10) {
                res += "0" + m;
            } else {
                res += m;
            }
            res += "m";
        }

        if (s >= 1) {
            if (s < 10) {
                res += "0" + s;
            } else {
                res += s;
            }
            res += "s";
        }

        return res;

    }

    public static class ByeByeThread extends Thread {
        // ByeBye Thread for graceful exit program.
        DatagramSocket conn;

        ByeByeThread(DatagramSocket conn) {
            this.conn = conn;
        }

        public void run() {
            try {
                Thread.sleep(200);
                byebye();

                this.conn.close(); // close listener

            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
                e.printStackTrace();
            }
        }
    }

    static class ServerReceiver extends Thread {
        // ServerReceiver Thread for all client
        DatagramSocket conn;

        ServerReceiver(DatagramSocket conn) {
            this.conn = conn;

        }

        public void run() {
            while (true) {
                try {

                    DatagramPacket buffer = new DatagramPacket(new byte[1024], 1024);
                    this.conn.receive(buffer);// read packet into buffer

                    System.out.printf("UDP message from %s\n", buffer.getAddress() + ":" + buffer.getPort()); // print
                                                                                                              // ip,
                                                                                                              // port

                    InetAddress clientAddress = buffer.getAddress(); // get sender ip
                    int clientPort = buffer.getPort(); // get sender port

                    handleMsg(new String(buffer.getData()), clientAddress, clientPort);// handle client msg using client
                                                                                       // address and port
                } catch (IOException e) {
                    System.out.println("client is disconnected");
                    break;
                }
            }

        }

        public void handleMsg(String msg, InetAddress addr, int port) {
            totalRequest++;// number of request is added

            /*
             * client packet form
             * 
             * (Option Number|Message)
             * 
             * 1|blah blah blah...
             * 1|hello world!
             * 
             * 2| => message is not required
             * 3| => message is not required
             * 4| => message is not required
             * 
             * 5| => maybe not arrived??
             * 
             */

            String[] req = msg.split("\\|"); // split client packet by '|' and takes option and convert to Integer.
            int requestOption = Integer.parseInt(req[0]);

            System.out.printf("Command %d\n\n\n", requestOption); // print Command #

            String requestData;
            try {
                requestData = req[1]; // get message parameter from packet.

            } catch (ArrayIndexOutOfBoundsException e) {
                requestData = "";
            }

            switch (requestOption) {
                case 1:
                    // Option 1
                    sendPacket(requestData.toUpperCase(), addr, port);
                    break;
                case 2:
                    // Option 2
                    String ip = addr + ":" + port;
                    sendPacket(ip, addr, port);
                    break;
                case 3:
                    // Option 3
                    sendPacket(totalRequest + "", addr, port);
                    break;
                case 4:
                    // Option 4
                    long elapsed = new Date().getTime() - startTime;
                    sendPacket(milliToTimeFormat(elapsed) + "", addr, port);
                    break;
                case 5:
                    // Option 5
                    this.conn.close();
                    this.stop();
                    break;
                default:
                    // Option default
                    this.conn.close();
                    this.stop();
                    break;
            }
        }

        void sendPacket(String packet, InetAddress addr, int port) {
            // send packet to server
            //using client ip, port (UDP)
            try {
                DatagramPacket msg = new DatagramPacket(packet.getBytes(), packet.getBytes().length, addr, port);
                this.conn.send(msg);
            } catch (IOException e) {

            }

        }

    }
}