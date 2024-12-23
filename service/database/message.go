package database

import (
	"database/sql"
)

type MessageWithAuthorAndAttachmentList []MessageWithAuthorAndAttachment

type Message struct {
	id           int64
	conversation int64
	author       string
	sendAt       string
	deliveredAt  *string
	seenAt       *string
	replyTo      *int64
	content      *string
	attachment   *string
}

type MessageWithAuthor struct {
	id           int64
	conversation int64
	author       User
	sendAt       string
	deliveredAt  *string
	seenAt       *string
	replyTo      *int64
	content      *string
	attachment   *string
}

type MessageWithAuthorAndAttachment struct {
	id           int64
	conversation int64
	author       UserWithImage
	sendAt       string
	deliveredAt  *string
	seenAt       *string
	replyTo      *int64
	content      *string
	attachment   *Image
}

func (m Message) GetId() int64 {
	return m.id
}

func (m MessageWithAuthor) GetId() int64 {
	return m.id
}

func (maa MessageWithAuthorAndAttachment) GetId() int64 {
	return maa.id
}

func (maa MessageWithAuthorAndAttachment) GetConversationId() int64 {
	return maa.conversation
}

func (maa MessageWithAuthorAndAttachment) GetAuthor() UserWithImage {
	return maa.author
}

func (maa MessageWithAuthorAndAttachment) GetTimestamps() [3]*string {
	return [3]*string{&maa.sendAt, maa.deliveredAt, maa.seenAt}
}

func (maa MessageWithAuthorAndAttachment) GetStatus() string {
	status := `sent`
	if maa.deliveredAt != nil {
		status = `delivered`
	}
	if maa.seenAt != nil {
		status = `seen`
	}
	return status
}

func (maa MessageWithAuthorAndAttachment) GetReplyTo() *int64 {
	return maa.replyTo
}

func (maa MessageWithAuthorAndAttachment) GetContent() *string {
	return maa.content
}

func (maa MessageWithAuthorAndAttachment) GetAttachment() *Image {
	return maa.attachment
}

func (db *appdbimpl) QueryRowMessage(query string, params ...any) (*Message, error) {
	message := Message{}
	if err := db.c.QueryRow(query, params).Scan(
		&message.id,
		&message.conversation,
		&message.author,
		&message.sendAt,
		&message.deliveredAt,
		&message.seenAt,
		&message.replyTo,
		&message.content,
		&message.attachment,
	); err != nil {
		return nil, err
	}

	return &message, nil
}

func (db *appdbimpl) QueryMessageWithAuthorAndAttachment(query string, params ...any) (MessageWithAuthorAndAttachmentList, error) {
	rows, err := db.c.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer SaveClose(rows)

	var tmpAttachmentUUID, tmpAttachmentExt, tmpAttachmentUrl, tmpAttachmentUploadedAt sql.NullString
	var tmpAttachmentWidth, tmpAttachmentHeight sql.NullInt64
	messages := make(MessageWithAuthorAndAttachmentList, 0)
	for rows.Next() {
		message := MessageWithAuthorAndAttachment{}
		message.author = UserWithImage{}

		if err := rows.Scan(
			&message.id,
			&message.conversation,
			&message.author.uuid,
			&message.author.username,
			&message.author.photo.uuid,
			&message.author.photo.extension,
			&message.author.photo.uploadedAt,
			&message.sendAt,
			&message.deliveredAt,
			&message.seenAt,
			&message.replyTo,
			&message.content,
			&tmpAttachmentUUID,
			&tmpAttachmentExt,
			&tmpAttachmentWidth,
			&tmpAttachmentHeight,
			&tmpAttachmentUrl,
			&tmpAttachmentUploadedAt,
		); err != nil {
			return nil, err
		}

		if tmpAttachmentUUID.Valid {
			message.attachment = &Image{
				tmpAttachmentUUID.String,
				tmpAttachmentExt.String,
				int(tmpAttachmentWidth.Int64),
				int(tmpAttachmentHeight.Int64),
				tmpAttachmentUrl.String,
				tmpAttachmentUploadedAt.String,
			}
		}

		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return messages, err
	}

	return messages, nil
}

func (db *appdbimpl) QueryRowMessageWithAuthorAndAttachment(query string, params ...any) (*MessageWithAuthorAndAttachment, error) {
	var tmpAttachmentUUID, tmpAttachmentExt, tmpAttachmentUrl, tmpAttachmentUploadedAt sql.NullString
	var tmpAttachmentWidth, tmpAttachmentHeight sql.NullInt64
	message := MessageWithAuthorAndAttachment{}
	message.author = UserWithImage{}

	if err := db.c.QueryRow(query, params...).Scan(
		&message.id,
		&message.conversation,
		&message.author.uuid,
		&message.author.username,
		&message.author.photo.uuid,
		&message.author.photo.extension,
		&message.author.photo.uploadedAt,
		&message.sendAt,
		&message.deliveredAt,
		&message.seenAt,
		&message.replyTo,
		&message.content,
		&tmpAttachmentUUID,
		&tmpAttachmentExt,
		&tmpAttachmentWidth,
		&tmpAttachmentHeight,
		&tmpAttachmentUrl,
		&tmpAttachmentUploadedAt,
	); err != nil {
		return nil, err
	}

	if tmpAttachmentUUID.Valid {
		message.attachment = &Image{
			tmpAttachmentUUID.String,
			tmpAttachmentExt.String,
			int(tmpAttachmentWidth.Int64),
			int(tmpAttachmentHeight.Int64),
			tmpAttachmentUrl.String,
			tmpAttachmentUploadedAt.String,
		}
	}

	return &message, nil
}
