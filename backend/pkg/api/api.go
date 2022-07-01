package api

import (
	"fmt"
	"net/http"

	mw "github.com/8run0/otel/backend/pkg/api/middleware"
	"github.com/8run0/otel/backend/pkg/svc"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	ErrCreateUser   = fmt.Errorf("api failed to create user")
	ErrDeleteUser   = fmt.Errorf("api failed to delete user")
	ErrUpdateUser   = fmt.Errorf("api failed to update user")
	ErrNoUserFound  = fmt.Errorf("api failed to get user")
	ErrNoUsersFound = fmt.Errorf("api failed to get users")

	ErrValidationFailed = fmt.Errorf("validation failed")
)

type Server struct {
	router *chi.Mux
}

type ErrorHandler struct {
}

type errorJSON struct {
	Message string `json:"message"`
}

func (eh ErrorHandler) HandleError(err error, w http.ResponseWriter, r *http.Request) {
	errJson := errorJSON{
		Message: err.Error(),
	}
	http.Error(w, errJson.Message, http.StatusInternalServerError)
}

func NewServer() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(mw.WithOTELTools)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})
	r.Mount("/posts", postsResource{}.Routes())
	r.Mount("/users", userResource{
		svc:          svc.NewUserService(),
		ErrorHandler: ErrorHandler{},
	}.Routes())
	return r
}
func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s Server) RunHttp(listenAddr string) error {
	return http.ListenAndServe(listenAddr, s.router)
}
