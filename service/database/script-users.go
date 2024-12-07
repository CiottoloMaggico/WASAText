package database

const qCreateUser = `
	INSERT INTO User (uuid, username, photo)  VALUES (?, ?, ?);
`
const qUpdateUser = `
	UPDATE User SET username = ?, photo = ? WHERE uuid = ?;
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
