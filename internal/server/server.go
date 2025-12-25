package server

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ne241099/daifugo-server/graph"
	"github.com/ne241099/daifugo-server/internal/sse"
)

// New は設定済みの Echo サーバーインスタンスを返します
// 必要な依存関係（ResolverやHub）は引数として受け取ります
func New(resolver *graph.Resolver, hub *sse.Hub) *echo.Echo {
	e := echo.New()

	// 1. ミドルウェアの設定
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS()) // フロントエンドとの通信用

	// 2. GraphQL サーバーの設定
	gqlServer := handler.NewDefaultServer(
		graph.NewExecutableSchema(
			graph.Config{Resolvers: resolver},
		),
	)

	// 3. ルーティングの定義

	// GraphQL エンドポイント
	e.POST("/query", func(c echo.Context) error {
		gqlServer.ServeHTTP(c.Response(), c.Request())
		return nil
	})
	e.GET("/query", func(c echo.Context) error {
		gqlServer.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	// Playground（ブラウザ確認用）
	e.GET("/", func(c echo.Context) error {
		playground.Handler("GraphQL Playground", "/query").ServeHTTP(c.Response(), c.Request())
		return nil
	})

	// SSE エンドポイント
	e.GET("/events", sse.NewHandler(hub))

	return e
}
