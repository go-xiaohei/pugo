package model

import "github.com/go-xiaohei/pugo/app/helper"

// Posts are posts list
type Posts []*Post

// implement sort.Sort interface
func (p Posts) Len() int {
	return len(p)
}

func (p Posts) Less(i, j int) bool {
	return p[i].dateTime.Unix() > p[j].dateTime.Unix()
}
func (p Posts) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// TopN get top N posts from list
func (p Posts) TopN(i int) []*Post {
	if i > len(p) {
		i = len(p)
	}
	return p[:i]
}

// Range get ranged[i:j] posts from list
func (p Posts) Range(i, j int) []*Post {
	if i > len(p)-1 {
		return nil
	}
	return p[i : j+1]
}

// TagPosts are list of posts belongs to a tag
type TagPosts struct {
	Posts
	Tag     *Tag
	destURL string
}

// SetDestURL set destUrl to tag post list
func (tp *TagPosts) SetDestURL(url string) {
	tp.destURL = url
}

// DestURL return compile file name of tag post list
func (tp *TagPosts) DestURL() string {
	return tp.destURL
}

// PagerPosts are list of posts by pagination
type PagerPosts struct {
	Posts
	Pager   *helper.Pager
	destURL string
	URL     string
}

// SetDestURL set destUrl to paged post list
func (pp *PagerPosts) SetDestURL(url string) {
	pp.destURL = url
}

// DestURL return compile file name of paged post list
func (pp *PagerPosts) DestURL() string {
	return pp.destURL
}
