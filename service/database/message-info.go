package database

type MessageInfoList []MessageInfo

type MessageInfo struct {
	message int64
	user    string
	status  int
	comment *string
}

func (i MessageInfo) GetMessageID() int64 {
	return i.message
}

func (i MessageInfo) GetUserUUID() string {
	return i.user
}

func (i MessageInfo) GetStatus() int {
	return i.status
}

func (i MessageInfo) GetComment() *string {
	return i.comment
}

func (db *appdbimpl) QueryMessageInfo(query string, params ...any) (MessageInfoList, error) {
	rows, err := db.c.Query(query, params...)

	if err != nil {
		return nil, err
	}
	defer SaveClose(rows)

	messageInfos := make(MessageInfoList, 0)
	for rows.Next() {
		messageInfo := MessageInfo{}

		if err := rows.Scan(
			&messageInfo.message,
			&messageInfo.user,
			&messageInfo.status,
			&messageInfo.comment,
		); err != nil {
			return nil, err
		}
		messageInfos = append(messageInfos, messageInfo)
	}

	if err := rows.Err(); err != nil {
		return messageInfos, err
	}
	return messageInfos, nil
}

func (db *appdbimpl) QueryRowMessageInfo(query string, params ...any) (*MessageInfo, error) {
	messageInfo := MessageInfo{}
	if err := db.c.QueryRow(query, params...).Scan(
		&messageInfo.message,
		&messageInfo.user,
		&messageInfo.status,
		&messageInfo.comment,
	); err != nil {
		return nil, err
	}

	return &messageInfo, nil
}
