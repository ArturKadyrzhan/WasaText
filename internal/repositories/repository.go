package repositories

import (
	"WasaText/cmd/database/models"
	"WasaText/internal/helpers"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"strings"
	"time"
)

type Repository struct {
	database *gorm.DB
}

func NewRepository(database *gorm.DB) *Repository {
	return &Repository{database: database}
}

func (r *Repository) CreateUser(user *models.User) (*models.User, error) {
	fmt.Println(user, "error etc")
	result := r.database.Create(user)
	if result.Error != nil {
		msg := result.Error
		return nil, msg
	}
	return user, nil
}

func (r *Repository) GetUser(user *models.User) (*models.User, error) {
	result := r.database.Model(&models.User{}).Where("username = ?", user.Username).Find(&user)
	if result.Error != nil {
		msg := result.Error
		return nil, msg
	}
	return user, nil

}

//обращение к датабэйз напиши сам буквально одна строка, Может быть несколько юзеров,
//функция принимает стринг в query, по нему база должна найти все совпадающие юзеры,
//"неважно  будет ли совпадать одна буква или все буквы

func (r *Repository) GetUsers(query string, userId uint) (*[]models.User, error) {
	var users []models.User
	if err := r.database.Where("LOWER(username) LIKE ? AND id != ?",
		"%"+strings.ToLower(query)+"%", userId).Find(&users).Error; err != nil {
		return nil, err
	}
	return &users, nil
}

//домашка

func (r *Repository) GetConversationsUsers(userId uint) (*[]models.User, error) {
	var users []models.User

	if err := r.database.Raw(`
        SELECT DISTINCT u.*FROM users u  JOIN conversations c ON (u.id = c.user1_id OR u.id = c.user2_id)
        WHERE (c.user1_id = ? OR c.user2_id = ?)
          AND u.id != ?`, userId, userId, userId).Scan(&users).Error; err != nil {
		return nil, err
	}

	return &users, nil
}

func (r *Repository) GetGroups(userId uint) (*[]models.Group, error) {
	var groups []models.Group

	if err := r.database.Raw(`
        SELECT DISTINCT g.*
        FROM groups g
        JOIN group_members gm ON g.id = gm.group_id
        WHERE gm.user_id = ?`, userId).Scan(&groups).Error; err != nil {
		return nil, err
	}

	return &groups, nil
}

func (r *Repository) CheckPrivateConversation(user1Id uint, user2Id uint) (*models.Conversation, error) {
	var conv models.Conversation

	result := r.database.Model(&models.Conversation{}).
		Where("(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)", user1Id, user2Id, user2Id, user1Id).
		Where("is_group = ?", false).
		First(&conv)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			user1 := models.User{ID: user1Id}
			user2 := models.User{ID: user2Id}

			conv = models.Conversation{
				User1ID: &user1.ID,
				User2ID: &user2.ID,
			}

			return r.CreateConversation(&conv)
		}
		return nil, result.Error
	}

	return &conv, nil
}

func (r *Repository) CreateMessage(message *models.Message) (*models.Message, error) {
	result := r.database.Create(&message).Find(&message)
	if result.Error != nil {
		return nil, result.Error
	}

	return message, nil
}

func (r *Repository) CreateConversation(conv *models.Conversation) (*models.Conversation, error) {
	result := r.database.Create(conv)
	if result.Error != nil {
		msg := result.Error
		return nil, msg
	}

	return conv, nil
}

func (r *Repository) GetPrivateMessages(user1ID uint, user2ID uint) (*[]models.Message, error) {
	var conversation models.Conversation
	err := r.database.Where("(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)", user1ID, user2ID, user2ID, user1ID).First(&conversation).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Println("Error fetching conversation:", err)
		return nil, err
	}

	var messages []models.Message
	err = r.database.
		Where("conversation_id = ?", conversation.ID).
		Order("created_at asc").
		Preload("Sender").
		Preload("Reactions").
		Find(&messages).Error

	return &messages, nil
}

func (r *Repository) GetGroupMessages(groupID uint) (*[]models.Message, error) {
	var conversation models.Conversation
	err := r.database.Where("group_id = ?", groupID).First(&conversation).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Println("Error fetching conversation:", err)
		return nil, err
	}

	var messages []models.Message
	err = r.database.
		Where("conversation_id = ?", conversation.ID).
		Order("created_at asc").
		Preload("Sender").
		Preload("Reactions").
		Find(&messages).Error

	return &messages, nil
}

