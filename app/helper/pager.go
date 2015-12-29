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

		layout string
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

// create Pager on a page number
func (p *PagerCursor) Page(i int) *Pager {
	if i < 1 {
		return nil
	}
	begin := (i - 1) * p.size
	if begin > p.all {
		return nil // no pager when begin number over all
	}
	pager := &Pager{
		Begin:   begin,
		Prev:    i - 1,
		Next:    i + 1,
		Current: i,
		Pages:   p.pages,
	}
	end := begin + p.size
	if end >= p.all {
		end = p.all
		pager.Next = 0 // no next
	}
	pager.End = end
	return pager
}

// set pager layout string,
// use to print url
func (pg *Pager) SetLayout(layout string) {
	pg.layout = layout
}

// prev url
func (pg *Pager) PrevUrl() string {
	if pg.Prev > 0 {
		return fmt.Sprintf(pg.layout, pg.Prev)
	}
	return ""
}

// next url
func (pg *Pager) NextUrl() string {
	if pg.Next > 0 {
		return fmt.Sprintf(pg.layout, pg.Next)
	}
	return ""
}
