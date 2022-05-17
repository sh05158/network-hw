/**
 * 20176342 Song Min Joon
 * P2POmokClient.java
 **/

// package submit_20176342_hw5;

import java.util.*;
import java.util.Timer;

import javax.swing.border.Border;
import javax.swing.event.MouseInputListener;
import javax.swing.text.DefaultCaret;

import java.awt.Container;
import java.awt.*;
import javax.swing.*;
import java.io.*;

import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.DatagramPacket;
import java.net.DatagramSocket;
import java.net.InetAddress;
import java.net.Socket;
import java.net.SocketException;
import java.net.UnknownHostException;

import java.awt.Event.*;
import java.awt.event.ActionListener;
import java.awt.event.KeyEvent;

import java.awt.event.KeyListener;

import java.awt.event.WindowAdapter;

import java.awt.event.WindowEvent;
import java.awt.event.MouseEvent;

import java.awt.event.*;
import javax.swing.JTextArea; //or JTextField



import java.util.List;


class clientclient {
	String nickname;
	String remote;

	clientclient(){

	}
	clientclient(String nickname, String remote){
		this.nickname = nickname;
		this.remote = remote;
	}
};

class TextAreaOutputStream extends OutputStream
{

// *************************************************************************************************
// INSTANCE MEMBERS
// *************************************************************************************************

private byte[]                          oneByte;                                                    // array for write(int val);
private Appender                        appender;                                                   // most recent action

public TextAreaOutputStream(JTextArea txtara) {
    this(txtara,1000);
    }

public TextAreaOutputStream(JTextArea txtara, int maxlin) {
    if(maxlin<1) { throw new IllegalArgumentException("TextAreaOutputStream maximum lines must be positive (value="+maxlin+")"); }
    oneByte=new byte[1];
    appender=new Appender(txtara,maxlin);
    }

/** Clear the current console text area. */
public synchronized void clear() {
    if(appender!=null) { appender.clear(); }
    }

public synchronized void close() {
    appender=null;
    }

public synchronized void flush() {
    }

public synchronized void write(int val) {
    oneByte[0]=(byte)val;
    write(oneByte,0,1);
    }

public synchronized void write(byte[] ba) {
    write(ba,0,ba.length);
    }

public synchronized void write(byte[] ba,int str,int len) {
    if(appender!=null) { appender.append(bytesToString(ba,str,len)); }
    }


static private String bytesToString(byte[] ba, int str, int len) {
    try { return new String(ba,str,len,"UTF-8"); } catch(UnsupportedEncodingException thr) { return new String(ba,str,len); } // all JVMs are required to support UTF-8
    }

// *************************************************************************************************
// STATIC MEMBERS
// *************************************************************************************************

    static class Appender
    implements Runnable
    {
    private final JTextArea             textArea;
    private final int                   maxLines;                                                   // maximum lines allowed in text area
    private final LinkedList<Integer>   lengths;                                                    // length of lines within text area
    private final List<String>          values;                                                     // values waiting to be appended

    private int                         curLength;                                                  // length of current line
    private boolean                     clear;
    private boolean                     queue;

    Appender(JTextArea txtara, int maxlin) {
        textArea =txtara;
        maxLines =maxlin;
        lengths  =new LinkedList<Integer>();
        values   =new ArrayList<String>();

        curLength=0;
        clear    =false;
        queue    =true;
        }

    synchronized void append(String val) {
        values.add(val);
        if(queue) { queue=false; EventQueue.invokeLater(this); }
        }

    synchronized void clear() {
        clear=true;
        curLength=0;
        lengths.clear();
        values.clear();
        if(queue) { queue=false; EventQueue.invokeLater(this); }
        }

    // MUST BE THE ONLY METHOD THAT TOUCHES textArea!
    public synchronized void run() {
        if(clear) { textArea.setText(""); }
        for(String val: values) {
            curLength+=val.length();
            if(val.endsWith(EOL1) || val.endsWith(EOL2)) {
                if(lengths.size()>=maxLines) { textArea.replaceRange("",0,lengths.removeFirst()); }
                lengths.addLast(curLength);
                curLength=0;
                }
            textArea.append(val);
            }
        values.clear();
        clear =false;
        queue =true;
        }

    static private final String         EOL1="\n";
    static private final String         EOL2=System.getProperty("line.separator",EOL1);
    }

} /* END PUBLIC CLASS */



	

	

