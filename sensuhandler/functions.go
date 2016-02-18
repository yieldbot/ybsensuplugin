// Library for all handler functions used by the Yieldbot Infrastructure
// teams in sensu.
//
// LICENSE:
//   Copyright 2015 Yieldbot. <devops@yieldbot.com>
//   Released under the MIT License; see LICENSE
//   for details.

// Package handler implements common data structures and functions for Yieldbot monitoring alerts and dashboards
package sensuhandler

import (
	"encoding/json"
	"github.com/yieldbot/sensuplugin/sensuutil"
	"io/ioutil"
	"os"
	"strings"
)

//cleanOutput will shorten the output to a more managable length
func CleanOutput(output string) string {
	return strings.Split(output, ":")[0]
}

// EventName generates a simple string for use by elasticsearch and internal logging of all monitoring alerts.
func EventName(client string, check string) string {
	return client + "_" + check
}

// AcquireMonitoredInstance sets the correct device that is being monitored. In the case of snmp trap collection, containers,
// or applicance monitoring the device running the sensu-client may not be the device actually being monitored.
func (e SensuEvent) AcquireMonitoredInstance() string {
	// var monitoredInstance string
	if e.Check.Source != "" {
		return e.Check.Source
	}
	return e.Client.Name
}

// func Set_time(t int) string {
//
// 	timeStamp := time.Unix(unixIntValue, 0)
// 	timestamp = time.Unix(timestamp, 0).Format(time.RFC822Z)
//
// }

func SetColor(status int) string {
	switch status {
	case 0:
		return "#33CC33"
	case 1:
		return "warning"
	case 2:
		return "#FF0000"
	case 3:
		return "#FF6600"
	default:
		return "#FF6600"
	}
}

// DefineSensuEnv sets the environment that the machine is running in based upon values
// dropped via Oahi during the Chef run.
func DefineSensuEnv(env string) string {
	switch env {
	case "prd":
		return "Prod "
	case "dev":
		return "Dev "
	case "stg":
		return "Stg "
	case "vagrant":
		return "Vagrant "
	default:
		return "Test "
	}
}

// DefineStatus converts the check result status from an integer to a string.
func DefineStatus(status int) string {
	switch status {
	case 0:
		return "OK"
	case 1:
		return "WARNING"
	case 2:
		return "CRITICAL"
	case 3:
		return "UNKNOWN"
	case 126:
		return "PERMISSION DENIED"
	case 127:
		return "CONFIG ERROR"
	case 129:
		return "GENERAL GOLANG ERROR"
	default:
		return "UNKNOWN"
	}
}

// CreateCheckName creates a monitor name that is easliy searchable in ES using different
// levels of granularity.
func CreateCheckName(check string) string {
	return strings.Replace(check, "-", ".", -1)
	// return fmtdCheck
}

// DefineCheckStateDuration calculates how long a monitor has been in a given state.
func DefineCheckStateDuration() int {
	return 0
}

// SetSensuEnv reads in the environment details provided by Oahi and drops them into a staticly defined struct.
func SetSensuEnv() *EnvDetails {
	envFile, err := ioutil.ReadFile(sensuutil.EnvironmentFile)
	if err != nil {
		sensuutil.EHndlr(err)
	}

	var envDetails EnvDetails
	err = json.Unmarshal(envFile, &envDetails)
	if err != nil {
		sensuutil.EHndlr(err)
	}
	return &envDetails
}

// AcquireSensuEvent reads in the check result generated by Snesu via stdin and drop it into a staticaly defined struct.
func (e SensuEvent) AcquireSensuEvent() *SensuEvent {
	results, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		sensuutil.EHndlr(err)
	}
	err = json.Unmarshal(results, &e)
	if err != nil {
		sensuutil.EHndlr(err)
	}
	return &e
}
