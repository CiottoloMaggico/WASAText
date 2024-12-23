package database

type UserWithImageList []UserWithImage

type User struct {
	uuid     string
	username string
	photo    string
}

type UserWithImage struct {
	uuid     string
	username string
	photo    Image
}

func (user User) GetUUID() string {
	return user.uuid
}

func (user User) GetUsername() string {
	return user.username
}

func (user User) GetPhoto() string {
	return user.photo
}

func (user UserWithImage) GetUUID() string {
	return user.uuid
}

func (user UserWithImage) GetUsername() string {
	return user.username
}

func (user UserWithImage) GetPhoto() Image {
	return user.photo
}

func (db *appdbimpl) QueryRowUser(query string, params ...any) (*User, error) {
	user := User{}

	if err := db.c.QueryRow(query, params...).Scan(
		&user.uuid,
		&user.username,
		&user.photo,
	); err != nil {
		return nil, err
	}

	return &user, nil
}

func (db *appdbimpl) QueryUserWithImage(query string, params ...any) (UserWithImageList, error) {
	rows, err := db.c.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer SaveClose(rows)

	res := make(UserWithImageList, 0)
	for rows.Next() {
		user := UserWithImage{}
		user.photo = Image{}

		if err := rows.Scan(
			&user.uuid,
			&user.username,
			&user.photo.uuid,
			&user.photo.extension,
			&user.photo.width,
			&user.photo.height,
			&user.photo.fullUrl,
			&user.photo.uploadedAt,
		); err != nil {
			return nil, err
		}

		res = append(res, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (db *appdbimpl) QueryRowUserWithImage(query string, params ...any) (*UserWithImage, error) {
	user := UserWithImage{}
	user.photo = Image{}

	if err := db.c.QueryRow(query, params...).Scan(
		&user.uuid,
		&user.username,
		&user.photo.uuid,
		&user.photo.extension,
		&user.photo.width,
		&user.photo.height,
		&user.photo.fullUrl,
		&user.photo.uploadedAt,
	); err != nil {
		return nil, err
	}

	return &user, nil
}
