package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
	"sync"
)

type ProxyHandlerNode struct {
	// 节点选项
	options *HandlerOptions
	// 节点处理器列表
	nodes []*ProxyHandlerNode
	// 匹配正则
	regxExpression *regexp.Regexp
	// 匹配计数器， 用于热度排序
	Popularity int
	// 目标服务器URL
	TargetURL *url.URL
	// 目标服务器路径(这段路径需要放到请求Path前面)
	TargetPath string
	// 操作锁
	locker sync.Mutex
	// 本地文件系统
	fileSystem *http.Handler
}

// 配置节点选项
func (node *ProxyHandlerNode) Options(options []*HandlerOptions) {
	node.locker.Lock()
	defer node.locker.Unlock()
	node.nodes = make([]*ProxyHandlerNode, 0)
	for i := 0; i < len(options); i++ {
		opt := options[i]
		if !opt.Enabled {
			continue
		}
		newNode := ProxyHandlerNode{}
		newNode.Popularity = 0
		newNode.options = opt
		var err error
		newNode.regxExpression, err = regexp.Compile(opt.Path)
		if err != nil {
			fmt.Printf("Regexp Compile Error：%v\n", err)
			continue
		}
		newNode.TargetURL, err = url.Parse(opt.Target)
		if err != nil {
			fmt.Printf("URL Parse Error：%v\n", err)
			continue
		}
		if newNode.options.Type == LocalFileSystem {
			var fs = http.FileServer(LocalFile(newNode.options.Target, false))
			newNode.fileSystem = &fs
		}
		newNode.TargetPath = newNode.TargetURL.Path
		newNode.TargetURL.Path = ""
		newNode.Options(opt.SubHandlers)
		node.nodes = append(node.nodes, &newNode)
	}
}

/*
 * 随着匹配次数热度爬升
 */
func popularityUp(nodes []*ProxyHandlerNode, node *ProxyHandlerNode, index int) {
	node.Popularity++
	if index > 0 {
		if node.Popularity > nodes[index-1].Popularity {
			nodes[index], nodes[index-1] = nodes[index-1], nodes[index]
		}
	}
}

/*
 * 匹配请求URL的处理器
 */
func (node *ProxyHandlerNode) MatchNode(url string) *ProxyHandlerNode {
	node.locker.Lock()
	defer node.locker.Unlock()
	for i := 0; i < len(node.nodes); i++ {
		var cnode = node.nodes[i]
		if cnode.regxExpression.Match([]byte(url)) {
			popularityUp(node.nodes, cnode, i)
			var ret = cnode.MatchNode(url)
			if ret != nil {
				return ret
			}
			return cnode
		}
	}
	return nil
}

/*
 * 代理请求
 */
func (node *ProxyHandlerNode) ProxyRequest(req *http.Request, res http.ResponseWriter) {
	// 替换url
	var requestPath = node.regxExpression.ReplaceAllString(req.URL.Path, "")
	// http重定向
	if len(requestPath) == 0 && len(req.URL.RawQuery) == 0 {
		// 把 /dir/app 重定向为 /dir/app/
		if !strings.HasSuffix(req.URL.Path, "/") {
			http.Redirect(res, req, req.URL.Path+"/", http.StatusTemporaryRedirect)
			return
		}
	}

	// 节点是本地文件系统
	if node.fileSystem != nil {
		req.URL.Path = requestPath
		(*node.fileSystem).ServeHTTP(res, req)
		return
	}

	// 修复路径
	req.URL.Path = node.TargetPath + requestPath

	// 代理
	proxy := httputil.NewSingleHostReverseProxy(node.TargetURL)
	proxy.ModifyResponse = func(r *http.Response) error {
		if node.options.DisableCache {
			r.Header.Set("Cache-Control", "no-store,no-cache")
			r.Header.Set("Pragma", "no-store,no-cache")
			r.Header.Set("Expires", "0")
			r.Header.Del("Last-Modified")
		}
		r.Header.Del("Server")
		return nil
	}
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		fmt.Printf("Serve Error %v {%v}\n", r.URL, err)
	}
	if node.options.DisableCache {
		req.Header.Set("Cache-Control", "no-store,no-cache")
		req.Header.Set("Pragma", "no-store,no-cache")
		req.Header.Del("if-modified-since")
		req.Header.Del("if-none-match")
	}
	proxy.ServeHTTP(res, req)
}
