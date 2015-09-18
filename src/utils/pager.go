package utils

import (
	"fmt"
	"html/template"
	"strconv"
)

// pager struct
type Pager struct {
	Current int
	All     int
	Pages   int
	Size    int
}

// create pager
func CreatePager(page, size, all int) *Pager {
	p := &Pager{
		Current: page,
		Size:    size,
		All:     all,
	}
	p.Pages = all / size
	if all%size > 0 {
		p.Pages++
	}
	return p
}

func (p *Pager) IsPrev() bool {
	return p.Current > 1
}

func (p *Pager) Prev() int {
	return p.Current - 1
}

func (p *Pager) IsNext() bool {
	return p.Current < p.Pages
}

func (p *Pager) Next() int {
	return p.Current + 1
}

// pager to HTML with number elements
func (p *Pager) HTML(layout string) template.HTML {
	tpl := ` <ul class="pager">`
	for i := 1; i <= p.Pages; i++ {
		if i == p.Current {
			tpl += `<li><a class="current" href="` + fmt.Sprintf(layout, i) + `">` + strconv.Itoa(i) + `</a></li>`
		} else {
			tpl += `<li><a href="` + fmt.Sprintf(layout, i) + `">` + strconv.Itoa(i) + `</a></li>`
		}
	}
	tpl += "</ul>"
	return template.HTML(tpl)
}

// pager to HTML with navigator elements
func (p *Pager) HTMLSimple(layout, lang string) template.HTML {
	tpl := `<div class="pager clear">`
	if p.Current > 1 {
		tpl += `<a class="prev" href="` + fmt.Sprintf(layout, p.Current-1) + `">PREV</a>`
	}
	if p.Current < p.Pages {
		tpl += `<a class="next" href="` + fmt.Sprintf(layout, p.Current+1) + `">NEXT</a>`
	}
	tpl += "</div>"
	return template.HTML(tpl)
}
