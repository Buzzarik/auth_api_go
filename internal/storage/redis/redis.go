package redis

import (
	"auth/internal/config"
	"context"
	"fmt"
	"time"
	rd "github.com/go-redis/redis/v8"
)

type Redis struct {
	client *rd.Client
};

func New(cnf config.ConfigRedis) (*Redis, error) {
	const op = "Redis.New";

	client := rd.NewClient(&rd.Options{
        Addr:     fmt.Sprintf("%s:%d", cnf.Host, cnf.Port),
        Password: cnf.Password,
        DB:       cnf.Db,
    });

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

	//NOTE: проверка на подключение
    _, err := client.Ping(ctx).Result();
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err);
    }

	return &Redis{client: client}, nil;
}

func (s *Redis) SetUser(phoneNumber string, userData interface{}) error{
	const op = "Redis.SetUser";
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel();

	//запись в redis
	err := s.client.HSet(ctx, phoneNumber, userData).Err();

	if err != nil {
		return fmt.Errorf("(HSet)%s: %w", op, err);
	}
	//установка тайминга на удержание записи
	err = s.client.Expire(ctx, phoneNumber, time.Minute*5).Err();
	if err != nil {
		return fmt.Errorf("%s: %w", op, err);
	}

	return nil;
}

func (s *Redis) GetUser(phoneNumber string) (map[string]string, error){
	const op = "Redis.GetUser";
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel();

	//запись в redis
	userData, err := s.client.HGetAll(ctx, phoneNumber).Result();
	
	if err == rd.Nil {
		return nil, nil;
	}

	if err != nil {
		return nil, fmt.Errorf("(HGetAll)%s: %w", op, err);
	}

	return userData, nil;
}