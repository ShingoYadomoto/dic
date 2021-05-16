package main

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

type (
	Result struct {
		Origin string
		Dist   string
	}

	Translator interface {
		Translate(string) ([]Result, error)
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

	EdictJ2ESearchXML struct {
		XMLName xml.Name `xml:"SearchDicItemResult"`

		// ErrorMessage  string `xml:"ErrorMessage"`
		// TotalHitCount int    `xml:"TotalHitCount"`
		// ItemCount int `xml:"ItemCount"`
		TitleList struct {
			DicItemTitleList []struct {
				ItemID string `xml:"ItemID"`
				Title  struct {
					Span string `xml:"span"`
				} `xml:"Title"`
				// LocID  string `xml:"LocID"`
			} `xml:"DicItemTitle"`
		} `xml:"TitleList"`
	}

	EdictJ2EGetXML struct {
		XMLName xml.Name `xml:"GetDicItemResult"`

		// ErrorMessage  string `xml:"ErrorMessage"`
		Head struct {
			Div struct {
				Span string `xml:"span"`
			} `xml:"div"`
		} `xml:"Head"`

		Body struct {
			Div struct {
				Div struct {
					DivList []string `xml:"div"`
				} `xml:"div"`
			} `xml:"div"`
		} `xml:"Body"`
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

var replacer = strings.NewReplacer(
	"\n", "",
	"\t", "",
)

func (xml EdictJ2ESearchXML) ItemIDList() []string {
	ret := make([]string, len(xml.TitleList.DicItemTitleList))
	for i, dicItemTitle := range xml.TitleList.DicItemTitleList {
		ret[i] = dicItemTitle.ItemID
	}
	return ret
}

func (xml EdictJ2EGetXML) Result() Result {
	return Result{
		Origin: replacer.Replace(xml.Head.Div.Span),
		Dist:   replacer.Replace(xml.Body.Div.Div.DivList[0]),
	}
}

func (d EdictJ2E) searchItemIDList(word string) ([]string, error) {
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

	xml := &EdictJ2ESearchXML{}
	err = decoder.Decode(xml)
	if err != nil {
		return nil, err
	}

	return xml.ItemIDList(), nil
}

func (d EdictJ2E) getResult(itemID string) (Result, error) {
	// Section1. create url
	// TODO: ↓ struct to query params using reflect ↓
	u, err := url.Parse(d.baseURL)
	if err != nil {
		return Result{}, err
	}
	u.Path = path.Join(u.Path, d.getItemPath)

	q := u.Query()
	q.Add("Dic", d.dic)
	q.Add("Item", itemID)
	q.Add("Loc", "")
	q.Add("Prof", d.prof)

	u.RawQuery = q.Encode()
	ru := u.String()
	// TODO: ↑ struct to query params using reflect ↑

	// Section2. http request
	resp, err := d.client.Get(ru)
	if err != nil {
		return Result{}, err
	}
	defer resp.Body.Close()

	// Section3. decode response xml
	decoder := xml.NewDecoder(resp.Body)

	xml := &EdictJ2EGetXML{}
	err = decoder.Decode(xml)
	if err != nil {
		return Result{}, err
	}

	return xml.Result(), nil
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
