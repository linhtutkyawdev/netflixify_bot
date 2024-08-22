package main

type Channel struct {
	ID       int64
	Title    string
	Password string
}

type Post struct {
	channel_id       int
	title            string
	rating           int
	description      string
	tags             string
	video_id         string
	video_path       string
	thumbnail_id     string
	thumbnail_path   string
	g_thumbnail_id   string
	g_thumbnail_path string
}
