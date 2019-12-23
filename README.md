# 短域名服务

一个简单的短域名生成跳转服务

![](http://source.qiniu.cnd.nsini.com/images/2019/11/48/bd/64/20191125-bfefea2da3a147e7616cfc58bd348c0b.jpeg?imageView2/2/w/1280/interlace/0/q/70)

## 演示

演示地址: [https://r.nsini.com](https://r.nsini.com)

该演示平台部署在 开普勒云平台上，若您有兴趣可能访问

- GitHub: [https://github.com/kplcloud/kplcloud](https://github.com/kplcloud/kplcloud)
- Demo: [https://kplcloud.nsini.com/about.html](https://kplcloud.nsini.com/about.html)

## 安装说明

平台后端基于[go-kit](https://github.com/go-kit/kit)、前端基于[ant-design](https://github.com/ant-design/ant-design)(版本略老)框架进行开发。

后端所使用到的依赖全部都在[go.mod](go.mod)里，前端的依赖在`package.json`，详情的请看`yarn.lock`，感谢开源社区的贡献。

前端: [https://github.com/icowan/shorter-view](https://github.com/icowan/shorter-view)

### 安装教程

该服务支持多种方式启动

- docker-compose: `$ cd install/docker-compose/ && docker-compose up`
- kubernetes: `$ kubectl apply -f install/kubernetes/`
- localhost: `$ make run`

### 依赖

- Golang 1.13+ [安装手册](https://golang.org/dl/)
- Docker 18.x+ [安装](https://docs.docker.com/install/)
- Mongo/Redis (主要用于存储短链信息)

## 快速开始

1. 克隆

```
$ mkdir -p $GOPATH/src/github.com/icowan
$ cd $GOPATH/src/github.com/icowan
$ git clone https://github.com/icowan/shorter.git
$ cd shorter
```

2. make 启动

```
$ make run
```

### 支持我

![微信](https://lattecake.oss-cn-beijing.aliyuncs.com/static%2Fimages%2Freward%2Fweixin-RMB-xxx.JPG)
![支付宝](https://lattecake.oss-cn-beijing.aliyuncs.com/static%2Fimages%2Freward%2Falipay-RMB-xxx.png)
