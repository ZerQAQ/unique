## TABLES

### TABLE log:

	id[int64](pk)

	password[string64(1-9a-f)] SHA-256加密后的密码 16进制数字符串 字母小写

### TABLE user:

	id[int64](pk)

	nick[string100]

	eemotionNum[int64]

	createdAt[int64]

### TABLE emotion:

	id[int64](pk) 心情id
	
	uid[int64] 用户id

	stars[int64] 星星数量 1~5

	type[int64] 心情类型 0:好心情 1:坏心情

	content[int64] 0~3 是否有语音和悦纳，分别用第0~1二进制位表示,例如有只语音的心情的content是01，即1

	brief[string20]

	text[string2000]

	accept[string2000]

	photoNum[int64] 0~9 照片数量

	createdAt[int64]

## DIR

- src
	- uid(int64)
		- emotionid(int64)
			- photo
				- num(int64).*
			- voice.*
	- head.*
- main.exe

## CODE

函数根据对应的URL命名，如果有可变字段，则在该字段前加_

例:

- POST user/photo 对应的函数是 postUserPhoto()

- GET src/photo/:eid/:num 对应的函数是 getSrcPhoto_Eid_Num()

全局变量首字母大写

用mylog(string)打印日志

## API:

根目录是/kuro

所有API返回的json数据中都有msg字段，retc字段和data字段

msg是服务器返回的信息 retc是返回代码

retc说明：
- 2 eemotion上传成功
- 1 正常
- -1 服务器错误
- -2 资源不存在/用户ID已存在
- -3 权限不足（skey错误或失效）/用户名或密码错误
- -4 数据格式错误 (不是合法的json)
- -5 在未上传emotion信息前上传了emotion语音或照片 / 悦纳的心情不是坏心情

下面的API下第一个代码段是请求体格式，第二个是回复体中data字段格式，除了POST/login外，POST/login返回的skey字段是没有data包裹的。

password字段是SHA256加密后的十六进制字符串 字母小写(测试的时候也可以是长度小于64的字符串)

skey是纯数字 长度在40以内

上传文件的格式：
```
POST http://localhost:8080/src/photo/2/1?skey=&filetype=png
Content-Type: file

二进制文件
```

### POST /user √
新建用户
```
{"id": [int64], "password" [string64], "nick" [string100], "growthPoint": [int64]}
```
```
{}
```

### POST /login √
登录
```
{"id": [int64], "password": [string64], "skeyLifeTime": [int64]}

skeyLifeTime是返回的skey的生命周期，单位秒，默认值是-1，即永远不失效
```
```
{"skey": [string]}
```

### POST /logout?skey= √
退出登录
```
```
```
{}
```

### GET /user?skey= √
返回用户信息
```

```
```
{
  "data": {
    "id": 113414123,
    "nick": "",
    "emotionNum": 0,
    "goodmoodNum": 0,
    "badmoodNum": 0,
    "acceptmoodNum": 0,
    "imageurl": "https://s1.ax1x.com/2020/03/20/8gHl79.jpg", //头像地址
    "growthPoint": 0
  },
  "msg": "ok",
  "retc": 1
}
```

### POST /user/photo?skey=&filetype= 上传头像 √
```
二进制文件
```
```
```

### POST /user?skey=&type=modify √
修改用户昵称
```
{"nick": [string100]}
```
```
```

考虑到包不一定具有时序性，建议发送完emotion包之后先sleep(0.1)

否则在emotion包到达前，语音和图片包都会被丢弃

全部发送成功之后返回的包里的retc字段是2

### POST /emotion?skey= √
```
{
	"stars": [int64],
	"type": [int64],
	"content": [int64],
	"text": string[2000], emotion的文字内容 2000字以内
	"photoNum": [int64],
}
```
```
```

### POST /emotion/:id?skey=&type=modify&key=stars √
```
{
	"stars": [int64]
}
```
```
{}
```

### POST /src/voice/:id?skey=&filetype= √
emotionid为:id的语音
```
二进制文件
```
```
```

### POST /src/photo/:id/:num?skey=&filetype= √
emotionid为:id的第num张图片
```
二进制文件
```
```
{notload: int64[], url: string} notload里面存着还未上传的照片，是1~photoNum的正整数
```

### GET /emotions?skey=&type=&content=&page=&rank=&search=&full= √

search
	模糊搜索给定的字符串，默认是空字符串，代表不搜索

page:
	页序号，一页最多20条信息

content、type:
	获取特定content和type的emotion，content和type的说明见TABLE emotion

full:
	默认等于0
	如果等于1，就返回完整的内容（text和accept）(文字内容和悦纳内容)

rank:
- 0 不排序
- 1 按照时间降序排序
- -1 按照时间升序排序
- 2 按照星星数量降序排序
- -2 按照星星升序排序

```
```
``` 
返回：
{
	"page": int64, 页序号
	"num": int64, emotionList的长度
	$emotionList
}

$emotionList 是长度为num的emotion数组，emotion的格式为：
{
	"id": int64,
	"stars": int64,
	"type": int64,
	"content": int64,
	"photoNum": int64,
	"brief": string[20], (心情文字的前20个字)
	"text": string[2000], full = 1 才有
	"accept": string[2000], full = 1 才有
	"createdAt": int64 (Unix时间戳 创建时间)
}
```

### GET/emotion/:id?skey= √

获取id为id的emotion的完整信息

我是这么想的 前端先获得emotionlist 然后根据里面的内容 可以知道想要那条id的具体内容 然后再通过这里获取完整内容

这样就不用发emotionlist的时候就发完整的text和accept内容了，这两个内容应该是长度最大的

当然emotionlist也有full字段 可以强制获取所有内容

```
```

```
{
	"id": int64,
	"stars": int64,
	"type": int64,
	"content": int64,
	"photoNum": int64,
	"brief": string[20], (心情文字的前20个字)
	"text": string[2000],
	"accept": string[2000],
	"createdAt": int64 (Unix时间戳 创建时间)
}
```

### GET/emotion?skey=&type=random&etype= √

随机获取一条用户信息

etype = 0 获取好心情

etype = 1 获取坏心情

```
```

```
{
	"id": int64,
	"stars": int64,
	"type": int64,
	"content": int64,
	"photoNum": int64,
	"brief": string[20], (心情文字的前20个字)
	"text": string[2000],
	"accept": string[2000],
	"createdAt": int64 (Unix时间戳 创建时间)
}
```



### GET /src/photo/:id/:num&skey= √

获取id为:id的心情的第:num张照片(从1开始计数)
```
```
```
二进制文件
```

### GET /src/voice/:id&skey= √

获取id为:id的心情的语音
```
```
```
二进制文件
```

### POST /emotion/:id?skey=&type=accept √
悦纳id为:id的心情
```
{"accept": string[2000]}
```
```
{}
```

### POST /emotion/:id?skey=&type=delete √
粉碎id为:id的心情
```
```
```
{}
```

### GET /src/motto √
随机获得一段格言
```
```
```
{"content": string, "author": string}
```