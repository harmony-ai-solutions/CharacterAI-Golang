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

func (p *Post) My(postsPage int, postsToLoad int, token string) (map[string]interface{}, error) {
	url := fmt.Sprintf("chat/posts/user/?scope=user&page=%d&posts_to_load=%d", postsPage, postsToLoad)
	return request(url, p.Session, token, http.MethodGet, nil, false, false)
}

func (p *Post) GetPosts(username string, postsPage int, postsToLoad int) (map[string]interface{}, error) {
	url := fmt.Sprintf("chat/posts/user/?username=%s&page=%d&posts_to_load=%d", username, postsPage, postsToLoad)
	return request(url, p.Session, "", http.MethodGet, nil, false, false)
}

func (p *Post) Upvote(postExternalID string, token string) (map[string]interface{}, error) {
	data := map[string]interface{}{"post_external_id": postExternalID}
	return request("chat/post/upvote/", p.Session, token, http.MethodPost, data, false, false)
}

func (p *Post) UndoUpvote(postExternalID string, token string) (map[string]interface{}, error) {
	data := map[string]interface{}{"post_external_id": postExternalID}
	return request("chat/post/undo-upvote/", p.Session, token, http.MethodPost, data, false, false)
}

func (p *Post) SendComment(postID, text, parentUUID, token string) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"post_external_id": postID,
		"text":             text,
		"parent_uuid":      parentUUID,
	}
	return request("chat/comment/create/", p.Session, token, http.MethodPost, data, false, false)
}

func (p *Post) DeleteComment(messageID int, postID, token string) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"external_id":      strconv.Itoa(messageID),
		"post_external_id": postID,
	}
	return request("chat/comment/delete/", p.Session, token, http.MethodPost, data, false, false)
}

func (p *Post) Create(postType, externalID, title, text, postVisibility, token string, data map[string]interface{}) (map[string]interface{}, error) {
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

	return request(postLink, p.Session, token, http.MethodPost, data, false, false)
}

func (p *Post) Delete(postID, token string) (map[string]interface{}, error) {
	data := map[string]interface{}{"external_id": postID}
	return request("chat/post/delete/", p.Session, token, http.MethodPost, data, false, false)
}

func (p *Post) GetTopics() (map[string]interface{}, error) {
	return request("chat/topics/", p.Session, "", http.MethodGet, nil, false, false)
}

func (p *Post) Feed(topic string, num, load int, sort, token string) (map[string]interface{}, error) {
	url := fmt.Sprintf("chat/posts/?topic=%s&page=%d&posts_to_load=%d&sort=%s", topic, num, load, sort)
	return request(url, p.Session, token, http.MethodGet, nil, false, false)
}
