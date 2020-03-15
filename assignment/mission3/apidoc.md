## TABLES

### TABLE log:

	id[int64](pk)

	password[string64(1-9a-f)] SHA-256加密后的密码 16进制数字符串 字母小写

### TABLE user:

	id[int64](pk)

	nick[string100]

	emtionNum[int64]

### TABLE motion:

	id[int64](pk) 心情id
	
	uid[int64](pk) 用户id

	timeUnix[int64] 发送时间

	stars[int64] 星星数量 1~5

	type[int64] 心情类型 0:好心情 1:坏心情-未悦纳 2:坏心情-已悦纳

	content[int64] 0~7 是否有文字照片和语音，分别用第0~2二进制位表示,例如只有照片和文字的心情的content是011，即3

	photoNum[int64] 0~9 照片数量

## API:

所有API返回的json数据中都有msg字段和retc字段，所以下面返回的格式里我只写除msg和retc之外的字段

msg是服务器返回的信息 retc是返回代码

retc说明：
- 2 emotion上传成功
- 1 正常
- -1 服务器错误
- -2 资源不存在
- -3 权限不足（skey错误或失效）

下面的API下第一个代码段是请求体格式，第二个是回复体格式

password字段是SHA256加密后的十六进制字符串 字母小写

### POST /user 
新建用户
```
{"id": [string100], "password" [string64]}
```
```
{}
```

### POST /login 
登录
```
{"id": [string100], "password": [string64], "skeyLifeTime": [int64]}

skeyLifeTime是返回的skey的生命周期，单位秒，默认值是-1，即永远不失效
```
```
{"skey": [string]}
```

### POST /logout 
退出登录
```
{"skey": [string]}
```
```
{}
```

### GET /user?skey= 
返回用户信息
```

```
```
{"nick": [string100], emotionNum: int64}
```

### POST /user?type=modify 
修改用户昵称
```
{"nick": [string100]}
```
```
```

下面四类请求要连着发，全部发完了才算创建成功

全部发送成功之后返回的包里的retc字段是2

### POST /motion?skey=
```
{
	"id": int64,
	"timeUnix": int64,
	"stars": int64,
	"type": int64,
	"content": int64,
	"photoNum": int64
}
```
```
```

### POST /src/text/:id 
motionid为:id的文字
```
字符串 不用json格式
```
```
```

### POST /src/voice/:id 
motionid为:id的语音
```
二进制文件
```
```
```

### POST /src/voice/:id/:num 
motionid为:id的第num张图片
```
二进制文件
```
```
```

### GET /motion?skey=&type=&content=&page=&rank=
获取id为:id的用户motion列表，可指定获取特定type和content的motion，可分页（一页数量最多20条，从0开始计数，-1代表返回所有数据），可排序（按照星星数量、日期等）

rank:
- 0 不排序
- 1 按照时间排序
- 2 按照星星数量排序

type、content:

- 见TABLE motion的说明

```
```
```
返回：
{
	"page": int64 页序号
	"num": int64, emotionList的长度
	$emotionList
}

$emotionList 是长度为num的emotion列表，emotion的格式为：

{
	"id": int64,
	"timeUnix": int64,
	"timeString": string[19],
	"stars": int64,
	"type": int64,
	"content": int64,
	"photoNum": int64
}
```

### GET /src/text/:id&skey= 
获取id为:id的心情文字
```
```
```
{"data": string}
```

### GET /src/photo/:id/:num&skey= 
获取id为:id的心情的第:num张照片
```
```
```
二进制文件
```

### GET /src/voice/:id&skey= 
获取id为:id的心情的语音
```
```
```
二进制文件
```

### POST /motion/:id?skey=&type=modify 
修改id为:id的心情的状态（悦纳）
```
{"type": int64} 把type的值修改成"type"
```
```
{}
```

### POST /motion/:id?skey=&type=delete 
删除id为:id的心情
```
```
```
{}
```

### GET /src/motto 
随机获得一段格言
```
```
```
{"data": string}
```