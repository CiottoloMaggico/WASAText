package database

import (
	"database/sql"
	"encoding/json"
)

type Image struct {
	Filename string         `json:"filename"`
	Size     int64          `json:"-"`
	Owner    sql.NullString `json:"-"`
	Width    int            `json:"width"`
	Height   int            `json:"height"`
}

func (i *Image) fullUrl() string {
	// TODO: fix with configuration url
	return "http://127.0.0.1:3000/media/" + i.Filename
}

func (i *Image) MarshalJSON() ([]byte, error) {
	type Alias Image
	return json.Marshal(&struct {
		FullUrl string `json:"fullUrl"`
		*Alias
	}{
		i.fullUrl(),
		(*Alias)(i),
	})
}

//type ConversationSummary struct {
//	Id            int     `json:"id"`
//	Name          string  `json:"name"`
//	Photo         Image   `json:"photo"`
//	ChatType      string  `json:"type"`
//	LatestMessage Message `json:"latestMessage"`
//}
