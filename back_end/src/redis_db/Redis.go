package redis_db

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"rabbit_gather/src/logger"
	"rabbit_gather/util"
	"time"
)

type DBConfig struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	ID       int    `json:"id"`
}

func init() {
	type Config struct {
		AllDatabase []DBConfig `json:"all_database"`
	}
	var config Config
	err := util.ParseJsonConfic(&config, "config/redis.config.json")
	if err != nil {
		panic(err.Error())
	}
	if config.AllDatabase == nil {
		panic("AllDatabase is nil")
	} else if len(config.AllDatabase) <= 0 {
		panic("AllDatabase is empty")
	}
	all_db_client = map[int]*ClientWrapper{}
	for _, dbConfig := range config.AllDatabase {
		rawClient := redis.NewClient(&redis.Options{
			Addr:     dbConfig.Addr,
			Password: dbConfig.Password,
			DB:       dbConfig.ID,
		})
		_, err := rawClient.Ping(context.Background()).Result()
		if err != nil {
			panic(err.Error())
		}
		all_db_client[dbConfig.ID] = &ClientWrapper{Client: *rawClient, DB_ID: dbConfig.ID}
	}
}

var all_db_client map[int]*ClientWrapper

type ClientWrapper struct {
	DB_ID int
	redis.Client
}

var log = logger.NewLoggerWrapper("ClientWrapper")

const DefaultCount = 100
const DefaultCursor = 0

func (c *ClientWrapper) DeleteAll(ctx context.Context, match string) error {
	pipe := c.Pipeline()
	defer func(pipe redis.Pipeliner) {
		err := pipe.Close()
		if err != nil {
			log.ERROR.Println("Error when Close pipe: ", err.Error())
		}
	}(pipe)
	defer func(pipe redis.Pipeliner, ctx context.Context) {
		_, err := pipe.Exec(ctx)
		if err != nil {
			log.ERROR.Println("Error when Exec pipe: ", err.Error())
		}
	}(pipe, ctx)

	iter := c.Client.Scan(ctx, DefaultCursor, match, DefaultCount).Iterator()
	for iter.Next(ctx) {
		pipe.Del(ctx, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return err
	}
	return nil
}

func GetClient(id int) *ClientWrapper {
	if cl, ok := all_db_client[id]; !ok {
		panic(fmt.Sprintf("Unknown db id: %d", id))
	} else {
		return cl
	}
}

func Close() error {
	var err error
	for i, wrapper := range all_db_client {
		e := wrapper.Close()
		if e != nil {
			if err == nil {
				err = fmt.Errorf("Error when close db: %d  %w", i, e)
			} else {
				err = fmt.Errorf("%s -> %w", err.Error(), fmt.Errorf("Error when close db: %d  %w", i, e))
			}
		}
	}
	return err
}

func (c *ClientWrapper) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	p, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = c.Client.Set(ctx, key, p, expiration).Result()
	return err
}

func (c *ClientWrapper) Get(ctx context.Context, key string, stk interface{}) error {
	p, err := c.Client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(p), stk)
}
