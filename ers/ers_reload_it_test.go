// +build integration

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
package ers

import (
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"path"
	"testing"

	"github.com/cgrates/cgrates/config"
	"github.com/cgrates/cgrates/engine"
	"github.com/cgrates/cgrates/utils"
)

var (
	reloadCfgPath string
	reloadCfg     *config.CGRConfig
	reloadRPC     *rpc.Client
	reloadTests   = []func(t *testing.T){
		testReloadITCreateCdrDirs,
		testReloadITInitConfig,
		testReloadITInitCdrDb,
		testReloadITResetDataDb,
		testReloadITStartEngine,
		testReloadITRpcConn,
		testReloadVerifyDisabledReaders,
		testReloadReloadConfig,
		testReloadVerifyFirstReload,
		testReloadITKillEngine,
	}
)

func TestERsReload(t *testing.T) {
	reloadCfgPath = path.Join(*dataDir, "conf", "samples", "ers_reload", "disabled")
	for _, test := range reloadTests {
		t.Run("TestERsReload", test)
	}
}

func testReloadITInitConfig(t *testing.T) {
	var err error
	if reloadCfg, err = config.NewCGRConfigFromPath(reloadCfgPath); err != nil {
		t.Fatal("Got config error: ", err.Error())
	}
}

// InitDb so we can rely on count
func testReloadITInitCdrDb(t *testing.T) {
	if err := engine.InitStorDb(reloadCfg); err != nil {
		t.Fatal(err)
	}
}

// Remove data in both rating and accounting db
func testReloadITResetDataDb(t *testing.T) {
	if err := engine.InitDataDb(reloadCfg); err != nil {
		t.Fatal(err)
	}
}

func testReloadITCreateCdrDirs(t *testing.T) {
	for _, dir := range []string{"/tmp/ers/in", "/tmp/ers/out",
		"/tmp/ers2/in", "/tmp/ers2/out"} {
		if err := os.RemoveAll(dir); err != nil {
			t.Fatal("Error removing folder: ", dir, err)
		}
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal("Error creating folder: ", dir, err)
		}
	}
}

func testReloadITStartEngine(t *testing.T) {
	if _, err := engine.StopStartEngine(reloadCfgPath, *waitRater); err != nil {
		t.Fatal(err)
	}
}

// Connect rpc client to rater
func testReloadITRpcConn(t *testing.T) {
	var err error
	reloadRPC, err = jsonrpc.Dial("tcp", reloadCfg.ListenCfg().RPCJSONListen) // We connect over JSON so we can also troubleshoot if needed
	if err != nil {
		t.Fatal("Could not connect to rater: ", err.Error())
	}
}

func testReloadVerifyDisabledReaders(t *testing.T) {
	if len(reloadCfg.ERsCfg().Readers) != 1 &&
		reloadCfg.ERsCfg().Readers[0].ID != utils.MetaDefault &&
		reloadCfg.ERsCfg().Enabled != false {
		t.Errorf("Unexpected active readers: <%+v>", utils.ToJSON(reloadCfg.ERsCfg().Readers))
	}
}

func testReloadReloadConfig(t *testing.T) {
	var reply string
	if err := reloadRPC.Call(utils.ConfigSv1ReloadConfig, &config.ConfigReloadWithArgDispatcher{
		Path:    path.Join(*dataDir, "conf", "samples", "ers_reload", "first_reload"),
		Section: config.ERsJson,
	}, &reply); err != nil {
		t.Error(err)
	} else if reply != utils.OK {
		t.Errorf("Unexpected reply received: <%+v>", reply)
	}
}

func testReloadVerifyFirstReload(t *testing.T) {
	var reply map[string]interface{}
	if err := reloadRPC.Call(utils.ConfigSv1GetJSONSection, &config.StringWithArgDispatcher{
		Section: config.ERsJson,
	}, &reply); err != nil {
		t.Error(err)
	} else if reply["Enabled"] != true {
		t.Errorf("Expecting: <true>, received: <%+v>", reply["Enabled"])
	} else if readers, canConvert := reply["Readers"].([]interface{}); !canConvert {
		t.Errorf("Cannot cast Readers to slice")
	} else if len(readers) != 3 { // 2 active readers and 1 default
		t.Errorf("Expecting: <2>, received: <%+v>", len(readers))
	}
}

func testReloadITKillEngine(t *testing.T) {
	if err := engine.KillEngine(*waitRater); err != nil {
		t.Error(err)
	}
}