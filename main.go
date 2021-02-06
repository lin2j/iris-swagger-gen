package main

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"log"
	"swagger.gen/common/message"
)

func main() {
	app := iris.New()
	// 演示三种请求2
	app.Get("/echo", handler)
	app.Post("/echo2", handler2)
	app.Delete("/echo3", handler2)
	// swagger 页面的 url，因为目录下有个 index.html 文件
	// 打开页面后会自动加载
	app.HandleDir("/swagger/", "./swagger")
	// 打印接口的基本信息
	for _, r := range app.GetRoutes() {
		log.Println(fmt.Sprintf("%-8s", r.Method) + r.Path)
	}
	if err := app.Run(iris.Addr(":8080")); err != nil {
		log.Fatal(err)
	}
}

func handler(ctx context.Context) {
	r := message.Request{}
	// 从url参数、表单中获取数据
	_ = ctx.ReadForm(&r)
	log.Println("receive message:", r.Msg)
	resp := message.Response{
		Msg: r.Msg,
	}
	_, _ = ctx.JSON(&resp)
}

func handler2(ctx context.Context) {
	r := message.Request{}
	// 从请求体中获取数据
	_ = ctx.ReadJSON(&r)
	log.Println("receive message:", r.Msg)
	resp := message.Response{
		Msg: r.Msg,
	}
	_, _ = ctx.JSON(&resp)
}
