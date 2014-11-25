// Copyright (c) Microsoft Open Technologies, Inc.
// All Rights Reserved.
// Licensed under the Apache License, Version 2.0.
// See License.txt in the project root for license information.
package common

import (
	"fmt"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
	"strings"
	"time"
	"log"
	powershell "github.com/MSOpenTech/packer-hyperv/packer/powershell"
)


type StepConfigureIp struct {
	ip string
}

func (s *StepConfigureIp) Run(state multistep.StateBag) multistep.StepAction {
//	driver := state.Get("driver").(Driver)
	ui := state.Get("ui").(packer.Ui)

	errorMsg := "Error configuring ip address: %s"
	vmName := state.Get("vmName").(string)

	ui.Say("Configuring ip address...")

	var script ScriptBuilder
	script.WriteLine("param([string]$vmName)")
	script.WriteLine("try {")
	script.WriteLine("  $adapter = Get-VMNetworkAdapter -VMName $vmName -ErrorAction SilentlyContinue")
	script.WriteLine("  $ip = $adapter.IPAddresses[0]")
	script.WriteLine("  if($ip -eq $null) {")
	script.WriteLine("    return $false")
	script.WriteLine("  }")
	script.WriteLine("} catch {")
	script.WriteLine("  return $false")
	script.WriteLine("}")
	script.WriteLine("$ip")

	count := 60
	var duration time.Duration = 1
	sleepTime := time.Minute * duration
	var ip string

	for count != 0 {
		powershell := new(powershell.PowerShellCmd)
		cmdOut, err := powershell.Output(script.String(), vmName);
		if err != nil {
			err := fmt.Errorf(errorMsg, err)
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}

		ip = strings.TrimSpace(string(cmdOut))

		if ip != "False" {
			break;
		}

		log.Println(fmt.Sprintf("Waiting for another %v minutes...", uint(duration)))
		time.Sleep(sleepTime)
		count--
	}

	if(count == 0){
		err := fmt.Errorf(errorMsg, "IP address assigned to the adapter is empty")
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	ui.Say("ip address is " + ip)
/*
	ui.Say("Adding to TrustedHosts (require elevated mode)")

	blockBuffer.Reset()
	blockBuffer.WriteString("start-process powershell -verb runas -argument ")
	blockBuffer.WriteString("\"Invoke-Command -scriptblock { Set-Item -path WSMan:\\localhost\\Client\\TrustedHosts '")
	blockBuffer.WriteString(ip)
	blockBuffer.WriteString("' -force}\"")

	err = driver.HypervManage(blockBuffer.String())

	if err != nil {
		err := fmt.Errorf(errorMsg, err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
*/
	s.ip = ip

	state.Put("ip", ip)


	return multistep.ActionContinue
}

func (s *StepConfigureIp) Cleanup(state multistep.StateBag) {
	// do nothing
}
