package main

import (
	web "github.com/iamlockon/shortu/web/api"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	web.Run()
}
