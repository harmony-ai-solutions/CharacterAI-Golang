package cai

import http "github.com/bogdanfinn/fhttp"

type User struct {
	Token   string
	Session *Session
}

func (u *User) Info() (map[string]interface{}, error) {
	return request("chat/user/", u.Session, u.Token, http.MethodGet, nil, false, false)
}

func (u *User) GetProfile(username string) (map[string]interface{}, error) {
	data := map[string]interface{}{"username": username}
	return request("chat/user/public/", u.Session, u.Token, http.MethodPost, data, false, false)
}

func (u *User) Followers() (map[string]interface{}, error) {
	return request("chat/user/followers/", u.Session, u.Token, http.MethodGet, nil, false, false)
}

func (u *User) Following() (map[string]interface{}, error) {
	return request("chat/user/following/", u.Session, u.Token, http.MethodGet, nil, false, false)
}

func (u *User) Recent() (map[string]interface{}, error) {
	return request("chat/characters/recent/", u.Session, u.Token, http.MethodGet, nil, false, false)
}

func (u *User) Characters() (map[string]interface{}, error) {
	return request("chat/characters/?scope=user", u.Session, u.Token, http.MethodGet, nil, false, false)
}

func (u *User) Update(username string, data map[string]interface{}) (map[string]interface{}, error) {
	data["username"] = username
	return request("chat/user/update/", u.Session, u.Token, http.MethodPost, data, false, false)
}
