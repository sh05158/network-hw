/**
 * 20176342 Song Min Joon
 * ChatTCPClient.c
 **/

#include <sys/select.h>
#define __SMARTSOCK_H__

#include <stdio.h>
#include <string.h>
#include <strings.h>
// socket & bind & listen & accept & connect
#include <sys/types.h>
#include <sys/socket.h>

// sockaddr_in
#include <netinet/in.h>

// read & write
#include <unistd.h>

// htonl
#include <arpa/inet.h>

// errno, perror
#include <errno.h>

// open
#include <fcntl.h>
#include <sys/stat.h>

//select
#include <sys/time.h>
#include <sys/types.h>
#include <unistd.h>
#include <sys/select.h>
#include <stdlib.h>
#include  <signal.h>
#include <time.h>
#include <netdb.h>

#define PORT      26342
#define HOST        "nsl2.cau.ac.kr"

#define MAX_USER     3
#define NIC_NAME_SIZE   9
#define NIC_NAME_MSG   (9+2)

#define BUF_SIZE  1024              //only message
#define MSG_SIZE  (BUF_SIZE)  //message + Nick Name
#define MSG_END    "\x01\x02\x03"
//클라이언트에 필요한 헤더 정의


char *nickChecked = "false";
char nickname[100];

void byebye();//print byebye
void processOption(int, int); //receive client 1~5 option and process properly function
void sendPacket(int, char[]); //function that send a packet to server
void handlePacket(char[]);//function that handling packets from server
void printRTT(struct timespec); //function that print RTT
void processMyMessage(char[], int);
void processCommandOption(char[], char[], int);

char *substr(int s, int e, char *str){
    char *new = (char*)malloc(sizeof(char)*(e-s+2));
    strncpy(new,str+s,e-s+1);
    new[e-s+1]=0;
    return new;
}

int indexOf(char *src, char *find){
    char src2[500];
    char find2[500];
    strcpy(src2, src);
    strcpy(find2, find);

    char *ptr = strstr(src2,find2);

    if(ptr!=NULL){
        return ptr-src2;
    }
    return -1;

}

struct timespec lastRequestTime;

void     INThandler(int); //function for graceful exit
//struct tm * lastRequestTime;

char** str_split(char* a_str, const char a_delim) //Implement split functions not in c
{
    char** result    = 0;
    size_t count     = 0;
    char* tmp        = a_str;
    char* last_comma = 0;
    char delim[2];
    delim[0] = a_delim;
    delim[1] = 0;

    /* Count how many elements will be extracted. */
    while (*tmp)
    {
        if (a_delim == *tmp)
        {
            count++;
            last_comma = tmp;
        }
        tmp++;
    }

    /* Add space for trailing token. */
    count += last_comma < (a_str + strlen(a_str) - 1);

    /* Add space for terminating null string so caller
       knows where the list of returned strings ends. */
    count++;

    result = malloc(sizeof(char*) * count);

    if (result)
    {
        size_t idx  = 0;
        char* token = strtok(a_str, delim);

        while (token)
        {
            *(result + idx++) = strdup(token);
            token = strtok(0, delim);
        }
        *(result + idx) = 0;
    }

    return result;
}

