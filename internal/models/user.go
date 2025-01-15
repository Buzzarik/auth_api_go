package models

import "time"

type User struct {
	ID          	int64     	`json:"id"`
	CreatedAt   	time.Time 	`json:"created_at"`
	Name        	string    	`json:"name"`
	PhoneNumber 	string    	`json:"phone_number"`
	HashPassword 	string 		`json:"hash_password"`
};