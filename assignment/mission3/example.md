## 1
```
POST http://39.97.228.101:8080/kuro/login
Accept: */*
Cache-Control: no-cache

{"id": 1, "password": "1"}
```

```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Mon, 16 Mar 2020 19:59:33 GMT
Content-Length: 51

{
  "msg": "ok",
  "retc": 1,
  "skey": "8672257192191626248"
}
```

## 2
```
POST http://39.97.228.101:8080/kuro/user?skey=947406850690932915&type=modify
Accept: */*
Cache-Control: no-cache

{"nick": "owo"}
```

```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Mon, 16 Mar 2020 13:30:28 GMT
Content-Length: 22

{
  "msg": "ok",
  "retc": 1
}
```

## 3
```
GET http://39.97.228.101:8080/kuro/user?skey=947406850690932915
Accept: application/json
```

```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Mon, 16 Mar 2020 19:53:04 GMT
Content-Length: 66

{
  "data": {
    "id": 113414123,
    "nick": "",
    "emotionNum": 0,
    "goodmoodNum": 0,
    "badmoodNum": 0,
    "acceptmoodNum": 0,
    "imageurl": "https://s1.ax1x.com/2020/03/20/8gHl79.jpg", //头像url
    "growthPoint": 0
  },
  "msg": "ok",
  "retc": 1
}


Response code: 200 (OK); Time: 24ms; Content length: 191 bytes

```

## 4

```cassandraql
GET http://localhost:8080/kuro/motto
Accept: application/json
```

```cassandraql
{
  "data": {
    "author": "惠特曼",
    "content": "做一个世界的水手,\n游遍每一个港口.\n"
  },
  "msg": "ok",
  "retc": 1
}
```

### 5 
```cassandraql
GET http://localhost:8080/kuro/emotion?type=random&skey=1&etype=0
```
```cassandraql
{
  "data": {
    "id": 6,
    "uid": 2,
    "stars": 3,
    "type": 0,
    "brief": "hello, orld",
    "content": 2,
    "photoNum": 0,
    "accept": "test accept content2",
    "text": "hello, orld",
    "createdAt": 0
  },
  "msg": "ok",
  "retc": 1
}
```