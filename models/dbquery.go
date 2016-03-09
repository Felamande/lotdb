package models

import (
	"errors"
	"fmt"

	"github.com/Felamande/lotdb/models/query"
	"github.com/Felamande/lotdb/settings"
)

func GetQueryResults(form QueryForm) ([][]int, error) {
	sum := form.Sum
	if sumOutOfRange(sum) {
		return nil, errors.New("没有该和值")
	}
	q, err := query.Connect(settings.DB.Type, settings.DB.Uri)
	if err != nil {
		fmt.Println(err)
		return nil, DatabaseError{settings.DB.Uri, err}
	}

	result, errs := q.Sum(form.Sum).Include(form.Include...).Exclude(form.Exclude...).Result()
	if len(errs) != 0 {
		return nil, errs[0]
	}

	return result, nil
}

func sumOutOfRange(sum int) bool {
	return sum > 165 || sum < 15
}
