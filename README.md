# DevOps
让你轻松快速拥有DevOps环境
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

# 设置工作目录
WORKDIR /app

# 从第一个构建阶段复制编译好的可执行文件到容器中
COPY --from=builder /go/src/devops-go .

# 复制HTML文件到Nginx的默认静态文件目录
COPY html /usr/share/nginx/html

# 暴露应用程序的端口
EXPOSE 22333

# 启动Nginx和Go程序
CMD ["sh", "-c", "nginx -g 'daemon off;' & ./devops-go"]
EOF
```

```shell
docker build -t devops-go:v1 .
```

```shell
docker run -d --name nginx-go \
    --restart=always \
    --privileged=true \
    -p 8080:80 \
    -p 22333:22333 \
    devops-go:v1
```
