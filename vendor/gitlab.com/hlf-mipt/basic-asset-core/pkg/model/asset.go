package model

// Asset describes basic details of what makes up a simple asset
// Insert struct field in alphabetic order => to achieve determinism across languages
// golang keeps the order when marshal to json but doesn't order automatically
type Asset struct {
	ID             string `json:"id"`
	AppraisedValue int    `json:"appraisedValue"`
	Color          string `json:"color"`
	Size           int    `json:"size"`
	Type           string `json:"type"`
	Owner          string `json:"owner"`
}
