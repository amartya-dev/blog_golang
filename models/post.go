package models

import (
	"os"
)

// Post represents a blog post
type Post struct {
	title       string
	description string
	image       os.File
}
