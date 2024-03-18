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

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/astaxie/beego"
	mq "github.com/eclipse/paho.mqtt.golang"
	"github.com/wengoldx/wing/invar"
	"github.com/wengoldx/wing/logger"
	"github.com/wengoldx/wing/secure"
)

// MQTT stub to manager MQTT connection.
//
// As usages you can connect remote MQTT broker and get client instance by follow usecases.
//
//
// UseCase 1 : Using nacos MQTT configs and connect without callbacks.
//
//	if err := mqtt.Singleton().GenClient(data); err != nil {
//		logger.E("Connect client err:", err)
//		return
//	}
//	client := mqtt.Singleton().Client
//
//
// UseCase 2 : Using nacos MQTT configs and connect with callbacks.
//
//	stub := mqtt.Singleton().SetHandlers(connectHandler, disconnectHandler)
//	if err := stub.GenClient(data); err != nil {
//		logger.E("Connect client err:", err)
//		return
//	}
//	client := stub.Client
//
//
// UseCase 3 : Using singleton stub set custom configs and connect.
//
//	stub := mqtt.Singleton()
//	// Here set your custom client configs
//	if err := stub.Connect(stub.GenConfigs()); err != nil {
//		logger.E("Connect client err:", err)
//		return
//	}
//	client := stub.Client
type MqttStub struct {
	Cfg               *ClientConfigs
	Client            mq.Client
	ConnectHandler    mq.OnConnectHandler
	DisconnectHandler mq.ConnectionLostHandler
}

// Singleton mqtt stub instance
var mqttStub *MqttStub

// Return Mqtt global Singleton
func Singleton() *MqttStub {
	if mqttStub == nil {
		mqttStub = &MqttStub{
			Cfg:               &ClientConfigs{},
			Client:            nil,
			ConnectHandler:    defConnectHandler,
			DisconnectHandler: defConnectLostHandler,
		}
	}
	return mqttStub
}

// Default connect handler, call Singleton().ConnectHandler to set custom
// handler before calling GenConfigs().
var defConnectHandler mq.OnConnectHandler = func(client mq.Client) {
	serve, opt := beego.BConfig.AppName, client.OptionsReader()
	logger.I("Server", serve, "connected mqtt as client:", opt.ClientID())
}

// Default disconnect handler, call Singleton().DisconnectHandler to set
// custom handler before calling GenConfigs().
var defConnectLostHandler mq.ConnectionLostHandler = func(client mq.Client, err error) {
	serve, opt := beego.BConfig.AppName, client.OptionsReader()
	logger.W("Server", serve, "disconnect mqtt client:", opt.ClientID())
}

// Generate mqtt config, default connection protocol using tcp, you can
// set mode 'tls' and cert files to using ssl protocol.
func (stub *MqttStub) GenConfigs(mode ...string) *mq.ClientOptions {
	options, protocol := mq.NewClientOptions(), "tcp://%s:%v"
	if len(mode) > 0 && mode[0] == "tls" {
		protocol = "ssl://%s:%v"
		if tlscfg := stub.newTLSConfig(); tlscfg != nil {
			options.SetTLSConfig(tlscfg)
		}
	}

	broker := fmt.Sprintf(protocol, stub.Cfg.Broker, stub.Cfg.Port)
	options.AddBroker(broker)
	options.SetClientID(stub.Cfg.ClientID)
	options.SetUsername(stub.Cfg.User.Account)
	options.SetPassword(stub.Cfg.User.Password)
	options.SetAutoReconnect(true)
	options.SetOnConnectHandler(stub.ConnectHandler)
	options.SetConnectionLostHandler(stub.DisconnectHandler)
	return options
}

// Generate mqtt client and connect with MQTT broker, the client using
// 'tcp' protocol and fixed id as format 'server@12345678'.
func (stub *MqttStub) GenClient(configs string, server ...string) error {
	svr := beego.BConfig.AppName
	if len(server) > 0 && server[0] != "" {
		svr = server[0]
	}

	if err := stub.parseConfig(configs, svr); err != nil {
		return err
	}

	opt := stub.GenConfigs() // using default tcp protocol
	if err := stub.Connect(opt); err != nil {
		logger.E("Generate", svr, "mqtt client err:", err)
		return err
	}
	return nil
}

