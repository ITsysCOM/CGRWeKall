/*
Real-time Online/Offline Charging System (OCS) for Telecom & ISP environments
Copyright (C) ITsysCOM GmbH

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>
*/

package console

import (
	"github.com/cgrates/cgrates/sla"
)

func init() {
	c := &CmdStatus{
		name:      "sla_info",
		rpcMethod: sla.SlaSv1Info,
	}
	commands[c.Name()] = c
	c.CommandExecuter = &CommandExecuter{c}
}

type CmdSlaInfo struct {
	name      string
	rpcMethod string
	rpcParams struct{}
	*CommandExecuter
}

func (self *CmdSlaInfo) Name() string {
	return self.name
}

func (self *CmdSlaInfo) RpcMethod() string {
	return self.rpcMethod
}

func (self *CmdSlaInfo) RpcParams(reset bool) interface{} {
	return self.rpcParams
}

func (self *CmdSlaInfo) PostprocessRpcParams() error {
	return nil
}

func (self *CmdSlaInfo) RpcResult() interface{} {
	var s map[string]interface{}
	return &s
}

func (self *CmdSlaInfo) ClientArgs() (args []string) {
	return
}
