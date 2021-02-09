package param

import (
	"context"
	"errors"

	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"

	"bcdpkg.in/go-project/pkg/hasura"
)

var (
	keyInCtx = "param:ctx:key"
)

type Database struct {
	DSN string
	Db  *gorm.DB
}

type OIDC struct {
	Addr   string
	Client *oidc.Provider
}

type Param struct {
	ServiceName    string
	HTTPListenAddr string
	CorsHosts      []string
	Client         Database
	GraphQLClient  *hasura.GqlClient
	Log            *logrus.Entry
	OIDC           OIDC
	Verbose        bool
}

func Inject(e *gin.Engine, p *Param) {
	e.Use(func(ctx *gin.Context) {
		ctx.Set(keyInCtx, p)
		c := context.WithValue(ctx.Request.Context(), keyInCtx, ctx)
		ctx.Request = ctx.Request.WithContext(c)
		ctx.Next()
	})
}

func Eject(ctx *gin.Context) *Param {
	v, _ := ctx.Get(keyInCtx)
	return v.(*Param)
}

func EjectParamForResolver(ctx context.Context) (*Param, error) {
	c := ctx.Value(keyInCtx)
	gCtx, ok := c.(*gin.Context)
	if !ok {
		err := errors.New("gin.Context has wrong type")
		return nil, err
	}
	v, _ := gCtx.Get(keyInCtx)
	return v.(*Param), nil
}