func (r *Repository) CheckGroupConversation(groupId uint) (*models.Conversation, error) {
	var conv models.Conversation

	result := r.database.Model(&models.Conversation{}).
		Where("group_id = ?", groupId).
		Where("is_group = ?", true).
		First(&conv)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			conv = models.Conversation{
				GroupID: &groupId,
				IsGroup: true,
			}

			return r.CreateConversation(&conv)
		}
		return nil, result.Error
	}

	return &conv, nil
}

func (r *Repository) MarkAsRead(convId, userId uint) (bool, error) {
	err := r.database.Model(&models.Message{}).
		Where("conversation_id = ? AND is_read = ? AND sender_id != ?", convId, false, userId).
		Update("is_read", true).Error
	if err != nil {
		log.Println("Error updating is_read:", err)
		return false, err
	}

	return true, nil
}

func (r *Repository) CreateGroup(payload *helpers.CreateGroupRequest, userId uint) (*models.Group, error) {
	var group models.Group

	group.CreatedBy = userId
	group.Name = payload.GroupName
	group.GroupPhotoURL = payload.GroupPhotoPath

	result := r.database.Create(&group)
	if result.Error != nil {
		msg := result.Error
		return nil, msg
	}

	return &group, nil
}

func (r *Repository) CreateGroupMembers(userId, addedById, groupId uint) (bool, error) {
	var groupMember models.GroupMember
	groupMember.GroupID = groupId
	groupMember.UserID = userId
	groupMember.AddedBy = addedById
	groupMember.AddedAt = time.Now()

	result := r.database.Create(&groupMember)
	if result.Error != nil {
		msg := result.Error
		return false, msg
	}

	return true, nil
}

func (r *Repository) UpdateUserProfile(user *models.User) (bool, error) {
	result := r.database.Model(&user).Where("id = ?", user.ID).Update("profile_photo_url", user.ProfilePhotoURL)
	if result.Error != nil {
		msg := result.Error
		return false, msg
	}

	return true, nil
}

func (r *Repository) GetGroupMembers(groupId uint) (*[]models.GroupMember, error) {
	var groupMember []models.GroupMember
	result := r.database.Model(&models.GroupMember{}).
		Where("group_id = ?", groupId).Preload("User").
		Find(&groupMember)
	if result.Error != nil {
		msg := result.Error
		return nil, msg
	}

	return &groupMember, nil
}

func (r *Repository) DeleteMessage(msgID uint) (bool, error) {
	err := r.database.Model(&models.Message{}).
		Where("id = ?", msgID).
		Delete(&models.Message{}).Error
	if err != nil {
		log.Println("Error updating is_read:", err)
		return false, err
	}

	return true, nil
}

func (r *Repository) DeleteGroupMember(userId uint, group models.Group) (bool, error) {
	result := r.database.Model(&models.GroupMember{}).
		Where("user_id = ? AND group_id = ?", userId, group.ID).
		Delete(&models.GroupMember{})

	if result.Error != nil {
		return false, result.Error
	}

	return true, nil
}

func (r *Repository) CommentMessage(payload *helpers.CommentMessage, userId uint) (bool, error) {
	var reaction models.Reaction

	result := r.database.Where("message_id = ? AND user_id = ?", payload.MessageId, userId).First(&reaction)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			reaction.MessageID = payload.MessageId
			reaction.UserID = userId
			reaction.Reaction = payload.Emoji

			if err := r.database.Create(&reaction).Error; err != nil {
				return false, err
			}
		}
		return false, result.Error

	}

	if err := r.database.Model(&models.Reaction{}).Where("message_id = ? AND user_id = ?", payload.MessageId, userId).Update("reaction", payload.Emoji).Error; err != nil {
		return false, err
	}

	return true, nil
}

func (r *Repository) UncommentMessage(payload *helpers.UncommentMessage, userId uint) (bool, error) {
	res := r.database.Model(&models.Reaction{}).
		Where("message_id = ? AND user_id =?", payload.MessageId, userId).
		Delete(&models.Reaction{})
	if res.Error != nil {
		log.Println("Error updating is_read:", res.Error)
		return false, res.Error
	}

	if res.RowsAffected == 0 {
		return false, nil
	}

	return true, nil
}
