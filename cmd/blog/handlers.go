package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type indexPage struct {
	Title         string
	FeaturedPosts []featuredPostData
	RecentPosts   []recentPostData
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
