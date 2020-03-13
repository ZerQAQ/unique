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
)

type article struct {
	Id int64 `json:id,string xorm:pk`
	Title string `json:title`
	Brief string `json:brief`
	Content string `json:content`
}

type contentBlock struct {
	Id int64 `xorm:pk`
	Content string `xorm:"varchar(1000)"`
}

var atcEng, cntBloEng *xorm.Engine

func errorHandle(c *gin.Context, err error){
	if err == nil {
		c.String(200, `{"msg": "ok"}`)
	} else{
		c.String(500, `{"msg": "operation fail, server error."}`)
	}
}

func newAtc(atc article){

}

func modAtc(atc article){

}

func delAtc(atc article){

}

func main(){

	atcEng, _ := xorm.NewEngine("mysql", "root:123456@/test?charset=utf8")
	cntBloEng, _ := xorm.NewEngine("mysql", "root:123456@/test?charset=utf8")

	atcEng.Sync2(new(article))
	cntBloEng.Sync2(new(contentBlock))

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "hello world")
	})
	
	r.POST("/blog", func(c *gin.Context) {

		atc := &article{}
		tp := c.DefaultQuery("type", "new")
		data, _ := ioutil.ReadAll(c.Request.Body)

		fmt.Printf("body:\n%v\n", string(data))

		if json.Unmarshal(data, atc) != nil {
			c.String(400, `{"msg": "body format error"}`)
		} else{
			if tp == "new"{
				ret, err := atcEng.Insert(atc)
				errorHandle(c, err)
				fmt.Printf("n req: ret:%v err:%v\n\n", ret, err);
			} else if tp == "modify"{
				ret, err := atcEng.Id(atc.Id).Update(atc)
				errorHandle(c, err)
				fmt.Printf("m req: ret:%v err:%v\n\n", ret, err);
			} else if tp == "delete"{
				ret, err := atcEng.Id(atc.Id).Delete(atc)
				errorHandle(c, err)
				fmt.Printf("d req: ret:%v err:%v\n\n", ret, err);
			}
		}
	})

	r.Run()
}
