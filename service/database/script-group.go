package database

const qSetGroupName = `
	UPDATE GroupConversation SET name = ? WHERE id = ?;
`

const qSetGroupPhoto = `
	UPDATE GroupConversation SET photo = ? WHERE id = ?;
`
