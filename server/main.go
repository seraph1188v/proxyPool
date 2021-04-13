package main

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/tencentyun/scf-go-lib/cloudfunction"
	events "github.com/tencentyun/scf-go-lib/events"
)

//收到的数据结构
type DefineEvent struct {
	URL     string `json:"url"`
	Content string `json:"content"`
}

//返回到代理端的数据结构
type RespEvent struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
	Data   string `json:"data"`
}

func hello(ctx context.Context, event events.APIGatewayRequest) (events.APIGatewayResponse, error) {

	var d DefineEvent
	var r RespEvent
	//读取url和http原始报文
	if err := json.Unmarshal([]byte(event.Body), &d); err != nil {
		log.Fatalf("%s \n", err)
	}
	}
	reqURL := d.URL

	rawReq, err := base64.StdEncoding.DecodeString(d.Content)
	if err != nil{
		log.Fatalf("%s \n", err)
	}
	}
	// 从原始报文创建 request 对象
	req, err := http.ReadRequest(bufio.NewReader(strings.NewReader(string(rawReq))))
	if err != nil {
		log.Fatalf("%s \n", err)
	}
	}
	req.RequestURI = ""
	url, err := url.Parse(reqURL)
	if err != nil {
		log.Fatalf("%s \n", err)
	}
	}
	req.URL = url
	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("%s \n", err)
	}
	}
	defer resp.Body.Close()

	// //获取返回的报文
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Fatalf("%s \n", err)
	}
	}
	//拼装返回客户端的数据
	res := base64.StdEncoding.EncodeToString(dump)
	r.Status = true
	r.Error = ""
	r.Data = res
	respByte, err := json.Marshal(r)
	if err != nil {
		log.Fatalf("%s \n", err)
	}
	}

	ret := events.APIGatewayResponse{}
	ret.IsBase64Encoded = true
	ret.StatusCode = 200
	ret.Body = string(respByte)
	return ret, nil
}

func main() {
	// Make the handler available for Remote Procedure Call by Cloud Function
	cloudfunction.Start(hello)
}
