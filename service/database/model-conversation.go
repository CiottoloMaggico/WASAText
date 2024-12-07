package database

import (
	"database/sql"
	"encoding/json"
	"errors"
)

type SerializedConversationSummary struct {
	Id            uint64   `json:"id"`
	Name          string   `json:"name"`
	Photo         Image    `json:"photo"`
	ChatType      string   `json:"chatType"`
	LatestMessage *Message `json:"latestMessage"`
	Seen          bool     `json:"read"`
}

type SerializedConversation struct {
	Id           uint64   `json:"id"`
	Name         string   `json:"name"`
	Photo        Image    `json:"photo"`
	ChatType     string   `json:"chatType"`
	Participants []string `json:"participants"`
}

// TODO: change request issuer with context
type Conversation interface {
	GetId() uint64
	GetName() string
	GetPhoto() Image
	GetRequestIssuer() User
	GetDB() *appdbimpl
	Type() string
	Validate() error
}

type DefaultConversation struct {
	Conversation
}

func (dc *DefaultConversation) MarshalJSON() ([]byte, error) {
	latestMessage, read, err := dc.GetLatestMessage()
	if err != nil {
		return nil, err
	}

	return json.Marshal(&SerializedConversationSummary{
		dc.GetId(),
		dc.GetName(),
		dc.GetPhoto(),
		dc.Type(),
		latestMessage,
		read,
	})
}

func (dc *DefaultConversation) MarshalDetailedJSON() ([]byte, error) {
	participants, err := dc.GetParticipants()
	if err != nil {
		return nil, err
	}
	participantsUUIDList := make([]string, 0, 200)
	for _, participant := range participants {
		participantsUUIDList = append(participantsUUIDList, participant.Uuid)
	}
	return json.Marshal(&SerializedConversation{
		dc.GetId(),
		dc.GetName(),
		dc.GetPhoto(),
		dc.Type(),
		participantsUUIDList,
	})
}

func (dc *DefaultConversation) GetParticipants() ([]User, error) {
	rows, err := dc.GetDB().c.Query(
		qGetConversationParticipants, dc.GetId(),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]User, 0, 100)
	for rows.Next() {
		user := User{}
		user.ProfileImage = Image{}
		image := &user.ProfileImage

		if err := rows.Scan(
			&user.Uuid,
			&user.Username,
			&image.Uuid,
			&image.Extension,
		); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (dc *DefaultConversation) SendMessage(author User, replyTo *uint64, attachment *Image, content *string) (*Message, error) {
	return dc.GetDB().NewMessage(dc, author, replyTo, attachment, content)
}

func (dc *DefaultConversation) GetMessages(pageSize int, pageNumber int) ([]Message, error) {
	rows, err := dc.GetDB().c.Query(
		qGetConversationMessagesPaginated,
		dc.GetId(), dc.GetRequestIssuer().Uuid, pageSize, pageNumber*pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]Message, 0, pageSize)
	for rows.Next() {
		var message Message
		message.Author = User{ProfileImage: Image{}}
		var attachmentUUID, attachmentExt sql.NullString
		if err := rows.Scan(
			&message.Id,
			&message.SendAt,
			&message.DeliveredAt,
			&message.SeenAt,
			&message.ReplyTo,
			&message.Content,
			&attachmentUUID,
			&attachmentExt,
			&message.Author.Uuid,
			&message.Author.Username,
			&message.Author.ProfileImage.Uuid,
			&message.Author.ProfileImage.Extension,
		); err != nil {
			return nil, err
		}

		if attachmentUUID.Valid {
			message.Attachment = &Image{attachmentUUID.String, attachmentExt.String, dc.GetDB()}
		}

		message.Conv, message.db = dc, dc.GetDB()
		messages = append(messages, message)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); errors.Is(err, sql.ErrNoRows) {
		return messages, nil
	} else if err != nil {
		return nil, err
	}

	return messages, nil
}

func (dc *DefaultConversation) GetMessage(messageId uint64) (*Message, error) {
	var message Message
	message.Author = User{ProfileImage: Image{}}
	var attachmentUUID, attachmentExt sql.NullString
	if err := dc.GetDB().c.QueryRow(qGetConversationMessageById, messageId, dc.GetId()).Scan(
		&message.Id,
		&message.SendAt,
		&message.DeliveredAt,
		&message.SeenAt,
		&message.ReplyTo,
		&message.Content,
		&attachmentUUID,
		&attachmentExt,
		&message.Author.Uuid,
		&message.Author.Username,
		&message.Author.ProfileImage.Uuid,
		&message.Author.ProfileImage.Extension,
	); err != nil {
		return nil, err
	}

	if attachmentUUID.Valid {
		message.Attachment = &Image{attachmentUUID.String, attachmentExt.String, dc.GetDB()}
	}

	message.Conv, message.db = dc, dc.GetDB()
	return &message, nil
}

func (dc *DefaultConversation) GetLatestMessage() (*Message, bool, error) {
	var tmpAttachmentUUID, tmpAttachmentExt sql.NullString
	var seen bool
	message := Message{Author: User{ProfileImage: Image{}}, Conv: Conversation(dc)}
	author := &message.Author

	if err := dc.GetDB().c.QueryRow(qGetLatestMessage, dc.GetId(), dc.GetRequestIssuer().Uuid).Scan(
		&message.Id,
		&message.SendAt,
		&message.DeliveredAt,
		&message.SeenAt,
		&message.ReplyTo,
		&message.Content,
		&tmpAttachmentUUID,
		&tmpAttachmentExt,
		&author.Uuid,
		&author.Username,
		&author.ProfileImage.Uuid,
		&author.ProfileImage.Extension,
		&seen,
	); errors.Is(err, sql.ErrNoRows) {
		return nil, true, nil
	} else if err != nil {
		return nil, false, err
	}

	if tmpAttachmentUUID.Valid {
		message.Attachment = &Image{tmpAttachmentUUID.String, tmpAttachmentExt.String, dc.GetDB()}
	}

	return &message, seen, nil
}
