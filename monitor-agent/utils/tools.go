package utils

import (
	"errors"
	"fmt"
	"futong-yw-monitor-center/monitor-agent/g"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

/**
 * @Author: Tao Jun
 * @Description: utils
 * @File:  tools
 * @Version: 1.0.0
 * @Date: 2021/11/19 下午6:27
 */
func GetOutboundIP(ip, port string) (outip string) {

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", ip, port), time.Duration(15)*time.Second)
	if err != nil {
		g.GetLog().Error("GetOutboundIP Err:%v\n", err)
		return
	}
	defer conn.Close()

	//localAddr := conn.LocalAddr().(*net.UDPAddr)
	localAddr := conn.LocalAddr().(*net.TCPAddr)
	//fmt.Println(localAddr.IP)
	//return localAddr.IP
	//fmt.Printf("%T,%s\n",localAddr.IP,localAddr.IP)
	tmp := fmt.Sprintf("%s", localAddr.IP)
	tmp = strings.TrimSpace(tmp)
	tmp = strings.Trim(tmp, "\n")
	return tmp
}

func IsExist(fileName string) bool {
	_, err := os.Stat(fileName)
	if err == nil {

		return true
	}
	if os.IsNotExist(err) {

		return false

	}
	return false
}
func HttpDownFile(url, fp string) error {

	if IsExist(fp) {
		return nil
	}

	client := &http.Client{}

	client = &http.Client{
		Timeout: 60 * time.Second,
	}

	request, _ := http.NewRequest("GET", url, nil)

	//request.Header.Set("Accept","text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	//request.Header.Set("Accept-Charset","GBK,utf-8;q=0.7,*;q=0.3")
	//request.Header.Set("Accept-Encoding","gzip,deflate,sdch")
	//request.Header.Set("Accept-Language","zh-CN,zh;q=0.8")
	//request.Header.Set("Cache-Control","max-age=0")
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("content-type", "application/json;charset=UTF-8")
	//request.Header.Set("User-Agent", userAgentSlice[rand.Intn(len(userAgentSlice))])

	//defer func() {
	//	if err := recover(); err != nil {
	//		log.Printf("跳过url:%s,err:%s \n", err, url)
	//		c <- struct{}{}
	//		//panic(fmt.Sprintf("err:%s\n",url))
	//
	//	}
	//}()

	response, err := client.Do(request)
	if err != nil {

		return err
	}

	if response.StatusCode != 200 {
		//log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		//c <- struct{}{}
		//wdownload.Done()
		return errors.New(fmt.Sprintf("status code error: %d %s", response.StatusCode, response.Status))
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		f := filepath.Dir(fp)
		return errors.New(fmt.Sprintf("body err: %s,dir: %s, url:%s\n", err.Error(), f, url))
	}

	//fp := string(filepath.Join("c:\\", "1"))

	err = ioutil.WriteFile(fp, body, 0777)
	if err != nil {
		return errors.New(fmt.Sprintf("%v fp:[%v]\n", err.Error(), fp))
	}
	//fmt.Printf("Download 成功: %+v\n", fp)
	return nil
}
