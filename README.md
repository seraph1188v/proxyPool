# 使用腾讯云函数实现代理

## 介绍

通过腾讯云函数实现http代理，客户端监听端口，接收来自浏览器或者其它程序的http请求，封装成自定义的格式，发送到API网关，服务端接收到之后还原出原始的http请求，发起请求并将响应返回给客户端。

## 实现效果

无尽的代理ip，全是白名单

## 使用方式
- 把server中的程序编译，打包放到腾讯云函数上
- 选择API网关触发，获取触发地址
- 修改client中的API网关地址
- 编译运行client