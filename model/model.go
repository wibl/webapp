package model

//Group is group of templates
type Group struct {
	ID    int64
	Title string
}

//Template is template
type Template struct {
	ID      int64
	GroupID int64
	Title   string
	Queue   string
	Body    string
}
