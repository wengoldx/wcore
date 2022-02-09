Commands for generating go files and micro files:
protoc --proto_path=${GOPATH}/src:. --micro_out=. --go_out=. proto/xxx/xxx.proto


client call method:

import (
	acc "github.com/wengoldx/wing/proto/account"
	srv "github.com/wengoldx/wing/proto/vcall"
	"github.com/wengoldx/wing/tool/client"
)

// Cli agent service helper
var Cli srv.AgentService

// Acc account service handler
var Acc acc.AccountService

// Instantiate the client of the accservices service accent handler and the client of the Vcall service agent handler
c := client.NewClient("service.name")
Cli = srv.NewAgentService("service.name", c)
Acc = acc.NewAccountService("service.name", c)

// use 
 Acc.ViaToken(context.TODO(), &acc.Token{Token: b64token})