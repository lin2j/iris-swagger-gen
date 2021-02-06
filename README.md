项目源代码在 `Github` 上：https://github.com/lin2j/iris-swagger-gen

# 下载 protoc 相关的工具

待补充

# 生成文档

将需要的生成的工具下载后，下面就来演示如何生成Swagger文档。

## 项目结构

未展开图

![iris-swagger-gen-3](https://www.lin2j.tech/upload/2021/02/iris-swagger-gen-3-075acafe42f1425cadba12ed38fb97ee.png)

展开图

![iris-swagger-gen-4](https://www.lin2j.tech/upload/2021/02/iris-swagger-gen-4-a9aad9b4ad434f1aa4b0a127e8609282.png)

- `common` 目录存放 `protoc`通过 `proto` 文件生成 `go` 文件。
- `proto` 目录存放请求用到的对象的描述，还有接口的描述。
- `swagger` 目录存放的是 `swagger` 页面的一些 `js`、`css` 、图片等文件，还有关键的 `json` 子目录下的 `json` 文件。`service.swagger.json` 文件是**通过工具生成**的。目录下的文件，除了 `index.html` 其他的都不用动。
- `proto.sh` 脚本的作用是读取 `proto` 文件，然后生成 `go` 文件到指定目录。
- `swagger.sh` 脚本的作用的读取 `proto` 文件和 `yaml` 文件，生成一个 `json` 文件到指定目录。

## 文件的内容

### `proto` 文件

`message.proto` 写的是一些请求参数和消息体的结构。

```protobuf
syntax = "proto3";

option go_package="message";
package message;

message Request {
  // 消息
  string Msg = 1;
}

message Response {
  // 响应内容
  string Msg = 1;
}
```

 `service.protp` 写的是请求接口的输入输出以及注释

```protobuf
syntax = "proto3";

import "message.proto";
package message;

service EchoService {
  /*
    Get 请求回显发送的请求内容
   */
  rpc Echo(message.Request) returns (message.Response) {}

  /*
    Post 请求回显发送的请求内容
   */
  rpc Echo2(message.Request) returns (message.Response) {}

  /*
    Delete 请求回显发送的请求内容
 */
  rpc Echo3(message.Request) returns (message.Response) {}
}
```

`service.yaml` 写的是请求接口的请求规则，比如请求方法、请求体等（还不太了解这些规则）

```yaml
type: google.api.service
config_version: 3

http:
  rules:
    - selector: message.EchoService.Echo
      get: /echo
    - selector: message.EchoService.Echo2
      post: /echo2
      # 不要漏掉了双引号
      body: "*"
    - selector: message.EchoService.Echo3
      delete: /echo3
      body: "*"

```

三个文件编写完毕，可以通过命令生成对应的 `go` 文件和 `swagger` 的 `json` 文件。这里直接写成脚本。

### `Shell` 脚本

`proto.sh` 

```shell
#!/bin/bash

mkdir -p common/message

#
protoc --proto_path=./proto --go_out=./common/message ./proto/message.proto
#
sed -i "s/,omitempty//g" ./common/message.pb.go
```

`swagger.sh`

```shell
#!/bin/bash

#CUR=$(pwd)
SOURCE_DIR="proto"
TARGET_DIR="swagger/json"
# 可以使用 -I 参数指定 GOOGLE_API 的目录
# https://github.com/googleapis/googleapis
# GOOGLE_API="googleapis-common-protos-1_3_1"
# 使用 yml 配置接口，可以不用引入 GOOGLE_API
# allow_delete_body=true  可以让 delete 请求支持从请求体获取数据
OUT_OPTS="allow_delete_body=true,logtostderr=true,grpc_api_configuration=$SOURCE_DIR/service.yml:$TARGET_DIR"

echo "protoc -I $SOURCE_DIR --swagger_out=$OUT_OPTS $SOURCE_DIR/service.proto"
protoc -I "$SOURCE_DIR" --swagger_out="$OUT_OPTS" "$SOURCE_DIR/service.proto"

VERSION=$(cat version)
sed -i "s/version not set/$VERSION/g" swagger/json/service.swagger.json
```

这里的路径使用的是相对路径，因此执行时，需要注意以下当前目录。

### 脚本的运行以及生成文件的配置

脚本写好后，执行以下。

```bash
sh proto.sh
sh swagger.sh
```

`sh proto.sh` 执行之后，会在 `common` 目录下新增一份 `go` 文件。

`sh swagger.sh` 执行之后，会在 `swagger/json` 下生成一份 `swagger.json` 文件。生成的文件需要配置到 `index.html` 中，才能够在页面上正常显示。

![iris-swagger-gen-5](https://www.lin2j.tech/upload/2021/02/iris-swagger-gen-5-1313d9d0de7e4d9bb7277f038494bb0c.png)

对应页面上的显示就是下图。

![iris-swagger-gen-6](https://www.lin2j.tech/upload/2021/02/iris-swagger-gen-6-8cfdc0c9711d402eb8c2bc4d702ae201.png)

### `Go` 代码

接下来是`go` 程序的编码，使用的是 `iris` 框架。将 `swagger` 目录映射到 `localhost:8080/swagger` 接口上，这样在浏览器可以看到 `swagger` 界面。

先拉取 `iris` 的依赖

```bash
go get github.com/kataras/iris/v12
```

```go
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
	// 从url参数、或者表单中获取数据
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

```

运行起来之后，会打印接口的基本信息，如下图。

![iris-swagger-gen-7](https://www.lin2j.tech/upload/2021/02/iris-swagger-gen-7-83ac9e29ff6a4993bd413ed3996609e3.png)

成功运行后，到浏览器访问 `localhost:8080/swagger` 接口，就可以看到`swagger`界面了

# 程序启动前运行脚本

开发过程中经常要增加或者修改接口的信息，如果每次修改之后，都要手动执行一遍 `swagger.sh` 未免显得繁琐。因此可以在借助 `Goland` 的功能，实现：在程序启动前，先执行一遍 `swagger.sh`  脚本。

如下图

![iris-swagger-gen-8](https://www.lin2j.tech/upload/2021/02/iris-swagger-gen-8-c0e5d785f0634b0eb80d987f0bf1e542.png)

`Run External tool` 的内容如下

![iris-swagger-gen-1](https://www.lin2j.tech/upload/2021/02/iris-swagger-gen-1-8224044f7fd1464399d0c1ccfeb01e61.png)

需要注意的是 `sh.exe` 是 `git bash` 的，`Working directory` 需要指定到脚本的所在目录，`Arguments` 配置为  `swagger.sh` 脚本。这样就可以了，配置完点击 `OK` 就好了。

接下来，尝试修改接口信息，可以是接口的注释说明，然后通过 `Goland` 直接启动程序成功后，再次访问 `localhost:8080/swagger` ，可以看到界面上的信息是修改后的最新信息。

# 结语

前段时间，因为项目需要，从 `Java` 转到 `Go` 开发，刚使用 `Go` 开发觉得还是不太舒服，现在也没觉得有多舒服，主要是`Java` 用惯了。刚开始的时候，连项目的目录结构我都是自己摸索的，仿照着 `Java` 的目录结构，我也为 `Go` 定制了一套属于自己的结构，以及各种编译打包的脚本。`swagger` 生成只是其中的一环。我也试过用网上的，直接在`Go`代码的请求接口增加注释来生成 `swagger` 文档，自己看着着实不太优雅。也可能是自己不熟悉，还不会用。所以综合之下，采用了通过 `proto` 生成文档的方法。也许将来我会发现更好的方式去生成接口文档，那么到时我再记录一下新的探索过程。