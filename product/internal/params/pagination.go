package params

import (
	"errors"
	"fmt"
	"strings"
)

type PaginationParams struct {
	Search string
	Sorts  []string
	Limit  int64
	Page   int64

	// to validate sort keys
	validSortKeys map[string]bool

	// local var. Used for sorting in the DB
	// sortMap        map[string]string
	sortDirections []string
	orderClause    string
}

func (pqr *PaginationParams) Validate() error {
	if pqr.Page < 1 {
		pqr.Page = 1
	}

	if pqr.Limit < 1 {
		pqr.Limit = 10
	}

	if len(pqr.Sorts) > 0 {
		newSorts := []string{}
		for _, sort := range pqr.Sorts {
			if strings.Contains(sort, ",") {
				newSorts = append(newSorts, strings.Split(sort, ",")...)
			} else {
				newSorts = append(newSorts, sort)
			}
		}

		// pqr.sortMap = make(map[string]string)
		pqr.sortDirections = make([]string, len(newSorts))
		for index, sortRaw := range newSorts {
			parts := strings.Split(sortRaw, ":")
			if len(parts) != 2 {
				return fmt.Errorf("%s is not valid sort format", sortRaw)
			}

			value := strings.ToLower(strings.TrimSpace(parts[0]))
			direction := strings.ToLower(strings.TrimSpace(parts[1]))

			if direction != "asc" && direction != "desc" {
				return errors.New("not valid sort direction")
			}

			if _, ok := pqr.validSortKeys[value]; !ok {
				return fmt.Errorf("%s is not valid sort key", value)
			}

			pqr.sortDirections[index] = fmt.Sprintf("%s %s", value, direction)

		}

		// if sortDirections is empty,
		// use the default sort with the first element of sorts (index 0)
		if len(pqr.sortDirections) > 0 {
			pqr.orderClause = strings.Join(pqr.sortDirections, ", ")
		} else {
			pqr.orderClause = pqr.Sorts[0] + " desc"
		}
	}

	return nil
}

func (pqr *PaginationParams) GetOrderClause() string {
	return pqr.orderClause
}

func (pqr *PaginationParams) SetValidSortKey(sortKeys ...string) {
	if pqr.validSortKeys == nil {
		pqr.validSortKeys = make(map[string]bool)
	}

	for _, sortKey := range sortKeys {
		pqr.validSortKeys[sortKey] = true
	}
}

func (pqr *PaginationParams) GetTotalPages(totalData int64) int64 {
	if pqr.Limit == 0 {
		return 0
	}
	totalPages := totalData / pqr.Limit
	if totalData%pqr.Limit != 0 {
		totalPages++
	}
	return totalPages
}
