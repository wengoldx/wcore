// Copyright (c) 2018-2028 Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2022/11/08   jidi           New version
// -------------------------------------------------------------------

package mqtt

type MqttConfig struct {
	Broker   string `json:"broker"`
	Port     int    `json:"port"`
	ClientID string `json:"client_id"`
	User     string `json:"user"`
	PWD      string `json:"pwd"`
	CAFile   string `json:"ca"`
	CerFile  string `json:"certificate"`
	KeyFile  string `json:"key"`
}
