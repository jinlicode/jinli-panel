package tools

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
)

var sockAddr = "/var/run/docker.sock"

// GetDockerImages 获取镜像map
func GetDockerImages() map[string]string {

	httpc := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", sockAddr)
			},
		},
	}

	// 获取容器列表
	resp, _ := httpc.Get("http://localhost/v1.40/images/json")

	resp.Header.Set("Content-Type", "application/json")

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var dataArr []map[string]interface{}

	json.Unmarshal(body, &dataArr)
	imagesMaps := make(map[string]string)

	for _, v := range dataArr {
		imagesMaps[v["RepoTags"].([]interface{})[0].(string)] = v["RepoTags"].([]interface{})[0].(string)
	}

	return imagesMaps
}

// GetDockerIP 获取容器ip
func GetDockerIP(containersName string) string {
	httpc := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", sockAddr)
			},
		},
	}

	resp, err := httpc.Get("http://localhost/v1.40/containers/" + containersName + "/json")
	if err != nil {
		return ""
	}

	resp.Header.Set("Content-Type", "application/json")

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	dataArr := make(map[string]interface{})

	json.Unmarshal(body, &dataArr)

	IPAddress := dataArr["NetworkSettings"].(map[string]interface{})["Networks"].(map[string]interface{})[containersName+"_net"].(map[string]interface{})["IPAddress"].(string)
	return IPAddress
}

// CheckDockerStatus 获取容器的运行状态
func CheckDockerStatus(containersName string) bool {
	httpc := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", sockAddr)
			},
		},
	}

	resp, err := httpc.Get("http://localhost/v1.40/containers/" + containersName + "/json")
	if err != nil {
		return false
	}

	resp.Header.Set("Content-Type", "application/json")

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	dataArr := make(map[string]interface{})

	json.Unmarshal(body, &dataArr)

	Running := dataArr["State"].(map[string]interface{})["Running"].(bool)
	return Running
}
