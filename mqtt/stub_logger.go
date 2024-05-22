// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------

package mqtt

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	mq "github.com/eclipse/paho.mqtt.golang"
	"github.com/wengoldx/wing/logger"
)

const (
	adapterMqtt = "mqtt"          // Adapter name of logger ouput by mqtt
	logTopicPre = "wengold/logs/" // mqtt logger publish topic prefix
)

// custom mqtt logger
type mqttLogger struct {
	Options *Options  // mqtt broker configs
	Stub    mq.Client // mqtt client instanse
	Topic   string    // publish topic
}

// Register mqtt logger as a beego logs, it will create
// single mqtt client to output logs only for prod mode
func SetupLogger(opts *Options) {
	if beego.BConfig.RunMode == "prod" && opts != nil {
		getMqttLogger := func() logs.Logger {
			return &mqttLogger{
				Topic:   logTopicPre + beego.BConfig.AppName,
				Options: opts,
			}
		}
		logs.Register(adapterMqtt, getMqttLogger)
		logs.SetLogger(adapterMqtt, "mqtt-logger")
	}
}

// Parse and return mqtt broker configs, it maybe nil returned
func GetOptions(data string, svr ...string) *Options {
	cfgs := &MqttConfigs{}
	if err := json.Unmarshal([]byte(data), &cfgs); err != nil {
		return nil
	}

	userkey := beego.BConfig.AppName
	if len(svr) > 0 && svr[0] != "" {
		userkey = svr[0]
	}

	if user, ok := cfgs.Users[userkey]; ok {
		return &Options{
			Host: cfgs.Broker.Host,
			Port: cfgs.Broker.Port,
			User: user,
		}
	}
	return nil
}

// Init mqtt logger topic and connect client with id of 'appname.logger'
func (w *mqttLogger) Init(config string) error {
	options, protocol := mq.NewClientOptions(), "tcp://%s:%v"
	broker := fmt.Sprintf(protocol, w.Options.Host, w.Options.Port)
	options.AddBroker(broker)
	options.SetClientID("logger." + beego.BConfig.AppName)
	options.SetUsername(w.Options.User.Account)
	options.SetPassword(w.Options.User.Password)
	options.SetAutoReconnect(true)

	w.Stub = mq.NewClient(options)
	if token := w.Stub.Connect(); token.Wait() && token.Error() != nil {
		w.Stub = nil

		// Delete mqtt logger from beego logs when connect failed
		logs.GetBeeLogger().DelLogger(adapterMqtt)
		logger.E("Setup mqtt logger err:", token.Error())
	}
	return nil
}

// Publish logs above warning level after mqtt client connected
func (w *mqttLogger) WriteMsg(when time.Time, msg string, level int) error {
	if w.Stub != nil && level <= logs.LevelWarning && msg != "" {
		msg = when.Format("2006/01/02 15:04:05.000") + " " + msg
		w.Stub.Publish(w.Topic, 0, false, msg)
	}
	return nil
}

// Disconnect mqtt client if living
func (w *mqttLogger) Destroy() {
	if w.Stub != nil && w.Stub.IsConnected() {
		w.Stub.Disconnect(0)
		w.Stub = nil
	}
}

// Do nothing here, none cache to output
func (w *mqttLogger) Flush() {}
