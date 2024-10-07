package main

import (
	"context"
	"fmt"
	"log"
	"training-golang/session-6-db-pgx/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dsn := "postgresql://postgres:admin@localhost:5432/go_db"
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalln(err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Connected to database")

	//query untuk mengambil row
	var u entity.User
	err = pool.QueryRow(ctx, "SELECT id, name FROM users order by id desc limit 1").Scan(&u.ID, &u.Name)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("user retrieved from database", u)

	//query untuk menjalankan perintah insert/update/delete
	_, err = pool.Exec(ctx, "insert into users (name,email,password,created_at,updated_at) values "+
		"('test','test@gmail.com','test123',NOW(),NOW())")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("user inserted to database")

	// //query untuk mengambil banyak row
	var users []entity.User
	rows, err := pool.Query(ctx, "SELECT id, name FROM users order by id desc")
	if err != nil {
		log.Fatalln(err)
	}
	for rows.Next() {
		var user entity.User
		rows.Scan(&user.ID, &user.Name)
		users = append(users, user)
	}
	fmt.Println("users retrieved from database", users)

}
