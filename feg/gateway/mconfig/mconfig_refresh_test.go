/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package mconfig provides gateway Go support for cloud managed configuration (mconfig)
package mconfig

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	fegmcfg "magma/feg/cloud/go/protos/mconfig"
	ltemcfg "magma/lte/cloud/go/protos/mconfig"
	orcprotos "magma/orc8r/cloud/go/protos"
	_ "magma/orc8r/cloud/go/protos/mconfig"

	"github.com/golang/protobuf/ptypes"
)

// JSON to recreate scenario with "static" mconfig file
const testMmconfigJsonV1 = `{
  "configs_by_key": {
    "mobilityd": {
      "@type": "type.googleapis.com/magma.mconfig.MobilityD",
      "logLevel": "INFO",
      "ipBlock": "192.168.128.0/24"
    },
    "mme": {
      "@type": "type.googleapis.com/magma.mconfig.MME",
      "mmeCode": 1,
      "mmeGid": 1,
      "logLevel": "INFO",
      "mcc": "001",
      "mnc": "01",
      "tac": 1,
      "enableDnsCaching": false,
      "relayEnabled": false,
      "nonEpsServiceControl": 0,
      "csfbMcc": "001",
      "csfbMnc": "01",
      "lac": 1
    },
    "enodebd": {
      "@type": "type.googleapis.com/magma.mconfig.EnodebD",
      "bandwidthMhz": 20,
      "specialSubframePattern": 7,
      "earfcndl": 44490,
      "logLevel": "INFO",
      "plmnidList": "00101",
      "pci": 260,
      "allowEnodebTransmit": false,
      "subframeAssignment": 2,
      "tac": 1
    },
    "control_proxy": {
      "@type": "type.googleapis.com/magma.mconfig.ControlProxy",
      "logLevel": "INFO"
    },
    "magmad": {
      "@type": "type.googleapis.com/magma.mconfig.MagmaD",
      "logLevel": "INFO",
      "checkinInterval": 60,
      "checkinTimeout": 10,
      "autoupgradeEnabled": false,
      "autoupgradePollInterval": 300,
      "package_version": "0.0.0-0"
    },
    "metricsd": {
      "@type": "type.googleapis.com/magma.mconfig.MetricsD",
      "logLevel": "INFO"
    },
    "pipelined": {
      "@type": "type.googleapis.com/magma.mconfig.PipelineD",
      "logLevel": "INFO",
      "ueIpBlock": "192.123.45.0/24",
      "natEnabled": true,
      "apps": []
    },
    "subscriberdb": {
      "@type": "type.googleapis.com/magma.mconfig.SubscriberDB",
      "logLevel": "INFO",
      "lteAuthOp": "EREREREREREREREREREREQ==",
      "lteAuthAmf": "gAA=",
      "relayEnabled": false
    },
    "dnsd": {
      "@type": "type.googleapis.com/magma.mconfig.DnsD",
      "logLevel": "INFO",
      "enableCaching": false,
      "localTTL": 0
    },
    "lighttpd": {
      "@type": "type.googleapis.com/magma.mconfig.LighttpD",
      "enableCaching": false,
      "logLevel": "INFO"
    },
    "directoryd": {
      "@type": "type.googleapis.com/magma.mconfig.DirectoryD",
      "logLevel": "INFO"
    },
    "policydb": {
      "@type": "type.googleapis.com/magma.mconfig.PolicyDB",
      "logLevel": "INFO"
    }
  }
}`

const testMmconfigJsonV2 = `{
	"configsByKey": {
		"s6a_proxy": {
			"@type": "type.googleapis.com/magma.mconfig.S6aConfig",
			"logLevel": "INFO",
			"server": {
				"protocol": "sctp",
				"address": "",
				"retransmits": 3,
				"watchdogInterval": 1,
				"retryCount": 5,
				"productName": "magma",
				"realm": "magma.com",
				"host": "magma.dagma.mai.com"
			}
		},
		"session_proxy": {
			"@type": "type.googleapis.com/magma.mconfig.SessionProxyConfig",
			"logLevel": "INFO",
			"gx": {
				"server": {
					 "protocol": "tcp",
					 "address": "",
					 "retransmits": 3,
					 "watchdogInterval": 1,
					 "retryCount": 5,
					 "productName": "magma",
		 			 "realm": "magma.com",
		 			 "host": "magma-fedgw.magma.com"
				}
			},
			"gy": {
				"server": {
					 "protocol": "tcp",
					 "address": "",
					 "retransmits": 3,
					 "watchdogInterval": 1,
					 "retryCount": 5,
					 "productName": "magma",
		 			 "realm": "magma.com",
		 			 "host": "magma-fedgw.magma.com"
				},
				"initMethod": "PER_KEY"
			}
		}
	}
}`

