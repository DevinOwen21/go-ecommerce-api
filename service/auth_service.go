package service

import (
	"database/sql"
	"errors"
	"go-ecommerce-api/dto"
	"go-ecommerce-api/model"
	"go-ecommerce-api/repository"
	"go-ecommerce-api/utils"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo}
}

func (s *AuthService) Register(req dto.RegisterRequest) (*dto.UserResponse, error) {
	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		return nil, err
	}

	_, err = s.userRepo.GetUserByEmail(req.Email)
	if err == nil {
		return nil, errors.New("User already exists")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}

	user := model.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	err = s.userRepo.CreateUser(&user)
	if err != nil {
		return nil, err
	}

	response := dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
	return &response, nil
}

func (s *AuthService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("Invalid email or password")
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.Password),
	)
	if err != nil {
		return nil, errors.New("Invalid email or password")
	}

	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	response := dto.LoginResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Token: token,
	}
	return &response, nil
}
