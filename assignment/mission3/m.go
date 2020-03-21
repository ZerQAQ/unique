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
	"strings"
	"time"
)

type log struct {
	Id int64 `json: "id" xorm:pk`
	Password string `json: "password" xorm:varchar(64)`
	CreatedAt int64 `xorm:"created" json:"-"`
}

type user struct {
	Id int64 `json:"id" xorm:pk`
	Nick string `json:"nick" xorm:varchar(100)`
	EmotionNum int64 `json:"emotionNum"`
	GoodmoodNum int64 `json:"goodmoodNum"`
	BadmoodNum int64 `json:"badmoodNum"`
	AcceptMoodNum int64 `json:"acceptmoodNum"`
	Imageurl string `json:"imageurl"`
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
	Accept string `json:"accept" xorm:varchar(2000)`
	Text string `json:"text" xorm:varchar(2000)`
	CreatedAt int64 `json:"createdAt" xorm:"created"`
	StringCreatedAt string `json:"stringCreatedAt" xorm:"varchar(50)"`
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

func maintainSession(){
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
	} else if cmd == NotBadEmotion{
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
	var dicd = make(map[string]interface{})
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

		fmt.Printf("\n id: %v pd: %v \n", newUser.Id, newUser.Password)
		i, errx := Sql.Insert(newUser)
		fmt.Printf("\n i: %v errx: %v \n", i, errx)

		var nick string
		if dicd["nick"] == nil { nick = "" } else { nick = dicd["nick"].(string) }
		var nu = user{Id: newUser.Id, Nick: nick, Imageurl:"https://s1.ax1x.com/2020/03/20/8gHl79.jpg"}
		nu.EmotionNum = 0
		nu.GoodmoodNum = 0
		nu.BadmoodNum = 0
		_, err := Sql.Insert(&nu)

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
	var mapd = make(map[string]interface{})
	json.Unmarshal(d, &mapd)
	_, ok := mapd["skeyLifeTime"]
	var lifetime int64
	if !ok { lifetime = 100 * 12 * 31 * 24 * 60 * 60 } else {
		lifetime = int64(mapd["skeyLifeTime"].(float64))
	}
	id := int64(mapd["id"].(float64))
	password := mapd["password"].(string)


	myLog(fmt.Sprintf("POST /login %v", string(d)))

	ok, _ = Sql.Where("id = ? and password = ?", id, password).Get(new(log))
	if ok {
		skey := newSession(id, lifetime)
		myLog(fmt.Sprintf("newskey: %v ", skey))
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

		myLog(fmt.Sprintf("body: %v", string(d)))


		var dicd = make(map[string]interface{})
		var newEmotion emotion
		json.Unmarshal(d, &dicd)
		json.Unmarshal(d, &newEmotion)

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
		newEmotion.StringCreatedAt = time.Now().Format("2006/01/02 15:04")
		_, err := Sql.Insert(&newEmotion)

		var u user
		u.Id = uid
		Sql.Get(&u)
		u.EmotionNum += 1

		if newEmotion.Type == 0 {
			u.GoodmoodNum += 1
		} else {
			u.BadmoodNum += 1
		}

		Sql.Id(uid).Update(&u)

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
			myLog(fmt.Sprintf("create notload. eid: %v nl: %v", newEmotion.Id, ul))
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
	myLog(fmt.Sprintf("postSrcVoice"))
	skey := c.DefaultQuery("skey", "null")
	filetype := c.DefaultQuery("filetype", "null")
	myLog(fmt.Sprintf("voice skey : %v", skey))
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
	myLog(fmt.Sprintf("eid: %v has: %v", eid, has))
	if !has {
		quickResp(NotUploading, c)
		return
	} else if ul.Voice == 1 {
		quickResp(OK, c)
		return
	} else{

		dir := fmt.Sprintf("src/%v/%v", uid, eid)
		os.MkdirAll(dir, os.ModePerm)
		path := fmt.Sprintf("src/%d/%d/voice.%v", uid, eid, filetype)
		os.Create(path)
		d, _ := ioutil.ReadAll(c.Request.Body)
		err := ioutil.WriteFile(path, d, 0666)
		if err != nil {
			quickResp(ServerError, c)
			return
		}
		ul.Voice = 1
		if uploadOK(ul) == 1 {
			delete(Uploading, ul.Id)
			quickResp(UploadSuccess, c)
			return
		} else {
			Uploading[ul.Id] = ul
			quickResp(OK, c)
			return
		}
	}
}

func postSrcPhoto_Id_Num(c *gin.Context){
	skey := c.DefaultQuery("skey", "null")
	filetype := c.DefaultQuery("filetype", "null")
	myLog(fmt.Sprintf("photo s: %v ft: %v", skey, filetype))
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
	myLog(fmt.Sprintf("eid: %v pn: %v n: %v", eid, ul.PhotoNum, num))
	if num > ul.PhotoNum {
		quickResp(FormatError, c)
		return
	}

	dir := fmt.Sprintf("src/%d/%d/photo", uid, eid)
	fmt.Printf("dir:\n%v\n", dir)
	os.MkdirAll(dir, os.ModePerm)
	path := fmt.Sprintf("src/%d/%d/photo/%d.%v", uid, eid, num, filetype)
	os.Create(path)
	d, _ := ioutil.ReadAll(c.Request.Body)
	err := ioutil.WriteFile(path, d, 0666)
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

	fmt.Printf("\n t: %v k: %v \n", tp, key)

	if skey == "null" || checkSession(skey) == -1 {
		quickResp(SkeyFail, c)
		return
	}
	eidr, _ := strconv.Atoi(c.Param("id"))
	eid := int64(eidr)
	uid := checkSession(skey)

	has, _ := Sql.Get(&emotion{Id: eid})

	if ! has {
		myLog("eid not exist")
		quickResp(NotExist, c)
		return
	}

	has, _ = Sql.Get(&emotion{Id: eid, Uid: uid})

	if !has {
		myLog("uid eid not exist")
		quickResp(SkeyFail, c)
		return
	}

	if tp == "delete" {
		em := emotion{Id:eid}
		Sql.Get(&em)

		_, err1 := Sql.Delete(emotion{Id: eid})
		addGrowth(uid, 1)

		u := user{Id:uid}
		Sql.Get(&u)

		if em.Type == 0 { u.GoodmoodNum -= 1 } else { u.BadmoodNum -= 1 }
		u.EmotionNum -= 1

		Sql.Id(u.Id).Update(&u)
		Sql.Id(u.Id).Cols("goodmood_num").Update(&u)
		Sql.Id(u.Id).Cols("badmood_num").Update(&u)
		Sql.Id(u.Id).Cols("emotion_num").Update(&u)
		

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

		var u user
		u.Id = uid
		Sql.Get(&u)
		u.AcceptMoodNum += 1
		Sql.Id(u.Id).Update(u)

		if acEmotion.Type != 1 {
			quickResp(NotBadEmotion, c)
			return
		}

		dr, _ := ioutil.ReadAll(c.Request.Body)
		var dicd = make(map[string]interface{})
		json.Unmarshal(dr, &dicd)
		acText := dicd["accept"].(string)

		acEmotion.Accept = acText
		acEmotion.Type = 0
		acEmotion.Content |= 2

		Sql.Id(acEmotion.Id).Update(&acEmotion)
		Sql.Id(acEmotion.Id).Cols("type").Update(&acEmotion)

		quickResp(OK, c)
		return
	} else if tp == "modify" {
		dr, _ := ioutil.ReadAll(c.Request.Body)
		var dicd = make(map[string]interface{})
		json.Unmarshal(dr, &dicd)
		//fmt.Printf("\n%v\n", key)
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
	if filetype == "null" {
		quickResp(FormatError, c)
		return
	}
	uid := checkSession(skey)
	dir := fmt.Sprintf("src/%v", uid)
	path := dir + "/head." + filetype
	os.MkdirAll(dir, os.ModePerm)
	delFile(dir, "head")
	os.Create(path)
	d, _ := ioutil.ReadAll(c.Request.Body)
	err := ioutil.WriteFile(path, d, 0666)
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
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token, Content-Disposition")
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

func getEmotion_Id(c *gin.Context) {
	skey := c.DefaultQuery("skey", "null")
	if skey == "null" || checkSession(skey) == -1 {
		quickResp(SkeyFail, c)
		return
	}

	eidr, _ := strconv.Atoi(c.Param("id"))
	eid := int64(eidr)
	uid := checkSession(skey)

	var rEmotion emotion
	rEmotion.Id = eid
	rEmotion.Uid = uid

	has, _ := Sql.Get(&rEmotion)

	if !has {
		quickResp(NotExist, c)
		return
	}

	fullResp(c, rEmotion)
}

func getEmotion(c *gin.Context) {
	skey := c.DefaultQuery("skey", "null")
	tp := c.DefaultQuery("type", "null")
	etpr, _ := strconv.Atoi(c.DefaultQuery("etype", "0"))
	etp := int64(etpr)

	myLog(fmt.Sprintf("random skey: %v", skey))

	fmt.Printf("\n s:%v t:%v \n", skey, tp)
	if skey == "null" || checkSession(skey) == -1 {
		quickResp(SkeyFail, c)
		return
	}

	uid := checkSession(skey)

	myLog(fmt.Sprintf("uid: %v", uid))

	if tp == "random" {
		var u user
		u.Id = uid
		Sql.Get(&u)

		var rEmotion emotion

		if u.EmotionNum == 0 {
			myLog("uEm = 0")
			quickResp(NotExist, c)
			return
		}

		var emotionN int64
		if etp == 0 { emotionN = u.GoodmoodNum } else { emotionN = u.BadmoodNum }

		has, _ := Sql.Where("uid = ? and type = ?", uid, etp).Limit(1, int(rand.Int63() % emotionN)).Get(&rEmotion)

		if !has {
			myLog("!has")
			quickResp(NotExist, c)
			return
		}

		fullResp(c, rEmotion)
		return
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
	r.Handle("GET", "/emotion/:id", getEmotion_Id)
	r.Handle("GET", "/emotion", getEmotion)

	r.GET("/emotions", getEmotions)
	r.GET("/src/text/:id", getSrcText_Id)
	r.GET("/src/photo/:id/:nu", getSrcPhoto_Id_Nu)
	r.GET("/src/voice/:id", getSrcVoice_Id)
	r.GET("/src/accept/:id", getSrcAccept_Id)

	r.Handle("POST", "/emotion", postEmotion)
	r.Handle("POST", "/src/voice/:id", postSrcVoice_Id)
	r.Handle("POST", "/src/photo/:id/:num", postSrcPhoto_Id_Num)
	r.Handle("POST", "/emotion/:id", postEmotion_Id)
	r.Handle("POST", "/user/photo", postUserPhoto)

	delFile("src/2", "head")
	go maintainSession()
	go sqlConnectKeepAlife()

	Sessions["1"] = 2
	SessionsLifetime["1"] = time.Now().Unix() * 2

	Router.Run()
}

/*
	BAIWEILIANG
*/
func readBinaryFile(path string) ([]byte, error) {
	t, err := ioutil.ReadFile(path)
	if err != nil {
		myLog(fmt.Sprintf("ERROR when openning %v \n", path))
	}
	return t, err
}

func matchFilePreffix(dirPath string, pre string) string {
	dir, err := ioutil.ReadDir(dirPath)
	if err != nil {return ""}
	for _, file := range dir {
		if strings.HasPrefix(file.Name(), pre) {
			return file.Name()
		}
	}
	return ""
}

//转换64位非负int的Itoa函数
func Itoa64(num int64) string {
	if num == 0 {return "0"}
	str := ""
	for ;num > 0; {
		str += strconv.Itoa(int(num % 10))
		num /= 10
	}
	strune := []rune(str)
	var ans string
	for l := len(strune) - 1;l >= 0;l-- {
		ans += string(strune[l])
	}
	return ans
}

//GET /emotions?skey=&type=&content=&page=&rank=&search=&full=
func getEmotions(c *gin.Context) {
	skey := c.Query("skey")
	uid := checkSession(skey)
	if uid == -1 {
		quickResp(SkeyFail, c)
		return
	}
	dataSourceName := "root:123456@/test?charset=utf8"
	Sql, err := xorm.NewEngine("mysql", dataSourceName)

	ty := c.DefaultQuery("type", "-1")	//-1表示忽略
	content := c.DefaultQuery("content", "-1")
	page := c.DefaultQuery("page", "1")	//Defult: p1
	rank := c.DefaultQuery("rank", "0")	//0 no sort;1 time up;2 star up
	search := c.DefaultQuery("search", "")	//""表示忽略
	full := c.DefaultQuery("full", "0")	//0表示不反回text和accept，1表示全部返回

	//填充sql语句: WHERE & ORDER
	where0, where1, where2, where3 := "", "", "", ""
	where0 = "uid=" + Itoa64(uid) + " "
	switch ty {
	case "0": where1 = "AND type=0 "
	case "1": where1 = "AND type=1 "
	default:
		where1 = ""
	}
	if content != "-1" && content != "" {where2 = "AND content=" + content + " "}
	if search != "" {where3 = "AND text LIKE '%" + search + "%'"}
	order := " ORDER BY"
	switch rank {
	case "1": order += " created_at ASC"
	case "-1": order += " created_at DESC"
	case "2": order += " stars ASC"
	case "-2": order += " stars DESC"
	default:
		order = ""
	}
	//分页
	limit := ""
	pageNum := 0
	if page != "0" {
		pageNum, _ := strconv.Atoi(page)
		if page == "" {pageNum = 1}
		offset := 1
		if pageNum > 0 { offset = (pageNum - 1) * 20 + 1}
		limit = " LIMIT " + strconv.Itoa(offset) + "," + "20"
	}

	//组装sql
	var sql string
	if full == "0" {sql = "SELECT id,stars,type,content,photo_num,brief,created_at,string_created_at FROM emotion WHERE "} else {
		sql = "SELECT id,stars,type,content,photo_num,brief,created_at,string_created_at,text,accept FROM emotion WHERE "
	}
	sql += where0 + where1 + where2 + where3 + order + limit + ";"

	type emotionList struct {
		Id      int64 `json:"id"`
		Stars   int64 `json:"stars"`
		Type int64 `json:"type"`
		Brief string `json:"brief" xorm:varchar(100)`
		Content int64 `json:"content"`
		PhotoNum int64 `json:"photoNum"`
		CreatedAt int64 `json:"createdAt"`
		StringCreatedAt string `json:"stringCreatedAt"`
	}
	type emotionListAll struct {
		Id      int64 `json:"id"`
		Stars   int64 `json:"stars"`
		Type int64 `json:"type"`
		Brief string `json:"brief" xorm:varchar(100)`
		Content int64 `json:"content"`
		PhotoNum int64 `json:"photoNum"`
		CreatedAt int64 `json:"createdAt"`
		StringCreatedAt string `json:"stringCreatedAt"`
		Text string `json:"text"`
		Accept string `json:"accept"`
	}

	//根据full值响应相应的json
	if full == "0" {
		list := make([]emotionList, 0)
		err := Sql.Sql(sql).Find(&list)
		if err != nil {
			quickResp(NotExist, c)
			return
		}
		type results struct {
			Page int64	`json:"page"`
			Num int64 `json:"num"`
			EmotionList []emotionList `json:"emotionList"`
		}
		if len(list) == 0 {pageNum = 1}
		respStruct := results{
			Page:        int64(pageNum),
			Num:         int64(len(list)),
			EmotionList: list,
		}
		fullResp(c, respStruct)
	} else {
		list := make([]emotionListAll, 0)
		err = Sql.Sql(sql).Find(&list)
		if err != nil {
			quickResp(NotExist, c)
			return
		}
		type results struct {
			Page int64	`json:"page"`
			Num int64 `json:"num"`
			EmotionList []emotionListAll `json:"emotionList"`
		}
		if len(list) == 0 {pageNum = 1}
		respStruct := results{
			Page:        int64(pageNum),
			Num:         int64(len(list)),
			EmotionList: list,
		}
		fullResp(c, respStruct)
	}
	return
}

//GET /src/text/:id?skey=
func getSrcText_Id(c *gin.Context) {
	skey := c.Query("skey")
	uid := checkSession(skey)
	if uid == -1 {
		quickResp(SkeyFail, c)
		return
	}
	dataSourceName := "root:123456@/test?charset=utf8"
	Sql, err := xorm.NewEngine("mysql", dataSourceName)
	type emotionList struct {
		Text string `json:"content" xorm:"varchar(2000)"`
		Accept string `json:"content" xorm:"varchar(2000)"`
	}
	var results []string
	err = Sql.Sql("SELECT text FROM emotion WHERE uid=?", uid).Find(&results)
	if err != nil {
		quickResp(NotExist, c)
		return
	}
	c.String(http.StatusOK, results[0])
}

//GET /src/photo/:id/:nu?skey=
func getSrcPhoto_Id_Nu(c *gin.Context) {
	skey := c.Query("skey")
	uid := checkSession(skey)
	if uid == -1 {
		quickResp(SkeyFail, c)
		return
	}
	id := c.Param("id")
	num := c.Param("nu")
	photoDir := "src/" + Itoa64(uid) + "/" + id + "/photo/"
	fileName := matchFilePreffix(photoDir, num)
	if fileName == "" {
		quickResp(NotExist, c)
		return
	}
	photoPath := photoDir + fileName
	photo, err := readBinaryFile(photoPath)
	if err != nil {
		quickResp(NotExist, c)
		return
	}
	suffix := strings.TrimPrefix(fileName, num + ".")
	if suffix == fileName {
		myLog("fail to add suffix of photo")
		suffix = ""
	}
	c.Data(http.StatusOK, "image/" + suffix, photo)
}

//GET /src/voice/:id?skey=
func getSrcVoice_Id(c *gin.Context) {
	skey := c.Query("skey")
	uid := checkSession(skey)
	if uid == -1 {
		quickResp(SkeyFail, c)
		return
	}
	id := c.Param("id")

	voicePath := "src/" + Itoa64(uid) + "/" + id + "/"
	fileName := matchFilePreffix(voicePath, "voice")
	if fileName == "" {
		quickResp(NotExist, c)
		return
	}

	voice, err := readBinaryFile(voicePath + fileName)
	if err != nil {
		quickResp(NotExist, c)
		return
	}
	suffix := strings.TrimPrefix(fileName, "voice.")
	if suffix == fileName {
		myLog("fail to add suffix of voice")
		suffix = ""
	}
	c.Data(http.StatusOK, "audio/" + suffix, voice)
}

//GET /src/accept/:id?skey=
func getSrcAccept_Id(c *gin.Context) {
	skey := c.Query("skey")
	uid := checkSession(skey)
	if uid == -1 {
		quickResp(SkeyFail, c)
		return
	}
	id := c.DefaultQuery("id", "")
	accept := ""
	err := Sql.Select("accept").Where("uid=?", uid).And("id=?", id).Find(&accept)
	if err != nil {
		quickResp(NotExist, c)
		return
	}
	c.String(http.StatusOK, accept)
}

//GET /user/photo?skey
func getUserPhoto(c *gin.Context) {
	skey := c.Query("skey")
	uid := checkSession(skey)
	if uid == -1 {
		quickResp(SkeyFail, c)
		return
	}
	path := "src/" + Itoa64(uid) + "/head"
	head, err := readBinaryFile(path)

	if err != nil {
		quickResp(NotExist, c)
		return
	}
	c.Data(http.StatusOK, "image", head)
}

func delFile(dir string, prefix string) bool {
	fileName := matchFilePreffix(dir, prefix)
	if fileName == "" {
		return false
	}
	err := os.Remove(dir + fileName)
	if err != nil {return false}
	return true
}