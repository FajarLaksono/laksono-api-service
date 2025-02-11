// Copyright (c) 2024 Fajar Laksono. All Rights Reserved.

package model

import (
	"fmt"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username"`
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	Email     string    `json:"email"`
}

type Users []User

func (data *User) ConvertToCreateUserResponse() *CreateUserResponse {
	result := &CreateUserResponse{
		ID:        data.ID,
		Username:  data.Username,
		Firstname: data.Firstname,
		Lastname:  data.Lastname,
		Email:     data.Email,
	}

	return result
}

func (data *User) ConvertToGetUserResponse() *GetUserResponse {
	fmt.Printf("element: %v\n", data)
	result := &GetUserResponse{
		ID:        data.ID,
		Username:  data.Username,
		Firstname: data.Firstname,
		Lastname:  data.Lastname,
		Email:     data.Email,
	}

	return result
}

func (data *Users) ConvertToGetUsersResponse() *GetUsersResponse {
	result := GetUsersResponse{}

	for _, element := range *data {
		result = append(result, *element.ConvertToGetUserResponse())
	}

	return &result
}

type CreateUserResponse struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username"`
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	Email     string    `json:"email"`
}

type CreateUserRequest struct {
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
}

func (data *CreateUserRequest) ConvertToUser(userID uuid.UUID) *User {
	result := &User{
		ID:        userID,
		Username:  data.Username,
		Firstname: data.Firstname,
		Lastname:  data.Lastname,
		Email:     data.Email,
	}

	return result
}

type GetUserResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	Email     string    `json:"email"`
}

type GetUsersResponse []GetUserResponse
