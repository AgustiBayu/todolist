package domain

import "context"

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}

func ToUserResponses(user []User) []UserResponse {
	var responses []UserResponse
	for _, u := range user {
		responses = append(responses, ToUserResponse(u))
	}
	return responses
}

func ToUserResponse(user User) UserResponse {
	return UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
}

type UserResponse struct {
	ID    int    `json:"user_id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRegisterRequest struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfrimPassword string `json:"confrim_password"`
}

type UserUsecase interface {
	Register(ctx context.Context, req UserRegisterRequest) error
	// Login(ctx context.Context, req UserLoginRequest) (UserResponse, error)
	// GetProfileById(ctx context.Context, userID int) (UserResponse, error)
	GetProfile(ctx context.Context) ([]UserResponse, error)
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	ReadById(ctx context.Context, userID int) (*User, error)
	ReadByAll(ctx context.Context) ([]User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, userID int) error
	FindByEmail(ctx context.Context, email string) (*User, error)
}
