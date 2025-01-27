package service

import (
	"WasaText/service/consts"
	models2 "WasaText/service/database/models"
	"WasaText/service/helpers"
	"WasaText/service/repositories"
	"fmt"
	"slices"
	"time"
)

type Service struct {
	Repository *repositories.Repository
}

func NewService(repository *repositories.Repository) *Service {
	return &Service{Repository: repository}
}

func (s *Service) CreateOrUpdateUser(user *models2.User) (string, error) {
	existUser, err := s.Repository.GetUser(user)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println(existUser, "oshibka itd")
	if existUser.ID == 0 {
		createdUser, err := s.Repository.CreateUser(user)
		if err != nil {
			return "", err
		}
		fmt.Println(createdUser, "user ne sozdalsya")
		return helpers.GenerateSessionToken(createdUser)
	}

	return helpers.GenerateSessionToken(existUser)
}

func (s *Service) GetUsers(query string, userId uint) (*[]models2.User, error) {
	return s.Repository.GetUsers(query, userId)
}

func (s *Service) GetConversations(userId uint) (map[string]interface{}, error) {
	res := make(map[string]interface{})
	users, err := s.Repository.GetConversationsUsers(userId)
	if err != nil {
		return nil, err
	}

	groups, err := s.Repository.GetGroups(userId)
	if err != nil {
		return nil, err
	}
	fmt.Println(groups, "!!!! ERROR")
	res["users"] = users
	res["groups"] = groups

	return res, nil

}

func (s *Service) SendMessage(userId uint, input *helpers.SendMessageRequest) (*models2.Message, error) {
	var conv *models2.Conversation
	var err error

	if !input.IsGroup {
		conv, err = s.Repository.CheckPrivateConversation(userId, input.ToUserId)
		if err != nil {
			return nil, err
		}
	} else {
		conv, err = s.Repository.CheckGroupConversation(input.GroupId)
		if err != nil {
			return nil, err
		}
	}
	var message models2.Message

	if input.PhotoPath == "" {
		message = models2.Message{
			ConversationID:   conv.ID,
			SenderID:         userId,
			Content:          input.Text,
			MessageType:      consts.MessageTypeText,
			CreatedAt:        time.Now(),
			RepliedMessageID: input.RepliedMessageId,
		}
	} else {
		message = models2.Message{
			ConversationID:   conv.ID,
			SenderID:         userId,
			Content:          input.PhotoPath,
			MessageType:      consts.MessageTypePhoto,
			CreatedAt:        time.Now(),
			RepliedMessageID: input.RepliedMessageId,
		}
	}

	return s.Repository.CreateMessage(&message)
}

func (s *Service) GetMessages(user1ID uint, payload *helpers.GetMessagesRequest) (*[]helpers.MessagesResponse, error) {
	var messages *[]models2.Message
	var err error

	if !payload.IsGroup {
		messages, err = s.Repository.GetPrivateMessages(user1ID, payload.UserOrGroupId)
		if err != nil {
			return nil, err
		}
	} else {
		messages, err = s.Repository.GetGroupMessages(payload.UserOrGroupId)
		if err != nil {
			return nil, err
		}
	}

	var response []helpers.MessagesResponse

	if messages != nil {
		messageMap := make(map[uint]models2.Message)
		for _, message := range *messages {
			messageMap[message.ID] = message
		}

		for _, message := range *messages {

			var emoji string
			repliedMessage := helpers.RepliedMessageResponse{}

			if message.Reactions.ID != 0 {
				emoji = message.Reactions.Reaction
			}

			if message.RepliedMessageID != nil {
				repliedMsg, found := messageMap[*message.RepliedMessageID]
				if found {
					repliedMessage = helpers.RepliedMessageResponse{
						MessageId: repliedMsg.ID,
						Content:   repliedMsg.Content,
					}
				}
			}

			if message.MessageType == consts.MessageTypeText {
				response = append(response, helpers.MessagesResponse{
					MessageId:      message.ID,
					Message:        message.Content,         // Message text
					UserId:         message.SenderID,        // Sender user ID
					ConversationId: message.ConversationID,  // Conversation ID
					Username:       message.Sender.Username, // Conversation ID
					CreatedAt:      message.CreatedAt,       // Conversation ID
					IsPhoto:        false,
					IsRead:         message.IsRead,
					Emoji:          emoji,
					RepliedMessage: repliedMessage,
				})
			} else {
				response = append(response, helpers.MessagesResponse{
					MessageId:      message.ID,
					Message:        message.Content,         // Message text
					UserId:         message.SenderID,        // Sender user ID
					ConversationId: message.ConversationID,  // Conversation ID
					Username:       message.Sender.Username, // Conversation ID
					CreatedAt:      message.CreatedAt,       // Conversation ID
					IsPhoto:        true,
					IsRead:         message.IsRead,
					Emoji:          emoji,
					RepliedMessage: repliedMessage,
				})
			}
		}
	}

	return &response, nil
}

