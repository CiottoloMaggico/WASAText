package translators

import (
	"github.com/ciottolomaggico/wasatext/service/models"
	"github.com/ciottolomaggico/wasatext/service/views"
)

func UserWithImageToView(user models.UserWithImage) views.UserWithImageView {
	return views.UserWithImageView{
		user.Uuid,
		user.Username,
		ImageToView(user.Image),
	}
}

func UserWithImageListToView(users []models.UserWithImage) []views.UserWithImageView {
	var res = make([]views.UserWithImageView, 0, cap(users))
	for _, user := range users {
		res = append(res, UserWithImageToView(user))
	}
	return res
}