func TestGatewayMconfigRefresh(t *testing.T) {
	stopRefreshTicker() // stop non-test config refresh

	// Create tmp mconfig test file
	tmpfile, err := ioutil.TempFile("", MconfigFileName)
	if err != nil {
		t.Fatal(err)
	}
	// Write V1 style marshaled configs
	if _, err = tmpfile.Write([]byte(testMmconfigJsonV1)); err != nil {
		t.Fatal(err)
	}
	mcpath := tmpfile.Name()
	if err = tmpfile.Close(); err != nil {
		t.Fatal(err)
	}
	t.Logf("Created gateway config file: %s", mcpath)

	defer func() {
		t.Logf("Remove temporary gateway config file: %s", mcpath)
		os.Remove(mcpath)
	}()

	// Start configs refresh ticker
	ticker := time.NewTicker(time.Millisecond * 50)
	go func() {
		for {
			<-ticker.C
			refreshErr := refreshConfigsFrom(mcpath)
			if refreshErr != nil {
				t.Error(refreshErr)
			}
		}
	}()

	time.Sleep(time.Millisecond * 120)
	pdcfg := &ltemcfg.PipelineD{}
	err = GetServiceConfigs("pipelined", pdcfg)
	if err != nil {
		t.Fatal(err)
	}
	expectedIpBlock := "192.123.45.0/24"
	if pdcfg.UeIpBlock != expectedIpBlock {
		t.Fatalf("pipelined Configs UeIpBlock Mismatch %s != %s", pdcfg.UeIpBlock, expectedIpBlock)
	}
	mc := GetGatewayConfigs()
	expectedIpBlock = "192.123.155.0/24"
	pdcfg.UeIpBlock = expectedIpBlock
	mc.ConfigsByKey["pipelined"], err = ptypes.MarshalAny(pdcfg)
	if err != nil {
		t.Fatal(err)
	}

	s6aBindAddr := ":12345"
	// Test marshaling of new configs
	mc.ConfigsByKey["s6a_proxy"], err = ptypes.MarshalAny(
		&fegmcfg.S6AConfig{
			LogLevel: 1,
			Server: &fegmcfg.DiamClientConfig{
				Protocol:         "sctp",
				Address:          "",
				Retransmits:      0x3,
				WatchdogInterval: 0x1,
				RetryCount:       0x7,
				LocalAddress:     s6aBindAddr,
				ProductName:      "magma",
				Realm:            "magma.com",
				Host:             "magma-fedgw.magma.com",
			}})

	mc.ConfigsByKey["session_proxy"], err = ptypes.MarshalAny(
		&fegmcfg.SessionProxyConfig{
			LogLevel: 1,
			Gx: &fegmcfg.GxConfig{
				Server: &fegmcfg.DiamClientConfig{
					Protocol:         "tcp",
					Address:          "",
					Retransmits:      0x3,
					WatchdogInterval: 0x1,
					RetryCount:       0x5,
					ProductName:      "magma",
					Realm:            "magma.com",
					Host:             "magma-fedgw.magma.com",
				},
			},
			Gy: &fegmcfg.GyConfig{
				Server: &fegmcfg.DiamClientConfig{
					Protocol:         "tcp",
					Address:          "",
					Retransmits:      0x3,
					WatchdogInterval: 0x1,
					RetryCount:       0x5,
					ProductName:      "magma",
					Realm:            "magma.com",
					Host:             "magma-fedgw.magma.com",
				},
				InitMethod: fegmcfg.GyInitMethod_PER_KEY,
			}})

	marshaled, err := orcprotos.MarshalIntern(mc)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile(mcpath, marshaled, os.ModePerm)
	if err != nil {
		t.Fatal(err)

	}
	time.Sleep(time.Millisecond * 120)
	err = GetServiceConfigs("pipelined", pdcfg)
	if err != nil {
		t.Fatal(err)
	}
	if pdcfg.UeIpBlock != expectedIpBlock {
		t.Fatalf("pipelined Configs UeIpBlock Mismatch %s != %s", pdcfg.UeIpBlock, expectedIpBlock)
	}

	sdcfg := &fegmcfg.SessionProxyConfig{}
	err = GetServiceConfigs("session_proxy", sdcfg)
	if err != nil {
		t.Fatal(err)
	}
	if sdcfg.Gy.InitMethod != fegmcfg.GyInitMethod_PER_KEY {
		t.Fatalf(
			"session_proxy Configs Gy.InitMethod Mismatch %d != %d",
			sdcfg.Gy.InitMethod, fegmcfg.GyInitMethod_PER_KEY)
	}

	// Test API's type enforcement/safety
	err = GetServiceConfigs("s6a_proxy", sdcfg)
	if err == nil {
		t.Fatal("Expected Error Getting s6a_proxy configs into SessionProxyConfig type")
	}

	s6acfg := &fegmcfg.S6AConfig{}
	err = GetServiceConfigs("s6a_proxy", s6acfg)
	if err != nil {
		t.Fatal(err)
	}
	if s6acfg.GetServer().GetLocalAddress() != s6aBindAddr {
		t.Fatalf(
			"s6a_proxy Configs Local Address Mismatch %s != %s",
			s6acfg.GetServer().GetLocalAddress(), s6aBindAddr)
		return
	}

	// test V2 - 'configsByKey' encoding version
	if err = ioutil.WriteFile(mcpath, []byte(testMmconfigJsonV2), os.ModePerm); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Millisecond * 120)
	const expectedHost = "magma.dagma.mai.com"
	err = GetServiceConfigs("s6a_proxy", s6acfg)
	if err != nil {
		t.Fatal(err)
	}
	if s6acfg.GetServer().GetHost() != expectedHost {
		t.Fatalf(
			"s6a_proxy Configs Host Mismatch %s != %s",
			s6acfg.GetServer().GetHost(), expectedHost)
	}
}
