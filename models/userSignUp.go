package models

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	Id       pgtype.Int4 `json:"id"`
	Name     pgtype.Text `json:"name"`
	Email    pgtype.Text `json:"email"`
	Password pgtype.Text `json:"password"`
}

type LoginClaims struct {
	Name     pgtype.Text `json:"name"`
	Email    pgtype.Text `json:"email"`
	Password pgtype.Text `json:"password"`
	Role     pgtype.Text `json:"role"`
	jwt.RegisteredClaims
}

// type UserLogin struct {
// 	Name     string `json:"name"`
// 	Email    string `json:"email"`
// 	Role 	 string `json:"role"`
// 	Password string `json:"password"`
// }

//DTO

type UserResponse struct {
	Id    int32  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		Id:    u.Id.Int32,
		Name:  u.Name.String,
		Email: u.Email.String,
		Role:  u.Email.String,
	}
}
