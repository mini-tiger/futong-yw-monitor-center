## 编译
```shell
# Linux 32位
GOOS=linux GOARCH=386 go build -o futongAgent.linux-386 main.go
# Linux 64位
GOOS=linux GOARCH=amd64 go build -o futongAgent.linux-amd64 main.go
# Windows 32位
GOOS=windows GOARCH=386 go build -o futongAgent.windows-386.exe main.go
# Windows 64位
GOOS=windows GOARCH=amd64 go build -o futongAgent.windows-amd64.exe main.go
# Linux ppc64le
GOOS=linux GOARCH=ppc64le go build -o futongAgent.linux-ppc64le main.go
```
