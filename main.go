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
	"bytes"
	"encoding/json"
	"os"
	"strings"

	"github.com/88250/gulu"
	"github.com/88250/lute"
	"github.com/88250/lute/util"
	"github.com/valyala/fasthttp"
)

var logger = gulu.Log.NewLogger(os.Stdout)

// handleTextBundle 处理 Markdown TextBundle。
// POST 请求 Body 传入 Markdown 原文；响应 Body 是 TextBundle 化后的 Markdown 文本。
func handleTextBundle(ctx *fasthttp.RequestCtx) {
	body := ctx.PostBody()
	engine := newLute()
	linkPrefixesStr := string(ctx.Request.Header.Peek("X-TextBundle-LinkPrefixes"))
	linkPrefixes := strings.Split(linkPrefixesStr, ",")
	md, links := engine.TextBundleStr("", util.BytesToStr(body), linkPrefixes)
	result := map[string]interface{}{
		"markdown":      md,
		"originalLinks": links,
	}
	resultBody, err := json.Marshal(result)
	if nil != err {
		ctx.Response.SetStatusCode(500)
		return
	}
	ctx.Response.Header.SetContentType("application/json; charset=utf-8")
	ctx.SetBody(resultBody)
}

// handleMarkdown2HTML 处理 Markdown 转 HTML。
// POST 请求 Body 传入 Markdown 原文；响应 Body 是处理好的 HTML。
func handleMarkdown2HTML(ctx *fasthttp.RequestCtx) {
	body := ctx.PostBody()

	if bytes.Contains(body, []byte("测试发布到链滴")) {
		logger.Infof("body: %s", body)
	}

	engine := newLute()

	CodeSyntaxHighlightLineNum := string(ctx.Request.Header.Peek("X-CodeSyntaxHighlightLineNum"))
	if "true" == CodeSyntaxHighlightLineNum {
		engine.SetCodeSyntaxHighlightLineNum(true)
	} else if "false" == CodeSyntaxHighlightLineNum {
		engine.SetCodeSyntaxHighlightLineNum(false)
	}

	CodeSyntaxHighlightDetectLang := string(ctx.Request.Header.Peek("X-CodeSyntaxHighlightDetectLang"))
	if "true" == CodeSyntaxHighlightDetectLang {
		engine.SetCodeSyntaxHighlightDetectLang(true)
	} else if "false" == CodeSyntaxHighlightDetectLang {
		engine.SetCodeSyntaxHighlightDetectLang(true)
	}

	ToC := string(ctx.Request.Header.Peek("X-ToC"))
	if "true" == ToC {
		engine.SetToC(true)
	} else if "false" == ToC {
		engine.SetToC(false)
	}

	Footnotes := string(ctx.Request.Header.Peek("X-Footnotes"))
	if "true" == Footnotes {
		engine.SetFootnotes(true)
	} else if "false" == Footnotes {
		engine.SetFootnotes(false)
	}

	AutoSpace := string(ctx.Request.Header.Peek("X-AutoSpace"))
	if "true" == AutoSpace {
		engine.SetAutoSpace(true)
	} else if "false" == AutoSpace {
		engine.SetAutoSpace(false)
	}

	FixTermTypo := string(ctx.Request.Header.Peek("X-FixTermTypo"))
	if "true" == FixTermTypo {
		engine.SetFixTermTypo(true)
	} else if "false" == FixTermTypo {
		engine.SetFixTermTypo(false)
	}

	HeadingID := string(ctx.Request.Header.Peek("X-HeadingID"))
	if "true" == HeadingID {
		engine.SetHeadingID(true)
	} else if "false" == HeadingID {
		engine.SetHeadingID(false)
	}

	IMADAOM := string(ctx.Request.Header.Peek("X-IMADAOM"))
	if "true" == IMADAOM {
		engine.SetInlineMathAllowDigitAfterOpenMarker(true)
	} else if "false" == IMADAOM {
		engine.SetInlineMathAllowDigitAfterOpenMarker(false)
	}

	ParagraphBeginningSpace := string(ctx.Request.Header.Peek("X-ParagraphBeginningSpace"))
	if "true" == ParagraphBeginningSpace {
		engine.SetChineseParagraphBeginningSpace(true)
	} else if "false" == ParagraphBeginningSpace {
		engine.SetChineseParagraphBeginningSpace(false)
	}

	html := engine.Markdown("", body)
	ctx.SetBody(html)
}

// handleMarkdownFormat 处理 Markdown 格式化。
// POST 请求 Body 传入 Markdown 原文；响应 Body 是格式化好的 Markdown 文本。
func handleMarkdownFormat(ctx *fasthttp.RequestCtx) {
	body := ctx.PostBody()
	engine := newLute()
	engine.ParseOptions.ImgPathAllowSpace = true
	formatted := engine.Format("", body)
	ctx.SetBody(formatted)
}

// handleHtml 处理 HTML 转 Markdown。
// POST 请求 Body 传入 HTML；响应 Body 是处理好的 HTML。
func handleHtml(ctx *fasthttp.RequestCtx) {
	body := ctx.PostBody()
	engine := newLute()
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
	defer gulu.Panic.Recover(nil)
	switch string(ctx.Path()) {
	case "/", "":
		handleMarkdown2HTML(ctx)
	case "/format":
		handleMarkdownFormat(ctx)
	case "/html":
		handleHtml(ctx)
	case "/textbundle":
		handleTextBundle(ctx)
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
		MaxRequestBodySize: 1024 * 1024 * 4, // 4MB
	}
	err := server.ListenAndServe(addr)
	if nil != err {
		logger.Fatalf("booting Lute HTTP server failed: %s", err.Error())
	}
}

func newLute() (ret *lute.Lute) {
	ret = lute.New()
	ret.ParseOptions.ImgPathAllowSpace = true
	return
}
