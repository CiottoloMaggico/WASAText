package database

type Image struct {
	uuid       string
	extension  string
	width      int
	height     int
	fullUrl    string
	uploadedAt string
}

func (im Image) GetUUID() string {
	return im.uuid
}

func (im Image) GetExtension() string {
	return im.extension
}

func (im Image) GetWidthAndHeight() [2]int {
	return [2]int{im.width, im.height}
}

func (im Image) GetFullUrl() string {
	return im.fullUrl
}

func (im Image) GetUploadedAt() string {
	return im.uploadedAt
}

func (im Image) Filename() string {
	return im.uuid + im.extension
}

func (db *appdbimpl) QueryRowImage(query string, params ...any) (*Image, error) {
	image := Image{}
	if err := db.c.QueryRow(query, params...).Scan(
		&image.uuid,
		&image.extension,
		&image.width,
		&image.height,
		&image.fullUrl,
		&image.uploadedAt,
	); err != nil {
		return nil, err
	}

	return &image, nil
}
