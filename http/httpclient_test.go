package http

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

type listReq struct {
	AccessToken string `json:"access_token"`
	Page        string `json:"page"`
	PageSize    string `json:"pageSize"`
}

type listRsp struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Count      int `json:"count"`
		Page       int `json:"page"`
		PageSize   int `json:"pageSize"`
		TotalPages int `json:"totalPages"`
		Records    []struct {
			ProjectID            string `json:"projectId"`
			PipelineID           string `json:"pipelineId"`
			PipelineName         string `json:"pipelineName"`
			PipelineDesc         string `json:"pipelineDesc"`
			InstanceFromTemplate bool   `json:"instanceFromTemplate"`
		} `json:"records"`
	} `json:"data"`
}

func TestHTTPClient(t *testing.T) {
	hc := NewHTTPClient(time.Second * 10)

	projectID := "kdpictxt"
	url := fmt.Sprintf("http://devops.apigw.o.oa.com/prod/apigw-user/pipelines/%s", projectID)

	fmt.Println("url:", url)

	req := listReq{
		AccessToken: "dTO4oVV82ri3IE57csUK9kaUPhFqFU",
		Page:        "1",
		PageSize:    "20",
	}
	fmt.Println("req:", req)

	rsp := listRsp{}
	err := hc.Do(http.MethodGet, url, &req, &rsp)
	t.Logf("rsp:\n%+v\n", rsp)
	if err != nil {
		panic(err)
	}
}
