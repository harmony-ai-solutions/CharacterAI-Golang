package cai

import http "github.com/bogdanfinn/fhttp"

type User struct {
	Token   string
	Session *Session
}

func (u *User) Info(token string) (map[string]interface{}, error) {
	return request("chat/user/", u.Session, token, http.MethodGet, nil, false, false)
}

func (u *User) GetProfile(username string, token string) (map[string]interface{}, error) {
	data := map[string]interface{}{"username": username}
	return request("chat/user/public/", u.Session, token, http.MethodPost, data, false, false)
}

func (u *User) Followers(token string) (map[string]interface{}, error) {
	return request("chat/user/followers/", u.Session, token, http.MethodGet, nil, false, false)
}

func (u *User) Following(token string) (map[string]interface{}, error) {
	return request("chat/user/following/", u.Session, token, http.MethodGet, nil, false, false)
}

func (u *User) Recent(token string) (map[string]interface{}, error) {
	return request("chat/characters/recent/", u.Session, token, http.MethodGet, nil, false, false)
}

func (u *User) Characters(token string) (map[string]interface{}, error) {
	return request("chat/characters/?scope=user", u.Session, token, http.MethodGet, nil, false, false)
}

func (u *User) Update(username string, token string, data map[string]interface{}) (map[string]interface{}, error) {
	data["username"] = username
	return request("chat/user/update/", u.Session, token, http.MethodPost, data, false, false)
}
