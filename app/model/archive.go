package model

// Archive is archive set for posts
type Archive struct {
	Year  int // each list by year
	Posts []*Post
}

// NewArchive converts posts to archive
func NewArchive(posts []*Post) []*Archive {
	archives := []*Archive{}
	var (
		last, lastYear int
	)
	for _, p := range posts {
		if len(archives) == 0 {
			archives = append(archives, &Archive{
				Year:  p.Created().Year(),
				Posts: []*Post{p},
			})
			continue
		}
		last = len(archives) - 1
		lastYear = archives[last].Year
		if lastYear == p.Created().Year() {
			archives[last].Posts = append(archives[last].Posts, p)
			continue
		}
		archives = append(archives, &Archive{
			Year:  p.Created().Year(),
			Posts: []*Post{p},
		})
	}
	return archives
}
