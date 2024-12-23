package database

type UserConversation struct {
	Id            int64             `json:"id"`
	Name          string            `json:"name"`
	Type          string            `json:"type"`
	Photo         Image             `json:"photo"`
	LatestMessage MessageWithAuthor `json:"latestMessage"`
	Read          bool              `json:"read"`
}

func (db *appdbimpl) QueryRowUserConversation(query string, params ...any) (*UserConversation, error) {
	userConversation := UserConversation{}
	userConversation.Photo, userConversation.LatestMessage = Image{}, MessageWithAuthor{author: User{}}

	if err := db.c.QueryRow(query, params...).Scan(
		userConversation.Id,
		userConversation.Name,
		userConversation.Type,
		userConversation.Photo.uuid,
		userConversation.Photo.extension,
		userConversation.Photo.width,
		userConversation.Photo.height,
		userConversation.Photo.fullUrl,
		userConversation.Photo.uploadedAt,
		userConversation.LatestMessage.id,
		userConversation.LatestMessage.sendAt,
		userConversation.LatestMessage.deliveredAt,
		userConversation.LatestMessage.seenAt,
		userConversation.LatestMessage.replyTo,
		userConversation.LatestMessage.content,
		userConversation.LatestMessage.attachment,
		userConversation.LatestMessage.author.uuid,
		userConversation.LatestMessage.author.username,
		userConversation.LatestMessage.author.photo,
	); err != nil {
		return nil, err
	}

	return &userConversation, nil
}

func (db *appdbimpl) QueryUserConversation(query string, params ...any) ([]UserConversation, error) {
	rows, err := db.c.Query(query, params...)

	if err != nil {
		return nil, err
	}

	userConversations := make([]UserConversation, 0)
	for rows.Next() {
		userConversation := UserConversation{}
		userConversation.Photo, userConversation.LatestMessage = Image{}, MessageWithAuthor{author: User{}}

		if err := rows.Scan(
			userConversation.Id,
			userConversation.Name,
			userConversation.Type,
			userConversation.Photo.uuid,
			userConversation.Photo.extension,
			userConversation.Photo.width,
			userConversation.Photo.height,
			userConversation.Photo.fullUrl,
			userConversation.Photo.uploadedAt,
			userConversation.LatestMessage.id,
			userConversation.LatestMessage.sendAt,
			userConversation.LatestMessage.deliveredAt,
			userConversation.LatestMessage.seenAt,
			userConversation.LatestMessage.replyTo,
			userConversation.LatestMessage.content,
			userConversation.LatestMessage.attachment,
			userConversation.LatestMessage.author.uuid,
			userConversation.LatestMessage.author.username,
			userConversation.LatestMessage.author.photo,
		); err != nil {
			return nil, err
		}

		userConversations = append(userConversations, userConversation)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userConversations, nil
}
