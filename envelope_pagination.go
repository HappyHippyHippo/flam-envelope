package rest

import (
	"fmt"
)

type EnvelopePagination struct {
	Search string `json:"search,omitempty"`
	Start  int64  `json:"start"`
	Count  int64  `json:"count"`
	Total  int64  `json:"total"`
	Prev   string `json:"prev,omitempty"`
	Next   string `json:"next,omitempty"`
}

func NewEnvelopePagination(
	search string,
	start,
	count,
	total int64,
) EnvelopePagination {
	report := EnvelopePagination{
		Search: search,
		Start:  start,
		Count:  count,
		Total:  total,
		Prev:   "",
		Next:   "",
	}

	if start > 0 {
		nStart := int64(0)
		if count < start {
			nStart = start - count
		}
		report.Prev = fmt.Sprintf(
			"?search=%s&start=%d&count=%d",
			search,
			nStart,
			count)
	}

	if start+count < total {
		report.Next = fmt.Sprintf(
			"?search=%s&start=%d&count=%d",
			search,
			start+count,
			count)
	}

	return report
}
