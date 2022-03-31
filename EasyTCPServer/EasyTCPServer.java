package EasyTCPServer;

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
    
    final public static int serverPort = 26342;
    public static int totalRequest = 0;
    public static long startTime;
    public static void main(String args[]){

        startTime = new Date().getTime();
        ServerSocket listener = null;
        Socket conn = null;
        try{
            listener = new ServerSocket(serverPort);

        } catch(IOException e){

        }
        System.out.printf("Server is ready to receive on port %s\n", serverPort);
        Runtime.getRuntime().addShutdownHook(new ByeByeThread(listener));

        try{
            while(true){
                conn = listener.accept();
                System.out.printf("Connection request from %s\n", conn.getInetAddress()+":"+conn.getPort());

                ServerReceiver th = new ServerReceiver(conn);
                th.start();
            }

        } catch(IOException e){
            
        }
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
        ServerSocket listener;

        ByeByeThread(ServerSocket listener){
            this.listener = listener;
        }

        public void run() {
            try {
                Thread.sleep(200);
                byebye();

                try{
                    this.listener.close();
                } catch(IOException e){

                }
            
                //some cleaning up code...

            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
                e.printStackTrace();
            }
        }
    }

    static class ServerReceiver extends Thread {
        Socket conn;
        DataInputStream is;
        DataOutputStream os;


        ServerReceiver(Socket conn){
            this.conn = conn;
            try{
                is = new DataInputStream(conn.getInputStream());
                os = new DataOutputStream(conn.getOutputStream());
            } catch(IOException e){
                
            }
        }

        public void run(){
            while(true){
                try{
                    int bufferSize = 1024;
                    byte[] buffer = new byte[bufferSize];

                    int count = is.read(buffer);

                    if(count == -1){
                        System.out.println("Client disconnected");
                        break;
                    }

                    String clientMsg = new String(buffer, 0, count);
                    handleMsg(clientMsg);
                }catch(IOException e){
                    System.out.println("client is disconnected");
                    break;
                }
            }
            
        }

        public void handleMsg(String msg){
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
                    sendPacket(requestData.toUpperCase());
                    break;
                case 2:
                    String ip = this.conn.getInetAddress()+":"+this.conn.getPort();
                    sendPacket(ip);
                    break;
                case 3:
                    sendPacket(totalRequest+"");
                    break;
                case 4:
                    long elapsed = new Date().getTime() - startTime;
                    sendPacket(milliToTimeFormat(elapsed)+"");
                    break;
                case 5:
                    try{
                        this.conn.close();
                        this.stop();
                    } catch(IOException e){

                    }
                    break;
                default:
                    try{
                        this.conn.close();
                        this.stop();
                    } catch(IOException e){

                    }
                    break;
            }
        }

        void sendPacket(String packet){
            try{
                os.write(packet.getBytes());
                os.flush();
            } catch(IOException e){

            }
            
        }

    }
}