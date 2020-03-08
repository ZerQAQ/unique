#include<bits/stdc++.h>

using namespace std;

void copyString(char* d, const char* s, int end = 0){
    memcpy(d, s, strlen(s) * sizeof(char));
    if(end){
        d[strlen(s) + 1] = '\0';
    }
}

int main(){
    const char *protocal = "HTML/1.1";
    char head[100];
    copyString(head, protocal);
    const char *temp = " 200 ok\n\n";
    copyString(head + strlen(protocal), temp, 1);
    printf("%s\n", head);
}