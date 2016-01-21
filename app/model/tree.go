package model

import (
	"path"
	"sort"
	"strings"
)

const (
	TreeIndex    = "index"
	TreePost     = "post"
	TreePage     = "page"
	TreeArchive  = "archive"
	TreePostList = "post-list"
	TreePostTag  = "post-tag"
	TreeTag      = "tag"
)

type Tree struct {
	Link     string
	I18n     string
	Type     string
	Sort     int
	children []*Tree
}

type treeSlice []*Tree

// implement sort.Sort interface
func (t treeSlice) Len() int           { return len(t) }
func (t treeSlice) Less(i, j int) bool { return t[i].Sort < t[j].Sort }
func (t treeSlice) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }

func NewTree() *Tree {
	return &Tree{
		Link: "/",
		Type: TreeIndex,
		I18n: "tree",
		Sort: 0,
	}
}

func (t *Tree) Children(link ...string) []*Tree {
	if len(link) == 0 {
		return t.children
	}
	if t2 := t.SubTree(link[0]); t2 != nil {
		return t2.children
	}
	return nil
}

func (t *Tree) SubTree(link string) *Tree {
	println("sub tree", link)
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
			return c.SubTree(strings.Join(linkData[1:], "/"))
		}
	}
	return nil
}

func (t *Tree) Print(prefix string) {
	if prefix == "" {
		prefix = "+"
	}
	println(prefix+":"+t.Link, t.I18n, "@"+t.Type)
	for _, c := range t.children {
		c.Print(prefix + "-")
	}
}

func (t *Tree) Add(link, linkType string, s int) {
	link = strings.TrimSuffix(link, path.Ext(link))
	linkData := strings.Split(strings.Trim(link, "/"), "/")
	if len(linkData) == 0 {
		return
	}
	isFind := false
	for _, c := range t.children {
		if c.Link == linkData[0] {
			c.Add(strings.Join(linkData[1:], "/"), linkType, s)
			isFind = true
			break
		}
	}
	if !isFind {
		tree := &Tree{
			Link: linkData[0],
			I18n: t.I18n + "." + linkData[0],
			Type: linkType,
			Sort: s,
		}
		if len(linkData) > 1 {
			tree.Add(strings.Join(linkData[1:], "/"), linkType, s)
			tree.Type = ""
		}
		t.children = append(t.children, tree)
		sort.Sort(treeSlice(t.children))
	}
}
