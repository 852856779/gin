# 指定基础镜像  
FROM golang:1.22.1
# 设置工作目录  
WORKDIR /go/src/app

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOPROXY=https://goproxy.cn,direct
ENV GOOS=linux
ENV GOARCH=amd64
RUN go mod init demo
RUN go mod tidy
# 将当前目录内容（除了go.mod和go.sum）复制到容器的/go/src/app目录下  
COPY . .   
# 安装项目依赖  
RUN go mod download  
# 安装Gin框架
RUN go get github.com/gin-gonic/gin
# 添加所有文件到工作目录
ADD . /go/src/app
# 运行命令
# CMD ["go", "run", "main.go"]