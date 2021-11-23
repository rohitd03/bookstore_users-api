package users

import (
	"fmt"

	"github.com/rohitd03/bookstore_users-api/datasources/mysql/users_db"
	"github.com/rohitd03/bookstore_users-api/logger"
	"github.com/rohitd03/bookstore_users-api/utils/errors"
)

const (
	queryInsertUser       = "INSERT INTO users(first_name, last_name, email, date_created, status, password) VALUES(?, ?, ?, ?, ?, ?)"
	queryGetUser          = " id, first_name, last_name, email, date_created, status from users WHERE id = ?;"
	queryUpdateUser       = "UPDATE users SET first_name=?, last_name=?, email=? WHERE id=?;"
	queryDeleteUser       = "Delete FROM users WHERE id = ?;"
	queryFindUserByStatus = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE status = ?;"
)

var (
	userDB = make(map[int64]*User)
)

func (user *User) Get() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		logger.Error("Error when trying to prepare get user statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Id)
	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); getErr != nil {
		logger.Error("Error when trying to get user by Id", getErr)
		return errors.NewInternalServerError("database error")
		// return mysql_utils.ParseError(getErr)
	}
	return nil
}

func (user *User) Save() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		// return errors.NewInternalServerError(err.Error())
		logger.Error("Error when trying to prepare save user statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	insertResult, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, &user.Status, &user.Password)
	if saveErr != nil {
		// return mysql_utils.ParseError(saveErr)
		logger.Error("Error when trying to save user", saveErr)
		return errors.NewInternalServerError("database error")
	}

	// Alternate way
	// insertResult, err := users_db.Client.Exec(queryInsertUser,user.FirstName, user.LastName, user.Email, user.DateCreated)

	userId, err := insertResult.LastInsertId()
	if err != nil {
		// return mysql_utils.ParseError(err)
		logger.Error("Error when trying to get last inserted id after creating a new user", err)
		return errors.NewInternalServerError("database error")
	}
	user.Id = userId
	// log.Println("Sab changa see")
	return nil
}

func (user *User) Update() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		// return errors.NewInternalServerError(err.Error())
		logger.Error("Error when trying to prepare update user statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()
	_, errUpdate := stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id)
	if errUpdate != nil {
		// return mysql_utils.ParseError(errUpdate)
		logger.Error("Error when trying to update user", errUpdate)
		return errors.NewInternalServerError("database error")
	}
	return nil
}

func (user *User) Delete() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		// return errors.NewInternalServerError(err.Error())
		logger.Error("Error when trying to prepare delete user statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()
	_, errDelete := stmt.Exec(user.Id)
	if errDelete != nil {
		// return mysql_utils.ParseError(errDelete)
		logger.Error("Error when trying to delete user", errDelete)
		return errors.NewInternalServerError("database error")
	}
	return nil
}

func (user *User) FindByStatus(status string) ([]User, *errors.RestErr) {
	stmt, err := users_db.Client.Prepare(queryFindUserByStatus)
	if err != nil {
		// return nil, errors.NewInternalServerError(err.Error())
		logger.Error("Error when trying to prepare find user by status statement", err)
		return nil, errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	rows, err := stmt.Query(status)
	if err != nil {
		// return nil, errors.NewInternalServerError(err.Error())
		logger.Error("Error when trying to find users by status", err)
		return nil, errors.NewInternalServerError("database error")
	}
	defer rows.Close()

	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
			// return nil, mysql_utils.ParseError(err)
			logger.Error("Error when trying to scan user row into user struct", err)
			return nil, errors.NewInternalServerError("database error")
		}
		results = append(results, user)
	}
	if len(results) == 0 {
		return nil, errors.NewNotFoundError(fmt.Sprintf("No users matching status %s", status))
	}
	return results, nil

}
