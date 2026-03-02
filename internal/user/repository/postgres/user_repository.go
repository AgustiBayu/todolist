package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"todolist/internal/domain"

	"github.com/redis/go-redis/v9"
)

type UserRepositoryImpl struct {
	DB    *sql.DB
	Redis *redis.Client
}

func NewUserRepository(DB *sql.DB, redis *redis.Client) domain.UserRepository {
	return &UserRepositoryImpl{
		DB:    DB,
		Redis: redis,
	}
}

func (r *UserRepositoryImpl) Create(ctx context.Context, user *domain.User) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var lastInsertId int
	SQL := `INSERT INTO users(name,email,password,created_at,updated_at)VALUES($1,$2,$3,now(),now()) RETURNING id`
	if err := tx.QueryRowContext(ctx, SQL, user.Name, user.Email, user.Password).Scan(&lastInsertId); err != nil {
		return err
	}
	user.ID = lastInsertId
	if err := tx.Commit(); err != nil {
		return err
	}
	cacheKey := fmt.Sprint("users::all")
	if errRedis := r.Redis.Del(ctx, cacheKey).Err(); errRedis != nil {
		fmt.Printf("failed deleted cache %s: %v/n", cacheKey, errRedis)
	}
	return nil
}

func (r *UserRepositoryImpl) ReadById(ctx context.Context, userID int) (*domain.User, error) {
	cacheKey := fmt.Sprintf("users::%d", userID)
	val, err := r.Redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var usr domain.User
		if err := json.Unmarshal([]byte(val), &usr); err == nil {
			return &usr, err
		}
	}
	SQL := `SELECT id,name,email,created_at,updated_at FROM users WHERE id = $1`
	var model UserModel
	if err := r.DB.QueryRowContext(ctx, SQL, userID).Scan(&model.ID, &model.Name, &model.Email, &model.CreatedAt, &model.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	userModel := model.ToUserDomain()
	jsonData, errModel := json.Marshal(userModel)
	if errModel != nil {
		return nil, errModel
	}
	r.Redis.Set(ctx, cacheKey, jsonData, 15*time.Minute)
	return userModel, nil
}
func (r *UserRepositoryImpl) ReadByAll(ctx context.Context) ([]domain.User, error) {
	cacheKey := fmt.Sprint("users::all")
	val, err := r.Redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var users []domain.User
		if err := json.Unmarshal([]byte(val), &users); err == nil {
			return users, nil
		}
	}
	SQL := `SELECT id,name,email,created_at,updated_at FROM users`
	rows, err := r.DB.QueryContext(ctx, SQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []domain.User
	for rows.Next() {
		var model UserModel
		if err := rows.Scan(&model.ID, &model.Name, &model.Email, &model.CreatedAt, &model.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, *model.ToUserDomain())
	}
	jsonData, errModel := json.Marshal(users)
	if errModel != nil {
		return nil, errModel
	}
	r.Redis.Set(ctx, cacheKey, jsonData, 15*time.Minute)
	return users, nil
}

func (r *UserRepositoryImpl) Update(ctx context.Context, user *domain.User) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	SQL := `UPDATE users SET name=$1, email=$2, updated_at = now() WHERE id = $3`
	result, err := tx.ExecContext(ctx, SQL, user.Name, user.Email, user.ID)
	if err != nil {
		return err
	}
	row, _ := result.RowsAffected()
	if row == 0 {
		return errors.New("data is not exist")
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	cacheKeyDetail := fmt.Sprintf("user:%d", user.ID)
	cacheList := fmt.Sprint("user::all")
	if errRedis := r.Redis.Del(ctx, cacheKeyDetail, cacheList).Err(); errRedis != nil {
		fmt.Printf("failed delete cache redis %s", errRedis)
	}
	return nil
}
func (r *UserRepositoryImpl) Delete(ctx context.Context, userID int) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	SQL := `DELETE users WHERE id = $1`
	result, err := tx.ExecContext(ctx, SQL, userID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("user not found")
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	cacheKeyDetail := fmt.Sprintf("user:%d", userID)
	cacheList := "user::all"

	if errRedis := r.Redis.Del(ctx, cacheKeyDetail, cacheList).Err(); errRedis != nil {
		fmt.Printf("failed delete cache redis %s", errRedis)
	}
	return nil
}
func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	cacheKey := fmt.Sprintf("email:%s", email)
	val, err := r.Redis.Get(ctx, cacheKey).Result()
	if err == nil {
		user := &domain.User{}
		if err := json.Unmarshal([]byte(val), user); err != nil {
			return user, nil
		}
	}
	SQL := `SELECT id, name, email,password FROM users WHERE email = $1 LIMIT 1`
	var model UserModel
	if err := r.DB.QueryRowContext(ctx, SQL, email).Scan(&model.ID, &model.Name, &model.Email, &model.Password); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	userDomain := model.ToUserDomain()
	jsonData, _ := json.Marshal(userDomain)
	r.Redis.Set(ctx, cacheKey, jsonData, 15*time.Minute)
	return userDomain, nil
}
