// Package cai
/*
Copyright © 2023 Harmony AI Solutions & Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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
