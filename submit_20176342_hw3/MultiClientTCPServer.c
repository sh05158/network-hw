/**
 * 20176342 Song Min Joon
 * MultiClientTCPServer.c
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
#include <ctype.h>
// open
#include <fcntl.h>
#include <sys/stat.h>

//select
#include <sys/time.h>
#include <sys/types.h>
#include <unistd.h>
#include <sys/select.h>

#define PORT      26342
#define HOST        "nsl2.cau.ac.kr"

#define MAX_USER     999

#define BUF_SIZE  1024
#define MSG_SIZE  (BUF_SIZE)
#define MSG_END    "\x01\x02\x03"

#include <unistd.h>
#include <stdlib.h>
#include <time.h>
#include  <signal.h>
#include <netdb.h>
//declare headers which needed in server


int totalRequests = 0; //total client request count
struct timespec startTime; // server running time

int uniqueID = 1; //unique ID of user
int userCount = 0; //user count
struct socket{ //socket struct
    int soc;
    int uniqueID;
    struct sockaddr_in* addr;
} socket_default = {-1,-1};


struct socket clientMap[MAX_USER+1]; //clientMap which have user socket information

void printClientNum(); //function that print how many user in this server
void broadCastToAll(int,char[]); //function that broadcast all user in this server
void byebye(); //function that print byebye
void handleError();
void handleMsg(struct socket soc, char[], int); //which handle client's message
void sendPacket(struct socket soc, char[]); //function that send packet to client( certain client socket parameter )
char* secToTimeFormat(int sec); //convert server running time to string format

void     INThandler(int); //signal handling for graceful exit

unsigned short serverPort=PORT;
int socNum=PORT;

int alarm_received = 0;


char** str_split(char* a_str, const char a_delim) //split function
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
    int s;

    signal(SIGINT, INThandler); //graceful exit
    clock_gettime(CLOCK_MONOTONIC_RAW, &startTime); //server start curr time

    struct sockaddr_in stAddr;
    socklen_t usocNumLen=sizeof(struct sockaddr);
    fd_set fdRead;

    int monitorSoc; //monitoring socket
    char cMSG[MSG_SIZE];


    socNum = socket(AF_INET, SOCK_STREAM, 0); //AF_INET = 2
    if(0 > socNum)
    {
        perror("socket : ");
        return -1;
    }

    int val = 1;

    setsockopt(socNum,SOL_SOCKET,SO_REUSEADDR,(char*)&val,sizeof val); //reuseable socket


    int sockfd;
    struct addrinfo hints, *servinfo, *p;
    struct sockaddr_in *h;
    int rv;

    memset(&hints, 0, sizeof hints);
    hints.ai_family = AF_UNSPEC; // use AF_INET6 to force IPv6
    hints.ai_socktype = SOCK_STREAM;

    char ip[100];

    if ( (rv = getaddrinfo( HOST , "26342" , &hints , &servinfo)) != 0) //convert host name to ip
    {
        fprintf(stderr, "getaddrinfo: %s\n", gai_strerror(rv));
        return 1;
    }

    // loop through all the results and connect to the first we can
    for(p = servinfo; p != NULL; p = p->ai_next)
    {
        h = (struct sockaddr_in *) p->ai_addr;
        strcpy(ip , inet_ntoa( h->sin_addr ) ); //record server ip
    }

    freeaddrinfo(servinfo); // all done with this structure

    printf("ip : %s \n",ip);


    bzero(&stAddr, sizeof(stAddr));
    stAddr.sin_family = AF_INET;
    stAddr.sin_addr.s_addr = inet_addr(ip);
    stAddr.sin_port = htons(serverPort);

    int serverSocket;
    serverSocket = bind(socNum, (struct sockaddr *)&stAddr,sizeof(stAddr));
    if(serverSocket < 0)
    {
        perror("bind : ");
        close(socNum);

        return -2;
    }
    serverSocket = listen(socNum, 999);

    if(serverSocket != 0)
    {
        perror("listen : ");
        close(socNum);

        return -3;
    }

    printf("Server is ready to receive on port %hu\n", serverPort); //print server start message

    struct itimerval it_val;  /* for setting itimer */


    if (signal(SIGALRM, (void (*)(int)) printClientNum) == SIG_ERR) { //set timer print Client num
        perror("Unable to catch SIGALRM");
        exit(1);
    }
    it_val.it_value.tv_sec =     60000/1000; //every 60 seconds
    it_val.it_value.tv_usec =    (60000*1000) % 1000000; //every 60 seconds
    it_val.it_interval = it_val.it_value;
    if (setitimer(ITIMER_REAL, &it_val, NULL) == -1) {
        perror("error calling setitimer()");
        exit(1);
    }


    while(1)
    {
        FD_ZERO(&fdRead);
        FD_SET(0, &fdRead);
        FD_SET(socNum, &fdRead);
        monitorSoc = socNum;

        for(int i = 0; i < uniqueID-1; ++i)
        {
            if(clientMap[i].soc == -1){
                continue;
            }
            FD_SET(clientMap[i].soc, &fdRead);
            if(monitorSoc < clientMap[i].soc){
                monitorSoc = clientMap[i].soc;
            }
        }
        s = select((monitorSoc+1), &fdRead, 0, 0, 0);
//
        if(s > 0) {
            //event receive
            if (0 != FD_ISSET(socNum, &fdRead)) {
                struct socket a;
                a.soc = accept(socNum, (struct sockaddr *) &stAddr, &usocNumLen);
                a.addr = &stAddr;
                a.uniqueID = uniqueID;

                struct socket curSoc;

                clientMap[uniqueID - 1] = a; //put user into client Map
                curSoc = a;

                if (curSoc.soc < 0) {
                    perror("Accept : ");
                    continue;
                }


                char *cIP = inet_ntoa(curSoc.addr->sin_addr);

                int port = (int)ntohs(curSoc.addr->sin_port);
                char p[20];



                sprintf(p,"%d",port);

                strcat(cIP,":");
                strcat(cIP,p);

                printf("Connection request from %s\n",cIP); //print client ip and port

                char buf[100];
                sprintf(buf, "Client %d connected. Number of connected clients = %d", uniqueID, ++userCount); //make a message
                broadCastToAll(1, buf);//send to everyone
                ++uniqueID;//add 1 to unique ID

            }
            for (int i = 0; i < uniqueID - 1; ++i) {
                if (clientMap[i].soc == -1) {
                    //skip already disconnected client
                    continue;
                }
                if (0 != FD_ISSET(clientMap[i].soc, &fdRead)) {
                    serverSocket = read(clientMap[i].soc, cMSG, MSG_SIZE);

                    if ((0 == strncmp(cMSG, MSG_END, sizeof(MSG_END))) || 0 == serverSocket) {
                        //if client sent MSG_END that client was disconnected
                        sprintf(cMSG, "Client %d disconnected. Number of connected clients = %d", clientMap[i].uniqueID,
                                --userCount);

                        close(clientMap[i].soc);
                        clientMap[i].soc = -1;

                        broadCastToAll(1, cMSG);
                        //broadcast to everyone him disconnected !!


                    } else {
                        //or it is normal client request
                        char res[BUF_SIZE] = "";

                        strncpy(res, cMSG, serverSocket);

                        char **clientMsg = str_split(res, '|');
                        int requestOption = atoi(*(clientMsg));
                        char *requestString = *(clientMsg + 1);

                        usleep(10000);

                        handleMsg(clientMap[i], requestString, requestOption);
                        //pass to handleMsg function to handle client message
                    }

                }
            }
        } else if(s == 0){
            //timer fired
            continue;
        } else {
            if (errno == EINTR) {
                /* We've been interrupted by another signal, and it might be
                 * because of the alarm(3) (using the SIGALRM) or any other
                 * signal we have received externally. */
                continue;
            }
        }








    }
    for(int uiCnt = 0; uiCnt < uniqueID; ++uiCnt)
    {
        close(clientMap[uiCnt].soc);
    }
    close(socNum);
    return 0;
}

