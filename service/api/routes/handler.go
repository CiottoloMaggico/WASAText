package routes

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type Handler func(w http.ResponseWriter, r *http.Request, params httprouter.Params, context RequestContext)
