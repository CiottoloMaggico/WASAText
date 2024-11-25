package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/ciottolomaggico/wasatext/service/utils/validators"
	"github.com/google/uuid"
)

type User struct {
	Uuid         string `json:"uuid"`
	Username     string `json:"username"`
	ProfileImage *Image `json:"image"`
}

func NewUser(username string, photo *Image) (*User, error) {
	rawUUID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to generate uuid: %w", err)
	}
	if ok, err := validators.UsernameIsValid(username); !ok {
		return nil, err
	}
	if photo == nil {
		photo = &Image{
			"default_user_image.jpg",
			1470000,
			sql.NullString{"", false},
			8000,
			8000,
		}
	}

	user := &User{rawUUID.String(), username, photo}
	return user, nil
}

func (u *User) Exists(appDB AppDatabase) (bool, error) {
	var tmpUUID string
	db := appDB.(*appdbimpl).c
	err := db.QueryRow(getUserUUIDQuery, u.Uuid).Scan(&tmpUUID)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("failed to check if user exists: %w", err)
	}
	return true, nil

}

func (u *User) Save(appDB AppDatabase) error {
	db := appDB.(*appdbimpl).c
	userExists, err := u.Exists(appDB)
	if err != nil {
		return err
	}
	// Create a new user row
	if !userExists {
		if _, err := db.Exec(userCreationQuery, u.Uuid, u.Username, u.ProfileImage.Filename); err != nil {
			return fmt.Errorf("failed to create the user: %w", err)
		}

		return nil
	}

	// Update the existing user row with new values
	if _, err := db.Exec(userUpdateQuery, u.Username, u.ProfileImage.Filename, u.Uuid); err != nil {
		return err
	}
	return nil
}

func (u *User) Delete(db AppDatabase) error {
	return errors.New("not implemented")
}

func (db *appdbimpl) GetUser(username string) (*User, error) {
	if ok, err := validators.UsernameIsValid(username); !ok {
		return nil, fmt.Errorf("invalid username: %w", err)
	}
	user := User{}
	user.ProfileImage = &Image{}
	image := user.ProfileImage

	if err := db.c.QueryRow(getUserQuery, username).Scan(
		&user.Uuid,
		&user.Username,
		&image.Filename,
		&image.Size,
		&image.Owner,
		&image.Width,
		&image.Height,
	); err != nil {
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	return &user, nil
}

func (db *appdbimpl) GetUserByUUID(UUID string) (*User, error) {
	if _, err := uuid.Parse(UUID); err != nil {
		return nil, fmt.Errorf("invalid uuid: %w", err)
	}
	user := User{}
	user.ProfileImage = &Image{}
	image := user.ProfileImage

	if err := db.c.QueryRow(getUserByUUIDQuery, UUID).Scan(
		&user.Uuid,
		&user.Username,
		&image.Filename,
		&image.Size,
		&image.Owner,
		&image.Width,
		&image.Height,
	); err != nil {
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	return &user, nil
}

func (db *appdbimpl) GetUsers(pageSize int, pageNumber int) ([]User, error) {
	rows, err := db.c.Query(
		getUsersPaginatedQuery,
		pageSize, pageNumber*pageSize,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]User, 0, pageSize)
	for rows.Next() {
		user := User{}
		user.ProfileImage = &Image{}
		image := user.ProfileImage

		if err := rows.Scan(
			&user.Uuid,
			&user.Username,
			&image.Filename,
			&image.Size,
			&image.Owner,
			&image.Width,
			&image.Height,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (db *appdbimpl) UsersCount() (int, error) {
	var count int
	err := db.c.QueryRow(
		getUsersCountQuery,
	).Scan(
		&count,
	)
	return count, err
}
