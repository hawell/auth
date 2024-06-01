package database

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
)

func addUser(t *sql.Tx, userId ObjectId, u NewUser) error {
	_, err := t.Exec("INSERT INTO User(Id, Email, Password, Status) VALUES (?, ?, ?, ?)", userId, u.Email, u.Password, u.Status)
	return err
}

func (db *Database) deleteUser(name string) error {
	res, err := db.db.Exec("DELETE FROM User WHERE Email = ?", name)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (db *Database) getUser(name string) (User, error) {
	res := db.db.QueryRow("SELECT Id, Email, Password, Status FROM User WHERE Email = ?", name)
	var u User
	err := res.Scan(&u.Id, &u.Email, &u.Password, &u.Status)
	return u, err
}

func deleteUsers(t *sql.Tx) error {
	_, err := t.Exec("DELETE FROM User")
	return err
}

func (db *Database) getVerification(userId ObjectId, verificationType VerificationType) (string, error) {
	res := db.db.QueryRow("SELECT Code FROM Verification WHERE User_Id = ? AND Type = ?", userId, verificationType)
	var code string
	err := res.Scan(&code)
	return code, err
}

func setVerification(t *sql.Tx, userId ObjectId, v Verification) error {
	_, err := t.Exec("REPLACE INTO Verification(Code, Type, User_Id) VALUES (?, ?, ?)", v.Code, v.Type, userId)
	return err
}

func deleteVerifications(t *sql.Tx) error {
	_, err := t.Exec("DELETE FROM Verification")
	return err
}

func setUserStatus(t *sql.Tx, userId ObjectId, status UserStatus) error {
	_, err := t.Exec("UPDATE User SET Status = ? WHERE Id = ?", status, userId)
	return err
}

func setUserPassword(t *sql.Tx, userId ObjectId, password string) error {
	_, err := t.Exec("UPDATE User SET Password = ? WHERE Id = ?", password, userId)
	return err
}

func parseError(err error) error {
	var mysqlErr *mysql.MySQLError
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	if errors.As(err, &mysqlErr) {
		switch mysqlErr.Number {
		case 1062:
			return ErrDuplicateEntry
		case 1452:
			return ErrInvalid
		default:
			return err
		}
	}
	return err
}

type actionFunction func(t *sql.Tx, userId ObjectId) error

func (db *Database) applyVerifiedAction(code string, verificationType VerificationType, action actionFunction) error {
	res := db.db.QueryRow("select U.Id, V.Type from Verification V left join User U on U.Id = V.User_Id WHERE Code = ?", code)
	var (
		userId     ObjectId
		storedType VerificationType
	)
	if err := res.Scan(&userId, &storedType); err != nil {
		return err
	}
	if storedType != verificationType {
		return errors.New("unknown verification type")
	}

	err := db.withTransaction(func(t *sql.Tx) error {
		if err := action(t, userId); err != nil {
			return err
		}
		if _, err := t.Exec("DELETE FROM Verification WHERE Code = ?", code); err != nil {
			return err
		}
		return nil
	})
	return err
}
