package main

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"
)

type (
	Translator interface {
		Translate(string) (string, error)
	}

	EdictJ2E struct {
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

	EdictJ2EXML struct {
		XMLName       xml.Name `xml:"SearchDicItemResult"`
		ErrorMessage  string   `xml:"ErrorMessage"`
		TotalHitCount int      `xml:"TotalHitCount"`
		ItemCount     int      `xml:"ItemCount"`

		TitleList struct {
			DicItemTitle []struct {
				ItemID string `xml:"ItemID"`
				Title  string `xml:"Title"`
				LocID  string `xml:"LocID"`
			} `xml:"DicItemTitle"`
		} `xml:"TitleList"`
	}
)

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

func (d EdictJ2E) Translate(origin string) (string, error) {
	_, err := d.searchItemIDs(origin)
	if err != nil {
		return "", err
	}
	return "", nil
}

func (d EdictJ2E) searchItemIDs(word string) (*EdictJ2EXML, error) {
	// Section1. create url
	// TODO: ↓ struct to query params using reflect ↓
	u, err := url.Parse(d.baseURL)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, d.searchItemPath)

	q := u.Query()
	q.Add("Dic", d.dic)
	q.Add("Scope", d.scope)
	q.Add("Match", d.match)
	q.Add("Merge", d.merge)
	q.Add("Prof", d.prof)
	q.Add("PageSize", fmt.Sprint(d.pageSize))
	q.Add("PageIndex", fmt.Sprint(d.pageIndex))
	q.Add("Word", word)

	u.RawQuery = q.Encode()
	ru := u.String()
	// TODO: ↑ struct to query params using reflect ↑

	// Section2. http request
	resp, err := d.client.Get(ru)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Section3. decode response xml
	decoder := xml.NewDecoder(resp.Body)

	xml := &EdictJ2EXML{}
	err = decoder.Decode(xml)
	if err != nil {
		return nil, err
	}

	return xml, nil
}
