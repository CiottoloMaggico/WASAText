package database

type Conversation struct {
	id int64
}

type Chat struct {
	Conversation
	user1 string
	user2 string
}

type Group struct {
	Conversation
	name   string
	author string
	photo  string
}

const MAX_GROUP_SIZE = 200

func (c Conversation) GetId() int64 {
	return c.id
}

func (c Chat) GetParticipantsUUID() [2]string {
	return [2]string{c.user1, c.user2}
}

func (c Group) GetName() string {
	return c.name
}

func (c Group) GetAuthor() string {
	return c.author
}

func (c Group) GetPhotoUUID() string {
	return c.photo
}

func (db *appdbimpl) QueryRowGroup(query string, params ...any) (*Group, error) {
	group := Group{}
	if err := db.c.QueryRow(
		query,
		params...,
	).Scan(
		&group.id,
		&group.name,
		&group.author,
		&group.photo,
	); err != nil {
		return nil, err
	}

	return &group, nil
}

func (db *appdbimpl) QueryRowChat(query string, params ...any) (*Chat, error) {
	chat := Chat{}
	if err := db.c.QueryRow(query, params...).Scan(
		&chat.id,
		&chat.user1,
		&chat.user2,
	); err != nil {
		return nil, err
	}

	return &chat, nil
}
