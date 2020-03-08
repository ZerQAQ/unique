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

const char* notFind = "HTTP/1.1 404 Not found\n\n404 Not found\0";
const char* STR_GET = "GET";

int split(const char *s, char (*ss)[100], const char *sp = " "){
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
                if(s[j] == '\n' || s[j] == '\r') continue;
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

void prtstr(char *s){
    int len = strlen(s);
    for(int i = 0; i < len; i++){
        printf("%d %c|", (int) s[i], s[i]);
    }
    printf("\n");
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
                _size = strlen(s);
            }
            memcpy(buf, s, _size * sizeof(char));
            size = _size;
        }
        ptr = buf;
        next = NULL;
    }
};

class bufferList{
public:
    static const int maxNoteNum = 1000;
    struct buffer *head, *tail;
    int len;

    bufferList(){
        head = tail = NULL;
        len = 0;
    }

    bool empty(){
        if(!len) return 1;
        else return 0;
    }

    void pop(){
        if(len){
            struct buffer *t = head;
            head = head->next;
            delete t;
            len--;
        }
    }

    buffer* front(){
        if(!len) return NULL;
        return head;
    }

    int push(struct buffer *buf){
        if(!len){
            tail = head = buf;
        }
        else{
            tail->next = buf;
            tail = buf;
        }
        len++;
    }

    void clear(){
        while(len){
           struct buffer *t = head;
           head = head->next;
           delete t;
           len--;
        }
        head = tail = NULL;
    }

    ~bufferList(){
        clear();
    }

    int write(FILE *f){
        if(len >= maxNoteNum) return -1; //这个判断放在外面 尽量每次都读完整的一个文件
        for(;;){
            struct buffer *newbuf = new buffer;
            int freadRet = fread(newbuf->buf, sizeof(char), buffer::maxLen, f);
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
                newbuf->size = freadRet;
                push(newbuf);
            }
        }
    }
};

