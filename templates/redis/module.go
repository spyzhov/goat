package redis

import "github.com/spyzhov/goat/templates"

func New() *templates.Template {
	return &templates.Template{
		ID:      "redis",
		Name:    "Redis",
		Package: "github.com/gomodule/redigo/redis",

		Environments: []*templates.Environment{
			{Name: "RedisConnect", Type: "string", Env: "REDIS_CONNECTION", Default: "localhost:6379"},
			{Name: "RedisDatabase", Type: "int", Env: "REDIS_DATABASE", Default: "0"},
			{Name: "RedisIdleConnections", Type: "int", Env: "REDIS_IDLE_CONNECTIONS", Default: "2"},
		},
		Properties: []*templates.Property{
			{Name: "Redis", Type: "*redis.Pool"},
		},
		Libraries: []*templates.Library{
			{Name: "github.com/gomodule/redigo/redis"},
		},
		Models: map[string]string{},

		TemplateSetter: func(config *templates.Config) (s string) {
			s = `
	if err = app.setRedis(); err != nil {
		logger.Panic("cannot connect to Redis", zap.Error(err))
		return nil, err
	}`
			return
		},
		TemplateSetterFunction: func(config *templates.Config) (s string) {
			s = `
// Redis connect
func (app *Application) setRedis() (err error) {
	app.Logger.Debug("Redis connect")
	app.Redis = &redis.Pool{
		MaxIdle:         app.Config.RedisIdleConnections,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", app.Config.RedisConnect, redis.DialDatabase(app.Config.RedisDatabase))
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
	conn := app.Redis.Get()
	defer app.Closer(conn, "Redis connection")
	_, err = conn.Do("PING")
	return
}`
			return
		},
		TemplateRunFunction: templates.BlankFunction,
		TemplateClosers: func(*templates.Config) (s string) {
			s = `
	defer app.Closer(app.Redis, "Redis connection")`
			return
		},

		Templates: templates.BlankFunctionMap,
	}
}
