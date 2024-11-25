package database

const getUsersCountQuery = `SELECT COUNT(*) FROM User;`
const getUserUUIDQuery = `SELECT uuid FROM User WHERE uuid = ?;`
const getUsersQuery = `
	SELECT User.uuid, User.username, Image.*
	FROM User, Image
	WHERE User.photoFilename = Image.filename
`
const getUserQuery = getUsersQuery + ` AND User.username = ?;`
const getUserByUUIDQuery = getUsersQuery + ` AND User.uuid = ?;`
const getUsersPaginatedQuery = getUsersQuery + ` LIMIT ? OFFSET ?;`
const userCreationQuery = `
	INSERT INTO User VALUES (?, ?, ?);
`
const userUpdateQuery = `
	UPDATE User SET username = ?, photoFilename = ? WHERE uuid = ?
`
