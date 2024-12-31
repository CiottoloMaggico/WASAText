package translators

import (
	"github.com/ciottolomaggico/wasatext/service/models"
	"github.com/ciottolomaggico/wasatext/service/views"
)

func UserWithImageToView(user models.UserWithImage) views.UserView {
	return views.UserView{
		user.Uuid,
		user.Username,
		ImageToView(user.Image),
	}
}

func UserWithImageListToView(users []models.UserWithImage) []views.UserView {
	var res = make([]views.UserView, 0, cap(users))
	for _, user := range users {
		res = append(res, UserWithImageToView(user))
	}
	return res
}
