package models

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Felamande/lotdb/settings"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3" //init sqlite3

	_ "github.com/go-sql-driver/mysql" //init mysql
)

func GetQueryResults(form QueryForm) ([][]int, error) {
	sum := form.Sum
	if !sumInRange(sum) {
		return nil, errors.New("没有该和值")
	}

	db, err := gorm.Open(settings.DB.Type, settings.DB.Uri)
	if err != nil {
		return nil, DatabaseError{settings.DB.Uri, err}
	}
	defer db.Close()

	tmp := db.Table("sum" + strconv.Itoa(sum)).Select("n1,n2,n3,n4,n5")

loop:
	for _, flt := range form.Filters {
		if flt.Value == 0 {
			continue
		}
		switch flt.Type {
		case typeInclude:

			if filterValOutOfRange(flt.Value) {
				return nil, fmt.Errorf("%d超出范围", flt.Value)
			}
			tmp = tmp.Where("? in (n1,n2,n3,n4,n5)", flt.Value)
		case typeExclude:
			if filterValOutOfRange(flt.Value) {
				continue loop
			}
			tmp = tmp.Where("not ? in (n1,n2,n3,n4,n5)", flt.Value)
		}

	}

	rows, err := tmp.Rows()

	if err != nil {
		return nil, DatabaseError{settings.DB.Uri, err}
	}

	var results [][]int
	for rows.Next() {
		var n1, n2, n3, n4, n5 int
		rows.Scan(&n1, &n2, &n3, &n4, &n5)
		results = append(results, []int{n1, n2, n3, n4, n5})
	}

	return results, nil
}

func sumInRange(sum int) bool {
	return sum <= 165 && sum >= 15
}

func filterValOutOfRange(val int) bool {
	return val > 35 || val < 0
}

type queryRe struct {
	N1 int
	N2 int
	N3 int
	N4 int
	N5 int
}
