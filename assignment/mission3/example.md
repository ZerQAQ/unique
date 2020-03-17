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
    "id": 2,
    "nick": "owo",
    "emotionNum": 0
  },
  "msg": "ok",
  "retc": 1
}
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