void epollAdd(int epfd, int fd, void *d = NULL){
    struct epoll_event e;
    e.events = EPOLLIN;
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

struct clientSock{
    int sock;
    class bufferList bufs;

    int bufLen, bufPointer;
	static const int bufMaxLen = 1024;
    char buf[bufMaxLen];

    clientSock(int _sock){
        sock = _sock;
        bufLen = bufPointer = 0;
    }

    int readLine(char* ret){
        int retPointer = 0;
		for(;;){
            printf("bp %d bl %d\n", bufPointer, bufLen);
            if(bufPointer == bufLen){
                bufPointer = bufLen = 0;

                int readRet = read(sock, buf, bufMaxLen);
                if(readRet == -1){
                    if(errno == EINTR){
                        continue;
                    }
                    else{
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

            while(bufPointer < bufLen){
                char ch = buf[bufPointer++];
                ret[retPointer++] = ch;
                if(ch == '\n' || retPointer == bufMaxLen - 1){
                    ret[retPointer] = '\0';
                    return retPointer;
                }
            }
        }
	}
};

class TCPserver{
public:
    static const int maxClientNum = 100;
    set<struct clientSock*> clients;
    struct epoll_event eventList[maxClientNum];
    int clientNum, serverSock, epollfd;
    struct sockaddr_in sockAddr;
    char IP[20];
    int port;

    TCPserver(const char* _IP, int _port){
        copyString(IP, _IP, 1);
        port = _port;

        clientNum = 0;
        serverSock = socket(AF_INET, SOCK_STREAM, IPPROTO_TCP);

        memset(&sockAddr, 0, sizeof(struct sockaddr_in));
        sockAddr.sin_family = AF_INET;  //使用IPv4地址
        sockAddr.sin_addr.s_addr = inet_addr(_IP);  //具体的IP地址
        sockAddr.sin_port = htons(_port);  //端口

        epollfd = epoll_create(maxClientNum);
        set_nonblocking(serverSock);
        printf("Server inited.\n");
    }

    void handleAccpet(){
        struct sockaddr_storage claddr;
        socklen_t addrlen = sizeof(struct sockaddr_storage);
        for (;;) {
            int clientfd = accept(serverSock, (struct sockaddr*)&claddr, &addrlen);
            if (clientfd < 0) {
                if (errno == EINTR)
                    continue;
                perror("accept ERROR: \n");
            }
            printf("\nclient connect: (fd = %d)\n", clientfd);

            if(clientNum < maxClientNum){
                set_nonblocking(clientfd);
                struct clientSock *cs = new clientSock(clientfd);
                clients.insert(cs);
                epollAdd(epollfd, clientfd, cs);
                clientNum++;
            }
            else{
                printf("too many client\n");
            }

            break;
        }   
    }

    void respondGetRequest(struct clientSock *clientptr, const char* protocal, const char* path){
        FILE *f = fopen(path, "rb");
        if(f == NULL){
            clientptr->bufs.push(new buffer(notFind));
        }
        else{
            char head[100];
            copyString(head, protocal);
            const char *temp = " 200 ok\n\n";
            copyString(head + strlen(protocal), temp, 1);
            printf("head:%s", head);

            clientptr->bufs.push(new buffer(head));
            clientptr->bufs.write(f);
            
            fclose(f);
        }
    }

    void closeClient(struct clientSock* clientptr){
        printf("\nclient closed (fd = %d)\n", clientptr->sock);
        epollDel(epollfd, clientptr->sock);
        if(close(clientptr->sock) < 0) perror("close error:\n");
        clientNum--;
        clients.erase(clientptr);
        delete clientptr;
    }

    void handleRead(void *ptr){
        struct clientSock* clientptr = (struct clientSock*) ptr;
        printf("\nhandling read (fd = %d)\n", clientptr->sock);
        for(;;){
            char _buf[clientSock::bufMaxLen];
            int readRet = clientptr->readLine(_buf);
            printf("readRet:%d\n", readRet);
            //printf("%s\n", _buf); //magic printf
            if(readRet <= 0){
                if(readRet == 0){
                    closeClient(clientptr);
                }
                else{
                    perror("read error");
                }
                break;
            }

            if(readRet >= 3){
                char __buf[4] = {_buf[0], _buf[1], _buf[2], '\0'};
                if(cmpString(__buf, "GET")){
                    printf("%s", _buf);
                    char strings[5][100];
                    int spLen = split(_buf, strings);
                    //prtstr(strings[1]);
                    //prtstr(strings[2]);
                    respondGetRequest(clientptr, strings[2], strings[1] + 1);
                }
            }
        }
        epollWrite(epollfd, clientptr->sock, 1, clientptr);
        printf("finish read handle (fd = %d)\n", clientptr->sock);
    }

    void handleWrite(void *ptr){
        struct clientSock *clientptr = (struct clientSock *) ptr;
        //printf("\nhandling write (fd = %d)\n", clientptr->sock);
        //printf("bufslen:%d\n", clientptr->bufs.len);
        long long byteNum = 0;
        while(!clientptr->bufs.empty()){
            //printf("in write while bufslen:%d\n", clientptr->bufs.len);
            //printf("head size:%d\n", clientptr->bufs.head->size);
            struct buffer *head = clientptr->bufs.head;
            if(!(head->size)){
                clientptr->bufs.pop();
            }
            else{
                //发送数据
                int writeRet = write(clientptr->sock, head->ptr, head->size);
                //printf("wret:%d\n", writeRet);
                if(writeRet <= 0){
                    if(errno == EINTR)
                        continue;
                    else if(errno == EAGAIN || errno == EWOULDBLOCK)
                        return;
                    else{
                        perror("write error:");
                        closeClient(clientptr);
                        return;
                    }
                }
                else{
                    byteNum += writeRet;
                    head->ptr += writeRet;
                    head->size -= writeRet;
                }
            }
        }
        epollWrite(epollfd, clientptr->sock, 0, clientptr);
        printf("%lld KBs sended.\n", (byteNum >> 10));
        printf("finish write handle (fd = %d)\n", clientptr->sock);
        closeClient(clientptr);
        return;
    }

    void handleError(void *ptr){
        perror("client error: \n");
        closeClient((struct clientSock *) ptr);
    }

    void work(){
        bind(serverSock, (struct sockaddr*) &sockAddr, sizeof(sockAddr));
        listen(serverSock, 128);
        epollAdd(epollfd, serverSock);
        printf("Server is listening %s:%d\n", IP, port);

        for(;;){
            int eventNum = epoll_wait(epollfd, eventList, maxClientNum, -1);
            //printf("evn: %d\n", eventNum);
            if(eventNum <= 0){
                if(eventNum < 0 && errno != EINTR){
                    perror("epoll_wait ERROR:\n");
                }
                continue;
            }
            for(int i = 0; i < eventNum; i++){
                struct epoll_event &e = eventList[i];
                if(e.data.ptr == NULL){
                    handleAccpet();
                }
                else{
                    if(e.events & (EPOLLIN | EPOLLHUP)){
                        handleRead(e.data.ptr);
                    }
                    if(e.events & EPOLLOUT){
                        handleWrite(e.data.ptr);
                    }
                    if(e.events & EPOLLERR){
                        handleError(e.data.ptr);
                    }
                }
            }
        }
    }

    void clear(){
        for(auto ptr : clients){
            closeClient(ptr);
        }
        clients.clear();
        clientNum = 0;
    }

    ~TCPserver(){
        clear();
        close(serverSock);
    }
};

int main(){
    TCPserver s("127.0.0.1", 80);
	s.work();
    return 0;
}