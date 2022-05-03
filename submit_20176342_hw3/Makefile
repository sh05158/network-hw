objects = EasyTCPClient MultiClientTCPServer
all: $(objects)

$(objects): %: %.c
	$(CC) $(CFLAGS) -o $@ $<

clean:
	$(RM) $(objects)