int main(int argc, char *argv[] )
{
    if( argc < 2 ){
        printf("Please check your nickname argument \n");
        byebye();
        exit(1);
    }

    if ( strlen(argv[1]) == 0){
        printf("please check your nickname argument! \n");
        byebye();
        exit(1);

    }

//    nickname = "user";
//    printf("argc %d ",argc);

    strcpy(nickname, argv[1]); //copy nickname

    signal(SIGINT, INThandler); //graceful exit
    int sockNum;
    struct sockaddr_in stAddr;
    int iLen;
    char cBuf[BUF_SIZE];
    fd_set fdRead;
    char cMSG[MSG_SIZE];
    unsigned short usPORT=PORT;
    //basic variables


    sockNum = 26342;

    /*** socket ***/

    sockNum = socket(AF_INET, SOCK_STREAM, 0);
    if(-1 == sockNum)
    {
        perror("socket:");
        return 100;
    }


    int sockfd;
    struct addrinfo hints, *servinfo, *p;
    struct sockaddr_in *h;
    int rv;

    memset(&hints, 0, sizeof hints);
    hints.ai_family = AF_UNSPEC; // use AF_INET6 to force IPv6
    hints.ai_socktype = SOCK_STREAM;

    char ip[100];

    if ( (rv = getaddrinfo( HOST , "26342" , &hints , &servinfo)) != 0) //hostname to ip convert
    {
        fprintf(stderr, "getaddrinfo: %s\n", gai_strerror(rv));
        return 1;
    }

    // loop through all the results and connect to the first we can
    for(p = servinfo; p != NULL; p = p->ai_next)
    {
        h = (struct sockaddr_in *) p->ai_addr;
        strcpy(ip , inet_ntoa( h->sin_addr ) );
    }

    freeaddrinfo(servinfo); // all done with this structure

//    printf("ip : %s \n",ip);



    hints.ai_family = AF_UNSPEC; // will get all the results regardless of IPv4 and IPv6
    hints.ai_socktype = SOCK_STREAM; // TCP stream sockets


    /** structure setting **/
    stAddr.sin_family = AF_INET;
    stAddr.sin_addr.s_addr = inet_addr(ip);
    stAddr.sin_port = htons(usPORT);

    iLen = sizeof(struct sockaddr_in);

    /** connect **/
    int iRtn = connect(sockNum, (struct sockaddr *)&stAddr, iLen);
    if(-1 == iRtn)
    {
        perror("connect:");
        close(sockNum);
        return 200;
    }

    sendPacket(sockNum, nickname);

    while(1)
    {
        FD_ZERO(&fdRead);
        FD_SET(0, &fdRead);
        FD_SET(sockNum, &fdRead);
        select(sockNum+1,&fdRead, 0, 0, 0);
        if(0 != (FD_ISSET(0, &fdRead) ) )
        {
            if(strcmp(nickChecked,"false")==0){
                continue;
            }


            char input[200] = "";


            do{
                iRtn = read(0, cBuf, BUF_SIZE);

            }while(iRtn == 0);

            cBuf[iRtn - 1] = 0;

            if(strlen(cBuf) < 1){
                continue;
            }
            processMyMessage(cBuf, sockNum);

        }
        if(0 != (FD_ISSET(sockNum, &fdRead) ))
        {
            iRtn = read(sockNum, cMSG, sizeof(cMSG));
            if(iRtn == 0)
            {
                //server connection is disconnected
                printf("Server does not respond\n");
                byebye();
                exit(1);
                break;
            }
            if( 0 == strncmp(cMSG, MSG_END, sizeof(MSG_END) ))
            {
                //client side exit process
                byebye();
                exit(1);
                break;

            }

            cBuf[iRtn - 1] = 0;

            char res[BUF_SIZE] = "";

            strncpy(res,cMSG,iRtn);

            if(strcmp(nickChecked,"true") == 0){
                if(strlen(res) != 0){
                    //input message is handled normally
                    handlePacket(res);
                }
            }
            else {
                if(strcmp(res,"duplicated") == 0){
                    printf("your nickname %s is duplicated. please use another nickname\n", nickname);
                    byebye();
                    exit(1);
                }
                else if(strcmp(res,"full") == 0){
                    printf("chatting room full. cannot connect\n");
                    byebye();
                    exit(1);
                }
                else {
                    nickChecked = "true";
                }
            }


//            res[0]='\0';


        }
    }



    /*** read & write ***/

    close(sockNum);
    return 0;
}
void processMyMessage(char* inputstr, int soc){
//    clock_gettime(CLOCK_MONOTONIC_RAW, &lastRequestTime);
    char requestString[BUF_SIZE] = "";
    char *arguments = "";
    if(inputstr[0] == '\\'){
        char command[30];

        char tempStr[300];
        strcpy(tempStr, inputstr);

        char** tokens = str_split(tempStr,'\\');

        if(*(tokens) == NULL){
            printf("Invalid Command \n");
            return;
        }

        if(inputstr[1] == ' '){
            printf("Invalid Command \n");
            return;
        }

        if(indexOf(inputstr," ")!=-1){
            //command which have space

            char tempCommand[500];

            strcpy(tempStr, inputstr);
            char** tokens2 = str_split(tempStr,' ');

            strcpy(tempCommand, *(tokens2));

            char** tokens3 = str_split(tempCommand,'\\');
//            command = *(tokens3);
            strcpy(command,*(tokens3));

            arguments = substr(indexOf(inputstr, " ")+1,strlen(inputstr), inputstr);

        }
        else {
//            command = *(tokens);
            strcpy(command,*(tokens));
        }

        processCommandOption(command, arguments, soc);

    }
    else {
        strcat(requestString , "1|");
        strcat(requestString , inputstr);

        sendPacket(soc, requestString);
    }
}
void processCommandOption(char* command, char* arguments, int soc){
    char requestString[500] = "2|";

    if(strcmp(command, "list") == 0){
        //if user command is list
        strcat(requestString, "1|");
        sendPacket(soc, requestString);

    } else if(strcmp(command, "dm") == 0){
        //if user command is dm
        int idx = indexOf(arguments, " ");
        if(idx != -1){
            char *toNickname = substr(0,idx-1,arguments);
            char *toMessage = substr(idx+1,strlen(arguments),arguments);

            char *commandCheck = strstr(toNickname,"\\");

            if(commandCheck){
                printf("Invalid Command \n");
                return;
            }

            strcat(requestString, "2|");
            strcat(requestString, toNickname);
            strcat(requestString, "|");
            strcat(requestString, toMessage);

            sendPacket(soc, requestString);

        }
        else {
            printf("Invalid Command \n");
        }

    } else if(strcmp(command, "exit") == 0){
        //if user command is exit
        strcat(requestString, "3|");
        sendPacket(soc, requestString);

        printf("connection is closed by server\n");
        byebye();
        close(soc);
        exit(1);

    } else if(strcmp(command, "ver") == 0){
        //if user command is ver
        strcat(requestString, "4|");
        sendPacket(soc, requestString);

    } else if(strcmp(command, "rtt") == 0){
        clock_gettime(CLOCK_MONOTONIC_RAW, &lastRequestTime);

        //if user command is rtt
        strcat(requestString, "5|");
        sendPacket(soc, requestString);

    } else {
        printf("Invalid Command \n");
    }

}

