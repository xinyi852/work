package common

type CreatePleskServerGroupData struct {
	Name     string `valid:"name" json:"name"`
	FillType string `valid:"fill_type" json:"fill_type"`
	// Servers  string `valid:"servers" json:"servers"`
}

type DeletePleskServerGroupData struct {
	GroupId string `valid:"group_id" json:"group_id"`
}

type UpdatePleskServerGroupData struct {
	GroupId  string `valid:"group_id" json:"group_id"`
	Name     string `valid:"name" json:"name"`
	FillType string `valid:"fill_type" json:"fill_type"`
}
