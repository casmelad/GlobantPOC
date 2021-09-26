package infrastructure

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/casmelad/GlobantPOC/pkg/domain/entities"
	_ "github.com/go-sql-driver/mysql"
)

var connectionString string

type config struct {
	User      string `env:"MYSQL_USER" envDefault:"root"`
	Password  string `env:"MYSQL_PASSWORD" envDefault:"BulkD3v_mysql"`
	Port      string `env:"MYSQL_PORT" envDefault:":3306"` //Como pasarlo hacia abajo?
	Host      string `env:"MYSQL_HOST" envDefault:""`
	DefaultDB string `env:"MYSQL_DEFAULTDB" envDefault:"Users"`
}

//revisar de donde tomar la configuración poara que no quede en duro
func Init() {

	cfg := config{}

	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	connectionString = fmt.Sprintf("%s:%s@tcp(%s%s)/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DefaultDB)

	fmt.Println(connectionString)
}

type MySqlUserRepository struct {
	db *sql.DB
}

func NewMySqlUserRepository() *MySqlUserRepository {
	Init()
	return &MySqlUserRepository{}
}

func (r *MySqlUserRepository) Add(usr entities.User) int {

	var err error

	if r.db, err = sql.Open("mysql", connectionString); err != nil {
		return 0
	}

	defer r.db.Close()

	result, err := r.db.Exec("INSERT INTO Users(Email, Name, LastName) VALUES (?, ?, ?)", usr.Email, usr.Name, usr.LastName)

	if err != nil {
		return 0
	}

	id, err := result.LastInsertId()

	if err != nil {
		//TODO log error
		return 0
	}

	return (int(id))
}

func (r *MySqlUserRepository) GetById(userId int) entities.User {

	var err error
	usr := entities.User{}
	r.db, err = sql.Open("mysql", connectionString)
	defer r.db.Close()
	r.db.QueryRow("SELECT Id, Email, Name, LastName FROM Users WHERE Id = ?", userId).Scan(&usr.Id, &usr.Email, &usr.Name, &usr.LastName)

	if err != nil {
		return entities.User{}
		//TODO log error
	}

	return usr

}

func (r *MySqlUserRepository) GetByEmail(email string) entities.User {
	var err error
	usr := entities.User{}
	r.db, err = sql.Open("mysql", connectionString)
	defer r.db.Close()
	r.db.QueryRow("SELECT Id, Email, Name, LastName FROM Users WHERE Email = ?", email).Scan(&usr.Id, &usr.Email, &usr.Name, &usr.LastName)

	if err != nil {
		return entities.User{}
		//TODO log error
	}

	return usr
}

func (r *MySqlUserRepository) GetAll() []entities.User {

	var err error
	usrs := []entities.User{}
	r.db, err = sql.Open("mysql", connectionString)
	defer r.db.Close()
	if err != nil {
		fmt.Println(err)
		//TODO log error
	}

	records, err := r.db.Query("SELECT Id, Email, Name, LastName FROM Users")

	if err != nil {
		fmt.Println(err)
	}

	for records.Next() {
		var user entities.User

		if err := records.Scan(&user.Id, &user.Email, &user.Name, &user.LastName); err != nil {
			log.Fatal(err)
		}

		usrs = append(usrs, user)
	}

	return usrs
}

func (r *MySqlUserRepository) Update(usr entities.User) int {
	var err error
	r.db, err = sql.Open("mysql", connectionString)
	r.db.Exec("UPDATE Users SET Name=?, LastName=? WHERE Id = ?", usr.Name, usr.LastName, usr.Id)

	if err != nil {
		return 0
		//TODO log error
	}

	return 1
}

func (r *MySqlUserRepository) Delete(userId int) int {
	var err error

	fmt.Println(userId)

	r.db, err = sql.Open("mysql", connectionString)

	defer r.db.Close()

	if err != nil {
		fmt.Println(err)
		return 0
		//TODO log error
	}

	_, err = r.db.Exec("DELETE FROM Users WHERE Id= ?", userId)

	if err != nil {
		fmt.Println(err)
		return 0
		//TODO log error
	}

	return 1
}
