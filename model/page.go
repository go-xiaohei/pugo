package model

import (
	"fmt"
	"github.com/go-xiaohei/pugo-static/parser"
	"gopkg.in/ini.v1"
	"html/template"
	"os"
	"time"
)

type Page struct {
	Title      string
	Slug       string
	Url        string
	HoverClass string
	Template   string
	Desc       string
	Created    Time
	Updated    Time
	Author     Author
	Raw        []byte
	rawType    string
    Meta map[string]string

	fileName string
	fileTime time.Time
}

func (p *Page) ContentHTML() template.HTML {
	if p.rawType == "markdown" {
		return template.HTML(markdown(p.Raw))
	}
	return template.HTML(p.Raw)
}

func NewPage(blocks []parser.Block, fi os.FileInfo) (*Page, error) {
	if len(blocks) < 2 {
		return nil, ErrPostBlockError
	}
	p := &Page{
		fileName: fi.Name(),
		fileTime: fi.ModTime(),
	}

	// parse first ini block
	iniF, err := ini.Load(blocks[0].Bytes())
	if err != nil {
		return nil, err
	}
	section := iniF.Section("DEFAULT")
	p.Title = section.Key("title").String()
	p.Slug = section.Key("slug").String()
	p.Desc = section.Key("desc").String()
	p.HoverClass = section.Key("hover").String()
	p.Template = section.Key("template").MustString("page.html")
    p.Meta = make(map[string]string)

	ct, err := time.Parse("2006-01-02", section.Key("date").String())
	if err != nil {
		return nil, err
	}
	p.Created = NewTime(ct)
	if upStr := section.Key("update_date").String(); upStr != "" {
		ut, err := time.Parse("2006-01-02", upStr)
		if err != nil {
			return nil, err
		}
		p.Updated = NewTime(ut)
	} else {
		// p.Updated = NewTime(p.fileTime)
	}
	p.Author = Author{
		Name:  section.Key("author").String(),
		Email: section.Key("author_email").String(),
		Url:   section.Key("author_url").String(),
	}

	// parse markdown block
	p.rawType = blocks[1].Type()
	p.Raw = blocks[1].Bytes()

	// build url
	p.Url = fmt.Sprintf("/%s", p.Slug)

    if len(blocks) > 2 {
        // parse meta block
        iniF, err = ini.Load(blocks[2].Bytes())
        if err != nil {
            return nil, err
        }
        keys := iniF.Section("DEFAULT").Keys()
        for _,k := range keys{
            p.Meta[k.Name()] = k.String()
        }
    }
	return p, nil
}
