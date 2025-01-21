package models

import "time"

type Token struct {
	IdUser 		int 		`json:"id_user"`
	Name 		string 		`json:"name"`
	PhoneNumber string 		`json:"phone_number"`
	Expiry 		time.Time 	`json:"expiry"`
	Hash 		string 		`json:"hash"`
	IdAPI 		int 		`json:"id_api"`
}