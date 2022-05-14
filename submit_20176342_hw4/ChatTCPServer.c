/**
 * 20176342 Song Min Joon
 * ChatTCPServer.c
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

#define MAX_USER     8

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
struct client{ //socket struct
    char nickname[100];
    char ip[30];
    char port[10];

    int soc;
    int uniqueID;
    int nickChecked;

    struct sockaddr_in addr;

} socket_default = {"","","",-1,-1,0};


struct client clientMap[1000]; //clientMap which have user socket information

void printClientNum(); //function that print how many user in this server
void broadCastToAll(int,char[]); //function that broadcast all user in this server
void byebye(); //function that print byebye
void handleError();
void handleMsg(struct client *soc, char[]); //which handle client's message
void sendPacket(struct client *soc, char[]); //function that send packet to client( certain client socket parameter )


void broadCastExceptMe(int route, char *msg, struct client *client){
    //send everyone except me with route and msg

//    printf("%s\n",msg);
    for(int i = 0; i < uniqueID; ++i){
        if(clientMap[i].soc == -1 || clientMap[i].soc == 0 || clientMap[i].uniqueID == (*client).uniqueID){
            continue;
        }
        char temp[1024] = "";
        sprintf(temp,"%d",route);
        strcat(temp,"|");
        strcat(temp,msg);
        sendPacket(&clientMap[i], temp);
    }
}



void     INThandler(int); //signal handling for graceful exit
struct client* getClientByNickname(char*);
int checkClientByNickname(char* nick);
char* getClientListString();

unsigned short serverPort=PORT;
int socNum=PORT;


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



void broadCastToAll(int route, char *msg){
    //send everyone with route and msg

//    printf("%s\n",msg);
    for(int i = 0; i < uniqueID; ++i){
        if(clientMap[i].soc == -1 || clientMap[i].soc == 0){
            continue;
        }
        char temp[1024] = "";
        sprintf(temp,"%d",route);
        strcat(temp,"|");
        strcat(temp,msg);
        sendPacket(&clientMap[i], temp);
    }
}

void handleRegisterNickname(struct client *soc, char *nick){
//    printf("handleRegsiterNickname\n");
    if( checkClientByNickname(nick) ){
        char temp[100] = "";
        temp[0]='\0';
        strcat(temp,"duplicated");
        sendPacket(soc, temp);
        close((*soc).soc);
        //close
    }
    else {
        char temp[100] = "";
        temp[0]='\0';
        strcat(temp,"success");
        (*soc).nickChecked = 1;
//        (*soc)->nickname = nick;
        strcpy(soc->nickname, nick);

        sendPacket(soc, temp);

//        char buf[100];
//        sprintf(buf, "Client %d connected. Number of connected clients = %d", uniqueID, ++userCount); //make a message
//        broadCastToAll(1, buf);//send to everyone

        char temp2[100];
        temp2[0]='\0';

        strcpy(temp2,(soc)->ip);
        strcat(temp2,":");
        strcat(temp2,(soc)->port);

        char temp3[100];
        temp3[0]='\0';

        sprintf(temp3,"[Welcome %s to CAU network class chat room at %s.]", (soc)->nickname, temp2);

        broadCastToAll(1, temp3);

        char temp4[100];
        temp4[0]='\0';

        sprintf(temp4,"[There are %d users connected.]", ++userCount);

        broadCastToAll(1, temp4);
    }
}

void handleMsg(struct client *soc, char *msg){

    char msgCopied[500];

    strcpy(msgCopied, msg);

    char **clientMsg = str_split(msg, '|');

    int requestOption = atoi(*(clientMsg));

//    char *requestString = *(clientMsg + 1);


    //client message handling

    char *s = NULL;
    char *t = NULL;
    char buffer[500] = "";
    struct timespec curr;
    char temp[500] = "";
    temp[0]='\0';
    buffer[0]='\0';
    totalRequests++;





    t = msgCopied;
    while (*t) {
        *t = toupper((unsigned char) *t);
        t++;
    }


    char *ret;

    ret = strstr(msgCopied, "I HATE PROFESSOR");
    if (ret){
        close((*soc).soc);
        char cMSG[500];
        cMSG[0]='\0';

        sprintf(cMSG, "%s is disconnected. There are %d users in the chat room", (*soc).nickname, --userCount);

        (*soc).soc = -1;

        broadCastToAll(6, cMSG);
        return;
    }

    switch(requestOption){
        case 1:
            //case 1
        {
            char *msg = *(clientMsg + 1);
            char formatMsg[500];

            strcpy(formatMsg,(*soc).nickname);

            strcat(formatMsg, "> ");
            strcat(formatMsg, msg);

            broadCastExceptMe(0, formatMsg, soc);

            break;
        }
        case 2:{
            //case 2
            int commandOption = atoi(*(clientMsg+1));

            switch(commandOption){
                case 1:{

                    strcpy(temp, "1|");
                    strcat(temp,getClientListString());
                    sendPacket(soc, temp);

                    break;
                }
                case 2:{

                    char *dmTarget = *(clientMsg + 2);

                    char *dmMessage = *(clientMsg + 3);

                    if( checkClientByNickname(dmTarget) ){
                        strcpy(temp, "2|");
                        strcat(temp,(*soc).nickname);
                        strcat(temp,"|");

                        if(dmMessage != NULL){
                            strcat(temp,dmMessage);
                        }

                        sendPacket(getClientByNickname(dmTarget ), temp );
                    }
                    break;
                }
                case 3:{

                    strcpy(temp, "3|");
                    sendPacket(soc, temp);

                    break;
                }
                case 4:{
                    strcpy(temp, "4|");
                    strcat(temp,"Chat TCP Version 0.1");
                    sendPacket(soc, temp);
                    break;
                }
                case 5:{
                    strcpy(temp, "5|");
                    sendPacket(soc, temp);
                    break;
                }
                default:
                    break;
            }
        }
        default:
            break;
    }

}

void byebye() {
    //byebye function
    printf("Bye bye~\n");
    exit(1);
}

void sendPacket(struct client *soc, char *msg){
    //send message to certain socket
    char packet[strlen(msg)+1];
    sprintf(packet,"%s",msg);

    printf("send Packet %s => %s \n",(*soc).nickname, msg);

    write((*soc).soc ,packet, strlen(msg));

}

void  INThandler(int sig)
{
    //signal handler
    signal(sig, SIG_IGN);
    byebye();
}




char* getClientListString(){
    static char returnString[1024] = "";
    returnString[0] = '\0';

    for(int i = 0; i < uniqueID; ++i){
        if(clientMap[i].soc == -1 || clientMap[i].soc == 0){
            continue;
        }
        strcat(returnString,"\n<");
        strcat(returnString,clientMap[i].nickname);
        strcat(returnString,", ");
        strcat(returnString,clientMap[i].ip);
        strcat(returnString,", ");
        strcat(returnString,clientMap[i].port);
        strcat(returnString,">");
    }

    return returnString;
}

int checkClientByNickname(char* nick){
    for(int i = 0; i < uniqueID; ++i){
        if(clientMap[i].soc == -1 || clientMap[i].soc == 0){
            continue;
        }

        if(strcmp(clientMap[i].nickname, nick) == 0){
            return 1;
        }
    }
    return 0;
}


struct client* getClientByNickname(char* nick){
    for(int i = 0; i < uniqueID; ++i){
        if(clientMap[i].soc == -1 || clientMap[i].soc == 0){
            continue;
        }

        if(strcmp(clientMap[i].nickname, nick) == 0){
            return &clientMap[i];
        }
    }
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

//    printf("ip : %s \n",ip);


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
                struct client *a =malloc(sizeof(struct client));
                (*a).soc = accept(socNum, (struct sockaddr *) &stAddr, &usocNumLen);
                (*a).addr = stAddr;
                (*a).uniqueID = uniqueID;
//                a->nickname = "";
                strcpy(a->nickname, "");
                struct client curSoc;

                curSoc = (*a);
                if (curSoc.soc < 0) {
                    perror("Accept : ");
                    continue;
                }

                if(userCount >= MAX_USER){
                    char temp[100] = "";
                    temp[0]='\0';
                    strcat(temp,"full");
                    sendPacket(a, temp);
                    close((*a).soc);
                    continue;
                }


                char *cIP = inet_ntoa(curSoc.addr.sin_addr);

                int port = (int)ntohs(curSoc.addr.sin_port);
                char p[20];


                sprintf(p,"%d",port);

                strcpy(a->ip, cIP);
                strcpy(a->port, p);
//                (*a).ip = cIP;
//
//                (*a).port = p;


                (*a).nickChecked = 0;


                strcat(cIP,":");
                strcat(cIP,p);

//                a.ip = cIP

                printf("Connection request from %s\n",cIP); //print client ip and port

//                char buf[100];
//                sprintf(buf, "Client %d connected. Number of connected clients = %d", uniqueID, ++userCount); //make a message
//                broadCastToAll(1, buf);//send to everyone

                clientMap[uniqueID - 1] = (*a); //put user into client Map


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
                        sprintf(cMSG, "%s is disconnected. There are %d users in the chat room", clientMap[i].nickname, --userCount);

                        close(clientMap[i].soc);
                        clientMap[i].soc = -1;

                        broadCastToAll(6, cMSG);
                        //broadcast to everyone him disconnected !!


                    } else {
                        //or it is normal client request


                        char res[BUF_SIZE] = "";
                        res[0]='\0';
                        strncpy(res, cMSG, serverSocket);


                        if(strcmp(clientMap[i].nickname,"")==0){
                            handleRegisterNickname(&clientMap[i], res);
                        }
                        else {
//                            char **clientMsg = str_split(res, '|');
//                            int requestOption = atoi(*(clientMsg));
//                            char *requestString = *(clientMsg + 1);

                            usleep(10000);

                            handleMsg(&clientMap[i], res);
                            //pass to handleMsg function to handle client message
                        }


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
