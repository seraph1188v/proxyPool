package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/tencentyun/scf-go-lib/events"
)

type Proxy struct {
}

//向API网关发送的数据结构
type DefineEvent struct {
	URL     string `json:"url"`
	Content string `json:"content"`
}

//接收到的数据结构
type RespEvent struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
	Data   string `json:"data"`
}

//API网关地址
var (
	proxyAddress string = ""
)

func (p *Proxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("接受请求 %s %s %s\n", req.Method, req.Host, req.RemoteAddr)

	reqDump, err := httputil.DumpRequest(req, true)
	reqDumpEncode := base64.StdEncoding.EncodeToString(reqDump)
	if err != nil {
		log.Fatalf("%s \n", err)
	}
	d := &DefineEvent{
		URL:     req.RequestURI,
		Content: reqDumpEncode,
	}
	dv, err := json.Marshal(&d)
	if err != nil {
		log.Fatalf("%s \n", err)
	}
	requestGW, err := http.NewRequest("POST", proxyAddress, bytes.NewReader(dv))
	if err != nil {
		log.Fatalf("%s \n", err)
	}

	client := &http.Client{}
	resp, err := client.Do(requestGW)
	if err != nil {
		log.Fatalf("%s \n", err)
	}
	defer resp.Body.Close()

	// 第二步： 把新请求复制到服务器端，并接收到服务器端返回的响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("%s \n", err)
	}

	var respbody events.APIGatewayResponse
	err = json.Unmarshal([]byte(body), &respbody)
	if err != nil {
		log.Fatalf("%s \n", err)
	}
	dumpBody := respbody.Body
	var r RespEvent
	err = json.Unmarshal([]byte(dumpBody), &r)
	if err != nil {
		log.Fatalf("%s \n", err)
	}
	retByte, err := base64.StdEncoding.DecodeString(r.Data)
	if err != nil {
		log.Fatalf("%s \n", err)
	}
	proxyresp, err := http.ReadResponse(bufio.NewReader(strings.NewReader(string(retByte))), req)
	if err != nil {
		log.Fatalf("%s \n", err)
	}
	for key, value := range proxyresp.Header {
		for _, v := range value {
			w.Header().Add(key, v)
		}
	}
	io.Copy(w, proxyresp.Body)
}

func main() {
	fmt.Println("Serve on :8080")
	http.Handle("/", &Proxy{})
	http.ListenAndServe("0.0.0.0:8080", nil)
}

// 代码运行之后，会在本地的 8080 端口启动代理服务。修改浏览器的代理为 127.0.0.1：:8080
// 再访问网站，可以验证代理正常工作，也能看到它在终端打印出所有的请求信息。
