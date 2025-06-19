package services

import (
	"app/models"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/util/gconv"
)

type ServiceService struct {
}

func (s *ServiceService) FetchRancherServices(ws *models.Workspace) (interface{}, error) {
	if ws.Category != models.WORKSPACE_CATEGORY_RANCHER {
		return nil, fmt.Errorf("workspace category is not rancher")
	}

	ctx := context.Background()
	client := gclient.New().Timeout(time.Second * 10).BasicAuth(ws.ApiKey, ws.ApiSecret)

	resp1, err := client.Get(ctx, fmt.Sprintf("%s/stacks?limit=500", ws.ApiBaseUrl))
	if err != nil {
		return nil, err
	}
	defer resp1.Close()
	if resp1.StatusCode > 300 {
		return nil, fmt.Errorf("fetch stacks error with status: %d", resp1.StatusCode)
	}

	body := resp1.ReadAll()
	stackBody := gconv.MapDeep(body)
	stackName := make(map[string]string)
	for _, v := range stackBody["data"].([]interface{}) {
		item := v.(map[string]interface{})
		id := item["id"]
		name := item["name"]
		if id == nil || name == nil {
			continue
		}
		stackName[id.(string)] = name.(string)
	}

	resp2, err := client.Get(ctx, fmt.Sprintf("%s/services?limit=2000", ws.ApiBaseUrl))
	if err != nil {
		return nil, err
	}
	defer resp2.Close()
	if resp2.StatusCode > 300 {
		return nil, fmt.Errorf("fetch services error with status: %d", resp2.StatusCode)
	}

	body = resp2.ReadAll()
	serviceBody := gconv.MapDeep(body)
	serviceList := []map[string]string{}
	for _, v := range serviceBody["data"].([]interface{}) {
		item := v.(map[string]interface{})
		name := item["name"]
		stackId := item["stackId"]
		if name == nil {
			continue
		}
		if item["kind"] != "service" {
			// 有可能不是在 rancher 中创建的服务
			continue
		}
		if item["system"] != nil && item["system"].(bool) {
			// 过滤掉系统服务
			continue
		}
		stack := stackName[stackId.(string)]
		serviceList = append(serviceList, map[string]string{
			"Name":   name.(string),
			"Stack":  stack,
			"HostName": fmt.Sprintf("%s.%s", name.(string), stack),
			"Label": fmt.Sprintf("%s/%s", stack, name.(string)),
		})
	}

	return serviceList, nil
}

func (s *ServiceService) FetchPortainerServices(ws *models.Workspace) (interface{}, error) {
	if ws.Category != models.WORKSPACE_CATEGORY_PORTAINER {
		return nil, fmt.Errorf("workspace category is not portainer")
	}

	ctx := context.Background()
	client := gclient.New().Timeout(time.Second * 10)

	resp1, err := client.Header(g.MapStrStr{"X-API-Key": ws.ApiKey}).Get(ctx, fmt.Sprintf("%s/docker/services", ws.ApiBaseUrl))
	if err != nil {
		return nil, err
	}
	defer resp1.Close()
	if resp1.StatusCode > 300 {
		return nil, fmt.Errorf("fetch services error with status: %d", resp1.StatusCode)
	}

	body := resp1.ReadAll()
	serviceBody, err := gjson.LoadContent(body)
	if err != nil {
		return nil, err
	}

	portainerServices := serviceBody.Array()
	serviceList := []map[string]string{}
	for _, v := range portainerServices {
		item := gjson.New(v)
		name := item.Get("Spec.Name").String()
		labelNames := []string{}

		labels := item.Get("Spec.Labels")
		stack := ""
		if !labels.IsEmpty() {
			value := labels.Map()["com.docker.stack.namespace"]
			if value != nil {
				stack = value.(string)
				labelNames = append(labelNames, stack)
			}
		}

		labelNames = append(labelNames, name)

		serviceList = append(serviceList, map[string]string{
			"Name": name,
			"Stack":  stack,
			"HostName": name,
			"Label": strings.Join(labelNames, "/"),
		})
	}
	return serviceList, nil
}