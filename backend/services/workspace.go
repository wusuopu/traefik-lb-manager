package services

import (
	"app/di"
	"app/models"
	"fmt"

	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/samber/lo"
)

type WorkspaceService struct {
}

func (w *WorkspaceService) GenerateTraefikConfig(ws *models.Workspace) error {
	payload, err := gjson.LoadJson([]byte(`{
		"http": {
			"routers": {
				"lets-encrypt-router": {
					"entryPoints":["web"],
					"priority":100000,
					"rule":"PathPrefix(` + "`" + "/.well-known/" + "`" + `)",
					"service":"lets-encrypt-service"
				}
			},
			"middlewares": {},
			"services": {
				"lets-encrypt-service": {
					"loadBalancer": {
						"servers":[
							{"url":""}
						]
					}
				}
			}
		},
		"tls": {
			"certificates": []
		}
	}`))
	if err != nil {
		return err
	}

	// 添加 let's encrypt 规则
	payload.Set("http.services.lets-encrypt-service.loadBalancer.servers.0.url", ws.ManagerBaseUrl)


	err = w.generateCertificates(ws, payload)
	if err != nil {
		return err
	}

	err = w.generateServers(ws, payload)
	if err != nil {
		return err
	}

	willUpdate := make(map[string]interface{})
	willUpdate["traefik_config"] = payload.MustToYamlString()
	willUpdate["traefik_json_config"] = payload.MustToJsonString()

	results := di.Container.DB.Model(&models.Workspace{}).Where("id = ?", ws.ID).Updates(willUpdate)
	if results.Error != nil {
		fmt.Printf("Error: %v\n", results.Error)
		return results.Error
	}

	return nil
}

func (w *WorkspaceService) generateServers(ws *models.Workspace, ret *gjson.Json) error {
	var servers []models.Server
	results := di.Container.DB.
		Model(&models.Server{}).
		Select("id", "Host").
		Where("workspace_id = ?", ws.ID).
		Where("enable = ?", 1).
		Find(&servers)
	if results.Error != nil {
		di.Container.Logger.Error(fmt.Sprintf("GenerateServers Find Server error: %s", results.Error.Error()))
		return results.Error
	}

	var services []models.Service
	results = di.Container.DB.
		Model(&models.Service{}).
		Where("workspace_id = ?", ws.ID).
		Find(&services)
	if results.Error != nil {
		di.Container.Logger.Error(fmt.Sprintf("GenerateServers Find Service error: %s", results.Error.Error()))
		return results.Error
	}

	var mds []models.Middleware
	results = di.Container.DB.
		Model(&models.Middleware{}).
		Where("workspace_id = ?", ws.ID).
		Find(&mds)
	if results.Error != nil {
		di.Container.Logger.Error(fmt.Sprintf("GenerateServers Find Middleware error: %s", results.Error.Error()))
		return results.Error
	}

	for _, s := range servers {
		err := w.generateRules(ws, ret, &s, services, mds)
		if err != nil {
			return err
		}
	}

  return nil
}
func (w *WorkspaceService) generateRules(ws *models.Workspace, ret *gjson.Json, server *models.Server, services []models.Service, mds []models.Middleware) error {
	var rules []models.Rule
	results := di.Container.DB.
		Model(&models.Rule{}).
		Where("workspace_id = ?", ws.ID).
		Where("server_id = ?", server.ID).
		Where("enable = ?", 1).
		Find(&rules)
	if results.Error != nil {
		di.Container.Logger.Error(fmt.Sprintf("GenerateRules Find Rule error: %s", results.Error.Error()))
		return results.Error
	}

	serviceIdsSet := gset.New()
	middlewareIdsSet := gset.New()

	serverRule := server.GetHostRules()
	for _, r := range rules {
		if r.Options == nil {
			continue
		}
		options, err := gjson.LoadJson(r.Options)
		if err != nil {
			di.Container.Logger.Error(fmt.Sprintf("Load Rule %d Options error: %s", r.ID, err))
			return err
		}
		baseKeyPath := fmt.Sprintf("http.routers._router__%d_%d", server.ID, r.ID)

		// set services
		serviceId := options.Get("service").Uint()
		_, ok := lo.Find(services, func(s models.Service) bool {
			return s.ID == serviceId
		})
		if !ok {
			di.Container.Logger.Error(fmt.Sprintf("Rule %d service(%d) is not exist", r.ID, serviceId))
			continue
		}
		ret.Set(fmt.Sprintf("%s.service", baseKeyPath), fmt.Sprintf("_service__%d", serviceId))
		serviceIdsSet.AddIfNotExist(serviceId)

		// set entryPoints
		entryPoints := options.Get("entryPoints").Array()
		for i, e  := range entryPoints {
			ret.Set(fmt.Sprintf("%s.entryPoints.%d", baseKeyPath, i), e.(string))
		}

		// set middlewares
		mdIds := options.Get("middlewares").Array()
		mdIndex := 0
		for _, mdId := range mdIds {
			id := gconv.Uint(mdId)
			_, ok := lo.Find(mds, func(m models.Middleware) bool {
				return m.ID == id
			})
			if !ok {
				di.Container.Logger.Error(fmt.Sprintf("Rule %d middleware(%d) is not exist", r.ID, id))
				continue
			}
			ret.Set(fmt.Sprintf("%s.middlewares.%d", baseKeyPath, mdIndex), fmt.Sprintf("_middleware__%d", id))
			middlewareIdsSet.AddIfNotExist(id)
			mdIndex++
		}

		// set priority
		priority := options.Get("priority").Uint()
		if priority > 0 {
			ret.Set(fmt.Sprintf("%s.priority", baseKeyPath), priority)
		}

		// set rule
		routeRule := ""
		if options.Get("advanceMode").Bool() {
			routeRule = options.Get("rule").String()
		} else {
			routeRule = fmt.Sprintf("PathPrefix(`%s`)",  options.Get("rule").String())
		}
		if serverRule == "" {
			ret.Set(fmt.Sprintf("%s.rule", baseKeyPath), routeRule)
		} else {
			ret.Set(fmt.Sprintf("%s.rule", baseKeyPath), fmt.Sprintf("%s && (%s)", serverRule, routeRule))
		}

	}

	for _, v := range serviceIdsSet.Slice() {
		service, ok := lo.Find(services, func(s models.Service) bool {
			return s.ID == v
		})
		if !ok {
			di.Container.Logger.Error(fmt.Sprintf("service(%d) is not exist", v))
			continue
		}
		w.generateService(ws, ret, &service)
	}

	for {
		if middlewareIdsSet.Size() == 0 {
			break
		}

		v := middlewareIdsSet.Pop()
		md, ok := lo.Find(mds, func(m models.Middleware) bool {
			return m.ID == v
		})
		if !ok {
			di.Container.Logger.Error(fmt.Sprintf("middleware(%d) is not exist", v))
			continue
		}
		w.generateMiddleware(ws, ret, &md, mds, middlewareIdsSet)
	}
  return nil
}

