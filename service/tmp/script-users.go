package database

const qCreateUser = `
	INSERT INTO User (uuid, username, photo)  VALUES (?, ?, ?);
`
const qSetUsername = `
	UPDATE User SET username = ? WHERE uuid = ?;
`

const qSetPhoto = `
	UPDATE User SET photo = ? WHERE uuid = ?;
`

const qGetUserByUsername = `
	SELECT * FROM ViewUsers WHERE username = ?;
`

const qGetUserByUUID = `
	SELECT * FROM ViewUsers WHERE uUuid = ?;
`

const qGetUsersPaginated = `
	SELECT * FROM ViewUsers LIMIT ? OFFSET ?;
`

const qGetUsersCount = `
	SELECT COUNT(*) FROM User;
`
