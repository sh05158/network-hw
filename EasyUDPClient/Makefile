JFLAGS = -g
JC = javac
OPT = -O3
#OPT = -g
WARN = -Wall

EasyUDPClient: EasyUDPClient.class

EasyUDPClient.class: EasyUDPClient.java
	$(JC) $(JFLAGS) EasyUDPClient.java

clean:
	$(RM) EasyUDPClient.class

clobber:
	$(RM) EasyUDPClient.class