# DevOps环境
# 大佬们，帮忙点个star，谢谢（加群：WeChat：Linux-LAMP）
## 本项目能让你轻松快速拥有DevOps环境

## 已经做成腾讯云镜像，直接运行下面命令（体验本项目）推荐使用（alpha版本）
```shell
docker run -d --name nginx-go-alpha \
    --restart=always \
    --privileged=true \
    -p 8090:80 \
    -p 22333:22333 \
    ccr.ccs.tencentyun.com/huanghuanhui/go:go-web-terminal-alpha
```
> 浏览器访问：ip+3333
## 已经做成腾讯云镜像，直接运行下面命令（体验本项目）（beta版本）
```shell
docker run -d --name nginx-go-beta \
    --restart=always \
    --privileged=true \
    -p 8090:80 \
    -p 22333:22333 \
    ccr.ccs.tencentyun.com/huanghuanhui/go:go-web-terminal-beta
```
> 浏览器访问：ip+8090

> 两者区别：alpha版本默认初始化的终端是空，beta版本默认使用容器内部的终端（推荐使用（alpha版本））
## 也可以按照下面的dockerfile自定义镜像
```shell
cat > dockerfile << 'EOF'
# 第一个构建阶段：编译Go应用程序
FROM golang:1.21.1-alpine3.18 as builder

# 在容器中切换到root用户
USER root

# 设置工作目录
WORKDIR /go/src

# 复制应用程序代码到容器中
COPY devops-go.go /go/src

# 配置Go模块和代理
RUN go env -w GO111MODULE=on && go env -w GOPROXY=https://goproxy.cn,direct

# 初始化Go模块并下载依赖
RUN go mod init devops-go 
RUN go mod tidy

# 编译Go应用程序为可执行文件
RUN go build -o devops-go

# 第二个构建阶段：创建最终的镜像
FROM nginx:1.25.2-alpine

# 将Alpine软件包仓库地址替换为国内源（阿里云源）
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 更新Alpine的软件包列表并安装OpenSSH客户端
RUN apk update && \
    apk add openssh-client && \
    rm -rf /var/cache/apk/*

# 设置工作目录
WORKDIR /root

# 从第一个构建阶段复制编译好的可执行文件到容器中
COPY --from=builder /go/src/devops-go .

# 复制HTML文件到Nginx的默认静态文件目录
COPY html /usr/share/nginx/html

# 暴露应用程序的端口
EXPOSE 22333

# 启动Nginx和Go程序
CMD ["sh", "-c", "nginx -g 'daemon off;' & ./devops-go "]
EOF
```

```shell
docker build -t devops-go:v1 .
```

```shell
docker run -d --name nginx-go \
    --restart=always \
    --privileged=true \
    -p 8090:80 \
    -p 22333:22333 \
    devops-go:v1
```

<img width="1547" alt="image" src="https://github.com/kubelsp/DevOps/assets/79031894/7983edac-4551-4519-a162-5f7221d30a63">
<img width="1547" alt="image" src="https://github.com/kubelsp/DevOps/assets/79031894/702d660f-38a2-4b9d-bf95-52720de23433">
<img width="1547" alt="image" src="https://github.com/kubelsp/DevOps/assets/79031894/48856fa8-c396-4df7-82b6-a1d331bdc27e">

<img width="1547" alt="image" src="https://github.com/kubelsp/DevOps/assets/79031894/fa7e6901-8a6d-40de-9eb5-3dee50e5ecbf">

