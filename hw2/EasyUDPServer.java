package hw2;

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
    
    final public static int serverPort = 26342;
    public static int totalRequest = 0;
    public static long startTime;
    public static void main(String args[]){

        startTime = new Date().getTime();
        DatagramSocket conn = null;

        try{
            conn = new DatagramSocket(serverPort);
        } catch(SocketException e){
            System.out.println("maybe port is already used.."+serverPort);
            System.exit(0);
        }

        System.out.printf("Server is ready to receive on port %s\n", serverPort);
        Runtime.getRuntime().addShutdownHook(new ByeByeThread(conn));


        ServerReceiver th = new ServerReceiver(conn);
        th.start();

        
    }

    public static void byebye(){
        System.out.println("Bye bye~");
    }

    public static String milliToTimeFormat(long ms){
        int ss = (int)Math.floor(ms / 1000);

        String res = "";

        int h = ss / 60 / 60;
        ss = ss - (h*60*60);

        int m = ss / 60;
        ss = ss - (m*60);

        int s = ss;

        if(h >= 1){
            if(h < 10){
                res += "0"+h;
            } else {
                res += h;
            }
            res+="h";
        }

        if(m >= 1){
            if(m < 10){
                res += "0"+m;
            } else {
                res += m;
            }
            res+="m";
        }

        if(s >= 1){
            if(s < 10){
                res += "0"+s;
            } else {
                res += s;
            }
            res+="s";
        }

        return res;

    }

    public static class ByeByeThread extends Thread{
        DatagramSocket conn;

        ByeByeThread(DatagramSocket conn){
            this.conn = conn;
        }

        public void run() {
            try {
                Thread.sleep(200);
                byebye();

               
                this.conn.close();
                
            
                //some cleaning up code...

            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
                e.printStackTrace();
            }
        }
    }

    static class ServerReceiver extends Thread {
        DatagramSocket conn;
        

        ServerReceiver(DatagramSocket conn){
            this.conn = conn;
            
        }

        public void run(){
            while(true){
                try{

                    DatagramPacket buffer = new DatagramPacket(new byte[1024],1024);
                    this.conn.receive(buffer);

                    System.out.printf("UDP message from %s\n", buffer.getAddress()+":"+buffer.getPort());

                    InetAddress clientAddress = buffer.getAddress();
                    int clientPort = buffer.getPort();

                    handleMsg(new String(buffer.getData()), clientAddress, clientPort);
                }catch(IOException e){
                    System.out.println("client is disconnected");
                    break;
                }
            }
            
        }

        public void handleMsg(String msg, InetAddress addr, int port){
            totalRequest++;
            String[] req = msg.split("\\|");
            int requestOption = Integer.parseInt( req[0] );

            System.out.printf("Command %d\n\n\n",requestOption);

            String requestData;
            try{
                requestData = req[1];
                
            } catch(ArrayIndexOutOfBoundsException e){
                requestData = "";
            }
            
            switch(requestOption){
                case 1:
                    sendPacket(requestData.toUpperCase(), addr, port);
                    break;
                case 2:
                    String ip = addr+":"+port;
                    sendPacket(ip, addr, port);
                    break;
                case 3:
                    sendPacket(totalRequest+"", addr, port);
                    break;
                case 4:
                    long elapsed = new Date().getTime() - startTime;
                    sendPacket(milliToTimeFormat(elapsed)+"", addr, port);
                    break;
                case 5:
                    this.conn.close();
                    this.stop();
                    break;
                default:
                    this.conn.close();
                    this.stop();
                    break;
            }
        }

        void sendPacket(String packet, InetAddress addr, int port){
            try{
                DatagramPacket msg = new DatagramPacket(packet.getBytes(), packet.getBytes().length, addr, port);
                this.conn.send(msg);
            } catch(IOException e){

            }
            
        }

    }
}