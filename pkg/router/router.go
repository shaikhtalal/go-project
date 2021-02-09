package router

import (
	tResolver "bcdpkg.in/go-project/pkg/todo/graph"
	tGenerated "bcdpkg.in/go-project/pkg/todo/graph/generated"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
)

//

func HandleHTTP(e *gin.Engine) {

	e.POST("/todo", TodoHandler())

	e.GET("/", PlaygroundHandler())

}

func PlaygroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL playground", "/query")
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func TodoHandler() gin.HandlerFunc {
	srv := handler.NewDefaultServer(tGenerated.NewExecutableSchema(tGenerated.Config{Resolvers: &tResolver.Resolver{}}))
	return func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
	}
}
