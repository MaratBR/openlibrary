package ao3import

type Ao3Book struct {
	Chapters   []Ao3Chapter
	ID         string
	Name       string
	Tags       map[string][]string
	AuthorName string
	Language   string
	Rating     string
	Summary    string
}

type Ao3Chapter struct {
	ID      string
	Title   string
	Content string
	Summary string
}
