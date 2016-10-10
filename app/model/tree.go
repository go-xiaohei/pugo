package model

import (
	"path"
	"path/filepath"
	"sort"
	"strings"
)

const (
	// TreeIndex is index page tree node
	TreeIndex = "index"
	// TreePost is a post node
	TreePost = "post"
	// TreePage is a page node
	TreePage = "page"
	// TreeArchive is node of archive page
	TreeArchive = "archive"
	// TreePostList is node of list page of posts
	TreePostList = "post-list"
	// TreePostTag is node of list posts belongs to a tag
	TreePostTag = "post-tag"
	// TreeTag is node of tag page, no used now
	TreeTag = "tag"
)

// Tree describe the position of one file in all compiled files
type Tree struct {
	Title    string
	Link     string
	I18n     string
	Type     string
	Sort     int
	Dest     string
	URL      string
	children []*Tree
}

type treeSlice []*Tree

// implement sort.Sort interface
func (t treeSlice) Len() int           { return len(t) }
func (t treeSlice) Less(i, j int) bool { return t[i].Sort < t[j].Sort }
func (t treeSlice) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }

func NewTree(dest string) *Tree {
	return &Tree{
		Link: "/",
		Type: TreeIndex,
		I18n: "tree",
		Sort: 0,
		Dest: dest,
		URL:  "",
	}
}

// Children return nodes of tree by url link
func (t *Tree) Children(link ...string) []*Tree {
	if len(link) == 0 {
		return t.children
	}
	if t2 := t.subTree(link[0]); t2 != nil {
		return t2.children
	}
	return nil
}

// IsValid return whether the node is compiled or not
func (t *Tree) IsValid() bool {
	return t.Type != ""
}

func (t *Tree) subTree(link string) *Tree {
	link = strings.TrimSuffix(link, path.Ext(link))
	linkData := strings.Split(strings.Trim(link, "/"), "/")
	if len(linkData) == 0 {
		return nil
	}
	for _, c := range t.children {
		if c.Link == linkData[0] {
			if len(linkData) == 1 {
				return c
			}
			return c.subTree(strings.Join(linkData[1:], "/"))
		}
	}
	return nil
}

// Print print tree nodes as readable string
func (t *Tree) Print(prefix string) {
	println(prefix+"/"+t.Link, t.I18n, t.Title, "@"+t.Type, t.URL)
	for _, c := range t.children {
		c.Print(prefix + "---")
	}
}

// Add add tree node
func (t *Tree) Add(link, title, linkType string, s int) {
	link = filepath.ToSlash(link)
	link = strings.TrimPrefix(link, t.Dest+"/")
	linkData := strings.SplitN(link, "/", 2)
	if len(linkData) == 0 {
		return
	}
	isFind := false
	for _, c := range t.children {
		if c.Link == linkData[0] {
			c.Add(strings.Join(linkData[1:], "/"), title, linkType, s)
			isFind = true
			break
		}
	}
	if !isFind {
		tree := &Tree{
			Title: title,
			Link:  linkData[0],
			I18n:  t.I18n + "." + linkData[0],
			Type:  linkType,
			Sort:  s,
			URL:   path.Join(t.URL, linkData[0]),
		}
		if len(linkData) > 1 {
			tree.Add(strings.Join(linkData[1:], "/"), title, linkType, s)
			tree.Type = ""
			tree.Title = ""
		}
		t.children = append(t.children, tree)
		sort.Sort(treeSlice(t.children))
	}
}
