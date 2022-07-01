package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	mw "github.com/8run0/otel/backend/pkg/api/middleware"
	"github.com/8run0/otel/backend/pkg/svc"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type userResource struct {
	svc svc.UserService
	ErrorHandler
}

func (us userResource) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", mw.SpanHttp("http:get_users", us.getUsers))
	r.Post("/", mw.SpanHttp("http:post_users", us.createUser))
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", mw.SpanHttp("http:get_user_by_id", us.getUserByID))
	})
	return r
}

func (u userResource) getUsers(w http.ResponseWriter, r *http.Request) {
	users, err := u.svc.GetUsers(r.Context())
	if err != nil {
		u.ErrorHandler.HandleError(err, w, r)

	}
	render.SetContentType(render.ContentTypeJSON)
	render.JSON(w, r, users)

}

func (u userResource) getUserByID(w http.ResponseWriter, r *http.Request) {

	idParam := chi.URLParam(r, "id")
	id, err := u.validateIDParam(idParam)
	if err != nil {
		u.ErrorHandler.HandleError(fmt.Errorf("%s : %w", ErrValidationFailed, err), w, r)
		return
	}
	ctx := r.Context()
	user, err := u.svc.GetUserByID(ctx, id)
	if err != nil {
		u.ErrorHandler.HandleError(fmt.Errorf("%s : %w", ErrNoUserFound, err), w, r)
		return
	}
	render.JSON(w, r, user)

}

var ErrInvalidID = fmt.Errorf("invalid id")

func (u userResource) validateIDParam(idParam string) (id int64, err error) {
	id, err = strconv.ParseInt(idParam, 10, 64)
	if id <= 0 || err != nil {
		return 0, ErrInvalidID
	}
	return id, nil
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (u userResource) createUser(w http.ResponseWriter, r *http.Request) {

	req := &CreateUserRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		u.ErrorHandler.HandleError(err, w, r)
		return
	}
	user, err := u.svc.CreateUser(r.Context(), svc.CreateUserRequest{
		Name:     req.Name,
		Password: req.Password,
	})
	if err != nil {
		u.ErrorHandler.HandleError(err, w, r)
		return
	}
	render.JSON(w, r, user)

}
