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

import (
	"errors"
	"fmt"
	http "github.com/bogdanfinn/fhttp"
	"strconv"
)

type Post struct {
	Token   string
	Session *Session
}

func (p *Post) GetPost(postID string) (map[string]interface{}, error) {
	url := fmt.Sprintf("chat/post/?post=%s", postID)
	return request(url, p.Session, "", http.MethodGet, nil, false, false)
}

func (p *Post) My(postsPage int, postsToLoad int) (map[string]interface{}, error) {
	url := fmt.Sprintf("chat/posts/user/?scope=user&page=%d&posts_to_load=%d", postsPage, postsToLoad)
	return request(url, p.Session, p.Token, http.MethodGet, nil, false, false)
}

func (p *Post) GetPosts(username string, postsPage int, postsToLoad int) (map[string]interface{}, error) {
	url := fmt.Sprintf("chat/posts/user/?username=%s&page=%d&posts_to_load=%d", username, postsPage, postsToLoad)
	return request(url, p.Session, "", http.MethodGet, nil, false, false)
}

func (p *Post) Upvote(postExternalID string) (map[string]interface{}, error) {
	data := map[string]interface{}{"post_external_id": postExternalID}
	return request("chat/post/upvote/", p.Session, p.Token, http.MethodPost, data, false, false)
}

func (p *Post) UndoUpvote(postExternalID string) (map[string]interface{}, error) {
	data := map[string]interface{}{"post_external_id": postExternalID}
	return request("chat/post/undo-upvote/", p.Session, p.Token, http.MethodPost, data, false, false)
}

func (p *Post) SendComment(postID, text, parentUUID string) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"post_external_id": postID,
		"text":             text,
		"parent_uuid":      parentUUID,
	}
	return request("chat/comment/create/", p.Session, p.Token, http.MethodPost, data, false, false)
}

func (p *Post) DeleteComment(messageID int, postID string) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"external_id":      strconv.Itoa(messageID),
		"post_external_id": postID,
	}
	return request("chat/comment/delete/", p.Session, p.Token, http.MethodPost, data, false, false)
}

func (p *Post) Create(postType, externalID, title, text, postVisibility string, data map[string]interface{}) (map[string]interface{}, error) {
	var postLink string
	switch postType {
	case http.MethodPost:
		postLink = "chat/post/create/"
		data["post_title"] = title
		data["topic_external_id"] = externalID
		data["post_text"] = text
	case "CHAT":
		postLink = "chat/chat-post/create/"
		data["post_title"] = title
		data["subject_external_id"] = externalID
		data["post_visibility"] = postVisibility
	default:
		return nil, errors.New("Invalid post_type")
	}

	return request(postLink, p.Session, p.Token, http.MethodPost, data, false, false)
}

func (p *Post) Delete(postID string) (map[string]interface{}, error) {
	data := map[string]interface{}{"external_id": postID}
	return request("chat/post/delete/", p.Session, p.Token, http.MethodPost, data, false, false)
}

func (p *Post) GetTopics() (map[string]interface{}, error) {
	return request("chat/topics/", p.Session, "", http.MethodGet, nil, false, false)
}

func (p *Post) Feed(topic string, num, load int, sort string) (map[string]interface{}, error) {
	url := fmt.Sprintf("chat/posts/?topic=%s&page=%d&posts_to_load=%d&sort=%s", topic, num, load, sort)
	return request(url, p.Session, p.Token, http.MethodGet, nil, false, false)
}
