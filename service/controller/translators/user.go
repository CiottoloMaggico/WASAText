package translators

import (
	"github.com/ciottolomaggico/wasatext/service/database"
	"github.com/ciottolomaggico/wasatext/service/views"
)

func UserWithImageToView(user database.UserWithImage) views.UserView {
	return views.UserView{
		user.GetUUID(),
		user.GetUsername(),
		ImageToView(user.GetPhoto()),
	}
}

func UserWithImageListToView(users database.UserWithImageList) []views.UserView {
	var res = make([]views.UserView, len(users))
	for _, user := range users {
		res = append(res, UserWithImageToView(user))
	}
	return res
}
