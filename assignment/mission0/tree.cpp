//
// Created by 10758 on 2020/3/3.
//

#include<bits/stdc++.h>

using namespace std;

template <class T>
class Stack{
public:
    T *a;
    int t, len;

    Stack(int n = 1e5){
        len = n;
        a = new T[n];
        t = -1;
    }

    ~Stack(){
        delete []a;
    }

    int empty(){
        if(t == -1) return 1;
        else return 0;
    }

    int push(T v){
        if(t == len) return -1;
        else{
            a[++t] = v;
            return t;
        }
    }

    T pop(){
        if(t == -1) return T();
        else{
            return a[t--];
        }
    }

    T top(){
        if(t == -1) return T();
        else{
            return a[t];
        }
    }

    void prt(){
        for(int i = 0; i <= t; i++){
            cout << a[i] << ' ';
        }
        cout << '\n';
    }
};

class Tree{
public:
    struct note{
        int v;
        note *ln, *rn;
    };

    note *root;
    int len;

    void build(note **o, int i, int *d){
        *o = new note;
        (*o)->v = d[i - 1];
        (*o)->ln = (*o)->rn = NULL;
        if(i * 2 <= len){
            build(&(*o)->ln, i * 2, d);
        }
        if(i * 2 + 1 <= len){
            build(&(*o)->rn, i * 2 + 1, d);
        }
    }

    Tree(int *d, int l){
        len = l;
        build(&root, 1, d);
    }

    void Rfir(note *o){
        cout << o->v << ' ';
        if(o->ln != NULL) Rfir(o->ln);
        if(o->rn != NULL) Rfir(o->rn);
    }

    void _Rfir(){
        Stack<note*> S;
        S.push(root);
        while(!S.empty()){
            note *o = S.pop();
            cout << o->v << ' ';
            if(o->rn != NULL){
                S.push(o->rn);
            }
            if(o->ln != NULL){
                S.push(o->ln);
            }
        }
    }

    void rootFir(){
        Rfir(root); cout << '\n';
        _Rfir(); cout << '\n';
    }

    void Rmid(note *o){
        if(o->ln != NULL) Rmid(o->ln);
        cout << o->v << ' ';
        if(o->rn != NULL) Rmid(o->rn);
    }

    void _Rmid(){
        Stack<note*> S;
        note *o = root;
        while(o){
            S.push(o);
            o = o->ln;
        }
        while(!S.empty()){
            note *o = S.pop();
            cout << o->v << ' ';
            if(o->rn != NULL){
                o = o->rn;
                while(o){
                    S.push(o);
                    o = o->ln;
                }
            }
        }
    }

    void rootmid(){
        Rmid(root); cout << endl;
        _Rmid(); cout << endl;
    }

    void Rlas(note *o){
        if(o->ln != NULL) Rlas(o->ln);
        if(o->rn != NULL) Rlas(o->rn);
        cout << o->v << ' ';
    }

    void _Rlas(){
        bool vis[len];
        for(int i = 0; i < len; i++) vis[i] = 0;
        Stack<note*> S;
        note *o = root;
        while(o){
            S.push(o);
            o = o->ln;
        }
        while(!S.empty()){
            note *o = S.top();
            if(o->rn && !vis[o->rn->v]){
                o = o->rn;
                while(o){
                    S.push(o);
                    o = o->ln;
                }
            }
            else{
                cout << o->v << ' ';
                vis[o->v] = 1;
                S.pop();
            }
        }
    }

    void rootLas(){
        Rlas(root); cout << endl;
        _Rlas(); cout << endl;
    }

    void bfs_r(){
        queue<note*> Q;
        Q.push(root);
        while(!Q.empty()){
            note* o = Q.front();
            Q.pop();
            cout << o->v << ' ';
            if(o->ln) Q.push(o->ln);
            if(o->rn) Q.push(o->rn);
        }
    }

    void bfs(){
        bfs_r(); cout << endl;
    }
};
