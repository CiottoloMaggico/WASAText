package api

type HTTPStatus struct {
	Status      int    `json:"code"`
	Description string `json:"description"`
}
