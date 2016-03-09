package models

type QueryForm struct {
	Sum     int   `form:"sum" json:"sum"`
	Include []int `form:"include" json:"include"`
	Exclude []int `form:"exclude" json:"exclude"`
}
