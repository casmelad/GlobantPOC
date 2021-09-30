package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/caarlos0/env/v6"
	"github.com/casmelad/GlobantPOC/pkg/users"
	_ "github.com/go-sql-driver/mysql" //only for implicit use
)

const (
	INSERTUSER         = "INSERT INTO Users(Email, Name, LastName) VALUES (?, ?, ?)"
	SELECTUSERBYID     = "SELECT Id, Email, Name, LastName FROM Users WHERE Id = ?"
	SELECTUSEERBYEMAIL = "SELECT Id, Email, Name, LastName FROM Users WHERE Email = ?"
	SELECTALLUSERS     = "SELECT Id, Email, Name, LastName FROM Users"
	UPDATEUSER         = "UPDATE Users SET Name=?, LastName=? WHERE Id = ?"
	DELETEUSER         = "DELETE FROM Users WHERE Id= ?"
)

type config struct {
	User      string `env:"MYSQL_USER" envDefault:"root"`
	Password  string `env:"MYSQL_PASSWORD" envDefault:"BulkD3v_mysql"`
	Port      string `env:"MYSQL_PORT" envDefault:":3306"`
	Host      string `env:"MYSQL_HOST" envDefault:""`
	DefaultDB string `env:"MYSQL_DEFAULTDB" envDefault:"Users"`
}

func initMySQLRepository() (*sql.DB, error) {

	cfg := config{}

	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	connectionString := fmt.Sprintf("%s:%s@tcp(%s%s)/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DefaultDB)
	db, err := sql.Open("mysql", connectionString)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

//MySQLRepository - is a mysql implementation of users repository
type MySQLRepository struct {
	db *sql.DB
}

//NewMySQLUserRepository - returns a MySQLRepository type pointer
func NewMySQLUserRepository() (*MySQLRepository, error) {

	db, err := initMySQLRepository()

	if err != nil {
		return nil, err
	}

	return &MySQLRepository{
		db: db,
	}, nil
}

//Add - adds a user to the repository
func (r *MySQLRepository) Add(ctx context.Context, usr users.User) (int, error) {

	stmt, err := r.db.Prepare(INSERTUSER)

	if err != nil {
		return 0, err
	}

	result, err := stmt.Exec(usr.Email, usr.Name, usr.LastName)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		//TODO log error
		return 0, err
	}

	return (int(id)), nil
}

//GetByID - retrieves a user from the repository based on the integer id
func (r *MySQLRepository) GetByID(ctx context.Context, userID int) (users.User, error) {

	usr := users.User{}

	err := r.db.QueryRow(SELECTUSERBYID, userID).
		Scan(&usr.ID, &usr.Email, &usr.Name, &usr.LastName)

	if err == sql.ErrNoRows {
		return usr, nil
	}

	return usr, err

}

//GetByEmail - retrieves a user from the repository based on the email address
func (r *MySQLRepository) GetByEmail(ctx context.Context, email string) (users.User, error) {

	usr := users.User{}
	row := r.db.QueryRow(SELECTUSEERBYEMAIL, email)
	err := row.Scan(&usr.ID, &usr.Email, &usr.Name, &usr.LastName)

	if err == sql.ErrNoRows {
		return usr, nil
	}

	return usr, err
}

//GetAll - retrieves all the users from the repository
func (r *MySQLRepository) GetAll(ctx context.Context) ([]users.User, error) {

	usrs := []users.User{}
	records, err := r.db.Query(SELECTALLUSERS)

	if err != nil {
		fmt.Println(err)
	}

	defer records.Close()

	for records.Next() {
		var user users.User

		if err := records.Scan(&user.ID, &user.Email, &user.Name, &user.LastName); err != nil {
			return []users.User{}, err
		}

		usrs = append(usrs, user)
	}

	return usrs, nil
}

//Update -  updates the information of a user
func (r *MySQLRepository) Update(ctx context.Context, usr users.User) error {

	stmt, err := r.db.Prepare(UPDATEUSER)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(usr.Name, usr.LastName, usr.ID)

	if err != nil {
		return err
	}

	if rows, err := result.RowsAffected(); rows == 0 || err != nil {
		return errors.New("no records were updated " + err.Error())
	}

	return nil
}

//Delete - deletes a user from the repository
func (r *MySQLRepository) Delete(ctx context.Context, userID int) error {

	stmt, err := r.db.Prepare(DELETEUSER)

	if err != nil {
		return err
	}

	result, err := stmt.Exec(userID)

	if err != nil {
		return err
	}

	if rows, err := result.RowsAffected(); rows == 0 || err != nil {
		return errors.New("no records were deleted " + err.Error())
	}

	return nil
}
