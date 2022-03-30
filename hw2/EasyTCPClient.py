#
# SimpleEchoTCPClient.py
#

from socket import *

serverName = 'localhost'
serverPort = 26342

clientSocket = socket(AF_INET, SOCK_STREAM)
clientSocket.connect((serverName, serverPort))

print("Client is running on port", clientSocket.getsockname()[1])

message = input('Input lowercase sentence: ')

clientSocket.send(message.encode())

modifiedMessage = clientSocket.recv(2048)
print('Reply from server:', modifiedMessage.decode())

clientSocket.close()
