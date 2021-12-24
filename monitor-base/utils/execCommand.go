package utils

import (
	"bytes"
	"os/exec"
)

/**
 * @Author: Tao Jun
 * @Description: utils
 * @File:  execCommand
 * @Version: 1.0.0
 * @Date: 2021/12/8 下午6:25
 */

func RunCommand(command string) (string, string) {
	cmd := exec.Command("/bin/bash", "-c", command)

	var out bytes.Buffer
	cmd.Stdout = &out

	var e bytes.Buffer
	cmd.Stderr = &e
	//cmd.Start()
	//buf, _ := cmd.CombinedOutput()
	if err := cmd.Run(); err != nil {
		return "", err.Error()
	}
	//if e.Len() != 0 && out.Len() == 0 {
	//	return e.String(), out.String()
	//}
	return e.String(), out.String()
}
