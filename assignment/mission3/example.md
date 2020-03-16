## 1
```
POST http://39.97.228.101:8080/kuro/login
Accept: */*
Cache-Control: no-cache

{"id": 1, "password": "1"}
```

```
HTTP/1.1 200 OK
Content-Type: text/plain; charset=utf-8
Date: Mon, 16 Mar 2020 13:29:01 GMT
Content-Length: 49

{"msg":"ok","retc":1,"skey":"947406850690932915"}
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
Date: Mon, 16 Mar 2020 13:31:14 GMT
Content-Length: 76

{
  "data": "{\"id\":1,\"nick\":\"owo\",\"emotionNum\":0}",
  "msg": "ok",
  "retc": 1
}
```