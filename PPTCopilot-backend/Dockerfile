FROM golang:1.20

WORKDIR /home/tmp

# Set environment variables
ENV GO111MODULE=on \
	CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64

# Install MySQL client and other dependencies
RUN apt-get update && \
	apt-get install -y --no-install-recommends default-mysql-client && \
	apt-get clean && \
	rm -rf /var/lib/apt/lists/*

WORKDIR /home/app

# 预先复制/缓存go.mod以预先下载依赖项，并且仅在后续构建中重新下载它们（如果它们发生变化）
COPY go.mod go.sum ./
# 下载bee工具以及依赖
RUN go env -w GO111MODULE=on && go env -w GOPROXY=https://goproxy.cn
RUN go install github.com/beego/bee/v2@latest && go mod download && go mod verify

COPY . .

# 如果没有docker-compose未传递，使用默认值host.docker.internal
ARG MYSQL_HOST=host.docker.internal
ARG MYSQL_PORT=3306
ENV MYSQL_HOST=${MYSQL_HOST}
ENV MYSQL_PORT=${MYSQL_PORT}

# 如果arg server_ip不为空，则替换配置文件中的server_ip
ARG SERVER_IP
# 运行env.py传递参数
RUN python3 env.py $SERVER_IP

CMD ["bee", "run", "-gendoc=true", "-downdoc=true"]
