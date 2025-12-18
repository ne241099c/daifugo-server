package main

import (
	"github.com/ne241099/daifugo-server/graph"
	"github.com/ne241099/daifugo-server/infra/inmem"
	"github.com/ne241099/daifugo-server/internal/sse"
	"github.com/ne241099/daifugo-server/usecase/user"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// middleware
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	// リポジトリ初期化
	userRepo := inmem.NewInmemUserRepository()

	// SSE Hub 作成
	hub := sse.NewHub()
	// SSE
	e.GET("/events", sse.NewHandler(hub)) // AIに聞く

	// GraphQL server
	gqlServer := handler.NewDefaultServer(
		graph.NewExecutableSchema(
			graph.Config{
				Resolvers: &graph.Resolver{
					Hub: hub,
					SignUpUseCase: &user.SignUpInteractor{
						UserRepository: userRepo,
					},
				},
			},
		),
	)

	// GraphQL エンドポイント
	e.POST("/query", func(c echo.Context) error {
		gqlServer.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	// Playground（開発用）
	e.GET("/playground", func(c echo.Context) error {
		playground.Handler("GraphQL Playground", "/query").
			ServeHTTP(c.Response(), c.Request())
		return nil
	})

	// サーバ起動
	e.Logger.Fatal(e.Start(":8080"))
}
