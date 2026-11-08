package redis

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"patients_categories/internal/app/config"

	"github.com/go-redis/redis/v8"
)

// servicePrefix — префикс сервиса, чтобы ключи в общем Redis
// не пересекались с ключами других приложений.
const servicePrefix = "patients_categories."

const jwtPrefix = "jwt."

type Client struct {
	client *redis.Client
}

func New(ctx context.Context, cfg config.RedisConfig) (*Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:        cfg.Host + ":" + strconv.Itoa(cfg.Port),
		Username:    cfg.User,
		Password:    cfg.Password,
		DB:          0,
		DialTimeout: cfg.DialTimeout,
		ReadTimeout: cfg.ReadTimeout,
	})

	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("cant ping redis: %w", err)
	}

	return &Client{client: redisClient}, nil
}

func (c *Client) Close() error {
	return c.client.Close()
}

func getJWTKey(token string) string {
	return servicePrefix + jwtPrefix + token
}

// WriteJWTToBlacklist сохраняет разлогиненный токен до момента истечения его срока.
func (c *Client) WriteJWTToBlacklist(ctx context.Context, jwtStr string, jwtTTL time.Duration) error {
	if jwtTTL <= 0 {
		return nil
	}
	return c.client.Set(ctx, getJWTKey(jwtStr), true, jwtTTL).Err()
}

// IsJWTBlacklisted сообщает, лежит ли токен в чёрном списке.
func (c *Client) IsJWTBlacklisted(ctx context.Context, jwtStr string) (bool, error) {
	err := c.client.Get(ctx, getJWTKey(jwtStr)).Err()
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
