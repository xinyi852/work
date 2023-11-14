package common

import (
	"encoding/xml"
)

type UsageResult struct {
	Data UsageData `xml:"data"`
	ID   int       `xml:"id"`
	// Name string    `xml:"name"`
	// GUID string    `xml:"guid"`
	// 其他字段...
}

type UsageData struct {
	UsageStat UsageStat `xml:"stat"`
	UsageDisk UsageDisk `xml:"disk_usage"`
}

type UsageStat struct {
	Traffic string `xml:"traffic"`
}

type UsageDisk struct {
	Httpdocs string `xml:"httpdocs"`
}

type UsageGet struct {
	Result UsageResult `xml:"result"`
}

type UsageSite struct {
	UsageGet UsageGet `xml:"get"`
}

type UsagePacket struct {
	XMLName   xml.Name  `xml:"packet"`
	Version   string    `xml:"version,attr"`
	UsageSite UsageSite `xml:"site"`
}
