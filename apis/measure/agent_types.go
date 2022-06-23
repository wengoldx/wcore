// Copyright (c) 2018-2028 Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2021/12/26   youhei         New version
// -------------------------------------------------------------------

package mea

type MeaAgent struct {
	Domain string // measure access url, such as http://192.168.1.100:3000
}

type ReqID struct {
	ReqID string `json:"reqid" description:"measure request key for seach both by redis and gate monitor"`
}

type ReqIDs struct {
	ReqIDs []string `json:"reqid" description:"measure request key for seach both by redis and gate monitor"`
}

// Measure body model data
type Measure struct {
	Sex       int    `json:"sex"                          description:"body model sex, 0:invalid, 1:male, 2:female"`
	Height    int    `json:"height"   validate:"gte=0"    description:"body model height in cm"`
	Weight    int    `json:"weight"   validate:"gte=0"    description:"body model weight in kg"`
	FrontURL  string `json:"fronturl"                     description:"img's url, get from upload img api"`
	SideURL   string `json:"sideurl"                      description:"img's url, get from upload img api"`
	Bust      int    `json:"bust"     validate:"gte=0"    description:"body model bust"`
	Waist     int    `json:"waist"    validate:"gte=0"    description:"body model waist"`
	Hipline   int    `json:"hipline"  validate:"gte=0"    description:"body model hipline"`
	Wrist     int    `json:"wrist"    validate:"gte=0"    description:"body model wrist"`
	NotifyURL string `json:"ntfurl"                       description:"ansync notifier url from measure to notify body result, only notify once"`
}

type BodyUpReq struct {
	ReqID string `json:"reqid" validate:"required" description:"measure request key for seach both by redis and gate monitor"`
	Measure
}

type BodyBasicResp struct {
	BodyID     int64    `json:"id"         description:"body model unique id of user"`
	Thumbnail  string   `json:"thumbnail"  description:"body model thumbnail image"`
	ReqID      string   `json:"reqid"      description:"measure request key for seach both by redis and gate monitor"`
	Sex        int      `json:"sex"        description:"body model sex, 1 male , 2 female"`
	Status     int      `json:"reqstate"   description:"measure state, rang in [1:success, 2:capturing, 3:waiting, 4:failed]"`
	Captures   []string `json:"ue4img"     description:"body model front and side captures"`
	CreateTime int      `json:"createtime" deacription:"body create time"`
}

type BodyDetailResp struct {
	BodyID           int64    `json:"id"               description:"body model unique id of user"`
	ReqID            string   `json:"reqid"            description:"measure request key for seach both by redis and gate monitor"`
	Status           int      `json:"reqstate"         description:"measure state, rang in [1:success, 2:capturing, 3:waiting, 4:failed]"`
	Sex              int      `json:"ismale"           description:"body model sex, 1 male , 2 female"`
	Thumbnail        string   `json:"img"              description:"body model thumbnail image"`
	Captures         []string `json:"ue4img"           description:"body model front and side captures"`
	Height           int      `json:"height"           description:"body model height"`
	Weight           int      `json:"weight"           description:"body model weight"`
	Bust             int      `json:"bust"             description:"body model bust"`
	Neck             int      `json:"neck"             description:"body model neck"`
	UpperNeck        int      `json:"upperneck"        description:"body model upperneck"`
	Shoulder         int      `json:"shoulder"         description:"body model shoulder"`
	Armlen           int      `json:"armlen"           description:"body model armlen"`
	Armcir           int      `json:"armcir"           description:"body model armcir"`
	Waist            int      `json:"waist"            description:"body model waist"`
	Hipcir           int      `json:"hipcir"           description:"body model hipcir"`
	Thighcir         int      `json:"thighcir"         description:"body model thighcir"`
	Knee             int      `json:"knee"             description:"body model knee"`
	Anklecir         int      `json:"anklecir"         description:"body model anklecir"`
	Hipline          int      `json:"hipline"          description:"body model hipline"`
	Wrist            int      `json:"wrist"            description:"body model wrist"`
	Clothlen         int      `json:"clothlen"         description:"body model clothlen"`
	Outsidelen       int      `json:"outsidelen"       description:"body model outsidelen"`
	Elbowcir         int      `json:"elbowcir"         description:"body model elbowcir"`
	BpShoulderdis    int      `json:"bpshoulderdis"    description:"body model bpshoulderdis"`
	BellyShoulderdis int      `json:"bellyshoulderdis" description:"body model bellyshoulderdis"`
	WaistShoulderdis int      `json:"waistshoulderdis" description:"body model waistshoulderdis"`
	Bpdis            int      `json:"bpdis"            description:"body model bpdis"`
	Hipheight        int      `json:"hipheight"        description:"body model hipheight"`
	WaistKneedis     int      `json:"waistkneedis"     description:"body model waistkneedis"`
	Bellycir         int      `json:"bellycir"         description:"body model bellycir"`
	Armscye          int      `json:"armscye"          description:"body model armscye"`
	CreateTime       int      `json:"createtime"       deacription:"body model create time"`
	ModifyTime       int      `json:"modifytime"       deacription:"body model modify time"`
}
