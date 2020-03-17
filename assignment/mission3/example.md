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
    "content": "我们曾如此渴望命运的波澜,\n到最后才发现,\n人生最曼妙的风景,\n竟是内心的淡定与从容.\n我们曾如此期待他人的认同,\n到最后才知道,\n世界是自己的,\n与他人毫无关系.\n--杨绛"
  },
  "msg": "ok",
  "retc": 1
}
```