func (w *WorkspaceService) generateService(ws *models.Workspace, ret *gjson.Json, service *models.Service) error {
	data := service.GetLBRuleMap()
	if data == nil {
		return nil
	}

	ret.Set(fmt.Sprintf("http.services._service__%d.loadBalancer", service.ID), data)
  return nil
}

func (w *WorkspaceService) generateMiddleware(ws *models.Workspace, ret *gjson.Json, md *models.Middleware, mds []models.Middleware, mdSet *gset.Set) error {
	keyPath := fmt.Sprintf("http.middlewares._middleware__%d.%s", md.ID, md.Category)
	value := ret.Get(keyPath)
	if value != nil {
		// 该 middleware 已经生成过
		return nil
	}

	data := md.GetRuleMap()
	if data == nil {
		return nil
	}

	if md.Category == "chain" {
		// 对于 chain middleware 还需要把包含的其他  middleware 添加进来
		chainIds, ok := data["middlewares"]
		chainNames := []string{}

		if ok {
			for _, v := range chainIds.([]interface{}) {
				id := gconv.Uint(v)
				_, ok := lo.Find(mds, func(m models.Middleware) bool {
					return m.ID == id
				})
				if ok {
					mdSet.AddIfNotExist(id)
					chainNames = append(chainNames, fmt.Sprintf("_middleware__%d", id))
				}
			}
		}
		data["middlewares"] = chainNames
	}

	ret.Set(keyPath, data)
  return nil
}

// 生成证书列表
func (w *WorkspaceService) generateCertificates(ws *models.Workspace, ret *gjson.Json) error {
	var certs []models.Certificate
	results := di.Container.DB.
		Model(&models.Certificate{}).
		Where("workspace_id = ?", ws.ID).
		Where("enable = ?", 1).
		Where("status = ?", models.CERTIFICATE_STATUS_COMPLETE).
		Find(&certs)

	if results.Error != nil {
		di.Container.Logger.Error(fmt.Sprintf("GenerateCertificates error: %s", results.Error.Error()))
		return results.Error
	}

	for i, cert := range certs {
		if cert.Cert == "" || cert.Key == "" {
			continue
		}
		name := fmt.Sprintf("/etc/traefik/ssl/%s__%d", cert.Domain, cert.ID)
		key := fmt.Sprintf("tls.certificates.%d.", i)
		ret.Set(key + "certFile", name + ".crt")
		ret.Set(key + "keyFile", name + ".key")
	}
	return nil
}