public class P2POmokClient extends JFrame {
	final static String serverName = "nsl2.cau.ac.kr"; // server host
    final static int serverPort = 56342;// server port
	static clientclient opponent;
	static Timer timer = new Timer();
	static TimerTask task;
	static Socket tcpconn;
	static DatagramSocket conn;

	static String targetServer;

	static boolean myTurn;
	static boolean isFirst;
	static Omok omok = null;

	static JPanel drawStone;

	static JTextArea ta = new JTextArea();



	public static void printf(String format, Object ... args){
		System.out.printf(format, args);
		ta.append(String.format(format,args).toString());
	}

	public static void println(String format){
		System.out.println(format);
		ta.append(format.toString()+"\n");
	}

	class MyMouseListener extends Frame implements MouseInputListener {

	
	@Override
	public void mouseClicked(java.awt.event.MouseEvent e) {
		// TODO Auto-generated method stub

		int x = e.getX();
		int y = e.getY();

		x-=175;
		y-=175;

		int detX = -1;
		int detY = -1;

		int startX = -22;
		int startY = -22;
		//x check
		for(int i = 0; i<10;i++){
			if(x >= startX && x<= startX+44){
				detY = i;
				break;
			}
			startX+=50;
		}

		//y check
		for(int i = 0; i<10;i++){
			if(y >= startY && y<= startY+44){
				detX = i;
				break;
			}
			startY+=50;
		}

		if(detX < 0 || detY < 0){
			println("Invalid Click");
			return;
		}
		else {
			// printf("(%d, %d)\n",detX,detY);
		}

		if(!myTurn || omok.isDone){
			println("Invalid Click");
			return;
		}

		HandleInput.processMyMessage("\\\\"+detX+" "+detY);


	}

	@Override
	public void mousePressed(java.awt.event.MouseEvent e) {
		// TODO Auto-generated method stub
		
	}

	@Override
	public void mouseReleased(java.awt.event.MouseEvent e) {
		// TODO Auto-generated method stub
		
	}

	@Override
	public void mouseEntered(java.awt.event.MouseEvent e) {
		// TODO Auto-generated method stub
		
	}

	@Override
	public void mouseExited(java.awt.event.MouseEvent e) {
		// TODO Auto-generated method stub
		
	}

	@Override
	public void mouseDragged(java.awt.event.MouseEvent e) {
		// TODO Auto-generated method stub
		
	}

	@Override
	public void mouseMoved(java.awt.event.MouseEvent e) {
		// TODO Auto-generated method stub
		
	}        
}

class DrawOmok extends JPanel{
	@Override
	protected void paintComponent(Graphics g) {
		super.paintComponent(g);

		g.setColor(new Color(212,163,60));
		g.fillRect(-0,-0,500,500);

		// g.setStroke(new BasicStroke(2));
		g.setColor(new Color(0,0,0));
		for (int x = 25; x <= 430; x += 50)
			for (int y = 25; y <= 430; y += 50)
				g.drawRect(x, y, 50, 50);


				if(omok == null){
					return;
				}
				int [][]arr = omok.omok;
	
				
				// g.setStroke(new BasicStroke(2));
				for (int x = 25, c = 0; x <= 475; x += 50, c++){
					for (int y = 25, cc = 0; y <= 475; y += 50, cc++){
						if(arr[c][cc] != 0){
							if(arr[c][cc] == 1){
								//black stone
								g.setColor(new Color(12,12,12));
								g.fillOval(y-18,x-18,35,35);
							}
							else if(arr[c][cc] == 2){
								//white stone
								g.setColor(new Color(222,222,222));
								g.fillOval(y-18,x-18,35,35);
							}
						}
					}
	
				}
	}

	public Dimension getPreferredSize() {
		return new Dimension(500, 500); // appropriate constants
	}
}

