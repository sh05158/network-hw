/**
 * 20176342 Song Min Joon
 * EasyTCPClient.c
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

void handleInput(int);
void byebye();//print byebye
void processOption(int, int); //receive client 1~5 option and process properly function
void sendPacket(int, char[]); //function that send a packet to server
void printOption();//print option
void handlePacket(char[]);//function that handling packets from server
void processMyMessage(int opt, char msg[]); // The function of processing packets from the server that I requested
void printRTT(struct timespec); //function that print RTT

void printDuration(char *msg); //Functions that parse and output server execution time from the server

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

int main()
{
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

    printf("ip : %s \n",ip);



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

    printOption(); //print 5 options 

    while(1)
    {
        FD_ZERO(&fdRead);
        FD_SET(0, &fdRead);
        FD_SET(sockNum, &fdRead);
        select(sockNum+1,&fdRead, 0, 0, 0);
        if(0 != (FD_ISSET(0, &fdRead) ) )
        {

            char input[200] = "";


            do{
                iRtn = read(0, cBuf, BUF_SIZE);

            }while(iRtn == 0);
            
            cBuf[iRtn - 1] = 0;
            
            processOption(sockNum, atoi(cBuf));

        }
        if(0 != (FD_ISSET(sockNum, &fdRead) ))
        {
//            printOption();
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

            if(strlen(res) != 0){
                //input message is handled normally
                handlePacket(res);
            }

//            res[0]='\0';


        }
    }



    /*** read & write ***/

    close(sockNum);
    return 0;
}

void printOption() {
    //print 5 options
    printf("<Menu>\n");
    printf("option 1) convert text to UPPER-case letters.\n");
    printf("option 2) ask the server what the IP address and port number of the client is.\n");
    printf("option 3) ask the server how many client requests(commands) it has served so far.\n");
    printf("option 4) ask the server program how long it has been running for since it started.\n");
    printf("option 5) exit client program\n");
    printf("Please select your option :\n");


}

void handleInput(int socID){
    printOption();
    printf("Please select your option :");
    int opt;
    scanf("%d",&opt);
    processOption(socID, opt);
}

void processOption(int socID, int opt) {

    clock_gettime(CLOCK_MONOTONIC_RAW, &lastRequestTime);
    //record a time for RTT Time


    char optString[10];
    char requestString[BUF_SIZE] = "";

    sprintf(optString,"%d",opt);

    switch (opt) {
        case 1:
            //Option 1
            printf("Input lowercase sentence: ");
            char input[200];


            scanf("%[^\n]%*c",input);
            fseek(stdin,0,SEEK_END);

            clock_gettime(CLOCK_MONOTONIC_RAW, &lastRequestTime);

            strcat(requestString, optString);
            strcat(requestString , "|");
            strcat(requestString,input);

            sendPacket(socID, requestString);

            break;

        case 2:
            //Option 2
            strcat(requestString, optString);
            strcat(requestString , "|");

            sendPacket(socID, requestString);

            break;


        case 3:
            //Option 3
            strcat(requestString, optString);
            strcat(requestString , "|");

            sendPacket(socID, requestString);

            break;


        case 4:
            //Option 4
            strcat(requestString, optString);
            strcat(requestString , "|");

            sendPacket(socID, requestString);

            break;


        case 5:
            //Option 5
            byebye();
            exit(1);
//            conn.Close();
//            os.Exit(0);
            break;


        default:
            //Option default
            byebye();
            exit(1);

//            conn.Close();
//            os.Exit(0);
            break;

    }
     printRTT(lastRequestTime);

}

void printRTT(struct timespec lastTime){
    //function for print RTT time
    struct timespec curr;

    clock_gettime(CLOCK_MONOTONIC_RAW, &curr);

    uint64_t delta_us = (curr.tv_sec - lastRequestTime.tv_sec) * 1000000 + (curr.tv_nsec - lastRequestTime.tv_nsec) / 1000;
    //calculate difference between two times
    printf("RTT = %lums \n\n\n",delta_us);
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

        0|1|BLAH BLAH BLAH...
        0|1|HELLO WORLD!

        0|2|127.0.0.1:6342    //my ip and port
        0|3|30 				  //count of server serves client requests
        0|4|01h30m12s		  //server running time

        1|client #3 is connected   //server broadcast message (server one-sided)
        1|client #2 is disconnected   //server broadcast message (server one-sided)

    */


    char** tokens = str_split(msg,'|');
    int route = atoi(*(tokens));
    int opt;

    switch(route){
        case 0:// route 0 is packet that related to client request
            opt = atoi( *(tokens+1) );
            processMyMessage(opt,*(tokens+2));
            printOption();

            break;
        case 1:// route 1 is packet that is not related to client request (server one-sided packet)
            printf("%s\n", *(tokens+1) );
            break;
        default:
            break;
    }
}

void processMyMessage(int opt, char msg[]) {
    //if i request server with option n, then process received packet with n
    char **temp;


    switch (opt) {
        case 1:
        //if i request option number 1
            printf("Reply from server: %s\n", msg);
            break;
        case 2:
            //if i request option number 2
            temp = str_split(msg, ':');
            char *ip = *(temp);
            char *port = *(temp+1);
            printf("Reply from server: client IP = %s, port = %s\n", ip,port);
            break;

        case 3:
        //if i request option number 3
            printf("Reply from server: requests served = %s\n", msg);
            break;

        case 4:
            //if i request option number 4
//            timeD, _ := time.ParseDuration(msg)
            printDuration(msg);
            break;
        default:
            break;

    }

//    printRTT(time.Since(lastRequestTime))
}

void printDuration(char msg[]) {
    //print server running time in proper form(HH:MM:ss)

    char *h;
    char *m;
    char *s;

    int hour = 0;
    int minute = 0;
    int second = 0;

    h = strstr(msg,"h");
    if(h){
        char **temp = str_split(msg, 'h');
        msg = *(temp+1);
        hour = atoi(*(temp));
    }

    m = strstr(msg,"m");
    if(m){
        char **temp = str_split(msg, 'm');
        msg = *(temp+1);
        minute = atoi(*(temp));
    }

    s = strstr(msg,"s");
    if(s){
        char **temp = str_split(msg, 's');
        msg = *(temp+1);
        second = atoi(*(temp));
    }

    printf("Reply from server: run time = %02d:%02d:%02d\n", hour, minute, second);
}

void  INThandler(int sig)
{
    //signal handler for graceful exit
    signal(sig, SIG_IGN);
    byebye();
    exit(1);
}

