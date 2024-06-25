package database

import "github.com/luisfernandogdutra/CursoGO/APIS/internal/entity"

type UserInterface interface {
	Create(user *entity.User) error
	FindByEmail(email string) (*entity.User, error)
}
