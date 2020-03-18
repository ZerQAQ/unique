package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

type log struct {
	Id int64 `json: "id" xorm:pk`
	Password string `json: "password" xorm:varchar(64)`
	CreatedAt time.Time `xorm:"created" json:"-"`
}

type user struct {
	Id int64 `json:"id" xorm:pk`
	Nick string `json:"nick" xorm:varchar(100)`
	EmotionNum int64 `json:"emotionNum"`
}

type emotion struct {
	Id      int64 `json:"id" xorm:pk`
	Uid     int64 `json:"uid"`
	Tid 	int64 `json:"-"'`
	Stars   int64 `json:"stars"`
	Type int64 `json:"type"`
	Brief string `json:"brief" xorm:varchar(100)`
	Content int64 `json:"content"`
	PhotoNum int64 `json:"photoNum"`
	CreatedAt time.Time `json:"createdAt" xorm:"created"`
}

type emotionText struct {
	Id int64 `json:"id"`
	Uid int64 `json:"uid"`
	Content string `json:"content" xorm:"varchar(2000)"`
}

type uploadStatus struct {
	Voice int64
	Photo [10]int64
}

var Router *gin.Engine
var Sql *xorm.Engine
var Sessions = make(map[string]int64)
var SessionsLifetime = make(map[string]int64)
var Uploading = make(map[int64]uploadStatus)


func myRand() string { return strconv.Itoa(rand.Int()) }

func checkSession(str string) int64{
	id, ok := Sessions[str]
	if !ok {return -1}
	if SessionsLifetime[str] < time.Now().Unix() {
		delSession(str)
		return -1
	} else {
		return id
	}
}

func newSession(id int64, life int64) string{
	skey := myRand()
	for checkSession(skey) == 0 { skey = myRand() }
	Sessions[skey] = id
	SessionsLifetime[skey] = time.Now().Unix() + life
	return skey
}

func delSession (str string) int64 {
	delete(Sessions, str)
	return 1
}

func maintainSeesion(){
	for true {
		for id := range Sessions{
			life := SessionsLifetime[id]
			if life < time.Now().Unix() {
				delete(Sessions, id)
				delete(SessionsLifetime, id)
			}
		}
		time.Sleep(3600 * time.Second) //一个小时清理一次skey
	}
}

func sqlConnectKeepAlife(){
	for true {
		Sql.Ping()
		time.Sleep(100 * time.Second) //100秒ping一次 保持连接鲜活
	}
}

func myLog(str string)  {
	fmt.Printf("[log][%v] %v\n", time.Now().String()[:19], str)
}

func idExist(id int64) bool{
	ok1, _ := Sql.Id(id).Get(new(log))
	ok2, _ := Sql.Id(id).Get(new(user))
	if ok1 != ok2{
		myLog("ERROR! user and log not match!\n")
	}
	return ok1 && ok2
}

var (
	OK = "ok"
	ServerError = "serverError"
	FormatError = "formatError"
	SkeyFail = "skeyFail"
	NotExist = "notExist"
	WrongLoginInfo = "wrongLoginInfo"
	NotUploading = "notUploading"
	UploadSuccess = "uploadSuccess"

	Mottos [][]string
	MottosLen int64
)

func quickResp(cmd string, c *gin.Context){
	if cmd == OK{
		c.JSON(200, gin.H{
			"msg": "ok",
			"retc": 1,
		})
	} else if cmd == ServerError{
		c.JSON(500, gin.H{
			"msg": "server error",
			"retc": -1,
		})
	} else if cmd == FormatError{
		c.JSON(400, gin.H{
			"msg": "format error",
			"retc": -4,
		})
	} else if cmd == SkeyFail{
		c.JSON(403, gin.H{
			"msg": "skey fail",
			"retc": -3,
		})
	} else if cmd == NotExist{
		c.JSON(404, gin.H{
			"msg": "source not exist",
			"retc": -2,
		})
	} else if cmd == WrongLoginInfo{
		c.JSON(403, gin.H{
			"msg": "wrong login info",
			"retc": -3,
		})
	} else if cmd == NotUploading{
		c.JSON(403, gin.H{
			"msg": "please POST emotion first",
			"retc": -5,
		})
	} else if cmd == UploadSuccess{
		c.JSON(200, gin.H{
			"msg": "upload success",
			"retc": 2,
		})
	}
}