	class DrawStone extends JPanel{
		@Override
		protected void paintComponent(Graphics g) {
			super.paintComponent(g);
			g.setColor(new Color(212,163,60));
			g.fillRect(-0,-0,500,500);

			// g.setStroke(new BasicStroke(2));
			g.setColor(new Color(0,0,0));
			for (int x = 25; x <= 430; x += 50)
				for (int y = 25; y <= 430; y += 50)
					g.drawRect(x, y, 50, 50);

			if(omok == null){
				return;
			}
			int [][]arr = omok.omok;

			
			// g.setStroke(new BasicStroke(2));
			for (int x = 25, c = 0; x <= 475; x += 50, c++){
				for (int y = 25, cc = 0; y <= 475; y += 50, cc++){
					if(arr[c][cc] != 0){
						if(arr[c][cc] == 1){
							//black stone
							g.setColor(new Color(12,12,12));
							g.fillOval(y-18,x-18,35,35);
						}
						else if(arr[c][cc] == 2){
							//white stone
							g.setColor(new Color(222,222,222));
							g.fillOval(y-18,x-18,35,35);
						}
					}
				}

			}
		}
	

		public Dimension getPreferredSize() {
			return new Dimension(500, 500); // appropriate constants
		}
	}

	public P2POmokClient()
	{
		setTitle("omok");
		setDefaultCloseOperation(JFrame.EXIT_ON_CLOSE);
		setVisible(true);
		setSize(1100,800);
		setLayout(null);
		setResizable(false);

		Container contentPane = getContentPane();

		contentPane.setLayout(null);
		contentPane.setSize(800,800);

		ta.setLocation(700,60);
		ta.setSize(300,600);
		ta.setRows(30);
		ta.setColumns(30);
		ta.setLineWrap(true);
		ta.setWrapStyleWord(true);

		DefaultCaret caret = (DefaultCaret)ta.getCaret();
		caret.setUpdatePolicy(DefaultCaret.ALWAYS_UPDATE);

        TextAreaOutputStream taos = new TextAreaOutputStream( ta, 60 );
        PrintStream ps = new PrintStream( taos );
        // System.setOut( ps );
        // System.setErr( ps );
		


		JScrollPane p = new JScrollPane( ta );
		p.setLocation(700,60);
		p.setSize(300,600);
		// uiPanel.revalidate();
		// uiPanel.repaint();

		contentPane.add(p);

		JTextArea input_ta = new JTextArea();
		input_ta.setLocation(700,670);
		input_ta.setSize(233,30);
		input_ta.setRows(1);
		input_ta.setColumns(20);

		input_ta.addKeyListener(new KeyListener(){
			

			@Override
			public void keyTyped(KeyEvent e) {
				// TODO Auto-generated method stub
				
			}

			@Override
			public void keyPressed(KeyEvent e) {
				// TODO Auto-generated method stub
				if(e.getKeyCode() == KeyEvent.VK_ENTER){
					//code to send message goes here
					if(input_ta.getText().length() >= 1){

						HandleInput.processMyMessage(input_ta.getText());
						
						input_ta.setText("");
					}
					
				}
			}

			@Override
			public void keyReleased(KeyEvent e) {
				// TODO Auto-generated method stub

				if(e.getKeyCode() == KeyEvent.VK_ENTER){
					//code to send message goes here

						
						input_ta.setText("");
					
				}
				
			}
		});

		Border border = BorderFactory.createLineBorder(Color.BLACK);
		input_ta.setBorder(border);
		contentPane.add(input_ta);

		JButton chat = new JButton();
		chat.setLocation(936,670);
		chat.setSize(64,30);
		chat.setText("Send");
		chat.addMouseListener(new MouseInputListener(){

			@Override
			public void mouseClicked(MouseEvent e) {
				// TODO Auto-generated method stub
				if(input_ta.getText().length() >= 1){

					HandleInput.processMyMessage(input_ta.getText());
					
					input_ta.setText("");
				}
				
			}

			@Override
			public void mousePressed(MouseEvent e) {
				// TODO Auto-generated method stub
				
			}

			@Override
			public void mouseReleased(MouseEvent e) {
				// TODO Auto-generated method stub
				
			}

			@Override
			public void mouseEntered(MouseEvent e) {
				// TODO Auto-generated method stub
				
			}

			@Override
			public void mouseExited(MouseEvent e) {
				// TODO Auto-generated method stub
				
			}

			@Override
			public void mouseDragged(MouseEvent e) {
				// TODO Auto-generated method stub
				
			}

			@Override
			public void mouseMoved(MouseEvent e) {
				// TODO Auto-generated method stub
				
			}
			
		});
		// chat.addActionListener(new ActionListener(){
		// 	public void actionPerformed(Action e){
		// 		JButton btn = (JButton) ((EventObject) e).getSource();
		// 		if(btn.getText().equals("Click"))
		// 			btn.setText("Hello");
		// 		else 
		// 			btn.setText("Click");
		// 	}
		// });
		contentPane.add(chat);

 
		JPanel omokPanel = new DrawOmok();
		omokPanel.setLocation(150,150);
		omokPanel.setSize(500,500);
		add(omokPanel);
		// contentPane.add(uiPanel);

		contentPane.addMouseListener(new MyMouseListener());



		addWindowListener(new java.awt.event.WindowAdapter() {
			public void windowClosing(WindowEvent winEvt) {
				HandleInput.processMyMessage("\\exit");
				println("bye bye~");
				System.exit(0);
			}
		});


		drawStone = new DrawStone();
		drawStone.setLocation(150,150);
		drawStone.setSize(500,500);

		add(drawStone);

		contentPane.revalidate();
        contentPane.repaint();
		


	}


    
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
		new P2POmokClient();
		if(args.length < 1){
			println("Please check your nickname argunemt");
			byebye();
			return;
		}

