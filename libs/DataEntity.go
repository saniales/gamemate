package libs

import "github.com/garyburd/redigo/redis" //Package to interact with Redis DB

type DataEntity interface {
    InsertIntoCache(pool *redis.Pool) error
    DeleteFromCache(pool *redis.Pool) error
}