func fullResp(c *gin.Context, d interface{}){
	c.JSON(200, gin.H{
		"msg": "ok",
		"retc": 1,
		"data": &d,
	})
}

func readStringFile(path string) []byte {
	f, err := os.Open(path)
	t, _ := ioutil.ReadAll(f)
	f.Close()
	if err != nil {
		myLog(fmt.Sprintf("ERROR when openning %v", path))
	}
	return t
}

func postUser(c *gin.Context)  {
	d, _ := ioutil.ReadAll(c.Request.Body)
	newUser := new(log)
	tp := c.DefaultQuery("type", "null")
	skey := c.DefaultQuery("skey", "null")
	skey = string(skey)

	if tp == "modify" {
		id := checkSession(skey)
		if id == -1 || !idExist(id){
			c.JSON(403, gin.H{
				"msg": "skey fail",
				"retc": -3,
			})
		} else {
			var newNick gin.H
			if json.Unmarshal(d, &newNick) == nil {
				Sql.Id(id).Update(&user{Nick: newNick["nick"].(string)})
				quickResp("ok", c)
			} else {quickResp("typeError", c)}
		}
		return
	}
	var dicd map[string]interface{}
	err := json.Unmarshal(d, &dicd)
	if err == nil {
		newUser.Id = int64(dicd["id"].(float64))
		newUser.Password = dicd["password"].(string)
		if newUser.Id == 0 || newUser.Password == ""{
			quickResp(FormatError, c)
			return
		}
		myLog(fmt.Sprintf("POST /user\n%v\n", string(d)))
		has, _ := Sql.Id(newUser.Id).Get(new(log))

		//fmt.Printf("has:%v", has)

		if has { //ID存在

			c.JSON(403, gin.H{
				"msg": "id has already exist",
				"retc": -2,
			})
			return
		}

		Sql.Insert(newUser)

		var nick string
		if dicd["nick"] == nil { nick = "" } else { nick = dicd["nick"].(string) }
		_, err := Sql.Insert(user{Id: newUser.Id, Nick: nick})

		if err == nil{ //ok
			quickResp(OK, c)
		} else { //服务器错误
			fmt.Print("ERROR:\n%v\n", err)
			quickResp(ServerError, c)
		}
	} else {
		myLog("Json error")
		myLog(fmt.Sprintf("%v", err))
		quickResp(FormatError, c)
	}
}

func postLogin(c *gin.Context)  {
	d, _ := ioutil.ReadAll(c.Request.Body)
	var mapd map[string]interface{}
	json.Unmarshal(d, &mapd)
	_, ok := mapd["skeyLifeTime"]
	if !ok {mapd["skeyLifeTime"] = float64(100 * 12 * 31 * 24 * 60 * 60)}
	id := int64(mapd["id"].(float64))
	password := mapd["password"].(string)
	lifetime := int64(mapd["skeyLifeTime"].(float64))

	ok, _ = Sql.Where("id = ? and password = ?", id, password).Get(new(log))
	if ok {
		skey := newSession(id, lifetime)
		c.JSON(200, gin.H{
			"msg": "ok",
			"retc": 1,
			"skey": skey,
		})
	} else {
		c.JSON(403, gin.H{
			"msg": "id or password wrong",
			"retc": -3,
		})
	}
}

func postLogout(c *gin.Context)  {
	skey := c.DefaultQuery("skey", "null")
	if skey == "null" || checkSession(skey) == -1 {
		quickResp(SkeyFail, c)
		return
	} else {
		delSession(skey)
		quickResp(OK, c)
		return
	}
}

func getUser(c *gin.Context)  {
	skey := c.DefaultQuery("skey", "null")
	if skey == "null" || checkSession(skey) == -1 {
		quickResp(SkeyFail, c)
		return
	} else {
		id := checkSession(skey)
		if !idExist(id) {
			myLog("ERROR! id not exist when getUser\n")
			quickResp(ServerError, c)
			return
		}
		var userData user
		Sql.Id(id).Get(&userData)
		fullResp(c, &userData)
		return
	}
}

