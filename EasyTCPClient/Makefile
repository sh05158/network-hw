JFLAGS = -g
JC = javac
OPT = -O3
#OPT = -g
WARN = -Wall

EasyTCPClient: EasyTCPClient.class

EasyTCPClient.class: EasyTCPClient.java
	$(JC) $(JFLAGS) EasyTCPClient.java

clean:
	$(RM) EasyTCPClient.class

clobber:
	$(RM) EasyTCPClient.class