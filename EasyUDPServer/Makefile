JFLAGS = -g
JC = javac
OPT = -O3
#OPT = -g
WARN = -Wall

EasyUDPServer: EasyUDPServer.class

EasyUDPServer.class: EasyUDPServer.java
	$(JC) $(JFLAGS) EasyUDPServer.java

clean:
	$(RM) EasyUDPServer.class

clobber:
	$(RM) EasyUDPServer.class