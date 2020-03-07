#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <unistd.h>
#include <arpa/inet.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <fcntl.h>
#include <bits/stdc++.h>
#include <sys/epoll.h>

using namespace std;

const char* notFind = "HTTP/1.1 404 Not found\n\n404 Not found";
const char* STR_GET = "GET";

int split(const char *s, char (*ss)[100], const char *sp = " \n"){
    int spLen = strlen(sp), sLen = strlen(s), ssPointer = 0;
    int last = 0;
    for(int i = 0; i < sLen; i++){
        bool isItv = 0;
        for(int j = 0; j < spLen; j++){
            if(s[i] == sp[j]){
                isItv = 1;
                break;
            }
        }
        if(i == sLen - 1){
            i++;
            isItv = 1;
        }
        if(isItv){
            int ssP = 0;
            for(int j = last; j < i; j++){
                ss[ssPointer][ssP++] = s[j];
            }
            ss[ssPointer++][ssP] = '\0';
            last = i + 1;
        }
    }
    return ssPointer;
}

bool cmpString(const char *a, const char *b){
    int lena = strlen(a), lenb = strlen(b);
    if(lena != lenb) return 0;
    for(int i = 0; i <= lena; i++) if(a[i] != b[i]) return 0;
    return 1;
}

void copyString(char* d, const char* s, int end = 0){
    memcpy(d, s, strlen(s) * sizeof(char));
    if(end){
        d[strlen(s) + 1] = '\0';
    }
}

void stringToCharPointer(string s, char* d){
    int len = s.length();
    for(int i = 0; i < len; i++){
        d[i] = s[i];
    }
    d[len + 1] = '\0';
}

// 设置socket为非阻塞
void set_nonblocking(int fd) {
    int flag = fcntl(fd, F_GETFL, 0);
    if (flag >= 0) {
        fcntl(fd, F_SETFL, flag | O_NONBLOCK);
    }
}

struct buffer{
    static const int maxLen = 10240;
    char buf[maxLen + 10], *ptr;
    int size;
    struct buffer* next;
    buffer(const char *s = NULL, int _size = -1){
        if(s){
            if(_size == -1){
                _size = strlen(s) + 1;
            }
            memcpy(buf, s, _size);
            size = _size;
            ptr = buf;
        }

        next = NULL;
        size = 0;
    }
};

class bufferList{
public:
    static const int maxNoteNum = 200; //防止被恶意攻击时占用太多内存
    queue<buffer> Q;

    bool empty(){
        return Q.empty();
    }

    void pop(){
        Q.pop();
    }

    buffer front(){
        return Q.front();
    }

    int push(buffer buf){
        if(Q.size() >= maxNoteNum) return -1;
        Q.push(buf);
        return 1;
    }

    int write(FILE *f){
        if(Q.size() >= maxNoteNum) return -1; //这个判断放在外面 尽量每次都读完整的一个文件
        for(;;){
            struct buffer newbuf;
            int freadRet = fread(newbuf.buf, sizeof(char), buffer::maxLen, f);
            if(freadRet <= 0){
                if(freadRet == 0){
                    return 1;
                }
                else{
                    printf("fread ERROR\n");
                    return -2;
                }
            }
            else{
                newbuf.size = freadRet;
                Q.push(newbuf);
            }
        }
    }
};

void epollAdd(int epfd, int fd, int tp = EPOLLIN, void *d = NULL){
    struct epoll_event e;
    e.events = tp;
    e.data.ptr = d;
    epoll_ctl(epfd, EPOLL_CTL_ADD, fd, &e);
}

void epollWrite(int epfd, int fd, int enabled, void *d){
    struct epoll_event e;
    e.events = EPOLLIN | (enabled ? EPOLLOUT: 0);
    e.data.ptr = d;
    epoll_ctl(epfd, EPOLL_CTL_MOD, fd, &e);
}

void epollDel(int epfd, int fd){
    epoll_ctl(epfd, EPOLL_CTL_DEL, fd, NULL);
}

class clientSocket{
public:
    int sock, epollList;
    class bufferList writeBufList;
    struct buffer writeBuf;

    int bufLen, bufPionter;
	static const int bufMaxLen = 1024;
    char buf[bufMaxLen * 2 + 10];
    class TCPserver* server;

    clientSocket(int _sock, int _epollList){
        epollList = _epollList;
        sock = _sock;
        set_nonblocking(sock);
    }

    ~clientSocket(){
        close(sock);
    }

    int readLine(char* ret){
        int retPointer = 0;
		for(;;){
            if(bufPionter == bufLen){
                bufPionter = bufLen = 0;

                int readRet = read(sock, buf, bufMaxLen);
                if(readRet == -1){
                    if(errno == EINTR){
                        continue;
                    }
                    else{
                        printf("socket read ERROR\n");
                        return -1;
                    }
                }
                else if(readRet == 0){
                    return 0;
                }
                else{
                    bufLen = readRet;
                }
            }

            while(bufPionter < bufLen){
                char ch = buf[bufPionter++];
                ret[retPointer++] = ch;
                if(ch == '\n'){
                    ret[retPointer] = '\0';
                    return retPointer;
                }
            }
        }
	}

