package model

import (
	"sort"
	"testing"

	"github.com/go-xiaohei/pugo/app/helper"
	. "github.com/smartystreets/goconvey/convey"
)

func preparePosts() []*Post {
	var p []*Post
	p = append(p, &Post{
		Title:     "abc",
		Slug:      "abc",
		Date:      "2016-01-15 12:20",
		TagString: []string{"a", "b", "c"},
	}, &Post{
		Title:     "xyz",
		Slug:      "xyz",
		Date:      "2016-09-12 12:20",
		TagString: []string{"a", "b"},
	}, &Post{
		Title:     "123",
		Slug:      "123",
		Date:      "2015-04-15 12:20",
		TagString: []string{"b"},
	}, &Post{
		Title:     "uvw",
		Slug:      "uvw",
		Date:      "2014-01-26 12:20",
		TagString: []string{"c"},
	})
	for _, pp := range p {
		pp.normalize()
	}
	return p
}

func TestArchive(t *testing.T) {
	Convey("Archives", t, func() {
		a := NewArchive(preparePosts())
		So(a.Data, ShouldHaveLength, 3)

		a.SetDestURL("/archive.html")
		So(a.DestURL(), ShouldEqual, "/archive.html")
	})
}

func TestPosts(t *testing.T) {
	Convey("Posts", t, func() {
		ps := Posts(preparePosts())
		sort.Sort(ps)
		So(ps[0].Date, ShouldEqual, "2016-09-12 12:20")
		So(ps[1].Title, ShouldEqual, "abc")
		So(ps[3].Slug, ShouldEqual, "uvw")

		Convey("PostsTopN", func() {
			ps2 := ps.TopN(2)
			So(ps2, ShouldHaveLength, 2)
			So(ps2[0].Date, ShouldEqual, "2016-09-12 12:20")
			So(ps2[1].Title, ShouldEqual, "abc")

			ps3 := ps.TopN(100)
			So(ps3, ShouldHaveLength, 4)
		})

		Convey("PostsRange", func() {
			ps2 := ps.Range(1, 3)
			So(ps2, ShouldHaveLength, 3)
			So(ps2[0].Date, ShouldEqual, "2016-01-15 12:20")
			So(ps2[1].Title, ShouldEqual, "123")

			ps3 := ps.Range(100, 200)
			So(ps3, ShouldBeNil)
		})
	})
}

func TestTagPosts(t *testing.T) {
	tp := make(map[string]*TagPosts)
	ps := preparePosts()
	for _, p := range ps {
		for _, tag := range p.Tags {
			if _, ok := tp[tag.Name]; !ok {
				tp[tag.Name] = &TagPosts{}
			}
			tp[tag.Name].Posts = append(tp[tag.Name].Posts, p)
			tp[tag.Name].Tag = tag
		}
	}
	for _, t2 := range tp {
		sort.Sort(t2.Posts)
	}
	Convey("TagPosts", t, func() {
		So(tp["a"].Posts, ShouldHaveLength, 2)
		So(tp["b"].Posts, ShouldHaveLength, 3)
		So(tp["a"].Tag, ShouldNotBeNil)

		tp["a"].SetDestURL("/tag/a.html")
		So(tp["a"].DestURL(), ShouldEqual, "/tag/a.html")
	})
}

func TestPagePosts(t *testing.T) {
	posts := preparePosts()
	var (
		ppMap  = make(map[int]*PagerPosts)
		cursor = helper.NewPagerCursor(3, len(posts))
		page   = 1
		layout = "posts/%d.html"
	)
	for {
		pager := cursor.Page(page)
		if pager == nil {
			break
		}
		currentPosts := posts[pager.Begin:pager.End]
		pager.SetLayout(layout)
		pp := &PagerPosts{
			Posts: currentPosts,
			Pager: pager,
			URL:   pager.URL(),
		}
		pp.SetDestURL(pager.URL())
		ppMap[pager.Current] = pp
		page++
	}

	Convey("PagedPosts", t, func() {
		So(ppMap[1].Posts, ShouldHaveLength, 3)
		So(ppMap[2].Posts, ShouldHaveLength, 1)
		So(ppMap[1].Pager, ShouldNotBeNil)

		So(ppMap[1].Pager.URL(), ShouldEqual, "posts/1.html")
		So(ppMap[1].DestURL(), ShouldEqual, "posts/1.html")
	})
}
