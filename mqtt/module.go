// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2024/03/15   youhei         New version
// -------------------------------------------------------------------

package mqtt

// ClientConfigs mqtt client configs to connect MQTT broker
type ClientConfigs struct {
	Broker   string       // remote MQTT broker address
	Port     int          // remote MQTT broker port number
	ClientID string       // current client unique id on broker
	User     *UserConfigs // user account and password to connect broker
	CAFile   string       // CA cert file for TSL
	CerFile  string       // certificate/key file for TSL
	KeyFile  string       // secure key file for TSL
}

// MqttConfigs mqtt configs pasered from nacos configs server
type MqttConfigs struct {
	Broker  string                  `json:"broker"`
	Port    int                     `json:"port"`
	Users   map[string]*UserConfigs `json:"svrcfg"`
	CAFile  string                  `json:"ca"`
	CerFile string                  `json:"certificate"`
	KeyFile string                  `json:"key"`
}

// UserConfigs mqtt client secure datas
type UserConfigs struct {
	Account  string `json:"user"`
	Password string `json:"pwd"`
}
