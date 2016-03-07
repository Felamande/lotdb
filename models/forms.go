package models

type QueryForm struct {
	Sum     int      `form:"sum" json:"sum"`
	Filters []Filter `form:"filters" json:"filters"`
}

type Filter struct {
	Type  string `form:"type" json:"type"`
	Value int    `form:"value" json:"value"`
}

const typeExclude = "exclude"
const typeInclude = "include"
