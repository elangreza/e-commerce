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
	sortDirections strings.Builder
	defaultSortKey string
}

func (pqr *PaginationParams) Validate() error {
	if pqr.Page < 1 {
		pqr.Page = 1
	}

	if pqr.Limit < 1 {
		pqr.Limit = 10
	}

	if len(pqr.Sorts) > 0 {
		newSorts := make([]string, 0)
		for _, sort := range pqr.Sorts {
			if strings.Contains(sort, ",") {
				newSorts = append(newSorts, strings.Split(sort, ",")...)
			} else {
				newSorts = append(newSorts, sort)
			}
		}

		for index, sortRaw := range newSorts {
			if sortRaw == "" {
				continue
			}

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

			pqr.sortDirections.WriteString(fmt.Sprintf("%s %s", value, direction))
			if len(newSorts) > 1 && index < len(newSorts)-1 {
				pqr.sortDirections.WriteString(", ")
			}
		}
	}

	return nil
}

func (pqr *PaginationParams) GetOrderClause() string {
	s := pqr.sortDirections.String()
	if len(s) > 0 {
		return s
	}

	if pqr.defaultSortKey != "" {
		return pqr.defaultSortKey + " desc"
	}

	return ""
}

// SetValidSortKey sets the valid sort keys for the pagination params.
// The first sort key will be used as the default sort key if no sort is provided.
// the default sort direction is desc.
// This method should be called before calling Validate().
func (pqr *PaginationParams) SetValidSortKey(sortKeys ...string) {
	if pqr.validSortKeys == nil {
		pqr.validSortKeys = make(map[string]bool)
	}

	for _, sortKey := range sortKeys {
		pqr.validSortKeys[sortKey] = true
	}

	if len(sortKeys) > 0 && pqr.defaultSortKey == "" {
		pqr.defaultSortKey = sortKeys[0]
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
