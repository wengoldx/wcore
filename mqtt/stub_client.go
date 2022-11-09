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
	"io/ioutil"

	"github.com/astaxie/beego"
	mq "github.com/eclipse/paho.mqtt.golang"
	"github.com/wengoldx/wing/logger"
)

// default connect handler, if you want to custom this, you can use GenConfig().SetOnConnectHandler()
var defConnectHandler mq.OnConnectHandler = func(client mq.Client) {
	opt := client.OptionsReader()
	logger.I("server", beego.BConfig.AppName, "already connected mqtt, client id is", opt.ClientID())
}

// default disconnect handler, if you want to custom this, you can use GenConfig().SetConnectionLostHandler()
var defConnectLostHandler mq.ConnectionLostHandler = func(client mq.Client, err error) {
	logger.I("disconnect, waitting for reconnect")
}

// generate mqtt config
// mode is default tcp, you can select the connection method, Currently, only "tcp" and "tls" are supported and
// you can use fixed clientid or random clientid
func GenConfig(mode string, config *MqttConfig) *mq.ClientOptions {
	opt := mq.NewClientOptions()

	protocol := "tcp://%s:%v"
	switch mode {
	case "tcp":
		protocol = "tcp://%s:%v"
	case "tls":
		{
			protocol = "ssl://%s:%v"
			tlsConfig := newTLSConfig(config.CAFile, config.CerFile, config.KeyFile)
			opt.SetTLSConfig(tlsConfig)
		}
	}
	broker := fmt.Sprintf(protocol, config.Broker, config.Port)
	opt.AddBroker(broker)
	opt.SetClientID(config.ClientID)
	opt.SetUsername(config.User)
	opt.SetPassword(config.PWD)
	opt.SetAutoReconnect(true)
	opt.SetOnConnectHandler(defConnectHandler)
	opt.SetConnectionLostHandler(defConnectLostHandler)
	return opt

}

func newTLSConfig(caPath string, cer ...string) *tls.Config {
	certpool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(caPath)
	if err != nil {
		logger.E("read ca file err", err)
		return nil
	}
	certpool.AppendCertsFromPEM(ca)
	tlsConfig := &tls.Config{
		RootCAs: certpool,
	}

	if len(cer) == 2 && cer[0] != "" && cer[1] != "" {
		// Import client certificate/key pair
		clientKeyPair, err := tls.LoadX509KeyPair(cer[0], cer[1])
		if err != nil {
			logger.E("load certificate and key err", err)
			return nil
		}
		tlsConfig.ClientAuth = tls.NoClientCert
		tlsConfig.ClientCAs = nil
		tlsConfig.InsecureSkipVerify = true
		tlsConfig.Certificates = []tls.Certificate{clientKeyPair}
	}
	return tlsConfig
}

func GenClient(opt *mq.ClientOptions) mq.Client {
	client := mq.NewClient(opt)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return client
}
