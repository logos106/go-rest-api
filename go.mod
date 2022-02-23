module github.com/saroopmathur/rest-api

go 1.17

require (
	github.com/0xAX/notificator v0.0.0-20210731104411-c42e3d4a43ee // indirect
	github.com/codegangsta/envy v0.0.0-20141216192214-4b78388c8ce4 // indirect
	github.com/codegangsta/gin v0.0.0-20211113050330-71f90109db02 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.0-20190314233015-f79a8a8ca69d // indirect
	github.com/golang/gddo v0.0.0-20210115222349-20d68f94ee1f // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/joho/godotenv v1.4.0 // indirect
	github.com/lib/pq v1.10.4 // indirect
	github.com/mattn/go-shellwords v1.0.12 // indirect
	github.com/rs/cors v1.8.0 // indirect
	github.com/russross/blackfriday/v2 v2.0.1 // indirect
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	github.com/urfave/cli v1.22.5 // indirect
	golang.org/x/crypto v0.0.0-20211215153901-e495a2d5b3d3 // indirect
)

replace github.com/saroopmathur/rest-api/router => ./router

replace github.com/saroopmathur/rest-api/models => ./models

replace github.com/saroopmathur/rest-api/db => ./db

replace github.com/saroopmathur/rest-api/middleware => ./middleware

replace github.com/saroopmathur/rest-api/handlers => ./handlers
