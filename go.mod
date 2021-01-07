module github.com/bjjb/urleen

go 1.14

replace github.com/bjjb/urleen/base62 => ./base62

require (
	github.com/bjjb/urleen/base62 v0.0.0-20191221021209-70772eec5343
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/go-redis/redis/v8 v8.4.4
	github.com/rs/cors v1.7.0
)
