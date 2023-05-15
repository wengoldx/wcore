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
	"github.com/wengoldx/wing/invar"
	"github.com/wengoldx/wing/logger"
	"io/ioutil"
	"strings"
)

// Swagger.json field keywords for version 2.0.0.
//
// `NOTICE`
//
//	DO NOT CHANGE THEME IF YOU KNOWE HOW TO CHANGE IT!
const (
	swaggerFile  = "./swagger/swagger.json"
	sfEnDescName = "description"
	sfGroupTags  = "tags"
	sfGroupName  = "name"
)

// A group informations
type Group struct {
	Name   string `json:"group"`  // group name as swagger controller path, like '/v3/acc'
	EnDesc string `json:"endesc"` // english description of group from swagger
	CnDesc string `json:"cndesc"` // chinese description of group manual update by user
}

// All routers of one server
type Routers struct {
	CnName string   `json:"cnname"` // backend server chinese name
	Groups []*Group `json:"groups"` // groups as swagger controllers
}

type SvrDesc struct {
	Server string            `json:"server"` // backend server english name
	CnName string            `json:"cnname"` // backend server chinese name
	Groups map[string]string `json:"groups"` // groups  chinese description
}

// Parse total routers and update description on chinese for local server routers,
// then marshal to string for next to push to nacos config server.
func UpdateRouters(data string) (string, error) {
	routers, err := loadSwaggerRouters()
	if err != nil {
		logger.E("Load local swagger, err:", err)
		return "", err
	}

	rsmap, nrs := parseNacosRouters(data)
	if nrs != nil {
		fetchChineseFields(nrs, routers)
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

// Parse total routers and update description on chinese for local server routers,
// then marshal to string for next to push to nacos config server.
func UpdateChineses(data string, descs []*SvrDesc) (string, error) {
	rsmap := make(map[string]*Routers)
	if err := json.Unmarshal([]byte(data), &rsmap); err != nil {
		return "", err
	}

	// check total routers map and input chinese values
	if len(rsmap) == 0 || len(descs) == 0 {
		return "", invar.ErrEmptyData
	}

	// fetch all routers and update chinese
	changed := false
	for _, svr := range descs {
		if routers, ok := rsmap[svr.Server]; ok {
			if routers.CnName != svr.CnName {
				routers.CnName, changed = svr.CnName, true
			}

			// update groups chinese descriptions
			for _, group := range routers.Groups {
				if cnname, ok := svr.Groups[group.Name]; ok {
					if group.CnDesc != cnname {
						group.CnDesc, changed = cnname, true
					}
				}
			}
		}
	}

	// check if exist chinese updated
	if !changed {
		return "", invar.ErrNotChanged
	}

	swagger, err := json.Marshal(rsmap)
	if err != nil {
		return "", err
	}

	logger.D("Updated routers chineses")
	return string(swagger), nil
}

// --------------------------------------------

// Load local server routers from swagger.json file.
func loadSwaggerRouters() (*Routers, error) {
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
	logger.I("Loaded swagger json, start parse routers")
	out := &Routers{}

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

			// append the group into groups array
			// logger.D("# Parsed group ["+gp.Name+"] \t desc:", gp.EnDesc)
			out.Groups = append(out.Groups, gp)
		}
	}

	logger.I("Finish parsed swagger routers")
	return out, nil
}

// Parse servers routers from nacos config data, then return all backend services
// routers map and local server swagger routers
func parseNacosRouters(data string) (map[string]*Routers, *Routers) {
	routers := make(map[string]*Routers)
	if data != "" && data != "{}" { // check data if empty
		if err := json.Unmarshal([]byte(data), &routers); err != nil {
			logger.E("Unmarshal nacos routers, err:", err)
			return nil, nil
		}
	}

	if len(routers) > 0 {
		svr := beego.BConfig.AppName
		if rs, ok := routers[svr]; ok {
			logger.D("Parsed nacos routers, found", svr)
			return routers, rs
		}
		logger.D("Parsed nacos routers, unexist", svr)
		return routers, nil
	}

	logger.D("Empty nacos routers, data:", data)
	return routers, nil
}

// Fetch the given routers and groups from src param and set chinese description to dest fileds.
func fetchChineseFields(src *Routers, dest *Routers) {
	dest.CnName = Condition(src.CnName != "", src.CnName, dest.CnName).(string)

	groups := make(map[string]string)
	if len(dest.Groups) > 0 {
		for _, group := range src.Groups {
			if group.CnDesc != "" {
				// logger.D("= Cached group ["+group.Name+"]\tchinese desc:", group.CnDesc)
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
