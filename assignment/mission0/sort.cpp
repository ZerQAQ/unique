#include<bits/stdc++.h>
using namespace std;

#include "tree.cpp"

class Heap{
public:
    static const int N = 1e5;
    int a[N], len;

    void up(int o){
        if(o == 1) return;
        if(a[o] < a[o / 2]){
            swap(a[o], a[o / 2]);
            up(o / 2);
        }
    }

    void down(int o){
        int mini = o;
        if(o * 2 <= len && a[o * 2] < a[mini]){
            mini = o * 2;
        }
        if(o * 2 + 1 <= len && a[o * 2 + 1] < a[mini]){
            mini = o * 2 + 1;
        }
        if(mini != o){
            swap(a[mini], a[o]);
            down(mini);
        }
    }

    Heap(int *d, int l){
        len = l;
        for(int i = 0; i < l; i++){
            a[i + 1] = d[i];
        }
        for(int i = l / 2; i >= 1; i--){
            down(i);
        }
    }

    int top(){
        if(len == 0) return -1;
        else return a[1];
    }

    int del(){
        if(len == 0) return -1;
        swap(a[1], a[len]);
        len--;
        down(1);
        return a[len + 1];
    }

    int add(int v){
        a[++len] = v;
        up(len);
        return 1;
    }
};

class Sorter{
public:
    int *a, n;
    int *res, *temp, *std;
    Sorter(int *d, int _n){
        n = _n;
        a = new int[n];
        res = new int[n];
        temp = new int[n];
        std = new int[n];
        for(int i = 0; i < n; i++){
            std[i] = res[i] = a[i] = d[i];
        }
        sort(std, std + n);
    }

    ~Sorter(){
        delete []a;
        delete []res;
        delete []temp;
        delete []std;
    }

    void cmp(){
        for(int i = 0; i < n; i++){
            if(std[i] != res[i]){
                cout << "WA\n";
                return;
            }
        }
        cout << "OK\n";
    }

    void init(){
        for(int i = 0; i < n; i++){
            res[i] = a[i];
        }
    }

    int search(int *d, int l, int r, int v){
        int mid, ans = -1;
        while(l <= r){
            mid = (l + r) / 2;
            if(d[mid] <= v){
                l = mid + 1;
            }
            else{
                ans = mid;
                r = mid - 1;
            }
        }
        return ans;
    }

    void bubble_sort(){
        init();
        /*
         * 两层循环
         * 第一层执行n次
         * 第二层执行n - 1,n - 2 ... 1次
         * 时间复杂度O(n^2)
         */
        for(int i = 0; i < n; i++){
            for(int j = n - 1; j > i; j--){
                if(res[j] < res[j - 1]){
                    swap(res[j], res[j - 1]);
                }
            }
        }
    }

    void insert_sort(){
        init();
        /*
         * 一层循环 执行n次
         * 二分查找时间复杂度是O(n)
         * 移动数组时间复杂度是O(n)
         * 所以总体时间复杂度是O(n^2)
         */
        int p = 0;
        for(int i = 0; i < n; i++){
            int ind = search(res, 0, p - 1, a[i]);
            if(ind == -1){
                res[p++] = a[i];
            }
            else{
                for(int j = p; j > ind; j--){
                    res[j] = res[j - 1];
                }
                res[ind] = a[i];
                p++;
            }
        }
    }

    void merge_sort_r(int l, int r){
        /*
         * 被第i次递归执行的所有函数的归并操作时间复杂度都是O(n)
         * 被递归执行的次数应该小于 [log(n)] + 1
         * 因而时间复杂度是O(nlogn)
         */
        if(r <= l)
            return;
        int m = (l + r) / 2;
        merge_sort_r(l, m);
        merge_sort_r(m + 1, r);
        int p1 = l, p2 = m + 1, p = l;
        while(p1 <= m && p2 <= r){
            if(res[p1] < res[p2]){
                temp[p++] = res[p1++];
            }
            else{
                temp[p++] = res[p2++];
            }
        }
        while(p1 <= m) temp[p++] = res[p1++];
        while(p2 <= r) temp[p++] = res[p2++];
        for(int i = l; i <= r; i++){
            res[i] = temp[i];
        }
    }

    void merge_sort(){
        init();
        merge_sort_r(0, n - 1);
    }

    void qsort_r(int l, int r){
        /*
         * 第i次被递归执行的所有函数的总交换次数不超过n
         * 递归执行次数小于 [logn] + 1
         * 总体时间复杂度O(nlogn)
         */
        if(r <= l) return;

        int i = l, j = r, t = res[l];
        while(i < j){
            while(res[j] >= t && i < j) j--;
            while(res[i] <= t && i < j) i++;
            if(i < j) swap(res[i], res[j]);
        }
        swap(res[l], res[j]);

        qsort_r(l, j - 1);
        qsort_r(j + 1, r);
    }

    void qsort(){
        init();
        qsort_r(0, n - 1);
    }

    void heap_sort(){
        /*
         * 建堆时间复杂度O(n)
         * 每次取出一个元素是O(logn)
         * 总体时间复杂度是O(nlogn)
         */
        Heap h(a, n);
        for(int i = 0; i < n; i++){
            res[i] = h.del();
        }
    }

    void prt(int *d = 0){
        if(d == 0) d = res;
        for(int i = 0; i < n; i++){
            cout << d[i] << ' ';
        }
        cout << '\n';
    }
};

typedef function<void(void)> func;

using namespace chrono;

double Time(func f){

    auto start = high_resolution_clock::now();
    f();
    auto end = high_resolution_clock::now();
    auto duration = duration_cast<nanoseconds>(end - start);

    return (int)duration.count();
}

int main(){

    const int N = 10000;
    srand((int)time(0));
    int t[N];
    for(int i = 0; i < N; i++){
        t[i] = rand();
    }
    Sorter s(t, N);
    //cout << "origin list: \n"; s.prt(s.a);

    int tm;

    tm = Time([&] () { s.bubble_sort(); });
    printf("bubble sort\n %.2lf kns ", (double)tm / (double)1000); s.cmp();
    tm = Time([&] () { s.insert_sort(); });
    printf("insert sort\n %.2lf kns ", (double)tm / (double)1000); s.cmp();
    tm = Time([&] () { s.merge_sort(); });
    printf("merge sort\n %.2lf kns ", (double)tm / (double)1000); s.cmp();
    tm = Time([&] () { s.qsort(); });
    printf("qsort\n %.2lf kns ", (double)tm / (double)1000); s.cmp();
    tm = Time([&] () { s.heap_sort(); });
    printf("heap sort\n %.2lf kns ", (double)tm / (double)1000); s.cmp();
    cout << endl;

    Stack<int> S(100);
    S.push(1);
    S.push(2);
    S.pop();
    S.push(3);
    S.push(4);
    S.prt();
    cout << endl;

    for(int i = 0; i < 8; i++) t[i] = i;

    Tree T(t, 8);
    T.rootFir(); //先序
    cout << endl;
    T.rootmid(); //中序
    cout << endl;
    T.rootLas(); //后序
    cout << endl;
    T.bfs(); //层序

    return 0;
}