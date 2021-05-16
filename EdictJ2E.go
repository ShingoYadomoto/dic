package main

import (
	"net/http"
	"time"
)

type EdictJ2E struct {
	client *http.Client

	baseURL        string
	searchItemPath string
	getItemPath    string

	dic       string
	scope     string
	match     string
	merge     string
	prof      string
	pageSize  int
	pageIndex int
}

func NewEdictJ2E(opts ...EdictJ2EOption) *EdictJ2E {
	dic := &EdictJ2E{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},

		baseURL:        BaseURLEJdict,
		searchItemPath: SearchItemPathEJdict,
		getItemPath:    GetItemPathEJdict,

		dic:       DicEdictJE,
		scope:     ScopeHeadword,
		match:     MatchTypeStartWith,
		merge:     MergeAnd,
		prof:      ProfXHTML,
		pageSize:  DefaultPageSize,
		pageIndex: DefaultPageIndex,
	}
	for _, opt := range opts {
		opt(dic)
	}

	return dic
}

func (d EdictJ2E) Translate(origin string) ([]Result, error) {
	itemIDList, err := d.searchItemIDList(origin)
	if err != nil {
		return nil, err
	}

	res := make([]Result, len(itemIDList))
	for i, itemID := range itemIDList {
		res[i], err = d.getResult(itemID)
		if err != nil {
			return nil, err
		}

	}
	return res, nil
}
