package users

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(input RegisterUserInput) (User, error)
	Login(input LoginUserInput) (User, error)
	GetUserByID(ID int) (User, error)
	ChangePassword(inputID GetUserDetailInput, input ChangePasswordUserInput) (User, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) Register(input RegisterUserInput) (User, error) {
	user := User{}
	user.Username = input.Username

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)

	if err != nil {
		return user, err
	}

	user.Password = string(passwordHash)

	savedUser, err := s.repository.Save(user)

	if err != nil {
		return savedUser, err
	}

	return savedUser, nil
}

func (s *service) Login(input LoginUserInput) (User, error) {
	username := input.Username
	password := input.Password

	user, err := s.repository.FindByUsername(username)
	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("username not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *service) GetUserByID(ID int) (User, error) {
	user, err := s.repository.FindByID(ID)
	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("no user found on with that ID")
	}

	return user, nil
}

func (s *service) ChangePassword(inputID GetUserDetailInput, input ChangePasswordUserInput) (User, error) {
	user, err := s.repository.FindByID(inputID.ID)
	if err != nil {
		return user, err
	}

	if user.ID != input.User.ID {
		return user, errors.New("not an owner of the user")
	}

	if input.NewPassword != input.ConfirmPassword {
		return user, errors.New("password not match")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.MinCost)

	if err != nil {
		return user, err
	}

	user.Password = string(passwordHash)

	updateUser, err := s.repository.Update(user)
	if err != nil {
		return updateUser, err
	}

	return updateUser, nil
}
