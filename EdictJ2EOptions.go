package main

const (
	BaseURLEJdict string = "https://public.dejizo.jp/NetDicV09.asmx"

	SearchItemPathEJdict string = "SearchDicItemLite"
	GetItemPathEJdict    string = "GetDicItemLite"

	DicEdictJE   string = "EdictJE"
	DicEJdict    string = "EJdict"
	DicWikipedia string = "wpedia"

	ScopeHeadword string = "HEADWORD"
	ScopeAnywhere string = "ANYWHERE"

	MatchTypeStartWith string = "STARTWITH"
	MatchTypeEnd       string = "ENDWITH"
	MatchTypeContain   string = "CONTAIN"
	MatchTypeExact     string = "EXACT"

	MergeAnd string = "AND"
	MergeOr  string = "OR"

	ProfXHTML string = "XHTML"

	DefaultPageSize int = 3

	DefaultPageIndex int = 0
)

type EdictJ2EOption func(*EdictJ2E)

func EdictJ2EMatchScope(s, e, c bool) EdictJ2EOption {
	sc := ScopeHeadword
	m := MatchTypeStartWith

	if s {
		sc = ScopeHeadword
		m = MatchTypeStartWith
	}
	if e {
		sc = ScopeHeadword
		m = MatchTypeEnd
	}
	if c {
		sc = ScopeAnywhere
		m = MatchTypeContain
	}

	return func(d *EdictJ2E) {
		d.scope = sc
		d.match = m
	}
}

func EdictJ2EMerge(merge string) EdictJ2EOption {
	return func(d *EdictJ2E) {
		d.merge = merge
	}
}

func EdictJ2EProf(prof string) EdictJ2EOption {
	return func(d *EdictJ2E) {
		d.prof = prof
	}
}

func EdictJ2EPageSize(pageSize int) EdictJ2EOption {
	return func(d *EdictJ2E) {
		d.pageSize = pageSize
	}
}

func EdictJ2EPageIndex(pageIndex int) EdictJ2EOption {
	return func(d *EdictJ2E) {
		d.pageIndex = pageIndex
	}
}
