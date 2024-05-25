# Core module of Wengold backend service

Golang后台开发的工具集，涵盖了Elastic分词系统、日志、MQTT、HTTP Rest4接口调用、Nacos配置中心、GRPC服务间调用、Wechat支付、Socket.IO等常用逻辑和模块的二次封装，使用beego能够快速便捷搭建后台微服务。

* `elastic` : 分词服务
* `invar` : 各种常用常量定义、error定义
* `logger` : beego日志二次封装，便捷使用
* `mqtt` : MQTT client功能操作封装
* `mvc` : 基于beego框架的Rest4接口功能操作封装
* `nacos` : Nacos注册配置中心终端功能操作封装
* `secure` : ASE、RSA、MD5、Base64、Hash等加解密，编解码封装
* `utils` : 各种常用工具集，包括钉钉通知、文件处理、队列、堆栈、短信、邮件、Task、Time等便捷使用接口
* `wechat` : 微信支付v3的整合、封装
* `wrpc` : 各种grpc的使用
* `wsio` : Socket.IO功能操作的封装
