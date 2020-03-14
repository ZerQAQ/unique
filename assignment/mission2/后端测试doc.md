# HOST:

http://39.97.228.101:8080

# API:

## POST /blog

上传一篇文章

### 例：

request:

	POST /blog

	{"title": "owo", "brief": "test", "content": "123456789987654321"}

respond:

	200 ok

	{"msg": "ok", "id": 4}

上传了一片新的文章，其id为4

## GET /blog/:id

获取文章

type = brief时，只获取content的前10个字符

type = full时，获取完整的content内容

### 例:

request:

	GET /blog/4&type=brief

respond:

	200 ok

	{"msg": "ok.", "data": {"Id":4,"Title":"owo","Brief":"test","Content":"1234567899","TimeString":"2020-03-15 03:12:41","TimeUnix":1584213161}}

## POST /blog/:id?type=modify

修改特定id的文章

### 例

request:

	POST /blog/4?type=modify

	{"title": "owo", "brief": "modify", "content": "modify content"}

respond:
	
	200 ok

	{"msg": "ok", "id": 5}

修改后文章id变成5

## POST /blog/:id?type=delete

删除特定id的文章

### 例子
requset:

	POST /blog/5?type=delete

respond:

	200 ok

	{"msg": "ok", "id": 5}

requset:

	GET /blog/5

respond:

	404 not found

	{"msg": "not exist."}

