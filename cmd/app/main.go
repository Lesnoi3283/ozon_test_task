package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"ozon_test_task/cfg"
	"ozon_test_task/internal/app/graph"
	"ozon_test_task/internal/app/graph/resolvers"
	"ozon_test_task/internal/app/middlewares"
	"ozon_test_task/pkg/authUtils"
	"ozon_test_task/pkg/database"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"
	"github.com/vektah/gqlparser/v2/ast"
	"go.uber.org/zap"
)

func main() {
	//load .env
	_ = godotenv.Load()

	//config set
	conf, err := cfg.Configure()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to configure: %w", err))
	}

	//logger set
	zCfg := zap.NewProductionConfig()
	level, err := zap.ParseAtomicLevel(conf.LogLevel)
	if err != nil {
		log.Fatalf("Cant parse log level, err: %v", err)
	}
	zCfg.Level = level
	zCfg.DisableStacktrace = true
	logger, err := zCfg.Build()
	if err != nil {
		log.Fatalf("logger was not started, err: %v", err)
	}
	sugar := logger.Sugar()

	resolver := &resolvers.Resolver{
		Cfg:    *conf,
		Logger: sugar,
	}

	//db set
	if conf.InMemoryStorage {
		sugar.Infof("Using in-memory storage")

		redis := redis.NewClient(&redis.Options{
			Addr:     conf.RedisAddress + ":" + conf.RedisPort,
			Password: conf.RedisPassword,
		})

		//check connection
		ctx := context.Background()
		err = redis.Ping(ctx).Err()
		if err != nil {
			sugar.Fatalf("Failed to connect to Redis: %v", err)
		}

		redisStorage := database.NewRepoRedis(redis)
		resolver.PostRepo = redisStorage
		resolver.UserRepo = redisStorage
		resolver.CommentRepo = redisStorage
	} else {
		sugar.Infof("Using database")

		postgres, err := sql.Open("postgres", conf.DBConnectionString)
		if err != nil {
			sugar.Fatalf("Failed to connect to postgres: %v", err)
		}
		postgresStorage := database.NewRepoPG(postgres)
		err = postgresStorage.InitDB()
		if err != nil {
			sugar.Fatalf("Failed to initialize postgres storage: %v", err)
		}
		resolver.PostRepo = postgresStorage
		resolver.UserRepo = postgresStorage
		resolver.CommentRepo = postgresStorage
	}

	//jwt manager set
	resolver.JWTManager = &authUtils.JWTHelper{}

	//middlewares set
	authMW := middlewares.GetAuthMiddleware(&authUtils.JWTHelper{}, resolver.UserRepo, sugar)

	//build GraphQL server
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	if conf.DebugMode {
		sugar.Infof("Debug mode enabled, playground is on")
		http.Handle("/playground", playground.Handler("GraphQL", "/query"))
	}
	http.Handle("/query", authMW(srv))

	sugar.Infof("connect to %s:%s for GraphQL", conf.ServerAddress, conf.ServerPort)
	sugar.Fatal(http.ListenAndServe(conf.ServerAddress+":"+conf.ServerPort, nil))
}