		String nickname = args[0];
		if(nickname.length() <= 0){
			println("Please check your nickname argunemt");
			byebye();
			return;

		}
		OutputStream os; 
		InputStream is;
		int localPort;
		String localIP;

        try {
            tcpconn = new Socket(serverName, serverPort); // create TCP Socket
            Runtime.getRuntime().addShutdownHook(new ByeByeThread(tcpconn)); // shutdown hook for graceful exit

			os = tcpconn.getOutputStream(); // output stream
			is = tcpconn.getInputStream(); // input stream

			InetAddress localAddr = tcpconn.getLocalAddress();
			localPort = tcpconn.getLocalPort();
			localIP = tcpconn.getLocalAddress().toString();

			// println("Client is running on "+tcpconn.getLocalAddress()+":"+tcpconn.getLocalPort());

			sendPacketTCP(os, nickname);

			String res = readPacketTCP(is);
			String[] responseArr = res.split("\\|");


			printf("welcome %s to p2p-omok server at %s. \n",nickname, tcpconn.getRemoteSocketAddress().toString());
			
			if(responseArr[0].equals("waiting")){
				printf("waiting for an opponent. \n");

				res = readPacketTCP(is);

				responseArr = res.split("\\|");

				if(responseArr[0].equals("matched")){
					printf("%s joined (%s). you play first. \n", responseArr[1], responseArr[2]);
					targetServer = responseArr[2];
					myTurn = true;
					isFirst = true;
					opponent = new clientclient(responseArr[1], responseArr[1]);

				}
			}
			else if(responseArr[0].equals("success")){
				//second client
				printf("%s is waiting for you (%s). \n", responseArr[1], responseArr[2]);
				printf("%s plays first. \n", responseArr[1]);
				targetServer = responseArr[2];

				myTurn = false;
				isFirst = false;

				opponent = new clientclient(responseArr[1], responseArr[2]);
			}

			tcpconn.close();


			try{
				conn = new DatagramSocket(localPort);
				
			} catch(SocketException e){
				println("please check your udp socket : "+localPort);
				byebye();
				return;
			}
	
			// printf("Client is running on %s \n",conn.getLocalAddress());
	
			omok = new Omok();
			omok.initOmok();
			omok.printOmok();
	
			if(myTurn){
				task = new TimerTask(){
					@Override
					public void run(){
						sendPacket(conn, "4|", targetServer.split(":")[0], Integer.parseInt(targetServer.split(":")[1]));
					}
				};
				timer.schedule(task, 10000);

			}
	
			HandlePacket receiver = new HandlePacket(conn);
			receiver.start();
	
			HandleInput sender = new HandleInput(conn, os, is, targetServer);
			sender.start();

        } catch (UnknownHostException e) {
            //server not connected
            println("Please check your server is running");
            // System.exit(0);
        } catch (IOException e) {
            //server not connected
            println("Please check your server is running");
            // System.exit(0);

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

		int route = Integer.parseInt(str.split("\\|")[0]);
		String[] msgArr = str.split("\\|");

		switch (route) {
			case 0:
				if (isMine) {
					//my chatting does not need to display!
					printf("%s\n", msgArr[1]);

					return;
				} else {
					printf("%s> %s\n", opponent.nickname, msgArr[1]);
				}
				break;

			case 1:
				int x = Integer.parseInt(msgArr[1].split(" ")[0]);
				int y = Integer.parseInt(msgArr[1].split(" ")[1]);
				if (isMine) {
					timer.cancel();
					if (isFirst) {
						omok.placeStone(x, y, 1);
					} else {
						omok.placeStone(x, y, 2);
					}
					myTurn = false;

				} else {

					task = new TimerTask(){
						@Override
						public void run(){
							sendPacket(conn, "4|", targetServer.split(":")[0], Integer.parseInt(targetServer.split(":")[1]));
						}
					};
					timer = new Timer();
					timer.schedule(task, 10000);

					if (isFirst) {
						omok.placeStone(x, y, 2);
					} else {
						omok.placeStone(x, y, 1);
					}
					myTurn = true;
				}

				omok.printOmok();

				int checkwin = omok.checkWin();

				if (checkwin == 1) {
					if (isFirst) {
						printf("you win\n");
					} else {
						printf("you lose\n");
					}
					omok.isDone = true;
				} else if (checkwin == 2) {
					if (isFirst) {
						printf("you lose\n");
					} else {
						printf("you win\n");
					}
					omok.isDone = true;
				}

		break;
		case 2:
			if (isMine) {
				//if i gave up the game
				printf("you lose\n");
			} else {
				//if opponent give up the game
				printf("%s give up this game ! \n", opponent.nickname);
				printf("you win\n");
			}
			omok.isDone = true;
			break;
		case 3:
			if (!omok.isDone) {
				if (isMine) {
					//if i gave up the game
					printf("you lose\n");
					// byebye();
					System.exit(0);
				} else {
					//if opponent give up the game
					printf("%s left the game ! \n", opponent.nickname);
					printf("you win\n");
				}
				omok.isDone = true;
			} else {
				if (isMine) {
					//if i gave up the game
					// byebye();
					System.exit(0);
				} else {
					//if opponent give up the game
					printf("%s left the game ! \n", opponent.nickname);
				}
			}

			break;

		case 4:
			if (!omok.isDone) {
				if (isMine) {
					//if i gave up the game
					printf("time out..(10s)\n");
					printf("you lose\n");

				} else {
					//if opponent give up the game
					printf("%s time out..(10s) \n", opponent.nickname);
					printf("you win\n");
				}
				omok.isDone = true;
			}
			break;
		}
	}

