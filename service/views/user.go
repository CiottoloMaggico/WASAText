package views

type UserView struct {
	Uuid     string `json:"uuid"`
	Username string `json:"username"`
	Photo    string `json:"photo"`
}

type UserWithImageView struct {
	Uuid     string    `json:"uuid"`
	Username string    `json:"username"`
	Photo    ImageView `json:"photo"`
}
