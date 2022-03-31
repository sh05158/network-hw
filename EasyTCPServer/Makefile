JFLAGS = -g
JC = javac
OPT = -O3
#OPT = -g
WARN = -Wall

EasyTCPServer: EasyTCPServer.class

EasyTCPServer.class: EasyTCPServer.java
	$(JC) $(JFLAGS) EasyTCPServer.java

clean:
	$(RM) EasyTCPServer.class

clobber:
	$(RM) EasyTCPServer.class