void printClientNum(){
    //print client connection number
    printf("Number of connected clients = %d\n",userCount);
}

void broadCastToAll(int route, char *msg){
    //send everyone with route and msg

    printf("%s\n",msg);
    for(int i = 0; i < uniqueID; ++i){
        if(clientMap[i].soc == -1 || clientMap[i].soc == 0){
            continue;
        }
        char temp[1024] = "";
        sprintf(temp,"%d",route);
        strcat(temp,"|");
        strcat(temp,msg);
        sendPacket(clientMap[i], temp);
    }
}

void handleMsg(struct socket soc, char *msg, int opt){
    //client message handling
    printf("Command %d\n\n",opt);
    char *s = NULL;
    char buffer[100] = "";
    struct timespec curr;
    char temp[100] = "";
    temp[0]='\0';
    buffer[0]='\0';
    totalRequests++;

    switch(opt){
        case 1:
            //case 1
            s = msg;
            while (*s) {
                *s = toupper((unsigned char) *s);
                s++;
            }
            strcpy(temp, "0|1|");
            strcat(temp,msg);
            sendPacket(soc, temp);
            break;
        case 2:{
            //case 2
            char *ip = inet_ntoa(soc.addr->sin_addr);

            int port = (int)ntohs(soc.addr->sin_port);
            char p[20];


            sprintf(p,"%d",port);

            strcat(ip,":");
            strcat(ip,p);

            strcpy(temp, "0|2|");
            strcat(temp,ip);

            sendPacket(soc, temp);
            break;
        }


        case 3:
            //case 3
            sprintf(buffer,"%d",totalRequests);
            strcpy(temp, "0|3|");
            strcat(temp,buffer);
            sendPacket(soc, temp);
            break;

        case 4: {
            //case 4
            char temp2[100] = "";
            clock_gettime(CLOCK_MONOTONIC_RAW, &curr);

            uint64_t delta_us = (curr.tv_sec - startTime.tv_sec) * 1000000 + (curr.tv_nsec - startTime.tv_nsec) / 1000;
            int time = (int) (delta_us / 1000000);

            strcpy(temp2, "0|4|");
            strcat(temp2, secToTimeFormat(time));

            sendPacket(soc, temp2);
            break;
        }
        case 5:
            //case 5
            close(soc.soc);
        default:
            //default
            close(soc.soc);
    }
}

