package models

import "time"

type Token struct {
	IdUser 		int64 		`json:"id_user"` //чтобы не выводить его в json
	Name 		string 		`json:"name"`
	PhoneNumber string 		`json:"phone_number"`
	Expiry 		time.Time 	`json:"expiry"`
	Hash 		string 		`json:"hash"`
	IdAPI 		int64 		`json:"id_api"`
}