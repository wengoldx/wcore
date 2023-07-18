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
	"fmt"
	"github.com/astaxie/beego"
	mq "github.com/eclipse/paho.mqtt.golang"
	"github.com/wengoldx/wing/logger"
	"github.com/wengoldx/wing/secure"
	"io/ioutil"
)

type MqttStub struct {
	Cfg    *MqttConfig
	Client mq.Client
}

// Singleton mqtt stub instance
var mqttStub *MqttStub

// Return Mqtt global Singleton
func Singleton() *MqttStub {
	if mqttStub == nil {
		mqttStub = &MqttStub{
			Cfg:    &MqttConfig{},
			Client: nil,
		}
	}
	return mqttStub
}

// Default connect handler, if you want to custom this, you can use GenConfig().SetOnConnectHandler()
var defConnectHandler mq.OnConnectHandler = func(client mq.Client) {
	serve, opt := beego.BConfig.AppName, client.OptionsReader()
	logger.I("Server", serve, "already connected mqtt as client:", opt.ClientID())
}

// Default disconnect handler, if you want to custom this, you can use GenConfig().SetConnectionLostHandler()
var defConnectLostHandler mq.ConnectionLostHandler = func(client mq.Client, err error) {
	serve, opt := beego.BConfig.AppName, client.OptionsReader()
	logger.W("Server", serve, "disconnect mqtt client:", opt.ClientID())
}

// Generate mqtt config, default connection mode using tcp, you can select the connection method,
// Currently, only support "tcp" or "tls", and fixed clientid or random clientid
func (stub *MqttStub) GenConfig(mode, svr string) *mq.ClientOptions {
	options, protocol := mq.NewClientOptions(), "tcp://%s:%v"
	if mode == "tls" {
		protocol = "ssl://%s:%v"
		tlsConfig := newTLSConfig(stub.Cfg.CAFile, stub.Cfg.CerFile, stub.Cfg.KeyFile)
		options.SetTLSConfig(tlsConfig)
	}

	broker := fmt.Sprintf(protocol, stub.Cfg.Broker, stub.Cfg.Port)
	options.AddBroker(broker)
	options.SetClientID(stub.Cfg.SvrCfg[svr].ClientID)
	options.SetUsername(stub.Cfg.SvrCfg[svr].User)
	options.SetPassword(stub.Cfg.SvrCfg[svr].PWD)
	options.SetAutoReconnect(true)
	options.SetOnConnectHandler(defConnectHandler)
	options.SetConnectionLostHandler(defConnectLostHandler)
	return options

}

func newTLSConfig(caPath string, cer ...string) *tls.Config {
	ca, err := ioutil.ReadFile(caPath)
	if err != nil {
		logger.E("Read CA file err:", err)
		return nil
	}

	certpool := x509.NewCertPool()
	certpool.AppendCertsFromPEM(ca)
	tlsConfig := &tls.Config{
		RootCAs: certpool,
	}

	if len(cer) == 2 && cer[0] != "" && cer[1] != "" {
		// import client certificate/key pair
		clientKeyPair, err := tls.LoadX509KeyPair(cer[0], cer[1])
		if err != nil {
			logger.E("Load certificate and key err", err)
			return nil
		}

		tlsConfig.ClientAuth = tls.NoClientCert
		tlsConfig.ClientCAs = nil
		tlsConfig.InsecureSkipVerify = true
		tlsConfig.Certificates = []tls.Certificate{clientKeyPair}
	}
	return tlsConfig
}

// New MQTT client from given options
func (stub *MqttStub) newClient(opt *mq.ClientOptions) (mq.Client, error) {
	client := mq.NewClient(opt)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		logger.E("Generate MQTT client, err:", token.Error())
		return nil, token.Error()
	}
	return client, nil
}

// `USAGE`
//
//	var MqClient mq.Client
//	if err := mqtt.Singleton().ParseConfig(data); err != nil {
//		logger.E("Parce mqtt config err:", err)
//		return
//	}
//	mqtt.Singleton().GenClient()
//	MqClient = mqtt.Singleton().Client

//	Parse all grpc certs from nacos config data, and cache to certs map
func (stub *MqttStub) ParseConfig(data string) error {
	cfg := &MqttConfig{}
	if err := json.Unmarshal([]byte(data), &cfg); err != nil {
		logger.E("Unmarshal mqtt setting, err:", err)
		return err
	}

	stub.Cfg = cfg
	return nil
}

// Generate mqtt client handler
func (stub *MqttStub) GenClient(server ...string) {
	svr := beego.BConfig.AppName
	if len(server) > 0 && server[0] != "" {
		svr = server[0]
	}

	if _, ok := stub.Cfg.SvrCfg[svr]; !ok {
		logger.E("Not found server mqtt config!")
		return
	}

	stub.Cfg.SvrCfg[svr].ClientID = secure.GenCode()
	opt := stub.GenConfig("tcp", svr)
	client, err := stub.newClient(opt)
	if err != nil {
		logger.E("Generate", svr, "mqtt client err:", err)
		return
	}
	stub.Client = client
}
