package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"io/ioutil"
	"math/rand"
	"strconv"
	"time"
)

type log struct {
	Id int64 `json: "id"`
	Password string `json: "password" `
}

type user struct {
	Id int64 `json:"id"`
	Nick string `json:"nick"`
	EmotionNum int64 `json:"emotionNum"`
	CreatedAt time.Time `json:"createdAt" xorm:"created"`
}

type emotion struct {
	Id      int64 `json:"id"`
	Uid     int64 `json:"uid"`
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
	v, ok := Sessions[str]
	if ok { return v } else { return -1 }
}

	_, ok := Sessions[str]
	if ok {
		return -1
	} else {
		Sessions[str] = id
		return 1
	}
}

func delSession (str string) int64 {
	delete(Sessions, str)
	return 1
}

func myLog(str string)  {
	fmt.Printf("[%v] %v", time.Now().String(), str)
}

func postUser(c *gin.Context)  {
	d, _ := ioutil.ReadAll(c.Request.Body)
	newUser := new(user)

	if json.Unmarshal(d, newUser) == nil {
		myLog(fmt.Sprintf("POST /user\n%v\n", string(d)))
		_, err := Sql.Insert(newUser)
		if err == nil{
			str, _ := json.Marshal(map[string]interface{}{
				"msg": "ok",
				"retc": 0,
			})
			c.String(200, string(str))
		}
	}
}

func getMotion_Id(){

}

func t1() int{
	type a struct {
		owo int
	}
	A := a{10}
	return A.owo
}
func t2() string{
	type a struct {
		qwq string
	}
	A := a{"10"}
	return A.qwq
}

func main() {
	Sql, err := xorm.NewEngine("mysql", "root:123456@/test?charset=utf8")
	Sql.DatabaseTZ, _ = time.LoadLocation("Asia/Shanghai")
	Sql.TZLocation, _ = time.LoadLocation("Asia/Shanghai")
	Router = gin.Default()
	if err != nil {fmt.Printf("%v\n", err)}
	Sql.Sync2(new(user))
	Sql.Sync2(new(log))
	Sql.Sync2(new(emotion))

	str, _ := json.Marshal(map[string]interface{}{
		"msg": "ok",
		"retc": 0,
	})
	fmt.Printf("%v\n", string(str))
	var body map[string]interface{}
	fmt.Printf("%v\n", body)
	json.Unmarshal(str, &body)
	fmt.Printf("%v\n", 1 + int(body["retc"].(float64)))
	fmt.Printf("%v\n", body["msg"].(string))
}