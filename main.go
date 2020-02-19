// Lute HTTP - HTTP Server for Lute.
// Copyright (c) 2019-present, b3log.org
//
// Lute HTTP is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package main

import (
	"os"

	"github.com/88250/gulu"
	"github.com/88250/lute"
	"github.com/valyala/fasthttp"
)

var logger = gulu.Log.NewLogger(os.Stdout)

// handleMarkdown2HTML 处理 Markdown 转 HTML。
// POST 请求 Body 传入 Markdown 原文；响应 Body 是处理好的 HTML。
func handleMarkdown2HTML(ctx *fasthttp.RequestCtx) {
	body := ctx.PostBody()

	engine := lute.New()

	CodeSyntaxHighlightLineNum := string(ctx.Request.Header.Peek("X-CodeSyntaxHighlightLineNum"))
	if "true" == CodeSyntaxHighlightLineNum {
		engine.CodeSyntaxHighlightLineNum = true
	} else if "false" == CodeSyntaxHighlightLineNum {
		engine.CodeSyntaxHighlightLineNum = false
	}

	CodeSyntaxHighlightDetectLang := string(ctx.Request.Header.Peek("X-CodeSyntaxHighlightDetectLang"))
	if "true" == CodeSyntaxHighlightDetectLang {
		engine.CodeSyntaxHighlightDetectLang = true
	} else if "false" == CodeSyntaxHighlightDetectLang {
		engine.CodeSyntaxHighlightDetectLang = false
	}

	ToC := string(ctx.Request.Header.Peek("X-ToC"))
	if "true" == ToC {
		engine.ToC = true
	} else if "false" == ToC {
		engine.ToC = false
	}

	Footnotes := string(ctx.Request.Header.Peek("X-Footnotes"))
	if "true" == Footnotes {
		engine.Footnotes = true
	} else if "false" == Footnotes {
		engine.Footnotes = false
	}

	AutoSpace := string(ctx.Request.Header.Peek("X-AutoSpace"))
	if "true" == AutoSpace {
		engine.AutoSpace = true
	} else if "false" == AutoSpace {
		engine.AutoSpace = false
	}

	FixTermTypo := string(ctx.Request.Header.Peek("X-FixTermTypo"))
	if "true" == FixTermTypo {
		engine.FixTermTypo = true
	} else if "false" == FixTermTypo {
		engine.FixTermTypo = false
	}

	ChinesePunct := string(ctx.Request.Header.Peek("X-ChinesePunct"))
	if "true" == ChinesePunct {
		engine.ChinesePunct = true
	} else if "false" == ChinesePunct {
		engine.ChinesePunct = false
	}

	IMADAOM := string(ctx.Request.Header.Peek("X-IMADAOM"))
	if "true" == IMADAOM {
		engine.InlineMathAllowDigitAfterOpenMarker = true
	} else if "false" == IMADAOM {
		engine.InlineMathAllowDigitAfterOpenMarker = false
	}

	html, err := engine.Markdown("", body)
	if nil != err {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		logger.Errorf("markdown text [%s] failed: %s\n", body, err.Error())
		return
	}
	ctx.SetBody(html)
}

// handleMarkdownFormat 处理 Markdown 格式化。
// POST 请求 Body 传入 Markdown 原文；响应 Body 是格式化好的 Markdown 文本。
func handleMarkdownFormat(ctx *fasthttp.RequestCtx) {
	body := ctx.PostBody()

	engine := lute.New()
	formatted, err := engine.Format("", body)
	if nil != err {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		logger.Errorf("format markdown text [%s] failed: %s\n", body, err.Error())
		return
	}
	ctx.SetBody(formatted)
}

// handleHtml 处理 HTML 转 Markdown。
// POST 请求 Body 传入 HTML；响应 Body 是处理好的 HTML。
func handleHtml(ctx *fasthttp.RequestCtx) {
	body := ctx.PostBody()

	engine := lute.New()
	html, err := engine.HTML2Markdown(gulu.Str.FromBytes(body))
	if nil != err {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		logger.Errorf("html [%s] to markdown failed: %s\n", body, err.Error())
		return
	}
	ctx.SetBody(gulu.Str.ToBytes(html))
}

// handle 处理请求分发。
func handle(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/", "":
		handleMarkdown2HTML(ctx)
	case "/format":
		handleMarkdownFormat(ctx)
	case "/html":
		handleHtml(ctx)
	default:
		ctx.SetStatusCode(fasthttp.StatusNotFound)
	}
}

// Lute 的 HTTP Server 入口点。
func main() {
	gulu.Log.SetLevel("info")

	addr := ":8249"
	logger.Infof("booting Lute HTTP on [%s]", addr)
	server := &fasthttp.Server{
		Handler:            handle,
		MaxRequestBodySize: 1024 * 1024 * 2, // 2MB
	}
	err := server.ListenAndServe(addr)
	if nil != err {
		logger.Fatalf("booting Lute HTTP server failed: %s", err.Error())
	}
}
