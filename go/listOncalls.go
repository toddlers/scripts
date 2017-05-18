package main

import (
	"flag"
	"fmt"
	"github.com/PagerDuty/go-pagerduty"
	"time"
)

type oncalls struct {
	name     string
	schedule string
}

func GetOnCalls(authtoken string, schedule, escalPolicy string) {

	var oc oncalls

	var scheduleOptions pagerduty.ListSchedulesOptions

	var scheduleID string

	client := pagerduty.NewClient(authtoken)

	if resp, err := client.ListSchedules(scheduleOptions); err == nil {
		for _, sched := range resp.Schedules {

			if schedule == sched.Name {
				scheduleID = sched.ID
				break
			}
		}
	} else {
		fmt.Println("Error : ", err)
	}

	var escalationPolicyId string

	escalPolicyOpts := pagerduty.ListEscalationPoliciesOptions{Query: escalPolicy}

	if escalResponse, err := client.ListEscalationPolicies(escalPolicyOpts); err == nil {
		escalationPolicyId = escalResponse.EscalationPolicies[0].ID
	} else {
		fmt.Println("Error :", err)
	}

	now := time.Now()

	onCallOptions := pagerduty.ListOnCallOptions{
		ScheduleIDs:         []string{scheduleID},
		Since:               now.String(),
		EscalationPolicyIDs: []string{escalationPolicyId},
	}

	oncallResponse, _ := client.ListOnCalls(onCallOptions)

	for _, p := range oncallResponse.OnCalls {
		oc.name = p.User.Summary
		oc.schedule = p.Schedule.Summary
	}

	fmt.Printf("%s\t%s\n", oc.name, oc.schedule)
}

func main() {
	schedule := flag.String("schedule", "", "Schedule ID")
	token := flag.String("token", "", "Pagerduty auth token")
	espol := flag.String("espol", "", "escalation policy")
	flag.Parse()
	GetOnCalls(*token, *schedule, *espol)
}
