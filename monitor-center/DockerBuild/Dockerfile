#xxxx 设置基础镜像,Harbor服务器pull镜像
FROM harbor.od.com/library/alpine:3.13.6_tz
#创建app目录
RUN mkdir -p /usr/src/datacenter_go
#设置工作目录
WORKDIR /usr/src/datacenter_go

#复制所有文件到 工作目录。 pwd当前目录是
#config.json jenkins 会拷贝到DockerBuild 下
COPY ./main ./config.json /usr/src/datacenter_go/

#COPY ./config.json /usr/src/datacenter_go/

# 时区 已经打包 到基础镜像
#RUN apk add -U tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && apk del tzdata
#暴露container的端口
EXPOSE 8001
#ENV mysqluser root
#运行命令
#CMD ["sh","-c","./main >server.log 2>&1"]
# 日志打印在 标准输出
CMD ["sh","-c","./main"]