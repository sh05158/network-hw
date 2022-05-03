#define _GNU_SOURCE
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <errno.h>
#include <netdb.h>
#include <unistd.h>
#include <signal.h>
#include <assert.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <netdb.h>

int main(int argc, char ** argv){

//get port
//int port = atoi(argv[1]);
    int server_port = 9034;
    char * name ="test";
    int namelength = strlen(name);



//set up server adress and socket
    int sock = socket(AF_INET,SOCK_STREAM,0);
    struct sockaddr_in server;
    memset(&server, 0, sizeof(server));
    server.sin_family = AF_INET;
    server.sin_addr.s_addr = inet_addr("127.0.0.1");
    server.sin_port = htons(server_port);


//connect
    if (connect(sock, (struct sockaddr *)&server, sizeof(server)) < 0) {
        perror("connect failed");
        exit(1);
    }
//set up client name

    char * buff = malloc(5000*sizeof(char));
//get the chatting
//char * other_message = malloc(5000*sizeof(char));
    while(1){
        printf("ENTER MESSAGE:\n");
        char message[5000];
        strcpy(message, name);
        strcat(message,": ");
        printf("%s", message);

        scanf("%[^\n]",buff);
        getchar();
        strcat(message,buff);

        int sent = send(sock , message , strlen(message) , MSG_DONTWAIT );
        if (sent == -1)
            perror("Send error: ");
        else
            printf("Sent bytes: %d\n", sent);

    }

    return 0;