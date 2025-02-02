package conversation

import (
	"chat/types"
	"chat/utils"
	"database/sql"
	"log"
)

func (c *Conversation) SaveConversation(db *sql.DB) bool {
	data := utils.ToJson(c.GetMessage())
	query := "INSERT INTO conversation (user_id, conversation_id, conversation_name, data) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE conversation_name = VALUES(conversation_name), data = VALUES(data)"

	stmt, err := db.Prepare(query)
	if err != nil {
		return false
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Println(err)
		}
	}(stmt)

	_, err = stmt.Exec(c.UserID, c.Id, c.Name, data)
	if err != nil {
		return false
	}
	return true
}
func GetConversationLengthByUserID(db *sql.DB, userId int64) int64 {
	var length int64
	err := db.QueryRow("SELECT COUNT(*) FROM conversation WHERE user_id = ?", userId).Scan(&length)
	if err != nil {
		return -1
	}
	return length
}

func LoadConversation(db *sql.DB, userId int64, conversationId int64) *Conversation {
	conversation := Conversation{
		UserID: userId,
		Id:     conversationId,
	}

	var data string
	err := db.QueryRow("SELECT conversation_name, data FROM conversation WHERE user_id = ? AND conversation_id = ?", userId, conversationId).Scan(&conversation.Name, &data)
	if err != nil {
		return nil
	}

	conversation.Message, err = utils.Unmarshal[[]types.ChatGPTMessage]([]byte(data))
	if err != nil {
		return nil
	}

	return &conversation
}

func LoadConversationList(db *sql.DB, userId int64) []Conversation {
	var conversationList []Conversation
	rows, err := db.Query("SELECT conversation_id, conversation_name FROM conversation WHERE user_id = ?", userId)
	if err != nil {
		return conversationList
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	for rows.Next() {
		var conversation Conversation
		err := rows.Scan(&conversation.Id, &conversation.Name)
		if err != nil {
			continue
		}
		conversationList = append(conversationList, conversation)
	}

	return utils.Reverse(conversationList)
}

func (c *Conversation) DeleteConversation(db *sql.DB) bool {
	_, err := db.Exec("DELETE FROM conversation WHERE user_id = ? AND conversation_id = ?", c.UserID, c.Id)
	if err != nil {
		return false
	}
	return true
}
