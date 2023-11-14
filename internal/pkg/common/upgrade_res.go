package common

type UpgradeResult struct {
	Status   string `xml:"status"`
	ErrCode  int    `xml:"errcode"`
	ErrText  string `xml:"errtext"`
	FilterID string `xml:"filter-id"`
	ID       int    `xml:"id"`
}

type SwitchSubscription struct {
	Result UpgradeResult `xml:"result"`
}

type WebSpace struct {
	SwitchSub SwitchSubscription `xml:"switch-subscription"`
}

type UpgradePacket struct {
	Version  string   `xml:"version,attr"`
	WebSpace WebSpace `xml:"webspace"`
}
