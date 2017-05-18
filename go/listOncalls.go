package main

import (
	"flag"
	"fmt"
	"github.com/PagerDuty/go-pagerduty"
	"log"
	"time"
)

type oncalls struct {
	name     string
	schedule string
}

type PagerDutyOptions struct {
	escalationPolicy string
	schedule         string
	client           *pagerduty.Client
}

func (pd *PagerDutyOptions) getScheduleId() string {

	var scheduleOptions pagerduty.ListSchedulesOptions

	var scheduleID string

	if resp, err := pd.client.ListSchedules(scheduleOptions); err == nil {
		for _, sched := range resp.Schedules {

			if pd.schedule == sched.Name {
				scheduleID = sched.ID
				break
			}
		}
	} else {
		fmt.Println("Error : ", err)
	}
	return scheduleID
}

func (pd *PagerDutyOptions) getEscalationPolicyId() string {

	var escalationPolicyId string

	escalPolicyOpts := pagerduty.ListEscalationPoliciesOptions{Query: pd.escalationPolicy}

	if escalResponse, err := pd.client.ListEscalationPolicies(escalPolicyOpts); err == nil {
		escalationPolicyId = escalResponse.EscalationPolicies[0].ID
	} else {
		fmt.Println("Error :", err)
	}
	return escalationPolicyId
}

func (pd *PagerDutyOptions) GetOnCalls() oncalls {

	var oc oncalls

	scheduleID := pd.getScheduleId()
	log.Println("Schedule ID is : ", scheduleID)
	escalationPolicyId := pd.getEscalationPolicyId()
	log.Println("Escalation Policy ID is : ", escalationPolicyId)

	now := time.Now()

	onCallOptions := pagerduty.ListOnCallOptions{
		ScheduleIDs:         []string{scheduleID},
		Since:               now.String(),
		EscalationPolicyIDs: []string{escalationPolicyId},
	}

	oncallResponse, _ := pd.client.ListOnCalls(onCallOptions)

	for _, p := range oncallResponse.OnCalls {
		oc.name = p.User.Summary
		oc.schedule = p.Schedule.Summary
	}

	return oc

}

func main() {
	schedule := flag.String("schedule", "", "Schedule ID")
	token := flag.String("token", "", "Pagerduty auth token")
	espol := flag.String("espol", "", "escalation policy")

	flag.Parse()

	if *schedule == "" || *token == "" || *espol == "" {
		fmt.Println("Missing input")
		fmt.Println("Usage:")
		flag.PrintDefaults()
		return
	}

	pdOpts := PagerDutyOptions{
		escalationPolicy: *espol,
		schedule:         *schedule,
		client:           pagerduty.NewClient(*token),
	}

	onCalls := pdOpts.GetOnCalls()
	log.Printf("Name: %s\t Schedule : %s\n", onCalls.name, onCalls.schedule)
}
