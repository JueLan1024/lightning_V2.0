package mysql

import (
	"context"
	"database/sql"
	"errors"
	"web_app/models"
)

// UsernameExists 判断用户名是否存在
func UsernameExists(ctx context.Context, username string) (err error) {
	sqlStr := `select count(user_id) from user where username=?`
	var count int
	if err := db.GetContext(ctx, &count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return ErrorUsernameExist
	}
	return
}

// CreateUser 创建用户信息
func CreateUser(user *models.User) (err error) {
	sqlStr := `INSERT INTO user (user_id, username, password) VALUES (?,?,?)`
	_, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password)
	return err
}

// GetUserByUsername 通过用户名获得用户信息
func GetUserByUsername(ctx context.Context, username string, user *models.User) (err error) {
	sqlStr := `SELECT user_id, username, password FROM user WHERE username = ?`
	if err = db.GetContext(ctx, user, sqlStr, username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrorUsernameNotFound
		}
	}
	return
}

// GetUserName 通过用户id获得用户名
func GetUserName(ctx context.Context, userID int64) (username string, err error) {
	sqlStr := `select username from user where user_id = ?`
	err = db.Get(&username, sqlStr, userID)
	if err == sql.ErrNoRows {
		return "", ErrorUserNotFound
	}
	if err != nil {
		return "", err
	}
	return username, nil
}
