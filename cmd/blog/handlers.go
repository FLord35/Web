package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type indexPage struct {
	Title         string
	FeaturedPosts []featuredPostData
	RecentPosts   []recentPostData
}

type loginPage struct {
	Title string
}

type adminPage struct {
	Title string
}

type postPage struct {
	Title string
	Post  []postData
}

type postData struct {
	Title    string `db:"title"`
	Subtitle string `db:"subtitle"`
	PostImg  string `db:"image_url"`
	Content  string `db:"content"`
}

type featuredPostData struct {
	PostID        string `db:"post_id"`
	Title         string `db:"title"`
	Subtitle      string `db:"subtitle"`
	BackgroundImg string `db:"image_url"`
	Author        string `db:"author"`
	AuthorImg     string `db:"author_url"`
	PublishDate   string `db:"publish_date"`
}

type recentPostData struct {
	PostID        string `db:"post_id"`
	Title         string `db:"title"`
	Subtitle      string `db:"subtitle"`
	BackgroundImg string `db:"image_url"`
	Author        string `db:"author"`
	AuthorImg     string `db:"author_url"`
	PublishDate   string `db:"publish_date"`
}

type newPostData struct {
	Title             string `json:"Title"`
	Subtitle          string `json:"Subtitle"`
	AuthorName        string `json:"AuthorName"`
	AuthorPhoto       string `json:"AuthorPhoto"`
	AuthorPhotoBase64 string `json:"AuthorPhotoBase64"`
	Date              string `json:"Date"`
	BigImage          string `json:"BigImage"`
	BigImageBase64    string `json:"BigImageBase64"`
	SmallImage        string `json:"SmallImage"`
	SmallImageBase64  string `json:"SmallImageBase64"`
	Content           string `json:"Content"`
}

func post(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		postIDStr := mux.Vars(r)["postID"]

		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			http.Error(w, "Invalid post id", http.StatusForbidden)
			log.Println(err)
			return
		}

		post, err := postByID(db, postID)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Order not found", 404)
				log.Println(err)
				return
			}

			http.Error(w, "Internal Server Error", 500)
			log.Println(err)
			return
		}

		ts, err := template.ParseFiles("pages/post.html")
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err)
			return
		}

		data := postPage{
			Title: "Escape",
			Post:  post,
		}

		err = ts.Execute(w, data)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err)
			return
		}

		log.Println("Request completed successfully")
	}
}

func index(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fPosts, err := featuredPosts(db)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err)
			return
		}

		rPosts, err := recentPosts(db)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err)
			return
		}

		ts, err := template.ParseFiles("pages/index.html")
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err.Error())
			return
		}

		data := indexPage{
			Title:         "Escape",
			FeaturedPosts: fPosts,
			RecentPosts:   rPosts,
		}

		err = ts.Execute(w, data)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err.Error())
			return
		}

		log.Println("Request completed successfully")
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("pages/login.html")
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}

	data := loginPage{
		Title: "Escape",
	}

	err = ts.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}

	log.Println("Request completed successfully")
}

func admin(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		ts, err := template.ParseFiles("pages/admin.html")
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err.Error())
			return
		}

		data := adminPage{
			Title: "Escape",
		}

		err = ts.Execute(w, data)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err.Error())
			return
		}

		log.Println("Request completed successfully")
	}
}

func featuredPosts(db *sqlx.DB) ([]featuredPostData, error) {
	const query = `
		SELECT
			post_id,
			title,
			subtitle,
			author,
			author_url,
			publish_date,
			image_url
		FROM
			post
		WHERE featured = 1
	`

	var fPosts []featuredPostData

	err := db.Select(&fPosts, query)
	if err != nil {
		return nil, err
	}

	return fPosts, nil
}

func recentPosts(db *sqlx.DB) ([]recentPostData, error) {
	const query = `
		SELECT
			post_id,
			title,
			subtitle,
			author,
			author_url,
			publish_date,
			image_url
		FROM
			post
		WHERE featured = 0
	`

	var rPosts []recentPostData

	err := db.Select(&rPosts, query)
	if err != nil {
		return nil, err
	}

	return rPosts, nil
}

func postByID(db *sqlx.DB, postID int) ([]postData, error) {
	const query = `
		SELECT
		    title,
		    subtitle,
		    image_url,
		    content
		FROM
			post
		WHERE post_id = ?
	`

	var post []postData

	err := db.Select(&post, query, postID)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func createPost(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		reqData, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading data", 500)
			log.Println(err.Error())
			return
		}

		var npd newPostData

		err = json.Unmarshal(reqData, &npd)
		if err != nil {
			http.Error(w, "Error decoding data", 500)
			log.Println(err.Error())
			return
		}

		APhotoB64Pure := npd.AuthorPhotoBase64[strings.IndexByte(npd.AuthorPhotoBase64, ',')+1:]
		APhotoB64, err := base64.StdEncoding.DecodeString(APhotoB64Pure)
		if err != nil {
			http.Error(w, "Error decoding image", 500)
			log.Println(err.Error())
			return
		}

		AuthorPhotoFile, err := os.Create("static/img/" + npd.AuthorPhoto)
		if err != nil {
			http.Error(w, "Error creating an image file", 500)
			fmt.Println(err.Error())
			return
		}

		_, err = AuthorPhotoFile.Write(APhotoB64)
		if err != nil {
			http.Error(w, "Error writing content to a file", 500)
			fmt.Println(err.Error())
			return
		}

		BImageB64Pure := npd.BigImageBase64[strings.IndexByte(npd.BigImageBase64, ',')+1:]
		BImageB64, err := base64.StdEncoding.DecodeString(BImageB64Pure)
		if err != nil {
			http.Error(w, "Error decoding image", 500)
			log.Println(err.Error())
			return
		}

		BigImageFile, err := os.Create("static/img/" + npd.BigImage)
		if err != nil {
			http.Error(w, "Error creating an image file", 500)
			fmt.Println(err.Error())
			return
		}

		_, err = BigImageFile.Write(BImageB64)
		if err != nil {
			http.Error(w, "Error writing content to a file", 500)
			fmt.Println(err.Error())
			return
		}

		SImageB64Pure := npd.SmallImageBase64[strings.IndexByte(npd.SmallImageBase64, ',')+1:]
		SImageB64, err := base64.StdEncoding.DecodeString(SImageB64Pure)
		if err != nil {
			http.Error(w, "Error decoding image", 500)
			log.Println(err.Error())
			return
		}

		SmallImageFile, err := os.Create("static/img/" + npd.SmallImage)
		if err != nil {
			http.Error(w, "Error creating an image file", 500)
			fmt.Println(err.Error())
			return
		}

		_, err = SmallImageFile.Write(SImageB64)
		if err != nil {
			http.Error(w, "Error writing content to a file", 500)
			fmt.Println(err.Error())
			return
		}

		err = writeToDB(db, npd)
		if err != nil {
			http.Error(w, "Error adding data to a database", 500)
			log.Println(err.Error())
			return
		}

		log.Println("Request completed successfully")
	}
}

func writeToDB(db *sqlx.DB, npd newPostData) error {
	const query = `
		INSERT INTO
			post
		(
			title,
			subtitle,
			author,
			author_url,
			publish_date,
			image_url,
			content
		)
		VALUES
		(
			?,
			?,
			?,
			CONCAT('../static/img/', ?),
			?,
			CONCAT('../static/img/', ?),
			?
		)
	`
	_, err := db.Exec(query, npd.Title, npd.Subtitle, npd.AuthorName, npd.AuthorPhoto, npd.Date, npd.BigImage, npd.Content)
	return err
}
