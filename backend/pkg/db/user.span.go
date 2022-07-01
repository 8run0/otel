package db

import (
	"context"

	"github.com/8run0/otel/backend/pkg/otel"
)

var _ userDatabaseImpl = &userDatabaseSpanner{}

type userDatabaseSpanner struct {
	next userDatabaseImpl
}

func (s userDatabaseSpanner) GetUsers(ctx context.Context) (users []User, err error) {
	tools := otel.GetToolsFromContext(ctx)
	ctx, span := tools.Tracer.Start(ctx, "userDatabase_GetUsers")
	defer span.End()
	return s.next.GetUsers(ctx)
}

func (s userDatabaseSpanner) GetUserByID(ctx context.Context, id int64) (user User, err error) {
	tools := otel.GetToolsFromContext(ctx)
	ctx, span := tools.Tracer.Start(ctx, "userDatabase_GetUserByID")
	defer span.End()
	return s.next.GetUserByID(ctx, id)
}

func (s userDatabaseSpanner) CreateUser(ctx context.Context, req *CreateUserRequest) (id int64, err error) {
	tools := otel.GetToolsFromContext(ctx)
	ctx, span := tools.Tracer.Start(ctx, "userDatabase_CreateUser")
	defer span.End()
	return s.next.CreateUser(ctx, req)
}

func (s userDatabaseSpanner) DeleteUser(ctx context.Context, req *DeleteUserRequest) (err error) {
	tools := otel.GetToolsFromContext(ctx)
	ctx, span := tools.Tracer.Start(ctx, "userDatabase_DeleteUser")
	defer span.End()
	return s.next.DeleteUser(ctx, req)
}

func (s userDatabaseSpanner) UpdateUser(ctx context.Context, req *UpdateUserRequest) (err error) {
	tools := otel.GetToolsFromContext(ctx)
	ctx, span := tools.Tracer.Start(ctx, "userDatabase_UpdateUser")
	defer span.End()
	return s.next.UpdateUser(ctx, req)
}
