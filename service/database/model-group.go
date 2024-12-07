package database

import (
	"errors"
	"fmt"
	"github.com/ciottolomaggico/wasatext/service/utils/validators"
	"github.com/google/uuid"
)

const MAX_GROUP_SIZE = 200

var InvalidId = errors.New("invalid id")
var NotMember = errors.New("this user isn't a group member")
var ExceedMaxGroupSize = fmt.Errorf("groups can handle at most %i participants", MAX_GROUP_SIZE)

// TODO: update request issuer with context
type GroupConversation struct {
	Id    uint64 `json:"id"`
	Name  string `json:"name"`
	Photo Image  `json:"photo"`

	requestIssuer User
	db            *appdbimpl
}

type UpdateGroupParams struct {
	Name  string
	Photo Image
}

func (gp UpdateGroupParams) Validate() error {
	if ok, err := validators.GroupNameIsValid(gp.Name); !ok {
		return err
	}
	if err := gp.Photo.Validate(); err != nil {
		return err
	}
	return nil
}

func (g GroupConversation) GetId() uint64 {
	return g.Id
}

func (g GroupConversation) GetName() string {
	return g.Name
}

func (g GroupConversation) GetPhoto() Image {
	return g.Photo
}

func (g GroupConversation) GetDB() *appdbimpl {
	return g.db
}

func (g GroupConversation) GetRequestIssuer() User {
	return g.requestIssuer
}

func (g GroupConversation) Type() string {
	return "group"
}

func (g GroupConversation) Validate() error {
	if ok, err := validators.GroupNameIsValid(g.Name); !ok {
		return err
	}
	if err := g.Photo.Validate(); err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) NewGroup(name string, photo Image, author User) (*GroupConversation, error) {
	// TODO: remove the feature to add participants in the same request of group creation, rather handle it with different requests in frontend
	if ok, err := validators.GroupNameIsValid(name); !ok {
		return nil, err
	}
	if err := photo.Validate(); err != nil {
		return nil, err
	}
	if err := author.Validate(); err != nil {
		return nil, err
	}

	tx, err := db.c.Begin()
	if err != nil {
		return nil, err
	}

	// Create the conversation
	var conversationId uint64
	if err := tx.QueryRow(qCreateConversation).Scan(&conversationId); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create a new group with the same id of the conversation
	if _, err := tx.Exec(qCreateGroup, conversationId, name, author.Uuid, photo.Uuid); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	group := &GroupConversation{
		conversationId, name, photo, author, db,
	}

	return group, nil
}

func (g *GroupConversation) Update(params UpdateGroupParams) error {
	if err := g.Validate(); err != nil {
		return err
	}
	if err := params.Validate(); err != nil {
		return err
	}

	if _, err := g.db.c.Exec(qUpdateGroup, params.Name, params.Photo.Uuid, g.Id); err != nil {
		return err
	}

	g.Name, g.Photo = params.Name, params.Photo
	return nil
}

func (g GroupConversation) AddParticipant(userUUID string) error {
	if _, err := uuid.Parse(userUUID); err != nil {
		return fmt.Errorf("invalid uuid: %w", err)
	}

	if _, err := g.db.c.Exec(qAddGroupToConversation, userUUID, g.Id); err != nil {
		return err
	}
	return nil
}

func (g GroupConversation) AddParticipants(participants []string) error {
	tx, err := g.db.c.Begin()
	if err != nil {
		return err
	}

	for _, participant := range participants {
		if err := g.AddParticipant(participant); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (g GroupConversation) RemoveParticipant(userUUID string) error {
	if ok := g.IsParticipant(userUUID); !ok {
		return NotMember
	}

	if _, err := g.db.c.Exec(qRemoveParticipant, userUUID, g.Id); err != nil {
		return err
	}
	return nil
}

func (g GroupConversation) IsParticipant(userUUID string) bool {
	if _, err := uuid.Parse(userUUID); err != nil {
		return false
	}
	var tmpUuid string
	if err := g.db.c.QueryRow(qGetParticipant, userUUID, g.Id).Scan(&tmpUuid); err != nil {
		return false
	}
	return true
}

func (db *appdbimpl) GetGroup(id uint64) (*GroupConversation, error) {
	if id < 0 {
		return nil, InvalidId
	}

	group := GroupConversation{}
	group.Photo = Image{}
	image := &group.Photo

	if err := db.c.QueryRow(qGetGroup, id).Scan(&group.Id, &group.Name, &image.Uuid, &image.Extension); err != nil {
		return nil, err
	}

	return &group, nil

}
