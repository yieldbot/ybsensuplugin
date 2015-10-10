// Library for all functions used in Yieldbot alert handlers and dashboard generators
//
// LICENSE:
//   Copyright 2015 Yieldbot. <devops@yieldbot.com>
//   Released under the MIT License; see LICENSE
//   for details.

// Package dracky implements common data structures and functions for Yieldbot monitoring alerts and dashboards
package dracky

import (
	"encoding/json"
	"github.com/yieldbot/dhuran"
	"io/ioutil"
	"os"
)

// Generate a simple string for use by elasticsearch and internal logging of all monitoring alerts.
func Event_name(client string, check string) string {
	return client + "_" + check
}

// Set the correct device that is being monitored. In the case of snmp trap collection, containers, or applicance
// monitoring the device running the sensu-client may not be the device actually being monitored.
func (e Sensu_Event) Acquire_monitored_instance() string {
	var monitored_instance string
	if e.Check.Source != "" {
		monitored_instance = e.Check.Source
	} else {
		monitored_instance = e.Client.Name
	}
	return monitored_instance
}

// func Set_time(t int) string {
//
// 	timeStamp := time.Unix(unixIntValue, 0)
// 	timestamp = time.Unix(timestamp, 0).Format(time.RFC822Z)
//
// }

// Set the environment that the machine is running in based upon values
// dropped via Oahi during the Chef run.
func Define_sensu_env(env string) string {
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

// Convert the check result status from an integer to a string.
func Define_status(status int) string {
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
	default:
		return "UNKNOWN"
	}
}

// Creates a monitor name that is easliy searchable in ES using different
// levels of granularity.
func Create_check_name(check string) string {
	return check
}

// Calculate how long a monitor has been in a given state.
func Define_check_state_duration() string {
	return " "
}

// Read in the environment details provided by Oahi and drop it into a staticly defined struct.
func Set_sensu_env() *Env_Details {
	env_file, err := ioutil.ReadFile(ENVIRONMENT_FILE)
	if err != nil {
		dhuran.Check(err)
	}

	var env_details Env_Details
	err = json.Unmarshal(env_file, &env_details)
	if err != nil {
		dhuran.Check(err)
	}
	return &env_details
}

// Read in the check result generated by Snesu via stdin and drop it into a staticaly defined struct.
func (e Sensu_Event) Acquire_sensu_event() *Sensu_Event {
	results, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		dhuran.Check(err)
	}
	err = json.Unmarshal(results, &e)
	if err != nil {
		dhuran.Check(err)
	}
	return &e
}
