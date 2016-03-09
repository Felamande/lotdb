package query

import (
	"github.com/jinzhu/gorm"
)

type query struct {
	root gorm.DB
	db   *gorm.DB
}

func Connect(dialect string, args ...interface{}) (*query, error) {
	root, err := gorm.Open(dialect, args)
	if err != nil {
		return nil, err
	}
	db := root.Select("n1,n2,n3,n4,n5")
	return &query{root, db}, nil

}

func (q *query) Include(num int) *query {
	q.db = q.db.Where("? in (n1,n2,n3,n4,n5)", num)
	return q
}
func (q *query) Exclude(nums ...int) *query {
	for _, n := range nums {
		q.db = q.db.Where("not ? in (n1,n2,n3,n4,n5)", n)
	}

	return q
}

func (q *query) If(exsit int, exclude ...int) *query {

	return q
}
