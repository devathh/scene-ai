package authuserpg

import "github.com/devathh/scene-ai/internal/modules/auth/domain/user"

func ToDomain(model *UserModel) *user.User {
	return user.From(
		model.ID,
		model.Firstname,
		model.Lastname,
		user.PasswordHash(model.PasswordHash),
		model.CreatedAt,
		model.UpdatedAt,
	)
}

func ToModel(domain *user.User) *UserModel {
	return &UserModel{
		ID:           domain.ID(),
		Firstname:    domain.Firstname(),
		Lastname:     domain.Lastname(),
		PasswordHash: domain.PasswordHash().String(),
		CreatedAt:    domain.CreatedAt(),
		UpdatedAt:    domain.UpdatedAt(),
	}
}
