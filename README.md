# Lute HTTP

## 简介

包装 [Lute](https://github.com/b3log/lute) 引擎以 HTTP 服务发布。

## 背景

该项目主要是为了：

* 让不同编程语言编写的应用能够复用标准的（符合 GFM/CommonMark） Markdown 渲染  
* 为 [Sym](https://github.com/b3log/symphony)、[Solo](https://github.com/b3log/solo)、[Pipe](https://github.com/b3log/pipe) 提供更好的 Markdown 渲染，解决各项目内建的 Markdown 渲染不统一的问题

## 使用

1. 安装 golang，然后获取并编译 `go get -u github.com/b3log/lute-http`
2. 编译成功后将在 ${GOPATH}/bin 下生成名为 `lute-http` 的可执行文件
3. 启动 lute-http 后再启动 Solo、Pipe、Sym 即可，如果成功的话 Solo 等的启动日志中会输出 `[Lute] is available`
4. 需要通过 [nohup](https://hacpai.com/man?cmd=nohup) 和 `&` 让进程在后台运行：`nohup ./lute-http > lute-http.log 2>&1 &`

## 授权

Lute HTTP 使用 [木兰宽松许可证, 第1版](http://license.coscl.org.cn/MulanPSL) 开源协议。

## 🙏 鸣谢

* [Lute](https://github.com/b3log/lute)：一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
* [fasthttp](https://github.com/valyala/fasthttp)：用 golang 写的高性能 HTTP 实现
