package translators

import (
	"github.com/ciottolomaggico/wasatext/service/models"
	"github.com/ciottolomaggico/wasatext/service/views"
)

func ImageToView(image models.Image) views.ImageView {
	return views.ImageView{
		image.Uuid,
		image.Width,
		image.Height,
		image.FullUrl,
	}
}