func (s *Service) CreateOrGetUser(user *models2.User) (*models2.User, error) {
	existUser, err := s.Repository.GetUser(user)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if existUser.ID == 0 {
		createdUser, err := s.Repository.CreateUser(user)
		if err != nil {
			return nil, err
		}

		return createdUser, nil
	}

	return existUser, nil
}

func (s *Service) MarkAsRead(userId uint, input *helpers.SendMessageRequest) (bool, error) {
	var conv *models2.Conversation
	var err error

	if !input.IsGroup {
		conv, err = s.Repository.CheckPrivateConversation(userId, input.ToUserId)
		if err != nil {
			return false, err
		}
	} else {
		conv, err = s.Repository.CheckGroupConversation(input.GroupId)
		if err != nil {
			return false, err
		}
	}

	return s.Repository.MarkAsRead(conv.ID, userId)
}

func (s *Service) CreateGroups(payload *helpers.CreateGroupRequest, addedById uint) (bool, error) {
	group, err := s.Repository.CreateGroup(payload, addedById)
	if err != nil {
		return false, err
	}

	_, err = s.Repository.CreateGroupMembers(addedById, addedById, group.ID)
	if err != nil {
		return false, err
	}

	for _, user := range payload.Users {
		_, err = s.Repository.CreateGroupMembers(user.ID, addedById, group.ID)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func (s *Service) UpdateUserProfile(userId uint, filepath string) (bool, error) {
	var user models2.User
	user.ID = userId
	user.ProfilePhotoURL = filepath
	return s.Repository.UpdateUserProfile(&user)
}

func (s *Service) AddUsersToGroup(userId uint, payload *helpers.AddUsersToGroup) (bool, error) {
	users, err := s.Repository.GetGroupMembers(payload.GroupId)
	if err != nil {
		return false, err
	}
	var userIds = make([]uint, len(*users))

	for _, user := range *users {
		userIds = append(userIds, user.UserID)
	}

	for _, v := range payload.Users {
		if !slices.Contains(userIds, v.ID) {
			_, err = s.Repository.CreateGroupMembers(v.ID, userId, payload.GroupId)
			if err != nil {
				return false, err
			}
		}
	}

	return true, nil
}

func (s *Service) DeleteMessage(msgID uint) (bool, error) {
	return s.Repository.DeleteMessage(msgID)
}

func (s *Service) LeaveGroup(userId uint, payload *helpers.GroupRequest) (bool, error) {
	return s.Repository.DeleteGroupMember(userId, payload.Group)
}

func (s *Service) CommentMessage(payload *helpers.CommentMessage, userId uint) (bool, error) {
	return s.Repository.CommentMessage(payload, userId)
}

func (s *Service) UncommentMessage(payload *helpers.UncommentMessage, userId uint) (bool, error) {
	return s.Repository.UncommentMessage(payload, userId)
}

func (s *Service) ForwardMessage(payload *helpers.ForwardMessage, userId uint) (bool, error) {
	if payload.IsPhoto {
		for _, user := range payload.Users {
			sendMessagePayload := helpers.SendMessageRequest{
				ToUserId:  user.ID,
				IsGroup:   false,
				PhotoPath: payload.Text,
			}
			_, err := s.SendMessage(userId, &sendMessagePayload)
			if err != nil {
				return false, err
			}
		}
	} else {
		for _, user := range payload.Users {
			sendMessagePayload := helpers.SendMessageRequest{
				ToUserId: user.ID,
				IsGroup:  false,
				Text:     payload.Text,
			}

			_, err := s.SendMessage(userId, &sendMessagePayload)
			if err != nil {
				return false, err
			}
		}
	}

	return true, nil
}