    void respondGetRequest(const char* protocal, const char* path){
        FILE *f = fopen(path, "rb");
        if(f == NULL){
            writeBufList.push(buffer(notFind, 37));
        }
        else{
            char head[100];
            copyString(head, protocal);
            const char* temp = " 200 OK\n\n";
            copyString(head + strlen(protocal), temp);
            int headLen = strlen(protocal) + strlen(temp);
            head[headLen] = '\0';

            writeBufList.push(buffer(head, headLen));
            writeBufList.write(f);
            
            fclose(f);
        }
    }

    void headleRead(){
        printf("reading...\n");
        for(;;){
            char _buf[bufMaxLen];
            int readRet = readLine(_buf);
            printf("readRet:%d\n", readRet);
            if(readRet <= 0){
                if(readRet == 0){
                    
                }
                break;
            }
            printf("%s", _buf);

            if(readRet >= 3){
                char __buf[4] = {_buf[0], _buf[1], _buf[2], '\0'};
                if(cmpString(__buf, "GET")){
                    printf("%s", _buf);
                    char strings[5][100];
                    int spLen = split(_buf, strings);
                    respondGetRequest(strings[2], strings[1] + 1);
                }
            }
            printf("end\n");
        }

        printf("adding write event\n");
        epollWrite(epollList, sock, 1, this);
    }

    int headleWrite(){
        printf("writing...\n");
        printf("wbs:%d\n", (int)writeBuf.size);
        printf("wblLen:%d\n", (int)writeBufList.Q.size());
        while(writeBufList.empty() == 0 || writeBuf.size){
            if(writeBuf.size){
                printf("wRet writing...\n");
                int wRet = write(sock, writeBuf.ptr, writeBuf.size);
                printf("wRet:%d\n", wRet);
                if(wRet < 0){
                    return -1;
                }
                else{
                    writeBuf.size -= wRet;
                    writeBuf.ptr += wRet;
                }
            }
            else{
                printf("poping writeBufList...\n");
                writeBuf = writeBufList.front();
                writeBufList.pop();
            }
        }
        printf("deleting write events...\n");
        epollWrite(epollList, sock, 0, this);
        return 1;
    }
   
};

class TCPserver
{
public:
    static const int maxClientNum = 100;
    int serverSock, port, epollList;
    struct sockaddr_in sockAddr;
    char IP[20];
    struct epoll_event eventList[maxClientNum + 10];
    class clientSocket *clientList[maxClientNum + 10];
    int clientListLen;

    void init(){
        clientListLen = 0;
        serverSock = socket(AF_INET, SOCK_STREAM, IPPROTO_TCP);

        sockAddr.sin_family = AF_INET;  //使用IPv4地址
        sockAddr.sin_addr.s_addr = inet_addr(IP);  //具体的IP地址
        sockAddr.sin_port = htons(port);  //端口

        epollList = epoll_create(maxClientNum);
        set_nonblocking(serverSock);
    }

    TCPserver(const char *IP_, int port_){
       memcpy(IP, IP_, (strlen(IP_) + 1) * sizeof(char));
       port = port_;
       init();
    }

    void headleAccept(){
        struct sockaddr_in clnt_addr;
        socklen_t clnt_addr_size = sizeof(clnt_addr);

        int clientSock = accept(serverSock, (struct sockaddr*)&clnt_addr, &clnt_addr_size);
        class clientSocket *cs = new class clientSocket(clientSock, epollList);
        epollAdd(epollList, clientSock, EPOLLIN, cs);
        clientList[clientListLen++] = cs;

        printf("clientsocket accepted\n");
    }

    void headleError(class clientSocket *cs){
        printf("deleting clinetSocket...\n");
        if(close(cs->sock) < 0){
            printf("clientsocket close ERROR\n");
        }
        epollDel(serverSock, cs->sock);
        delete cs;
    }

    void work(){
        bind(serverSock, (struct sockaddr*) &sockAddr, sizeof(sockAddr));
        listen(serverSock, 128);
        epollAdd(epollList, serverSock, EPOLLIN);
        printf("Server is listening on %s:%d\n", IP, port);

        for(;;){
            int waitRet = epoll_wait(epollList, eventList, maxClientNum, -1);
            if(waitRet <= 0){
                if(waitRet < 0 && errno != EINTR){
                    perror("epoll_wait ERROR\n");
                }
                continue;
            }

            class clientSocket *cs;
            for(int i = 0; i < waitRet; i++){
                struct epoll_event e = eventList[i];
                if(e.data.ptr == NULL){
                    headleAccept();
                }
                else{
                    if(e.events & (EPOLLIN | EPOLLHUP)){
                        ((class clientSocket *) (e.data.ptr))->headleRead();
                    }
                    if(e.events & EPOLLOUT){
                        ((class clientSocket *) (e.data.ptr))->headleWrite();
                    }
                    if(e.events & EPOLLERR){
                        headleError((class clientSocket *) (e.data.ptr));
                    }
                }
            }
        }
    }

    ~TCPserver(){
		close(serverSock);
	}
};


int main(){
    TCPserver s("127.0.0.1", 80);
	for(;;) s.work();
    return 0;
}