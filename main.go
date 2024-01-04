package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/davecgh/go-spew/spew"
	_ "github.com/mattn/go-sqlite3"
	"github.com/maxsei/bob-chinook/models/models"
	"github.com/stephenafamo/bob"
)

const DbFilename = "./Chinook_Sqlite.sqlite"

func DownloadDb() {
	resp, err := http.Get("https://raw.githubusercontent.com/lerocha/chinook-database/master/ChinookDatabase/DataSources/Chinook_Sqlite.sqlite")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		panic("non ok status on downloading database")
	}
	f, err := os.OpenFile(DbFilename, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if _, err := io.Copy(f, resp.Body); err != nil {
		os.Remove(DbFilename)
		return
	}
}

func main() {
	if _, err := os.Stat(DbFilename); errors.Is(err, os.ErrNotExist) {
		DownloadDb()
	}
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared", DbFilename))
	if err != nil {
		panic(err)
	}
	defer db.Close()
	{
		// Query to get all table names from sqlite_master
		rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table';")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		// Iterate through the result set and print table names
		fmt.Println("Table Names:")
		for rows.Next() {
			var tableName string
			err := rows.Scan(&tableName)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(tableName)
		}
	}

	bdb := bob.NewDB(db)
	// experiment1(bdb)
	// experiment1_5(bdb)
	// experiment1(bdb)
	// experiment2(bdb)
	experiment3(bdb)
}

func experiment1(bdb bob.DB) {
	ctx := context.Background()
	artist, err := models.FindArtist(ctx, bdb, 1)
	if err != nil {
		panic(err)
	}
	albums, err := artist.ArtistIdAlbums(ctx, bdb).All()
	if err != nil {
		panic(err)
	}
	spew.Dump(artist)
	spew.Dump(albums)
}

func experiment1_5(bdb bob.DB) {
	ctx := context.Background()
	// artist, err := models.FindArtist(ctx, bdb, 1)
	// if err != nil {
	// 	panic(err)
	// }
	artist := &models.Artist{
		ArtistId: 1,
	}
	albums, err := artist.ArtistIdAlbums(ctx, bdb).All()
	if err != nil {
		panic(err)
	}
	spew.Dump(artist)
	spew.Dump(albums)
}

func experiment2(bdb bob.DB) {
	ctx := context.Background()
	mod := models.ThenLoadArtistArtistIdAlbums()
	artists, err := models.Artists.Query(ctx, bdb, mod).One()
	if err != nil {
		panic(err)
	}
	spew.Dump(artists)
}

func experiment3(bdb bob.DB) {
	ctx := context.Background()
	mod := models.ThenLoadArtistArtistIdAlbums(models.ThenLoadAlbumAlbumIdTracks())
	artists, err := models.Artists.Query(ctx, bdb, mod).One()
	if err != nil {
		panic(err)
	}
	spew.Dump(artists)
}
