package bg

import (
	"fmt"
	"strconv"
	"strings"
)

/**
 * @Author: Tao Jun
 * @Description: g
 * @File:  version
 * @Version: 1.0.0
 * @Date: 2021/11/21 上午10:13
 */

const AgentName = "ft-agent"

type AgentVersion uint64

var Version string

// agent 端版本
var AgentVer AgentVersion

func init() {
	SetAgentVersion(Version)
}
func SetAgentVersion(v interface{}) {
	var vv int
	switch v.(type) {
	case string:
		vv, _ = strconv.Atoi(v.(string))
		break
	case int, int32, int64, uint64, uint, uint32:
		break
	}

	//a := AgentVersion(uint64(vv))
	AgentVer = AgentVersion(uint64(vv))
}

func GetAgentVer() {
	fmt.Println(AgentVer)
}

// req response
type AgentVerResp struct {
	AgentUrl    string // agent 下载地址?
	Update      bool
	AgentNewVer AgentVersion
	Pattern     string
}

type AgentVerReq struct {
	Os       string       `json:"os"`
	Arch     string       `json:"arch"`
	HostId   string       `json:"hostId"`
	AgentVer AgentVersion `json:"agentVer"` // agent端的agent版本
}

func GetPatternStartParams(Pattern string, getconfurl string) string {
	switch strings.ToLower(Pattern) {
	case "ssh":
		return fmt.Sprintf("./%s -i=false -w=false -o=true -c=%s", AgentName, getconfurl)
	case "web":
		return fmt.Sprintf("./%s -w=true -i=false -o=false", AgentName)
	case "agent":
		return fmt.Sprintf("./%s -i=true -o=false -w=false", AgentName)
	default:
		return ""
	}

}

////center 端版本
//var Linux64Ver AgentVersion = 100 // 测试用写死
//var Linux32Ver AgentVersion
