package models

import ("time")

//LoginSession Session for logged used with token to access API
type LoginSession struct {
	Model
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type"`
	RefreshToken string    `json:"refresh_token"`
	Expiry       time.Time `json:"expiry"`

	SessionID    string `json:"session_id" gorm:"column:session_id" `
	SessionValue string `json:"session_value" gorm:"column:session_value" `
	UserID       string `json:"user_id" gorm:"column:user_id" `
}

//WriteToDB Write model to DB
func (l *LoginSession) WriteToDB() {
	db := OpenDB()
	db.Create(&l)
}
