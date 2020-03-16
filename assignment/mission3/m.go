package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type log struct {
	Id int64 `json: "id" xorm:pk`
	Password string `json: "password" xorm:varchar(64)`
	CreatedAt time.Time `xorm:"created"`
}

type user struct {
	Id int64 `json:"id" xorm:pk`
	Nick string `json:"nick" xorm:varchar(100)`
	EmotionNum int64 `json:"emotionNum"`
}

type emotion struct {
	Id      int64 `json:"id" xorm:pk`
	Uid     int64 `json:"uid" xorm:pk`
	Stars   int64 `json:"stars"`
	Type int64 `json:"type"`
	Content int64 `json:"content"`
	PhotoNum int64 `json:"photoNum"`
	CreatedAt time.Time `json:"createdAt" xorm:"created"`
}

var Router *gin.Engine
var Sql *xorm.Engine
var Sessions = make(map[string]int64)
var SessionsLifetime = make(map[string]int64)

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

func myLog(str string)  {
	fmt.Printf("[%v] %v", time.Now().String(), str)
}

func idExist(id int64) bool{
	ok1, _ := Sql.Id(id).Get(new(log))
	ok2, _ := Sql.Id(id).Get(new(user))
	if ok1 != ok2{
		myLog("ERROR! user and log not match!\n")
	}
	return ok1 && ok2
}

var OK = "ok"
var ServerError = "serverError"
var TypeError = "typeError"
var SkeyFail = "skeyFail"
var Mottos []string
var MottosLen int64

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
	} else if cmd == TypeError{
		c.JSON(400, gin.H{
			"msg": "format error",
			"retc": -4,
		})
	} else if cmd == SkeyFail{
		c.JSON(403, gin.H{
			"msg": "skey fail",
			"retc": -3,
		})
	}
}

func fullResp(c *gin.Context, d *gin.H){
	c.JSON(200, gin.H{
		"msg": "ok",
		"retc": 1,
		"data": &d,
	})
}

func readFile(path string) []byte {
	f, err := os.Open(path)
	t, _ := ioutil.ReadAll(f)
	f.Close()
	if err != nil {
		myLog(fmt.Sprintf("ERROR when openning %v \n", path))
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

	if json.Unmarshal(d, newUser) == nil {
		myLog(fmt.Sprintf("POST /user\n%v\n", string(d)))
		has, _ := Sql.Id(newUser.Id).Get(new(log))

		fmt.Printf("has:%v", has)

		if has { //ID存在

			c.JSON(403, gin.H{
				"msg": "id has already exist",
				"retc": -2,
			})
			return
		}

		Sql.Insert(newUser)
		_, err := Sql.Insert(user{Id: newUser.Id})
		if err == nil{ //ok
			quickResp(OK, c)
		} else { //服务器错误
			fmt.Print("ERROR:\n%v\n", err)
			quickResp(ServerError, c)
		}
	} else {quickResp("typeError", c)}
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
		str, _ := json.Marshal(userData)
		c.JSON(200, gin.H{
			"msg": "ok",
			"retc": 1,
			"data": string(str),
		})
		return
	}
}

func getMotto(c *gin.Context)  {
	fmt.Printf("getmotto\n")
	fullResp(c, &gin.H{
		"content": Mottos[rand.Int63() % MottosLen],
	})
	fmt.Printf("ret\n")
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
	r := Router.Group("/kuro")
	json.Unmarshal(readFile("motto.json"), &Mottos)
	MottosLen = int64(len(Mottos))

	r.Handle("POST", "/user", postUser)
	r.Handle("POST", "/login", postLogin)
	r.Handle("GET", "/user", getUser)
	r.Handle("GET", "/motto", getMotto)

	Router.Run()
}