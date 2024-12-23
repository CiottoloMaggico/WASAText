package translators

import (
	"github.com/ciottolomaggico/wasatext/service/database"
	"github.com/ciottolomaggico/wasatext/service/views"
)

func ImageToView(image database.Image) views.ImageView {
	measures := image.GetWidthAndHeight()
	return views.ImageView{
		image.GetUUID(),
		measures[0],
		measures[1],
		image.GetFullUrl(),
	}
}
