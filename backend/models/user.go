package models

type User struct {
	Username string `json:"username"`
	FaceData []byte `json:"faceData"`
}

var Users = make(map[string]User)
