package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type t1 struct {
	Id int64
	T string
}


func main()  {
	t := time.Now()
	fmt.Print(t.String()[:19], " ", t.Unix())
}