void byebye() {
    //byebye function
    printf("Bye bye~\n");
    exit(1);
}

char* secToTimeFormat(int sec) {
    // convert Date.now(millisecond) format to HH:MM:ss
//    int ss = (int) Math.floor(ms / 1000);

    static char res[100] = "";
    res[0] = '\0';

    int ss = sec;

    int h = ss / 60 / 60;
    ss = ss - (h * 60 * 60);

    int m = ss / 60;
    ss = ss - (m * 60);

    int s = ss;

    if (h >= 1) {
        if (h < 10) {
            strcat(res, "0");
            char temp[100] = "";
            sprintf(temp,"%d",h);
            strcat(res,temp);
        } else {
            char temp[100] = "";
            sprintf(temp,"%d",h);
            strcat(res,temp);
        }

        strcat(res,"h");
    }

    if (m >= 1) {
        if (m < 10) {
            strcat(res, "0");
            char temp[100] = "";
            sprintf(temp,"%d",m);
            strcat(res,temp);
        } else {
            char temp[100] = "";
            sprintf(temp,"%d",m);
            strcat(res,temp);
        }

        strcat(res,"m");
    }

    if (s >= 1) {
        if (s < 10) {
            strcat(res, "0");
            char temp[100] = "";
            sprintf(temp,"%d",s);
            strcat(res,temp);
        } else {
            char temp[100] = "";
            sprintf(temp,"%d",s);
            strcat(res,temp);
        }

        strcat(res,"s");
    }

    return res;

}

void sendPacket(struct socket soc, char *msg){
    //send message to certain socket
    char packet[strlen(msg)+1];
    sprintf(packet,"%s",msg);

    write(soc.soc ,packet, strlen(msg));

}

void  INThandler(int sig)
{
    //signal handler
    signal(sig, SIG_IGN);
    byebye();
}

