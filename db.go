package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

var dbCon *sql.DB

func conDB() {
	_ = mysql.Config{
		User:   "root",
		Passwd: "123",
		Net:    "tcp",
		Addr:   "localhost:3306",
		DBName: "recordings",
	}
	var err error
	dbCon, err = sql.Open("mysql", "root:123@/recordings")

	if err != nil {
		log.Fatal(err)
	}

	pingErr := dbCon.Ping()

	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Printf("DB Connected\n")

}

func allAlbums() ([]album, error) {
	var albums []album

	rows, err := dbCon.Query("SELECT * FROM album ")
	if err != nil {
		return nil, fmt.Errorf("allAlbums  %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var alb album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("allAlbums %v", err)
		}

		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("allAlbums %v", err)
	}

	return albums, nil
}

func albumsByArtist(name string) ([]album, error) {
	var albums []album

	rows, err := dbCon.Query("SELECT * FROM album WHERE artist = ?", name)
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}

	defer rows.Close()

	for rows.Next() {
		var alb album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
		}

		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}

	return albums, nil
}

func albumById(id int64) (album, error) {
	var alb album

	row := dbCon.QueryRow("SELECT * FROM album WHERE id = ?", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("albumsById %d: no such album", id)
		}
		return alb, fmt.Errorf("albumsById %d: %v", id, err)
	}
	return alb, nil
}

func addAlbum(alb album) (int64, error) {
	result, err := dbCon.Exec("INSERT INTO album(title,artist,price) VALUES (?,?,?)", alb.Title, alb.Artist, alb.Price)
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}

	return id, nil
}
