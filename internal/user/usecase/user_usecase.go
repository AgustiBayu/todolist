package usecase

import (
	"context"
	"todolist/internal/domain"
	"todolist/internal/exception"
)

type UserUsecaseImpl struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(userRepo domain.UserRepository) domain.UserUsecase {
	return &UserUsecaseImpl{
		userRepo: userRepo,
	}
}

// func (u *UserUsecaseImpl) Register(ctx context.Context, req domain.UserRegisterRequest) error {

// }
// func (u *UserUsecaseImpl) Login(ctx context.Context, req domain.UserLoginRequest) (domain.UserResponse, error) {

// }
// func (u *UserUsecaseImpl) GetProfileById(ctx context.Context, userID int) (domain.UserResponse, error) {
// 	user, err := u.userRepo.ReadById(ctx, userID)
// 	if err != nil {
// 		return domain.UserResponse{}, exception.InternalServerError("failed fetch data users: " + err.Error())
// 	}
// 	if user == nil {
// 		return domain.UserResponse{}, exception.NotFound("user not found")
// 	}
// 	return domain.ToUserResponse(*user), nil
// }

func (u *UserUsecaseImpl) GetProfile(ctx context.Context) ([]domain.UserResponse, error) {
	users, err := u.userRepo.ReadByAll(ctx)
	if err != nil {
		return nil, exception.InternalServerError("failed fetch data users: " + err.Error())
	}
	if len(users) == 0 {
		return nil, exception.NotFoundError("no user registered yet")
	}
	return domain.ToUserResponses(users), nil
}
