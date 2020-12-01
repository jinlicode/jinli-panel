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
