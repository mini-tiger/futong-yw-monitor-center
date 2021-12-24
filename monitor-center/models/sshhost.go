package models

import (
	"context"
	"errors"
	"fmt"
	"futong-yw-monitor-center/monitor-base/bg"
	"futong-yw-monitor-center/monitor-center/g"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"path"
	"strings"
	"time"
)

/**
 * @Author: Tao Jun
 * @Description: models
 * @File:  sshhost
 * @Version: 1.0.0
 * @Date: 2021/11/29 下午3:40
 */

type SSH_Host struct {
	// 链接句柄
	SSHClient *ssh.Client `json:"-"`

	// 资源配置表
	Username string `json:"username"`
	Password string `json:"password"`
	Port     int    `json:"port"`
	Ip       string `json:"ip"`
	HostId   string `json:"hostId"`
	OsType   string `json:"osType"`   // linux or windows
	AuthType uint8  `json:"authType"` // pass or pub
}

func (l *SSH_Host) GetConn(ctx context.Context) (err error) {
	var Auth []ssh.AuthMethod
	if l.AuthType == 1 { // pass
		Auth = []ssh.AuthMethod{ssh.Password(l.Password)}
	} else { // pub
		signer, err := ssh.ParsePrivateKey([]byte(l.Password))
		if err != nil {
			return err
		}
		Auth = []ssh.AuthMethod{
			//证书验证
			ssh.PublicKeys(signer),
			//密码验证
			//ssh.Password("xxxx"),
		}
	}

	clientConfig := &ssh.ClientConfig{
		User:    l.Username,
		Auth:    Auth,
		Timeout: 20 * time.Second,
		// 这个是问你要不要验证远程主机，以保证安全性。这里不验证
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr := fmt.Sprintf("%s:%d", l.Ip, l.Port)
	l.SSHClient, err = ssh.Dial("tcp", addr, clientConfig)
	return err
}

func (l *SSH_Host) ExecShell(ctx context.Context, cmd string, sync bool) ([]byte, error) {
	defer func() {
		if err := recover(); err != nil {
		}
	}()

	var data []byte
	// create session
	session, err := l.SSHClient.NewSession()
	if err != nil {
		return data, err
	}
	defer session.Close()

	if sync {
		data, err = session.Output(cmd)
	} else {
		err = session.Start(cmd)
	}

	if err != nil {
		return data, err
	}
	return data, nil
}
func (l *SSH_Host) FileTransfer(ctx context.Context, srcfile string, remotefile string) error {
	// create sftp client
	sftpClient, err := sftp.NewClient(l.SSHClient)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	srcFile, err := os.Open(srcfile)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := sftpClient.Create(remotefile)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	buf := make([]byte, 1024*1024)
	for {
		select {
		case <-ctx.Done():
			return errors.New("FileTransfer timeout")
		default:
			n, err := srcFile.Read(buf)
			if err != nil {
				if err == io.EOF {
					err = nil
				}
				return err
			}
			if n == 0 {
				return nil
			}
			n1, err := dstFile.Write(buf[0:n])
			if n1 != n || err != nil {
				return err
			}
		}
	}
}

func (l *SSH_Host) CmdExec() {
	err := l.GetConn(context.TODO())
	if err != nil {
		g.GetLog().Error("SSH Conn host:%s,err:%s\n", l.Ip, err)
		return
	}
	defer l.SSHClient.Close()
	checkBinCmdStr := fmt.Sprintf("chmod 777 /usr/local/bin/%s;echo $?", bg.AgentName)
	bytes, err := l.ExecShell(context.TODO(), checkBinCmdStr, true)
	if err != nil {
		g.GetLog().Error("SSH ExecShell chmod host:%s,err:%s\n", l.Ip, err)
		return
	}
	//fmt.Println(strings.TrimSpace(string(bytes)))
	sshCmdStr := bg.GetPatternStartParams("ssh", g.GetConfig().GetConfUrl)
	sshCmdStr = fmt.Sprintf("nohup %s &", sshCmdStr)
	g.GetLog().Debug("IP:%s host:%s ssh cmd:%v\n", l.Ip, l.HostId, sshCmdStr)
	if (strings.TrimSpace(string(bytes))) == "0" {
		// 直接执行

		_, err := l.ExecShell(context.TODO(),
			fmt.Sprintf("cd /usr/local/bin/; %s", sshCmdStr), true)

		if err != nil {
			g.GetLog().Error("SSH ExecShell ssh Pattern host:%s,err:%s\n", l.Ip, err)
			return
		}
		g.GetLog().Debug("SSH ExecShell ssh Pattern host:%s,success\n", l.Ip)
	} else {
		srcFile := path.Join(g.CurrentDir, string(os.PathSeparator), "packages", bg.AgentName)
		remoteFile := fmt.Sprintf("/usr/local/bin/%s", bg.AgentName)
		// 传送文件 后 执行
		err := l.FileTransfer(context.TODO(),
			srcFile, remoteFile)
		time.Sleep(1 * time.Second)
		if err != nil {
			g.GetLog().Error("SFTP host:%s,err:%s\n", l.Ip, err)
			return
		} else {
			g.GetLog().Info("SFTP host:%s success\n", l.Ip)
			//sshCmdStr := bg.GetPatternStartParams("ssh", g.GetConfig().GetConfUrl)
			_, err := l.ExecShell(context.TODO(),
				fmt.Sprintf("cd /usr/local/bin/ && chmod 777 %s ; %s", bg.AgentName, sshCmdStr),
				true)
			if err != nil {
				g.GetLog().Error("SSH ExecShell ssh Pattern host:%s,err:%s\n", l.Ip, err)
				return
			}
			g.GetLog().Debug("SSH ExecShell ssh Pattern host:%s,success\n", l.Ip)
		}
	}
}
