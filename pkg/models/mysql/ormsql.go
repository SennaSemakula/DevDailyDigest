/**
package ormsql

Users model to implement data access layer (ORM) to provide user operations on data
**/
package ormsql

import (
	"github.com/Pioneersltd/DevDailyDigest/v1/pkg/models"
	"database/sql"
	"errors"
	"fmt"
)

var ErrNoRecord = errors.New("no matching record found")

type Model interface {
	Insert() (int, error)
	Get() interface{}
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(email, country, city, language, level, frameworks, resume, git, date, jobStatus, referral string, recruit bool) (int, error) {

	// Check if user exists
	// TODO: Figure out a better way for duplicate return type. 2 is not acceptable

	if user, _ := m.Get(email); user != nil {
		query := fmt.Sprintf(`UPDATE USERS
				 SET email=?, git=?, recruit=?, resume=?, language=?, frameworks=?, level=?, country=?, city=?, date=?, jobstatus=?, referral=?
				 WHERE email=%q
				 `, email)

		_, err := m.DB.Exec(query, email, git, recruit, resume, language, frameworks, level, country, city, date, jobStatus, referral)

		if err != nil {
			return 0, err
		}

		return 2, nil
	} else {
		query := `INSERT INTO USERS (email, git, recruit, resume, language, frameworks, level, country, city, date, jobstatus, referral)
		VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`
		_, err := m.DB.Exec(query, email, git, recruit, resume, language, frameworks, level, country, city, date, jobStatus, referral)

		if err != nil {
			return 0, err
		}
	}

	return 1, nil
}

func (m *UserModel) Get(email string) (*models.User, error) {
	query := `SELECT email FROM USERS
			  WHERE email = ?`

	row := m.DB.QueryRow(query, email)
	user := &models.User{}

	err := row.Scan(&user.Email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return user, nil
}

func (m *UserModel) Delete(email string) error {
	query := `DELETE FROM USERS
			  WHERE email = ?`

	m.DB.QueryRow(query, email)

	return nil
}

func (m *UserModel) GetAll() ([]*models.User, error) {
	query := `SELECT * FROM USERS`

	rows, err := m.DB.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	users := []*models.User{}

	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(&user.Email, &user.Git, &user.Recruit, &user.Resume, &user.Language, &user.Frameworks, &user.Level, &user.Country, &user.City, &user.Date, &user.JobStatus, &user.Referral)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	// Check if there were any errors during iterations
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
