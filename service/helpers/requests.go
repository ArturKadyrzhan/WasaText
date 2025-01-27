package helpers

import (
	models2 "WasaText/service/database/models"
)

type SendMessageRequest struct {
	Text             string `json:"text"`
	ToUserId         uint   ` json:"toUserId"`
	IsGroup          bool   ` json:"isGroup"`
	GroupId          uint   `json:"groupId"`
	PhotoPath        string `json:"photoPath"`
	RepliedMessageId *uint  `json:"repliedMessageId"`
}

type GetMessagesRequest struct {
	UserOrGroupId uint `json:"id"`
	IsGroup       bool `json:"isGroup"`
}

type CreateGroupRequest struct {
	GroupName      string         `json:"groupName"`
	GroupPhotoPath string         `json:"groupPhotoPath"`
	Users          []models2.User `json:"users"`
}

type AddUsersToGroup struct {
	GroupId uint           ` json:"groupId"`
	Users   []models2.User `json:"users"`
}

type DeleteMessage struct {
	MessageId uint ` json:"id"`
}

type CommentMessage struct {
	MessageId uint   ` json:"messageId"`
	Emoji     string ` json:"emoji"`
}

type GroupRequest struct {
	Group models2.Group `json:"group"`
}

type UncommentMessage struct {
	MessageId uint `json:"id"`
}

type ForwardMessage struct {
	MessageId uint           `json:"id"`
	Text      string         `json:"text"`
	IsPhoto   bool           `json:"isPhoto"`
	Users     []models2.User `json:"users"`
}
