package neo4j

import (
	"github.com/pkg/errors"
	"github.com/simonhylander/booker"
	"time"
)

type User struct {
	firstname string
	lastname string
	username string
	password string
	email string
	active bool
}

func (u User) FindByUsername(username string) (booker.User, error) {
	var user booker.User

	users := []User{
		{
			firstname: "Simon",
			lastname: "hylander",
			username: "simon",
			password: "$2a$10$ZHKJU/bOqkEpMUczOQ5swORcP2qQowiWEcX9hQmZuy0eMwKiu24Su",
			email: "hylandersimon@gmail.com",
			active: true,
		},
	}

	for _, user := range users {
		if (user.username == username) {

			return booker.User{
				Base:               booker.Base{
					ID: 1,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					DeletedAt: time.Now(),
				},
				FirstName:         user.firstname,
				LastName:           user.lastname,
				Username:           user.username,
				Password:           user.password,
				Email:              user.email,
				Mobile:             "",
				Phone:              "",
				Address:            "",
				Active:             user.active,
				LastLogin:          time.Time{},
				LastPasswordChange: time.Time{},
				Token:              "", // TODO
				/*Role:              nil,
				RoleID:             0,
				CompanyID:          0,
				LocationID:         0,*/
			}, nil
		}
	}

	return user, errors.New("User not found")
}

// TODO:
func (u User) Update(user booker.User) error {
	return nil
	//return db.Update(&user)
}

/*import (
	"github.com/go-pg/pg/v9/orm"
	"github.com/simonhylander/booker"
)

// User represents the client for user table
type User struct{}

// View returns single user by ID
func (u User) View(db orm.DB, id int) (gorsk.User, error) {
	var user gorsk.User
	sql := `SELECT "user".*, "role"."id" AS "role__id", "role"."access_level" AS "role__access_level", "role"."name" AS "role__name" 
	FROM "users" AS "user" LEFT JOIN "roles" AS "role" ON "role"."id" = "user"."role_id" 
	WHERE ("user"."id" = ? and deleted_at is null)`
	_, err := db.QueryOne(&user, sql, id)
	return user, err
}

// FindByUsername queries for single user by username
func (u User) FindByUsername(db orm.DB, uname string) (gorsk.User, error) {
	var user gorsk.User
	sql := `SELECT "user".*, "role"."id" AS "role__id", "role"."access_level" AS "role__access_level", "role"."name" AS "role__name" 
	FROM "users" AS "user" LEFT JOIN "roles" AS "role" ON "role"."id" = "user"."role_id" 
	WHERE ("user"."username" = ? and deleted_at is null)`
	_, err := db.QueryOne(&user, sql, uname)
	return user, err
}

// FindByToken queries for single user by token
func (u User) FindByToken(db orm.DB, token string) (gorsk.User, error) {
	var user gorsk.User
	sql := `SELECT "user".*, "role"."id" AS "role__id", "role"."access_level" AS "role__access_level", "role"."name" AS "role__name" 
	FROM "users" AS "user" LEFT JOIN "roles" AS "role" ON "role"."id" = "user"."role_id" 
	WHERE ("user"."token" = ? and deleted_at is null)`
	_, err := db.QueryOne(&user, sql, token)
	return user, err
}

// Update updates user's info
func (u User) Update(db orm.DB, user gorsk.User) error {
	return db.Update(&user)
}
*/