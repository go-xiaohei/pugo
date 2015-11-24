package model

type Tag struct {
	Name string
	Url  string
}

func NewTag(name string) Tag {
	return Tag{
		Name: name,
		Url:  "/tags/" + name + ".html",
	}
}
