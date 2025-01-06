package views

type ImageView struct {
	Uuid    string `json:"uuid"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	FullUrl string `json:"fullUrl"`
}
