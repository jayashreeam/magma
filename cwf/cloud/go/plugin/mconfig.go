/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package plugin

import (
	"log"
	"strings"

	"magma/cwf/cloud/go/cwf"
	"magma/cwf/cloud/go/plugin/models"
	fegmconfig "magma/feg/cloud/go/protos/mconfig"
	ltemconfig "magma/lte/cloud/go/protos/mconfig"
	merrors "magma/orc8r/cloud/go/errors"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/configurator"

	"github.com/go-openapi/swag"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

const (
	DefaultUeIpBlock = "192.168.128.0/24"
)

var networkServicesByName = map[string]ltemconfig.PipelineD_NetworkServices{
	"metering":           ltemconfig.PipelineD_METERING,
	"dpi":                ltemconfig.PipelineD_DPI,
	"policy_enforcement": ltemconfig.PipelineD_ENFORCEMENT,
}

type Builder struct{}

func (*Builder) Build(
	networkID string,
	gatewayID string,
	graph configurator.EntityGraph,
	network configurator.Network,
	mconfigOut map[string]proto.Message,
) error {
	// we only build an mconfig if carrier_wifi network configs exist
	inwConfig, found := network.Configs[cwf.CwfNetworkType]
	if !found || inwConfig == nil {
		return nil
	}
	nwConfig := inwConfig.(*models.NetworkCarrierWifiConfigs)
	gwConfig, err := graph.GetEntity(cwf.CwfGatewayType, gatewayID)
	if err == merrors.ErrNotFound {
		return nil
	}
	if err != nil {
		return errors.WithStack(err)
	}
	if gwConfig.Config == nil {
		return nil
	}

	vals, err := buildFromConfigs(nwConfig, gwConfig.Config.(*models.GatewayCwfConfigs))
	if err != nil {
		return errors.WithStack(err)
	}
	for k, v := range vals {
		mconfigOut[k] = v
	}
	return nil
}

func buildFromConfigs(nwConfig *models.NetworkCarrierWifiConfigs, gwConfig *models.GatewayCwfConfigs) (map[string]proto.Message, error) {
	ret := map[string]proto.Message{}
	if nwConfig == nil {
		return ret, nil
	}
	pipelineDServices, err := getPipelineDServicesConfig(nwConfig.NetworkServices)
	if err != nil {
		return ret, err
	}
	allowed_ues, err := getPipelineDAllowedUes(gwConfig.AllowedUes)
	if err != nil {
		return ret, err
	}

	eapAka := nwConfig.EapAka
	aaa := nwConfig.AaaServer
	if eapAka != nil {
		mc := &fegmconfig.EapAkaConfig{LogLevel: protos.LogLevel_INFO}
		protos.FillIn(eapAka, mc)
		ret["eap_aka"] = mc
	}
	if aaa != nil {
		mc := &fegmconfig.AAAConfig{LogLevel: protos.LogLevel_INFO}
		protos.FillIn(aaa, mc)
		ret["aaa_server"] = mc
	}
	ret["pipelined"] = &ltemconfig.PipelineD{
		LogLevel:      protos.LogLevel_INFO,
		UeIpBlock:     DefaultUeIpBlock, // Unused by CWF
		NatEnabled:    false,
		DefaultRuleId: swag.StringValue(nwConfig.DefaultRuleID),
		RelayEnabled:  true,
		Services:      pipelineDServices,
		AllowedUes:    allowed_ues,
	}
	ret["sessiond"] = &ltemconfig.SessionD{
		LogLevel:     protos.LogLevel_INFO,
		RelayEnabled: true,
	}
	ret["redirectd"] = &ltemconfig.RedirectD{
		LogLevel: protos.LogLevel_INFO,
	}
	return ret, err
}
func getPipelineDAllowedUes(allowed_ues models.AllowedUes) ([]*ltemconfig.PipelineD_AllowedUE, error) {
	ues := make([]*ltemconfig.PipelineD_AllowedUE, 0, len(allowed_ues))
	for _, entry := range allowed_ues {
		ues = append(ues, &ltemconfig.PipelineD_AllowedUE{Ip: entry.IP.String(), Key: int32(entry.Key)})
	}
	return ues, nil
}
func getPipelineDServicesConfig(networkServices []string) ([]ltemconfig.PipelineD_NetworkServices, error) {
	apps := make([]ltemconfig.PipelineD_NetworkServices, 0, len(networkServices))
	for _, service := range networkServices {
		mc, found := networkServicesByName[strings.ToLower(service)]
		if !found {
			log.Printf("CWAG: unknown network service name %s", service)
		} else {
			apps = append(apps, mc)
		}
	}
	if len(apps) == 0 {
		apps = []ltemconfig.PipelineD_NetworkServices{
			ltemconfig.PipelineD_METERING,
			ltemconfig.PipelineD_DPI,
			ltemconfig.PipelineD_ENFORCEMENT,
		}
	}
	return apps, nil
}
