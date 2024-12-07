package database

const qGetImage = `SELECT uuid, extension FROM Image WHERE uuid = ?`

const qCreateImage = `
	INSERT INTO Image (uuid, extension) VALUES (?, ?);
`
