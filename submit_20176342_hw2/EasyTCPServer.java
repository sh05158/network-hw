/**
 * 20176342 Song Min Joon
 * EasyTCPServer.java
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
import java.net.UnknownHostException;
import java.sql.Time;

public class EasyTCPServer {

    final public static int serverPort = 26342;// for server port
    public static int totalRequest = 0; // total Request count global variable for server.
    public static long startTime;// for saving server start time

    public static void main(String args[]) {

        startTime = new Date().getTime(); //save server start time(curr time)
        ServerSocket listener = null;
        Socket conn = null;
        try {
            listener = new ServerSocket(serverPort); // create listener for tcp socket

        } catch (IOException e) {

        }
        System.out.printf("Server is ready to receive on port %s\n", serverPort);
        Runtime.getRuntime().addShutdownHook(new ByeByeThread(listener)); // shutdown hook for graceful exit

        try {
            while (true) {
                // listener is waiting for tcp connection of clients.
                conn = listener.accept();
                System.out.printf("Connection request from %s\n", conn.getInetAddress() + ":" + conn.getPort());

                // when client is connect to server, make individual sub-thread to communicate
                // each client.
                ServerReceiver th = new ServerReceiver(conn);
                th.start();
            }

        } catch (IOException e) {

        }
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
        ServerSocket listener;

        ByeByeThread(ServerSocket listener) {
            this.listener = listener;
        }

        public void run() {
            try {
                Thread.sleep(200);
                byebye();//print bye bye~

                try {
                    this.listener.close(); // close listener
                } catch (IOException e) {

                }

            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
                e.printStackTrace();
            }
        }
    }

    static class ServerReceiver extends Thread {
        // ServerReceiver Thread for each client
        Socket conn;
        DataInputStream is;
        DataOutputStream os;

        ServerReceiver(Socket conn) {
            this.conn = conn;
            try {
                is = new DataInputStream(conn.getInputStream());
                os = new DataOutputStream(conn.getOutputStream());
            } catch (IOException e) {

            }
        }

        public void run() {
            while (true) {
                try {
                    int bufferSize = 1024;
                    byte[] buffer = new byte[bufferSize];

                    int count = is.read(buffer);

                    if (count == -1) {
                        System.out.println("Client disconnected");
                        break;
                    }

                    // read packet from server and handle message
                    String clientMsg = new String(buffer, 0, count);

                    handleMsg(clientMsg);

                } catch (IOException e) {
                    System.out.println("client is disconnected");
                    break;
                }
            }

        }

        public void handleMsg(String msg) {
            totalRequest++; // number of request is added

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
                    sendPacket(requestData.toUpperCase());
                    break;
                case 2:
                    // Option 2
                    String ip = (this.conn.getInetAddress() + ":" + this.conn.getPort());
                    sendPacket(ip);
                    break;
                case 3:
                    // Option 3
                    sendPacket(totalRequest + "");
                    break;
                case 4:
                    // Option 4
                    long elapsed = new Date().getTime() - startTime;
                    sendPacket(milliToTimeFormat(elapsed) + "");
                    break;
                case 5:
                    // Option 5
                    try {
                        this.conn.close();
                        this.stop();
                    } catch (IOException e) {

                    }
                    break;
                default:
                    // Option default
                    try {
                        this.conn.close();
                        this.stop();
                    } catch (IOException e) {

                    }
                    break;
            }
        }

        void sendPacket(String packet) {
            // send packet to server
            try {
                os.write(packet.getBytes());
                os.flush();
            } catch (IOException e) {

            }

        }

    }
}