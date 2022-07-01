package svc

import (
	"context"

	"github.com/8run0/otel/backend/pkg/otel"
)

var _ userServiceImpl = &userServiceSpanner{}

type userServiceSpanner struct {
	next userServiceImpl
}

func (s *userServiceSpanner) GetUsers(ctx context.Context) (users []User, err error) {
	tools := otel.GetToolsFromContext(ctx)
	ctx, span := tools.Tracer.Start(ctx, "userService_GetUsers")
	defer span.End()
	return s.next.GetUsers(ctx)
}

func (s *userServiceSpanner) GetUserByID(ctx context.Context, id int64) (user User, err error) {
	tools := otel.GetToolsFromContext(ctx)
	ctx, span := tools.Tracer.Start(ctx, "userService_GetUserByID")
	defer span.End()
	return s.next.GetUserByID(ctx, id)
}

func (s *userServiceSpanner) CreateUser(ctx context.Context, req CreateUserRequest) (id int64, err error) {
	tools := otel.GetToolsFromContext(ctx)
	ctx, span := tools.Tracer.Start(ctx, "userService_CreateUser")
	defer span.End()
	return s.next.CreateUser(ctx, req)
}

func (s *userServiceSpanner) DeleteUser(ctx context.Context, req DeleteUserRequest) (err error) {
	tools := otel.GetToolsFromContext(ctx)
	ctx, span := tools.Tracer.Start(ctx, "userService_DeleteUser")
	defer span.End()
	return s.next.DeleteUser(ctx, req)
}

func (s *userServiceSpanner) UpdateUser(ctx context.Context, req UpdateUserRequest) (err error) {
	tools := otel.GetToolsFromContext(ctx)
	ctx, span := tools.Tracer.Start(ctx, "userService_UpdateUser")
	defer span.End()
	return s.next.UpdateUser(ctx, req)
}