// New client from given options and connect with broker
func (stub *MqttStub) Connect(opt *mq.ClientOptions) error {
	client := mq.NewClient(opt)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		logger.E("Connect mqtt client, err:", token.Error())
		return token.Error()
	}

	stub.Client = client
	return nil
}

// Set mqtt client connect and disconnect handler, it called must before calling GenConfigs()
func (stub *MqttStub) SetHandlers(conn mq.OnConnectHandler, disc mq.ConnectionLostHandler) *MqttStub {
	stub.ConnectHandler, stub.DisconnectHandler = conn, disc
	return stub
}

// Publish empty message topic, it same use for just notify
func (stub *MqttStub) Notify(topic string, Qos ...byte) error {
	return stub.Publish(topic, nil, Qos...)
}

// Publish indicate topic message which formated as bytes array, and can set Qos to 0 ~ 2
func (stub *MqttStub) Publish(topic string, data interface{}, Qos ...byte) error {
	if stub.Client == nil {
		logger.E("Abort publish topic:", topic, "on nil client!!")
		return invar.ErrInvalidClient
	}

	payload := []byte{}
	if data != nil {
		if tmpd, err := json.Marshal(data); err != nil {
			return err
		} else {
			payload = tmpd
		}
	}

	qosv := byte(0)
	if len(Qos) > 0 && Qos[0] > 0 && Qos[0] <= 2 {
		qosv = Qos[0]
	}

	token := stub.Client.Publish(topic, qosv, false, payload)
	if token.Wait() && token.Error() != nil {
		logger.E("Publish topic:", topic, "err:", token.Error())
		return token.Error()
	}

	logger.I("Published topic:", topic)
	return nil
}

// Subscribe given topic and set callback
func (stub *MqttStub) Subscribe(topic string, hanlder mq.MessageHandler, Qos ...byte) error {
	if stub.Client == nil {
		logger.E("Abort subscribe topic:", topic, "on nil client!!")
		return invar.ErrInvalidClient
	}

	qosv := byte(0)
	if len(Qos) > 0 && Qos[0] > 0 && Qos[0] <= 2 {
		qosv = Qos[0]
	}

	token := stub.Client.Subscribe(topic, qosv, hanlder)
	if token.Wait() && token.Error() != nil {
		logger.E("Subscribe topic:", topic, "err:", token.Error())
		return token.Error()
	}
	logger.I("Subscribed topic:", topic)
	return nil
}

// Load and create secure configs for TLS protocol to connect.
func (stub *MqttStub) newTLSConfig() *tls.Config {
	ca, err := os.ReadFile(stub.Cfg.CAFile)
	if err != nil {
		logger.E("Read CA file err:", err)
		return nil
	}

	certpool := x509.NewCertPool()
	certpool.AppendCertsFromPEM(ca)
	tlsConfig := &tls.Config{RootCAs: certpool}

	// Import client certificate/key pair
	if stub.Cfg.CerFile != "" && stub.Cfg.KeyFile != "" {
		keyPair, err := tls.LoadX509KeyPair(stub.Cfg.CerFile, stub.Cfg.KeyFile)
		if err != nil {
			logger.E("Load cert and key err:", err)
			return nil
		}

		tlsConfig.ClientAuth = tls.NoClientCert
		tlsConfig.ClientCAs = nil
		tlsConfig.InsecureSkipVerify = true
		tlsConfig.Certificates = []tls.Certificate{keyPair}
	}
	return tlsConfig
}

// Parse mqtt broker and all user datas from nacos config center
func (stub *MqttStub) parseConfig(data, svr string) error {
	cfgs := &MqttConfigs{}
	if err := json.Unmarshal([]byte(data), &cfgs); err != nil {
		logger.E("Unmarshal mqtt settings, err:", err)
		return err
	}

	// Create client configs and fix the id as 'server@123456789'
	stub.Cfg = &ClientConfigs{
		Broker:   cfgs.Broker,
		Port:     cfgs.Port,
		ClientID: svr + "@" + secure.GenCode(),
		CAFile:   cfgs.CAFile,
		CerFile:  cfgs.CerFile,
		KeyFile:  cfgs.KeyFile,
	}

	if user, ok := cfgs.Users[svr]; ok {
		stub.Cfg.User = user
		return nil
	}
	return errors.New("Not found mqtt user: " + svr)
}
