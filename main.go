package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
)

var rootHandler = &ProxyHandlerNode{}

func handleRequestAndRedirect(res http.ResponseWriter, req *http.Request) {
	//
	var node = rootHandler.MatchNode(req.URL.Path)
	if node != nil {
		node.ProxyRequest(req, res)
		return
	}
	fmt.Printf("Not Handle Request %v\n", req.RequestURI)
	res.WriteHeader(http.StatusBadGateway)
	res.Write([]byte("Siteweb Core Proxy: No suitable forwarding processor matched."))
}

func main() {
	// 初始化处理器
	rootHandler.Options(handlers)
	// 跳过tsl证书验证
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	http.HandleFunc("/", handleRequestAndRedirect)
	// http
	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