    public static void byebye() {
        //print bye bye
        println("Bye bye~");
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

            println("disconnected by server");
            System.exit(0);

            return "";
        }

    }

    public static void sendPacket(DatagramSocket conn, String requestString, String ip, int port) {
        // use datagram Packet to send packet to server

        try {
            InetAddress serverIP = InetAddress.getByName(ip);

            DatagramPacket dp = new DatagramPacket(requestString.getBytes(), requestString.getBytes().length, serverIP,
                    port);

            conn.send(dp);
			emit(conn, requestString, true);

			// println("send packet : "+requestString+" / "+ip+ " : "+port);

			

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

			println("read Packet : "+res);
            return res;

        } catch (IOException e) {

            return null;
        }

    }

	static class HandlePacket extends Thread {
		DatagramSocket conn;

        HandlePacket(DatagramSocket conn) {
            this.conn = conn;

        }

        public void run() {
            while (true) {
                try {

                    DatagramPacket buffer = new DatagramPacket(new byte[1024], 1024);
                    this.conn.receive(buffer);// read packet into buffer

                    // printf("UDP message from %s\n", buffer.getAddress() + ":" + buffer.getPort()); // print
                    //                                                                                           // ip,
                    //                                                                                           // port

                    InetAddress clientAddress = buffer.getAddress(); // get sender ip
                    int clientPort = buffer.getPort(); // get sender port

                    emit(conn, new String(buffer.getData(),0 , buffer.getLength()), false);// handle client msg using client
                                                                                       // address and port
                } catch (IOException e) {
                    println("client is disconnected");
                    break;
                }
            }

        }
	}

	static class HandleInput extends Thread {
        //HandleInput thread
        static DatagramSocket conn = null;
        static OutputStream os;
        static InputStream is;
		static String targetServer;
		static String targetIP;
		static String targetPort;

        HandleInput(DatagramSocket conn2, OutputStream os, InputStream is, String targetServer) {
            this.conn = conn2;
            this.os = os;
            this.is = is;

			this.targetIP = targetServer.split(":")[0];
			this.targetPort = targetServer.split(":")[1];
        }

        public void run() {
            while (true) {
                try {
                    //infinite loop
                    Scanner sc = new Scanner(System.in);
                    String inputstr = sc.nextLine();
					if(inputstr.length() >= 1){
						processMyMessage(inputstr);
					}
                } catch (NoSuchElementException e) {
                    System.exit(0);
                }

            }

        }

		public static void processMyMessage(String inputstr){
			/*

		define route

		0|message    normal chat
		1|x y		 place stone
		2|gg
		3|exit

			*/
			if(conn == null){
				return;
			}
			if (inputstr.length() >= 2 && inputstr.substring(0, 2).equals("\\\\")) {
				//place stone
				if (omok.isDone) {
					printf("This game is already done. \n");
					return;
				}

				if (omok.checkWin() != 0) {
					printf("This game is already done. \n");
					return;
				}

				// inputstr = inputstr[:len(inputstr)-1];
				if (!myTurn) {
					printf("Not your turn \n");
					return;
				}
				
				if(inputstr.split("\\\\\\\\").length < 2){
					printf("Invalid command \n");
					return;
				}
				String temp = inputstr.split("\\\\\\\\")[1];
				String[] tempArr = temp.split(" ");

				int x,y;
				x = -1;
				y = -1;

				int xyCount = 0;
				for (String v : tempArr) {
					if (v.matches("-?\\d+")) {
						if (xyCount == 0) {
							x = Integer.parseInt(v);
							xyCount++;
						} else if (xyCount == 1) {
							y = Integer.parseInt(v);
							xyCount++;
							break;
						}
					}

				}

				if (x == -1 || y == -1) {
					printf("Invalid Command \n");
					return;
				}

				if (omok.canPlace(x, y)) {
					sendPacket(conn, "1|"+x+" "+y, targetIP, Integer.parseInt(targetPort));
				} else {
					printf("Invalid move! \n");
				}

			} else if (inputstr.substring(0, 1).equals("\\")) {
				//command

				// inputstr = inputstr[:len(inputstr)-1]

				String command = inputstr.split("\\\\")[1];

				if(command.equals("gg")){
					if (omok.isDone) {
						printf("This game is already done. \n");
						return;
					}

					if (omok.checkWin() != 0) {
						printf("This game is already done. \n");
						return;
					}

					sendPacket(conn, "2|", targetIP, Integer.parseInt(targetPort));
				} else if(command.equals("exit")){
					sendPacket(conn, "3|", targetIP, Integer.parseInt(targetPort));

				} else {
					printf("Invalid Command \n");
					return;
				}
			} else {
				//normal chat
					sendPacket(conn, "0|"+inputstr, targetIP, Integer.parseInt(targetPort));
			}
		}
    }

	static class Omok {
		int[][] omok;
		boolean isDone;
		Omok(){
			this.omok = new int[10][10];
			this.initOmok();
			this.isDone = false;
		}
	
		void initOmok(){
			for(int i = 0; i<10; i++){
				for( int j = 0; j<10; j++){
					this.omok[i][j]=0;
				}
			}
		}
	
		void printOmok(){
			printf("  y 0 1 2 3 4 5 6 7 8 9  \n");
			printf("x -----------------------\n");
			for (int i = 0; i<10; i++) {
				printf("%d |", i);
		
				for (int j = 0; j<10; j++) {
					int val = omok[i][j];
					if (val == 0) {
						printf(" +");
					} else if (val == 1) {
						printf(" O");
					} else if (val == 2) {
						printf(" @");
					}
				}
		
				printf(" |\n");
			}
			printf("  -----------------------\n");
		}
		void placeStone(int x, int y, int num){
			if (omok[x][y] != 0) {
				printf("Unknown err (%d, %d) is already %d\n", x, y, omok[x][y]);
				return;
			}
			omok[x][y] = num;
			
			drawStone.revalidate();
			drawStone.repaint();
		}
		int checkWin(){
			int lastChecked = 0;
			int checkSum = 0;
		
			for(int i = 0; i<10; i++){
				lastChecked = 0;
				checkSum = 0;
	
				for( int j = 0; j<10; j++){
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
		
			for(int i = 0; i<10; i++){
				lastChecked = 0;
				checkSum = 0;
	
				for( int j = 0; j<10; j++){
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
}