func postEmotion(c *gin.Context){
	skey := c.DefaultQuery("skey", "null")
	if skey == "null" || checkSession(skey) == -1 {
		quickResp(SkeyFail, c)
		return
	} else {
		uid := checkSession(skey)
		d, _ := ioutil.ReadAll(c.Request.Body)
		var dicd map[string]interface{}
		var newEmotion emotion
		json.Unmarshal(d, dicd)
		json.Unmarshal(d, newEmotion)

		text := dicd["text"].(string)
		content := int64(dicd["content"].(float64))
		if len(text) > 20 {
			newEmotion.Brief = text[:20]
		} else {
			newEmotion.Brief = text
		}
		var newEmotionText emotionText
		newEmotionText.Content = text
		newEmotionText.Uid = uid
		_, err1 := Sql.Insert(newEmotionText)

		newEmotion.Tid = newEmotionText.Id
		newEmotion.Uid = uid
		_, err2 := Sql.Insert(newEmotion)

		if err1 != nil || err2 != nil {
			quickResp(ServerError, c)
		}

		if newEmotion.PhotoNum > 0 || content & 1 == 1 {
			var ul uploadStatus
			if newEmotion.PhotoNum > 0 {
				for i := 1; i <= int(newEmotion.PhotoNum); i++ {
					ul.Photo[i] = int64(i)
				}
			}
			if content & 1 == 0 {
				ul.Voice = 1
			}
			Uploading[newEmotion.Id] = ul
		}

		eid := newEmotion.Id

		fullResp(c, gin.H{
			"id": eid,
		})
	}
}

func uploadOK(s uploadStatus) int64 {
	for i := 1; i <= 9; i++{
		if s.Photo[i] != 0 { return 0 }
	}
	if s.Voice == 1 { return 1 } else { return 0 }
}

func postSrcVoice_Id(c *gin.Context)  {
	skey := c.DefaultQuery("skey", "null")
	filetype := c.DefaultQuery("filetype", "null")
	if skey == "null" || checkSession(skey) == -1 {
		quickResp(SkeyFail, c)
		return
	}
	if filetype == "null" {
		quickResp(FormatError, c)
		return
	}
	uid := checkSession(skey)
	eidr, _ := strconv.Atoi(c.Param("id"))
	eid := int64(eidr)
	ul, has := Uploading[eid]
	if !has {
		quickResp(NotUploading, c)
	} else if ul.Voice == 1 {
		quickResp(OK, c)
	} else{
		f, _ := c.FormFile("file")
		path := fmt.Sprintf("src/%d/%d/voice.%v", uid, eid, filetype)
		c.SaveUploadedFile(f, path)
		if uploadOK(ul) == 1 {
			quickResp(UploadSuccess, c)
		} else {
			quickResp(OK, c)
		}
	}
}
/*
func postSrcPhoto_Id(c *gin.Context){
	skey := c.DefaultQuery("skey", "null")
	filetype := c.DefaultQuery("filetype", "null")
	if skey == "null" || checkSession(skey) == -1 {
		quickResp(SkeyFail, c)
		return
	}
	if filetype == "null" {
		quickResp(FormatError, c)
		return
	}
	uid := checkSession(skey)
	eidr, _ := strconv.Atoi(c.Param("id"))
	eid := int64(eidr)

}
 */



func getMotto(c *gin.Context)  {
	k := rand.Int63() % MottosLen
	fullResp(c, &gin.H{
		"content": Mottos[k][0],
		"author": Mottos[k][1],
	})
	fmt.Printf("ret\n")
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}

func main() {
	rand.Seed(time.Now().Unix())
	Sql, _ = xorm.NewEngine("mysql", "root:123456@/test?charset=utf8")
	Sql.DatabaseTZ, _ = time.LoadLocation("Asia/Shanghai")
	Sql.TZLocation, _ = time.LoadLocation("Asia/Shanghai")
	Sql.Sync2(new(user))
	Sql.Sync2(new(log))
	Sql.Sync2(new(emotion))
	Sql.Sync2(new(emotionText))
	Router = gin.Default()
	Router.Use(Cors())
	r := Router.Group("/kuro")
	json.Unmarshal(readStringFile("motto.json"), &Mottos)
	MottosLen = int64(len(Mottos))

	r.Handle("POST", "/user", postUser)
	r.Handle("POST", "/login", postLogin)
	r.Handle("POST", "/logout", postLogout)
	r.Handle("GET", "/user", getUser)
	r.Handle("GET", "/motto", getMotto)

	r.Handle("POST", "/emotion", postEmotion)
	r.Handle("POST", "/src/voice/:id", postSrcVoice_Id)

	go maintainSeesion()
	go sqlConnectKeepAlife()

	Router.Run()
}