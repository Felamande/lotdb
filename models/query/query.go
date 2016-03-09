package query

import (
	"strconv"

	"github.com/Felamande/lotdb/settings"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3" //init sqlite3

	_ "github.com/go-sql-driver/mysql" //init mysql
)

type query struct {
	db   *gorm.DB
	errs []error
}

func Connect(dialect string, args ...interface{}) (*query, error) {
	db, err := gorm.Open(dialect, args...)
	if err != nil {
		return nil, err
	}
	db = db.Select("n1,n2,n3,n4,n5")
	return &query{db, []error{}}, nil

}

func (q *query) Sum(sum int) *query {
	if len(q.errs) != 0 {
		return q
	}

	q.db = q.db.Table("sum" + strconv.Itoa(sum))
	return q
}

func (q *query) Include(nums ...int) *query {
	if len(q.errs) != 0 {
		return q
	}
	for _, num := range nums {
		if num == 0 {
			continue
		}
		q.db = q.db.Where("? in (n1,n2,n3,n4,n5)", num)
	}

	return q
}
func (q *query) Exclude(nums ...int) *query {
	if len(q.errs) != 0 {
		return q
	}

	for _, num := range nums {
		if filterValOutOfRange(num) {
			continue
		}
		q.db = q.db.Where("not ? in (n1,n2,n3,n4,n5)", num)
	}

	return q
}

func (q *query) If(exsit int, exclude ...int) *query {
	if len(q.errs) != 0 {
		return q
	}
	return q
}
func (q *query) Result() ([][]int, []error) {
	defer q.db.Close()

	if len(q.errs) != 0 {
		return nil, q.errs
	}
	rows, err := q.db.Rows()

	if err != nil {
		q.errs = append(q.errs, DatabaseError{settings.DB.Uri, err})
		return nil, q.errs
	}

	var results [][]int
	for rows.Next() {
		var n1, n2, n3, n4, n5 int
		rows.Scan(&n1, &n2, &n3, &n4, &n5)
		results = append(results, []int{n1, n2, n3, n4, n5})
	}

	return results, nil
}

func filterValOutOfRange(val int) bool {
	return val > 35 || val <= 0
}

func (q *query) Count() int {
	var i int
	q.db.Count(&i)
	return i
}
