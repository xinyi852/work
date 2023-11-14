package common

import (
	"encoding/xml"
)

type PlanResult struct {
	Status string `xml:"status"`
	ID     int    `xml:"id"`
	Name   string `xml:"name"`
	GUID   string `xml:"guid"`
	// 其他字段...
}

type Get struct {
	Result PlanResult `xml:"result"`
}

type ServicePlan struct {
	Get Get `xml:"get"`
}

type PlanPacket struct {
	XMLName     xml.Name    `xml:"packet"`
	Version     string      `xml:"version,attr"`
	ServicePlan ServicePlan `xml:"service-plan"`
}
