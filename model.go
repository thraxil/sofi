package main

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type Image struct {
	Id         int       `db:"id"`
	AHash      string    `db:"ahash"`
	Ext        string    `db:"ext"`
	Url        string    `db:"url"`
	InsertedAt time.Time `db:"inserted_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

type Tag struct {
	Id         int       `db:"id"`
	Slug       string    `db:"slug"`
	Tag        string    `db:"tag"`
	InsertedAt time.Time `db:"inserted_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

func count_images(db *sqlx.DB) int {
	var count int
	db.Get(&count, "SELECT COUNT(*) FROM images")
	return count
}

func newest_images(db *sqlx.DB, page int) []Image {
	// page size is 20
	images := []Image{}
	db.Select(&images, "SELECT * FROM images ORDER BY id DESC LIMIT 20 OFFSET $1", (page-1)*20)
	return images
}

func random_image(db *sqlx.DB) Image {
	image := Image{}
	db.Get(&image, "SELECT * FROM images ORDER BY RANDOM() LIMIT 1")
	return image
}

func get_image_by_id(db *sqlx.DB, id int) Image {
	image := Image{}
	db.Get(&image, "SELECT * FROM images WHERE id=$1", id)
	return image
}

func getAllTags(db *sqlx.DB) []Tag {
	tags := []Tag{}
	db.Select(&tags, "SELECT * FROM tags order by lower(tag)")
	return tags
}

func get_tag_by_slug(db *sqlx.DB, slug string) Tag {
	tag := Tag{}
	db.Get(&tag, "SELECT * FROM tags WHERE slug=$1", slug)
	return tag
}

func get_images_by_tag(db *sqlx.DB, tag Tag) []Image {
	images := []Image{}
	db.Select(&images, "SELECT images.* FROM images JOIN image_tags ON images.id=image_tags.image_id WHERE image_tags.tag_id=$1", tag.Id)
	return images
}

func get_image_tags(db *sqlx.DB, image Image) []Tag {
	tags := []Tag{}
	db.Select(&tags, "SELECT tags.* FROM tags JOIN image_tags ON tags.id=image_tags.tag_id WHERE image_tags.image_id=$1", image.Id)
	return tags
}
