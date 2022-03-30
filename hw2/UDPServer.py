#
# SimpleEchoUDPServer.py
#

from socket import *

serverPort = 26342
serverSocket = socket(AF_INET, SOCK_DGRAM)
serverSocket.bind(('', serverPort))

print("Server is ready to receive on port", serverPort)

while True:
    message, clientAddress = serverSocket.recvfrom(2048)
    print('UDP message from', clientAddress)
    modifiedMessage = message.decode().upper()
    serverSocket.sendto(modifiedMessage.encode(), clientAddress)
