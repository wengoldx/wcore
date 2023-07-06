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
	"fmt"
	"github.com/astaxie/beego"
	mq "github.com/eclipse/paho.mqtt.golang"
	"github.com/wengoldx/wing/logger"
	"io/ioutil"
)

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
func GenConfig(mode string, config *MqttConfig) *mq.ClientOptions {
	options, protocol := mq.NewClientOptions(), "tcp://%s:%v"
	if mode == "tls" {
		protocol = "ssl://%s:%v"
		tlsConfig := newTLSConfig(config.CAFile, config.CerFile, config.KeyFile)
		options.SetTLSConfig(tlsConfig)
	}

	broker := fmt.Sprintf(protocol, config.Broker, config.Port)
	options.AddBroker(broker)
	options.SetClientID(config.ClientID)
	options.SetUsername(config.User)
	options.SetPassword(config.PWD)
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

// Generate MQTT client from given options
func GenClient(opt *mq.ClientOptions) (mq.Client, error) {
	client := mq.NewClient(opt)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		logger.E("Generate MQTT client, err:", token.Error())
		return nil, token.Error()
	}
	return client, nil
}
