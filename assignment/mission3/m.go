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
	GrowthPoint int64 `json:"growthPoint"`
}

type emotion struct {
	Id      int64 `json:"id"`
	Uid     int64 `json:"uid"`
	Stars   int64 `json:"stars"`
	Type int64 `json:"type"`
	Brief string `json:"brief" xorm:varchar(100)`
	Content int64 `json:"content"`
	PhotoNum int64 `json:"photoNum"`
	Accept string `json:"-" xorm:varchar(2000)`
	Text string `json:"-" xorm:varchar(2000)`
	CreatedAt time.Time `json:"createdAt" xorm:"created"`
}

type uploadStatus struct {
	Id int64
	Voice int64
	Photo [10]int64
	PhotoNum int64
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
	fmt.Printf("[mylog][%v] %v\n", time.Now().String()[:19], str)
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
	NotBadEmotion = "notBadEmotion"

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
	} else if cmd == NotUploading{
		c.JSON(403, gin.H{
			"msg": "not bad emotion",
			"retc": -5,
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
		myLog(fmt.Sprintf("POST /user %v", string(d)))
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

	myLog(fmt.Sprintf("POST /login %v", string(d)))

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

/*
func GetUserPhoto(c *gin.Context) {
	skey := c.DefaultQuery("skey", "null")
	if skey == "null" || checkSession(skey) == -1 {
		quickResp(SkeyFail, c)
		return
	}
	uid := checkSession(skey)
	path := fmt.Sprintf("src/%v/head", uid)
	_, err := os.Stat(path)
}
 */

func postEmotion(c *gin.Context){
	skey := c.DefaultQuery("skey", "null")
	if skey == "null" || checkSession(skey) == -1 {
		quickResp(SkeyFail, c)
		return
	} else {
		uid := checkSession(skey)
		d, _ := ioutil.ReadAll(c.Request.Body)

		myLog(fmt.Sprintf("body:\n%v\n", string(d)))


		var dicd map[string]interface{}
		var newEmotion emotion
		json.Unmarshal(d, &dicd)
		json.Unmarshal(d, &newEmotion)

		fmt.Printf("dicd:\n%v\n", dicd)

		text := dicd["text"].(string)
		content := int64(dicd["content"].(float64))
		if len(text) > 20 {
			newEmotion.Brief = text[:20]
		} else {
			newEmotion.Brief = text
		}
		newEmotion.Text = text
		newEmotion.Uid = uid
		newEmotion.Id = 0
		fmt.Printf("emo:\n%v", newEmotion)
		_, err := Sql.Insert(&newEmotion)

		if err != nil {
			quickResp(ServerError, c)
		}

		if newEmotion.Type == 0 {
			addGrowth(uid, 3)
		}

		if newEmotion.PhotoNum > 0 || content & 1 == 1 {
			var ul uploadStatus
			ul.PhotoNum = newEmotion.PhotoNum
			if newEmotion.PhotoNum > 0 {
				for i := 1; i <= int(newEmotion.PhotoNum); i++ {
					ul.Photo[i] = int64(i)
				}
			}
			if content & 1 == 0 {
				ul.Voice = 1
			}
			ul.Id = newEmotion.Id
			Uploading[newEmotion.Id] = ul
		}

		fullResp(c, gin.H{
			"id": newEmotion.Id,
		})


	}
}

func uploadOK(s uploadStatus) int64 {
	for i := 1; i <= 9; i++{
		if s.Photo[i] != 0 { return 0 }
	}
	if s.Voice == 1 { delete(Uploading, s.Id); return 1 } else { return 0 }
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
		if f == nil {
			quickResp(FormatError, c)
			return
		}
		dir := fmt.Sprintf("src/%v/%v", uid, eid)
		os.MkdirAll(dir, os.ModePerm)
		path := fmt.Sprintf("src/%d/%d/voice.%v", uid, eid, filetype)
		err := c.SaveUploadedFile(f, path)
		if err != nil {
			quickResp(ServerError, c)
			return
		}
		ul.Voice = 1
		if uploadOK(ul) == 1 {
			quickResp(UploadSuccess, c)
		} else {
			Uploading[ul.Id] = ul
			quickResp(OK, c)
		}
	}
}

func postSrcPhoto_Id_Num(c *gin.Context){
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
	numr, _ := strconv.Atoi(c.Param("num"))
	eid := int64(eidr)
	num := int64(numr)
	ul, has := Uploading[eid]

	if !has {
		quickResp(NotUploading, c)
		return
	}
	if num > ul.PhotoNum {
		quickResp(FormatError, c)
		return
	}

	f, _ := c.FormFile("file")
	if f == nil {
		quickResp(FormatError, c)
		return
	}
	dir := fmt.Sprintf("src/%d/%d/photo", uid, eid)
	fmt.Printf("dir:\n%v\n", dir)
	os.MkdirAll(dir, os.ModePerm)
	path := fmt.Sprintf("src/%d/%d/photo/%d.%v", uid, eid, num, filetype)
	err := c.SaveUploadedFile(f, path)
	if err != nil {
		quickResp(ServerError, c)
		return
	}
	ul.Photo[num] = 0
	if uploadOK(ul) == 1 {
		quickResp(UploadSuccess, c)
		return
	} else {
		Uploading[ul.Id] = ul
		var notload []int64
		for i := 1; i <= 9; i++ {
			if ul.Photo[i] != 0 {
				notload = append(notload, int64(i))
			}
		}
		fullResp(c, gin.H{
			"notLoad": notload,
			"url": fmt.Sprintf("src/photo/%v/%v", eid, num),
		})
	}
}

func addGrowth(uid int64, v int64) {
	var u user
	u.Id = uid
	Sql.Get(&u)
	u.GrowthPoint += v
	Sql.Id(uid).Update(u)
}

func postEmotion_Id(c *gin.Context){
	skey := c.DefaultQuery("skey", "null")
	tp := c.DefaultQuery("type", "null")
	key := c.DefaultQuery("key", "null")
	if skey == "null" || checkSession(skey) == -1 {
		quickResp(SkeyFail, c)
		return
	}
	eidr, _ := strconv.Atoi(c.Param("id"))
	eid := int64(eidr)
	uid := checkSession(skey)

	has, _ := Sql.Get(&emotion{Id: eid})

	if ! has {
		quickResp(NotExist, c)
		return
	}

	has, _ = Sql.Get(&emotion{Id: eid, Uid: uid})

	if !has {
		quickResp(SkeyFail, c)
		return
	}

	if tp == "delete" {
		_, err1 := Sql.Delete(emotion{Id: eid})
		addGrowth(uid, 1)

		path := fmt.Sprintf("src/%v/%v", uid, eid)
		err2 := os.RemoveAll(path)
		if err1 != nil || err2 != nil {
			quickResp(ServerError, c)
			return
		}
		quickResp(OK, c)
		return
	} else if tp == "accept" {
		var acEmotion emotion
		acEmotion.Id = eid
		acEmotion.Uid = uid
		addGrowth(uid, 3)
		Sql.Get(&acEmotion)

		fmt.Print(acEmotion)
		fmt.Printf("\n%v\n",acEmotion)

		if acEmotion.Type != 1 {
			quickResp(NotBadEmotion, c)
			return
		}

		d, _ := ioutil.ReadAll(c.Request.Body)
		acText := string(d)

		acEmotion.Accept = acText
		acEmotion.Type = 0
		acEmotion.Content |= 2

		Sql.Id(acEmotion.Id).Update(&acEmotion)
		Sql.Id(acEmotion.Id).Cols("type").Update(&acEmotion)

		quickResp(OK, c)
		return
	} else if tp == "modify" {
		dr, _ := ioutil.ReadAll(c.Request.Body)
		var dicd map[string]interface{}
		json.Unmarshal(dr, &dicd)
		if key == "stars" {
			has, _ := Sql.Get(&emotion{Id: eid, Uid: uid})
			if !has {
				quickResp(NotExist, c)
				return
			}
			_, err := Sql.Id(eid).Update(&emotion{Stars: int64(dicd["stars"].(float64))})
			if err != nil {
				quickResp(ServerError, c)
				return
			}
			quickResp(OK, c)
		} else {
			quickResp(FormatError, c)
		}
	} else {
		quickResp(FormatError, c)
		return
	}
}

func getMotto(c *gin.Context)  {
	k := rand.Int63() % MottosLen
	fullResp(c, &gin.H{
		"content": Mottos[k][0],
		"author": Mottos[k][1],
	})
	fmt.Printf("ret\n")
}

func postUserPhoto(c *gin.Context)  {
	skey := c.DefaultQuery("skey", "null")
	filetype := c.DefaultQuery("filetype", "")
	if skey == "null" || checkSession(skey) == -1 {
		quickResp(SkeyFail, c)
		return
	}
	uid := checkSession(skey)
	dir := fmt.Sprintf("src/%v", uid)
	os.MkdirAll(dir, os.ModePerm)

	f, _ := c.FormFile("file")
	if f == nil {
		quickResp(FormatError, c)
		return
	}
	path := dir + fmt.Sprintf("head.%v", filetype)
	err := c.SaveUploadedFile(f, path)
	if err != nil {
		quickResp(ServerError, c)
		return
	}
	quickResp(OK, c)
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
	r.Handle("POST", "/src/photo/:id/:num", postSrcPhoto_Id_Num)
	r.Handle("POST", "/emotion/:id", postEmotion_Id)

	go maintainSeesion()
	go sqlConnectKeepAlife()

	Sessions["1"] = 2
	SessionsLifetime["1"] = time.Now().Unix() * 2

	Router.Run()
}