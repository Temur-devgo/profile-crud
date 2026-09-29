package repository

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"sarkor/internal/models"
)

type UsersRepo struct {
	conn *sql.DB
}

func NewUsersRepo(db *sql.DB) *UsersRepo {
	return &UsersRepo{conn: db}
}

const insertUser = `INSERT INTO users (login, password, username, age) VALUES (?, ?, ?, ?)`

func (u *UsersRepo) AddNewUser(ctx *gin.Context, userData models.UserRegReq, password []byte) error {
	stmt, err := u.conn.Prepare(insertUser)
	if err != nil {
		return err
	}

	if _, err = stmt.ExecContext(ctx, userData.Login, password, userData.Name, userData.Age); err != nil {
		return err
	}

	return nil
}

const selectUserByName = `SELECT id, username, age FROM users WHERE username = ?`

func (u *UsersRepo) GetUserByName(ctx *gin.Context, userName string) (models.UserDb, error) {
	var res models.UserDb

	stmt, err := u.conn.Prepare(selectUserByName)
	if err != nil {
		return res, err
	}

	if err = stmt.QueryRowContext(ctx, userName).Scan(&res.ID, &res.Name, &res.Age); err != nil {
		return res, err
	}

	return res, nil
}

const insertUserPhone = `INSERT INTO user_phones (user_id, phone, description, is_fax) VALUES (?, ?, ?, ?)`

func (u *UsersRepo) AddUserPhone(ctx *gin.Context, userID int, phoneReq models.UserPhoneReq) error {
	stmt, err := u.conn.Prepare(insertUserPhone)
	if err != nil {
		return err
	}

	if _, err = stmt.ExecContext(ctx, userID, phoneReq.Phone, phoneReq.Description, phoneReq.IsFax); err != nil {
		return err
	}

	return nil
}

const selectUserPhone = "SELECT id, user_id, phone, description, is_fax FROM user_phones WHERE phone LIKE ?"

func (u *UsersRepo) GetUserPhone(ctx *gin.Context, phone string) ([]models.UserPhoneDb, error) {
	var res []models.UserPhoneDb

	stmt, err := u.conn.Prepare(selectUserPhone)
	if err != nil {
		return res, err
	}

	rows, err := stmt.QueryContext(ctx, "%"+phone+"%")
	if err != nil {
		return res, err
	}

	for rows.Next() {
		var phoneData models.UserPhoneDb
		if err = rows.Scan(&phoneData.ID, &phoneData.UserID, &phoneData.Phone, &phoneData.Description, &phoneData.IsFax); err != nil {
			return res, err
		}

		res = append(res, phoneData)
	}

	return res, nil
}

const updateUserPhone = `UPDATE user_phones SET phone = ?, description = ?, is_fax = ? WHERE id = ? AND user_id = ?`

func (u *UsersRepo) UpdateUserPhone(ctx *gin.Context, userID int, req models.UserPhoneUpdateReq) error {
	stmt, err := u.conn.Prepare(updateUserPhone)
	if err != nil {
		return err
	}

	if _, err = stmt.ExecContext(ctx, req.Phone, req.Description, req.IsFax, req.ID, userID); err != nil {
		return err
	}

	return nil
}

const deleteUserPhone = `DELETE FROM user_phones WHERE phone = ? AND user_id = ?`

func (u *UsersRepo) RemoveUserPhone(ctx *gin.Context, userID int, phone string) error {
	stmt, err := u.conn.Prepare(deleteUserPhone)
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, phone, userID)
	if err != nil {
		return err
	}

	rowsCnt, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsCnt == 0 {
		return fmt.Errorf("phone not fount: %s", phone)
	}

	return nil
}

const userExist = `SELECT id, password FROM users WHERE login = ?`

func (u *UsersRepo) GetUserByLogin(ctx *gin.Context, login string) (*models.UserDb, error) {
	stmt, err := u.conn.Prepare(userExist)
	if err != nil {
		return nil, err
	}

	var res models.UserDb
	if err = stmt.QueryRowContext(ctx, login).Scan(&res.ID, &res.PassHash); err != nil {
		return nil, err
	}

	if res.ID == 0 {
		return nil, nil
	}

	return &res, nil
}
