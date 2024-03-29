// Copyright (c) 2018-Now Dunyu All Rights Reserved.
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
	"reflect"

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
//	if err := mqtt.GenClient(data); err != nil {
//		logger.E("Connect client err:", err)
//		return
//	}
//
//
// UseCase 2 : Using nacos MQTT configs and connect with callbacks.
//
//	stub := mqtt.SetOptions(2)
//	stub.ConnectHandler = ConnectHandler
//	if err := mqtt.GenClient(data); err != nil {
//		logger.E("Connect client err:", err)
//		return
//	}
//
//
// UseCase 3 : Using singleton stub set custom configs and connect.
//
//	stub := mqtt.Singleton()
//	// Here use stub.Client to set your client configs
//	if err := stub.Connect(stub.GenConfigs()); err != nil {
//		logger.E("Connect client err:", err)
//		return
//	}
type MqttStub struct {
	Options           *Options
	Client            mq.Client
	ConnectHandler    mq.OnConnectHandler      // Connect callback handler
	DisconnectHandler mq.ConnectionLostHandler // Disconnect callback handler
	ReconnectHandler  mq.ReconnectHandler      // Reconnect callback handler
	MessageHandler    mq.MessageHandler        // Default publish message callback handler
	qos               byte                     // The default qos for publish or subscribe
	remain            bool                     // The default remain flag
}

// Singleton mqtt stub instance
var mqttStub *MqttStub

// Default connect handler, change it before call GenConfigs().
var connHandler mq.OnConnectHandler = func(client mq.Client) {
	serve, opt := beego.BConfig.AppName, client.OptionsReader()
	logger.I("Server", serve, "connected mqtt as client:", opt.ClientID())
}

// Default disconnect handler, change it before call GenConfigs().
var lostHandler mq.ConnectionLostHandler = func(client mq.Client, err error) {
	serve, opt := beego.BConfig.AppName, client.OptionsReader()
	logger.W("Server", serve, "disconnect mqtt client:", opt.ClientID())
}

// Return Mqtt global singleton
func Singleton() *MqttStub {
	if mqttStub == nil {
		mqttStub = &MqttStub{
			Options: &Options{}, Client: nil,
			ConnectHandler: connHandler, DisconnectHandler: lostHandler,
			qos: byte(0), remain: false,
		}
	}
	return mqttStub
}

// Generate mqtt client and connect with MQTT broker, the client using
// 'tcp' protocol and fixed id as format 'server@12345678'.
func GenClient(configs string, server ...string) error {
	stub, svr := Singleton(), beego.BConfig.AppName
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

// Set default qos and remain flag
func SetOptions(qos byte, remain ...bool) *MqttStub {
	stub := Singleton()
	if len(remain) > 0 {
		stub.remain = remain[0]
	}
	stub.qos = qos
	return stub
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

	broker := fmt.Sprintf(protocol, stub.Options.Host, stub.Options.Port)
	options.AddBroker(broker)
	options.SetClientID(stub.Options.ClientID)
	options.SetUsername(stub.Options.User.Account)
	options.SetPassword(stub.Options.User.Password)
	options.SetAutoReconnect(true)

	// set callback handlers if exist
	options.SetOnConnectHandler(stub.ConnectHandler)
	options.SetConnectionLostHandler(stub.DisconnectHandler)
	options.SetReconnectingHandler(stub.ReconnectHandler)
	options.SetDefaultPublishHandler(stub.MessageHandler)
	return options
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

// Publish empty message topic, it same use for just notify
func (stub *MqttStub) Notify(topic string, Qos ...byte) error {
	return stub.PublishOptions(topic, nil, stub.remain, Qos...)
}

// Publish indicate topic message, the Qos can be set current call in 0 ~ 2
func (stub *MqttStub) Publish(topic string, data any, Qos ...byte) error {
	return stub.PublishOptions(topic, data, stub.remain, Qos...)
}

// Publish indicate topic message with input remain flag and Qos options,
//
// Notice that the data will encode as json bytes array if value type is Struct,
// Pointer or map, or instead nil data to empty bytes array.
func (stub *MqttStub) PublishOptions(topic string, data any, remain bool, Qos ...byte) error {
	if stub.Client == nil {
		logger.E("Abort publish topic:", topic, "on nil client!!")
		return invar.ErrInvalidClient
	}

	var payload any
	if data != nil {
		switch reflect.ValueOf(data).Kind() {
		case reflect.Struct, reflect.Pointer, reflect.Map:
			if buffer, err := json.Marshal(data); err != nil {
				return err
			} else {
				payload = buffer
			}
		default:
			payload = data
		}
	} else {
		payload = []byte{} // Instead nil to empty bytes
	}

	qosv := stub.qos
	if len(Qos) > 0 && Qos[0] > 0 && Qos[0] <= 2 {
		qosv = Qos[0]
	}

	token := stub.Client.Publish(topic, qosv, remain, payload)
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

	qosv := stub.qos
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
	opts := stub.Options
	ca, err := os.ReadFile(opts.CAFile)
	if err != nil {
		logger.E("Read CA file err:", err)
		return nil
	}

	certpool := x509.NewCertPool()
	certpool.AppendCertsFromPEM(ca)
	tlsConfig := &tls.Config{RootCAs: certpool}

	// Import client certificate/key pair
	if opts.CerFile != "" && opts.KeyFile != "" {
		keyPair, err := tls.LoadX509KeyPair(opts.CerFile, opts.KeyFile)
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
	runmode := beego.BConfig.RunMode
	if broker, ok := cfgs.Brokers[runmode]; !ok {
		return errors.New("Notfound broker on run mode:" + runmode)
	} else if user, ok := cfgs.Users[svr]; !ok {
		return errors.New("Notfound mqtt user: " + svr)
	} else {
		stub.Options = &Options{
			Host: broker.Host, Port: broker.Port, User: user,
			ClientID: svr + "@" + secure.GenCode(),
			CAFile:   cfgs.CAFile,
			CerFile:  cfgs.CerFile,
			KeyFile:  cfgs.KeyFile,
		}
	}
	return nil
}
