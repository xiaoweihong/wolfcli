package model

type Org struct {
	Id            string
	Name          string
	OrgLevel      int64
	OrgCode       int64
	SuperiorOrgId string
	Status        int64
	Orgs          []Org
}
