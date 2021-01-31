package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"log"
	"swagger.gen/common/message"
)

func main() {
	app := iris.New()
	app.Get("/echo", handler)
	app.Post("/echo2", handler2)
	// swagger 页面的 url，因为目录下有个 index.html 文件
	// 打开页面后会自动加载
	app.HandleDir("/swagger/", "./swagger")
	if err := app.Run(iris.Addr(":8080")); err != nil {
		log.Fatal(err)
	}
}

func handler(ctx context.Context) {
	r := message.Request{}
	// 处理 Get 请求，从url参数中获取数据
	_ = ctx.ReadForm(&r)
	log.Println("receive message:", r.Msg)
	resp := message.Response{
		Msg: r.Msg,
	}
	_, _ = ctx.JSON(&resp)
}

func handler2(ctx context.Context) {
	r := message.Request{}
	// 处理 Post 请求，从请求体中获取数据
	_ = ctx.ReadJSON(&r)
	log.Println("receive message:", r.Msg)
	resp := message.Response{
		Msg: r.Msg,
	}
	_, _ = ctx.JSON(&resp)
}
