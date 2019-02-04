package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"fmt"
	"strings"
	"math/rand"
	"time"
	"sync"
)

/*
	使用curl命令模拟网页请求
	curl http://localhost:8080/  | iconv -f utf-8 -t gbk
	curl --data "users=david,nancy" http://localhost:8080/import  | iconv -f utf-8 -t gbk
	curl http://localhost:8080/lucky  | iconv -f utf-8 -t gbk
*/

//抽奖者名单
var userList []string
//使用互斥锁
var mu sync.Mutex

//控制器
type lotteryController struct {
	Ctx iris.Context
}

//传控制器到iris的实际web应用
func newApp() *iris.Application {
	app := iris.New()
	mu = sync.Mutex{}
	mvc.New(app.Party("/")).Handle(&lotteryController{})
	return app
}

func main() {
	app := newApp()
	userList = []string{}

	app.Run(iris.Addr(":8080"))
}

//首页返回当前参与抽奖的用户数
func (c *lotteryController) Get() string {
	count := len(userList)
	return fmt.Sprintf("当前总共参与抽奖的用户数：%d\n", count)
}

//POST http://localhost:8080/import
//params: users
func (c *lotteryController) PostImport() string {
	strUsers := c.Ctx.FormValue("users")
	users := strings.Split(strUsers, ",")
	count1 := len(userList)

	//互斥锁加锁、解锁
	mu.Lock()
	defer mu.Unlock()

	for _, u := range users {
		u = strings.TrimSpace(u)
		if len(u) > 0 {
			userList = append(userList, u)
		}
	}
	count2 := len(userList)
	return fmt.Sprintf("当前总共参与抽奖的用户数： %d，成功导入的用户数:%d \n", count2, count2-count1)
}

//Get http://localhost:8080/lucky
func (c *lotteryController) GetLucky() string {

	//互斥锁加锁、解锁
	mu.Lock()
	defer mu.Unlock()

	count := len(userList)
	if count > 1 {
		seed := time.Now().UnixNano()
		index := rand.New(rand.NewSource(seed)).Int31n(int32(count))
		user := userList[index]
		userList = append(userList[0:index], userList[index+1:]...)
		return fmt.Sprintf("当前剩余用户数:%s,剩余用户数:%d\n", user, count-1)
	} else if count == 1 {
		user := userList[0]
		return fmt.Sprintf("当前中奖用户数： %s,剩余用户数: %d\n", user, count-1)
	} else {
		return fmt.Sprintf("已经没有参与用户，请通过/ import导入用户\n")
	}
}
