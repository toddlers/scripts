package main

import (
	"flag"
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

	var scheduleId string

	if resp, err := pd.client.ListSchedules(scheduleOptions); err == nil {
		for _, sched := range resp.Schedules {

			if pd.schedule == sched.Name {
				scheduleId = sched.ID
				break
			}
		}
	} else {
		log.Println("Error : ", err)
	}
	return scheduleId
}

func (pd *PagerDutyOptions) getEscalationPolicyId() string {

	var escalationPolicyId string

	escalationPolicyOptions := pagerduty.ListEscalationPoliciesOptions{Query: pd.escalationPolicy}

	if resp, err := pd.client.ListEscalationPolicies(escalationPolicyOptions); err == nil {
		escalationPolicyId = resp.EscalationPolicies[0].ID
	} else {
		log.Println("Error :", err)
	}
	return escalationPolicyId
}

func (pd *PagerDutyOptions) GetOnCalls() oncalls {

	var oc oncalls

	scheduleId := pd.getScheduleId()
	log.Println("Schedule ID is : ", scheduleId)
	escalationPolicyId := pd.getEscalationPolicyId()
	log.Println("Escalation Policy ID is : ", escalationPolicyId)

	now := time.Now()

	onCallOptions := pagerduty.ListOnCallOptions{
		ScheduleIDs:         []string{scheduleId},
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
		log.Println("Missing input")
		log.Println("Usage:")
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
