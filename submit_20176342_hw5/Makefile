objects = ChatTCPClient ChatTCPServer
all: $(objects)

$(objects): %: %.c
	$(CC) $(CFLAGS) -o $@ $<

clean:
	$(RM) $(objects)