package helper

import "fmt"

type (
	// PagerCursor creates Pager with each page number
	PagerCursor struct {
		all   int
		size  int
		pages int
	}
	// Pager contains pagination data when on a page number
	Pager struct {
		Begin   int
		End     int
		Prev    int
		Next    int
		Current int
		Pages   int
		All     int

		layout string
	}
	PagerItem struct {
		Page int
		Link string
	}
)

// NewPagerCursor with size and all count
func NewPagerCursor(size, all int) *PagerCursor {
	pc := &PagerCursor{
		all:  all,
		size: size,
	}
	if all%size == 0 {
		pc.pages = all / size
	} else {
		pc.pages = all/size + 1
	}
	return pc
}

// Page creates Pager on a page number
func (pg *PagerCursor) Page(i int) *Pager {
	if i < 1 {
		return nil
	}
	begin := (i - 1) * pg.size
	if begin > pg.all {
		return nil // no pager when begin number over all
	}
	pager := &Pager{
		Begin:   begin,
		Prev:    i - 1,
		Next:    i + 1,
		Current: i,
		Pages:   pg.pages,
		All:     pg.all,
	}
	end := begin + pg.size
	if end >= pg.all {
		end = pg.all
		pager.Next = 0 // no next
	}
	pager.End = end
	return pager
}

// SetLayout sets pager layout string,
// use to print url
func (pg *Pager) SetLayout(layout string) {
	pg.layout = layout
}

// PrevURL returns prev url
func (pg *Pager) PrevURL() string {
	if pg.Prev > 0 {
		return fmt.Sprintf(pg.layout, pg.Prev)
	}
	return ""
}

// NextURL returns next url
func (pg *Pager) NextURL() string {
	if pg.Next > 0 {
		return fmt.Sprintf(pg.layout, pg.Next)
	}
	return ""
}

// URL returns page current url
func (pg *Pager) URL() string {
	return fmt.Sprintf(pg.layout, pg.Current)
}

// PageItems returns each page item in this pager
func (pg *Pager) PageItems() []*PagerItem {
	var items []*PagerItem
	for i := 1; i <= pg.Pages; i++ {
		item := &PagerItem{
			Page: i,
			Link: fmt.Sprintf(pg.layout, i),
		}
		items = append(items, item)
	}
	return items
}