void printRTT(struct timespec lastTime){
    //function for print RTT time
    struct timespec curr;

    clock_gettime(CLOCK_MONOTONIC_RAW, &curr);

    uint64_t delta_us = (curr.tv_sec - lastRequestTime.tv_sec) * 1000000 + (curr.tv_nsec - lastRequestTime.tv_nsec) / 1000000;
    //calculate difference between two times
    printf("RTT = %lums   \n\n",delta_us);
}

void sendPacket(int socID, char buffer[]) {
    //send packet to server

    char cMSG[strlen(buffer)];

    sprintf(cMSG, "%s", buffer);

    write(socID, cMSG, strlen(buffer));
}

void byebye() {
    //print byebye
    printf("Bye bye~\n");
}

void handlePacket(char msg[]){
//    printf("server raw message : %s\n",msg);

    //handle server packet


    /*
        server packet form

        (Route|Option Number|Message) client request packet
        (Route|Message) server one sided packet

        Route 0 => normal message 0|nickname|message
        Route 1 => list   just print server message 1|message
        Route 2 => dm     message 2|nickname|message
        Route 3 => disconnect  message 3|
        Route 4 => show version 4|message(Version)
        Route 5 => show rtt 5|message(rtt)
        Route 6 => show message 6|message(disconnected)

    */
    char** tokens = str_split(msg,'|');
    int route = atoi(*(tokens));
    int opt;

    switch(route){
        case 0:// route 0 is normal message packet
            printf("%s\n", *(tokens+1) );

            break;
        case 2:
            //print dm (from, message)
            printf("from: %s> %s\n", *(tokens+1), *(tokens+2) );
            break;
        case 3:
            //disconnect, do nothing
//            printf("%s\n", *(tokens+1) );
            break;
        case 4:
            //print version string
            printf("%s\n", *(tokens+1) );
            break;
        case 5://
//            //print RTT
            printRTT(lastRequestTime);
            break;
        case 1:// route 1 is packet that is not related to client request (server one-sided packet)
            //print list of users< nickname, IP, port >
            printf("%s\n", *(tokens+1) );
            break;
        case 6:
            //server another user disconnect message
            printf("%s\n", *(tokens+1) );
            break;
        default:
            break;
    }
}

void  INThandler(int sig)
{
    //signal handler for graceful exit
    signal(sig, SIG_IGN);
    byebye();
    exit(1);
}

