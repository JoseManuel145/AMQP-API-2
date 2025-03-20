package entities

type Report struct {
	ID      int    `json:"Id"`
	Title   string `json:"Title"`
	Content string `json:"Content"`
	Status  string `json:"Status"`
}
