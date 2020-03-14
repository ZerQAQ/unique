package	main

import (
	"encoding/json"
	_ "encoding/json"
	"fmt"
	_ "fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"io/ioutil"
	"strconv"
	"time"
)

func min(a int64, b int64) int64{
	if a < b {return a} else {return b}
}

type article struct {
	Id int64 `json:id,string xorm:pk`
	Title string `json:title`
	Brief string `json:brief`
	Content string `json:content xorm:"varchar(1000)"`
	TimeString string `json:timeString xorm:"varchar(50)"`
	TimeUnix int64 `json:timeUnix`
}

type contentBlock struct {
	Id int64 `xorm:pk`
	Content string `xorm:"varchar(1000)"`
}

var atcEng, cntBloEng *xorm.Engine

func errorHandle(c *gin.Context, v int64){
	if v > 0 {
		c.String(200, fmt.Sprintf(`{"msg": "ok", "id": %d}`, v))
	} else if v == 0 {
		c.String(500, `{"msg": "operation fail, server error."}`)
	} else if v == -1 {
		c.String(404, `{"msg": "not exist."}`)
	}
}
var blockLen = 10

func newAtc(atc *article) int64{
	atc.Id = 0
	atc.TimeString = time.Now().String()[:19]
	atc.TimeUnix = time.Now().Unix()
	if len(atc.Brief) > 100 {atc.Brief = atc.Brief[:100]}
	blen := (len(atc.Content) - 1) / blockLen + 1
	clen := len(atc.Content)
	var idList []int64
	for i := 0; i < blen; i++{
		t := contentBlock{Content: atc.Content[blockLen * i: min(int64(blockLen * (i + 1)), int64(clen))]}
		cntBloEng.Insert(&t)
		idList = append(idList, t.Id)
	}
	t, err1 := json.Marshal(idList)
	atc.Content = string(t)
	_, err2 := atcEng.Insert(atc)
	if err1 == nil && err2 == nil{
		return atc.Id
	}else {return 0}
}

func delAtc(atc *article) int64{
	//fmt.Printf("atcid: %v\n", atc.Id)
	has, _ := atcEng.Get(atc)
	if !has {return -1}
	var clist []int64
	json.Unmarshal([]byte(atc.Content), &clist)
	//fmt.Printf("id:%v\ncnt:%v\n\n", atc.Id, atc.Content)
	//fmt.Printf("clist:%v\n\n", clist)
	for _, i := range clist{
		fmt.Print(i)
		cntBloEng.Id(i).Delete(new(contentBlock))
	}
	_, err := atcEng.Id(atc.Id).Delete(new(article))
	if err == nil{
		return atc.Id
	} else {return 0}
}

func modAtc(atc *article) int64{
	//fmt.Printf("modifying %v ...\n", atc.Id)
	ret1 := delAtc(&article{Id: atc.Id})
	ret2 := newAtc(atc)
	if ret1 > 1 {return ret2} else {return ret1}
}

func getAtc(id int64, tp string, ret *article) int64{
	has, _ := atcEng.Id(id).Get(ret)
	if !has {return -1}
	var idList []int64
	json.Unmarshal([]byte(ret.Content), &idList)
	ret.Content = ""

	if tp == "full" {
		for _, id := range idList{
			temp := &contentBlock{}
			cntBloEng.Id(id).Get(temp)
			ret.Content += temp.Content
		}
	} else{
		temp := &contentBlock{}
		cntBloEng.Id(idList[0]).Get(temp)
		ret.Content += temp.Content
	}
	return 1
}

func main(){

	atcEng, _ = xorm.NewEngine("mysql", "root:123456@/test?charset=utf8")
	cntBloEng, _ = xorm.NewEngine("mysql", "root:123456@/test?charset=utf8")

	atcEng.Sync2(new(article))
	cntBloEng.Sync2(new(contentBlock))

	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	blogRouter := router.Group("/blog")
	router.GET("/", func(c *gin.Context) {
		c.String(200, "hello world")
	})

	blogRouter.POST("", func(c *gin.Context) {

		atc := &article{}
		data, _ := ioutil.ReadAll(c.Request.Body)

		fmt.Printf("body:\n%v\n", string(data))
		if json.Unmarshal(data, atc) != nil {
			c.String(400, fmt.Sprintf(`{"msg": "body format error", "id": %d}`, atc.Id))
		} else {
			errorHandle(c, newAtc(atc))
		}
	})

	blogRouter.POST("/:id", func(c *gin.Context) {
		atc := &article{}
		tp := c.DefaultQuery("type", "new")
		data, _ := ioutil.ReadAll(c.Request.Body)

		idr := c.Param("id")
		id, _ := strconv.Atoi(idr)

		fmt.Printf("body:\n%v\n", string(data))
		if tp == "modify"{
			if json.Unmarshal(data, atc) != nil {
				c.String(400, `{"msg": "body format error"}`)
			} else { atc.Id = int64(id); errorHandle(c, modAtc(atc)) }
		} else if tp == "delete"{
			atc.Id = int64(id)
			errorHandle(c, delAtc(&article{Id: int64(id)}))
		}
	})

	blogRouter.GET("/:id", func(c *gin.Context) {
		tp := c.DefaultQuery("type", "brief")
		idr := c.Param("id")
		fmt.Printf("tp: %v", tp)
		id, _ := strconv.Atoi(idr)
		atc := new(article)

		if getAtc(int64(id), tp, atc) == -1 {
			c.String(404, `{"msg": "not exist."}`)
		} else{
			ret, _ := json.Marshal(*atc)
			c.String(200, `{"msg": "ok.", "data": ` + string(ret) + `}`)
		}
	})

	router.Run()
}
