// Copyright (c) 2018-2028 Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2023/04/18   tangxiaoyu     New version
// -------------------------------------------------------------------

package comm

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/wengoldx/wing/logger"
	"io/ioutil"
	"reflect"
	"strings"
)

// Swagger.json field keywords for version 2.0.0.
//
// `NOTICE`
//
//	DO NOT CHANGE THEME IF YOU KNOWE HOW TO CHANGE IT!
const (
	swaggerFile  = "./swagger/swagger.json"
	sfServerName = "basePath"
	sfPathName   = "paths"
	sfMethodGet  = "get"
	sfMethodPost = "post"
	sfEnDescName = "description"
	sfGroupTags  = "tags"
	sfGroupName  = "name"
)

// A router informations
type Router struct {
	Router string `json:"router"` // restful router full path start /
	Method string `json:"method"` // restful router http method, such as GET, POST...
	Group  string `json:"group"`  // beego controller keyworld
	EnDesc string `json:"endesc"` // english description of router from swagger
	CnDesc string `json:"cndesc"` // chinese description of router manual update by user
}

// A group informations
type Group struct {
	Name   string `json:"group"`  // group name as swagger controller path, like '/v3/acc'
	EnDesc string `json:"endesc"` // english description of group from swagger
	CnDesc string `json:"cndesc"` // chinese description of group manual update by user
}

// All routers of one server
type Routers struct {
	Server  string    `json:"server"`  // backend server name
	CnName  string    `json:"cnanem"`  // server name of chinese
	Groups  []*Group  `json:"groups"`  // groups as swagger controllers
	Routers []*Router `json:"routers"` // routers parsed from swagger.json file
}

// Load local server routers from swagger.json file.
func LoadSwaggerRouters() (*Routers, error) {
	buff, err := ioutil.ReadFile(swaggerFile)
	if err != nil {
		logger.E("Load swagger routers, err:", err)
		return nil, err
	}

	routers := make(map[string]interface{})
	if err := json.Unmarshal(buff, &routers); err != nil {
		logger.E("Unmarshal swagger routers err:", err)
		return nil, err
	}
	logger.I("Make routers, and unmarshal swagger json")

	out := &Routers{}
	if basePath, ok := routers[sfServerName]; ok && basePath != nil {
		out.Server = basePath.(string) // parse server name
		logger.I("Parsed server name:", out.Server)

		// parse routers by path keyword
		if ps, ok := routers[sfPathName]; ok && ps != nil {
			paths := ps.(map[string]interface{})

			for path, pathvals := range paths {
				router := &Router{Router: path} // parse router path

				// parse http method, HERE only support GET or POST methods
				var mvs interface{}
				pvs := pathvals.(map[string]interface{})
				if pmg, ok := pvs[sfMethodGet]; ok && pmg != nil {
					router.Method, mvs = "GET", pmg
				} else if pmp, ok := pvs[sfMethodPost]; ok && pmp != nil {
					router.Method, mvs = "POST", pmp
				} else {
					logger.W("Invalid method of path:", path)
					continue
				}

				// parse beego controller group name
				method := mvs.(map[string]interface{})
				if gps, ok := method[sfGroupTags]; ok && gps != nil {
					groups := reflect.ValueOf(gps)
					router.Group = groups.Index(0).Interface().(string)
				}

				// parse router path english description
				if desc, ok := method[sfEnDescName]; ok && desc != nil {
					router.EnDesc = desc.(string)
				}

				// append the router into routers array
				logger.D("> Parsed ["+router.Method+"]\tpath:", path, "\tdesc:", router.EnDesc)
				out.Routers = append(out.Routers, router)
			}
		}

		// parse groups by tags keyword
		if gps, ok := routers[sfGroupTags]; ok && gps != nil {
			groups := gps.([]interface{}) // parse all group array

			for _, group := range groups {
				t := group.(map[string]interface{})
				gp := &Group{}

				// parse group name value
				if gpn, ok := t[sfGroupName]; ok && gpn != nil {
					gp.Name = gpn.(string)
				}

				// parse group english description
				if gpd, ok := t[sfEnDescName]; ok && gpd != nil {
					gp.EnDesc = strings.TrimRight(gpd.(string), "\n")
				}

				logger.D("# Parsed group ["+gp.Name+"] \t desc:", gp.EnDesc)
				out.Groups = append(out.Groups, gp)
			}
		}
	}

	logger.I("Finished parse, out:", out)
	return out, nil
}

// Parse servers routers from nacos config data, then return all backend services
// routers map and local server swagger routers
func ParseNacosRouters(data string) (map[string]*Routers, *Routers) {
	routers := make(map[string]*Routers)
	if data != "" && data != "{}" { // check data if empty
		if err := json.Unmarshal([]byte(data), &routers); err != nil {
			logger.E("Unmarshal swagger routers, err:", err)
			return nil, nil
		}
	}

	svr := beego.BConfig.AppName
	if rs, ok := routers[svr]; ok {
		logger.D("Parsed routers and found", svr)
		return routers, rs
	}

	logger.D("Parsed routers, but unexist", svr)
	return routers, nil
}

// Fetch the given routers and groups from src param and set chinese description to dest fileds.
func FetchChineseFields(src *Routers, dest *Routers) {
	dest.CnName = Condition(src.CnName != "", src.CnName, dest.CnName).(string)

	/* -------------------------------- */
	/* cache router chinese description */
	/* -------------------------------- */
	routers := make(map[string]string)
	if len(dest.Routers) > 0 {
		for _, router := range src.Routers {
			if router.CnDesc != "" {
				logger.D("- Cached router ["+router.Router+"]\tchinese desc:", router.CnDesc)
				routers[router.Router] = router.CnDesc
			}
		}
	}

	// set chinese description to dest routers
	if len(routers) > 0 {
		for _, router := range dest.Routers {
			router.CnDesc = routers[router.Router]
		}
	}

	/* -------------------------------- */
	/* cache groups chinese description */
	/* -------------------------------- */
	groups := make(map[string]string)
	if len(dest.Groups) > 0 {
		for _, group := range src.Groups {
			if group.CnDesc != "" {
				logger.D("= Cached group ["+group.Name+"]\tchinese desc:", group.CnDesc)
				groups[group.Name] = group.CnDesc
			}
		}
	}

	// set chinese description to dest groups
	if len(groups) > 0 {
		for _, group := range dest.Groups {
			group.CnDesc = groups[group.Name]
		}
	}
}

// Parse total routers and update description on chinese for local server routers,
// then marshal to string and push to nacos config server.
func UpdateRouters(data string) (string, error) {
	routers, err := LoadSwaggerRouters()
	if err != nil {
		logger.E("Load local swagger, err:", err)
		return "", err
	}

	rsmap, nrs := ParseNacosRouters(data)
	if nrs != nil {
		FetchChineseFields(nrs, routers)
	}

	svr := beego.BConfig.AppName
	if rsmap != nil {
		rsmap[svr] = routers
	} else {
		rsmap = make(map[string]*Routers)
		rsmap[svr] = routers
	}

	swagger, err := json.Marshal(rsmap)
	if err != nil {
		logger.E("Marshal routers, err:", err)
		return "", err
	}

	logger.D("Updated routers apis for", svr)
	return string(swagger), nil
}
