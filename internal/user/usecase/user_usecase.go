package usecase

import (
	"context"
	"todolist/internal/domain"
	"todolist/internal/exception"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecaseImpl struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(userRepo domain.UserRepository) domain.UserUsecase {
	return &UserUsecaseImpl{
		userRepo: userRepo,
	}
}

func (u *UserUsecaseImpl) Register(ctx context.Context, req domain.UserRegisterRequest) error {
	if req.Password != req.ConfrimPassword {
		return exception.BadRequestError("password and confrim password do not match")
	}
	hasPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return exception.InternalServerError("failed bcrypt passoword")
	}
	mail, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return exception.InternalServerError("database connection error")
	}
	if mail != nil {
		return exception.ConflictError("email is already exists")
	}
	user := domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hasPass),
	}
	err = u.userRepo.Create(ctx, &user)
	if err != nil {
		return exception.InternalServerError("register is failed")
	}
	return nil
}

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
