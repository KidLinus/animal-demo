package kingdom

import (
	"context"
	"fmt"
	"io"
	"runtime/debug"
	"time"

	"github.com/creasty/defaults"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func (api *API) Gin(allowOrigins []string) *gin.Engine {
	router := gin.New()
	router.MaxMultipartMemory = 15 << 20 // 15 MiB
	router.Use(api.middlewareLogger)
	router.Use(api.middlewareCrash)
	router.Use(api.middlewareSession)
	router.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Accept-Encoding", "Content-Length", "Content-Type", "Host", "Origin", "Access-Control-Allow-Credentials", "Access-Control-Allow-Headers", "Access-Control-Allow-Methods", "Access-Control-Allow-Origin"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "Content-Disposition", "Access-Control-Allow-Credentials", "Access-Control-Allow-Headers", "Access-Control-Allow-Methods", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.Handle("GET", "/animal", endpoint(api.AnimalSearch))
	router.Handle("GET", "/animal/:id", endpoint(api.AnimalGet))
	router.Handle("GET", "/animal/:id/coi", endpoint(api.AnimalGetCOI))
	router.Handle("GET", "/animal/:id/parents", endpoint(api.AnimalGetParentTree))
	router.Handle("GET", "/animal/coi", endpoint(api.AnimalsGetCOI))
	router.Handle("GET", "/animal/parents", endpoint(api.AnimalsGetParentTree))
	return router
}

func (api *API) middlewareLogger(ctx *gin.Context) {
	start := time.Now()
	ctx.Next()
	fmt.Printf("[%s] %d %s %s\n", time.Since(start), ctx.Writer.Status(), ctx.Request.Method, ctx.Request.URL.Path)
}

func (api *API) middlewareCrash(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("[CRASH] %s\n%s\n", err, debug.Stack())
		}
	}()
	ctx.Next()
}

func (api *API) middlewareSession(ctx *gin.Context) {
	var session context.Context = ctx
	ctx.Set("context", session)
}

func endpoint[IN, OUT any](fn func(context.Context, IN) (OUT, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		input := new(IN)
		defaults.MustSet(input)
		ctx.ShouldBindUri(input)
		ctx.ShouldBindQuery(input)
		ctx.ShouldBind(input)
		if err := binding.Validator.ValidateStruct(input); err != nil {
			errors := []gin.H{}
			for _, e := range err.(validator.ValidationErrors) {
				errors = append(errors, gin.H{"field": e.Field(), "validation": e.Tag(), "error": e.Error()})
			}
			ctx.JSON(400, gin.H{"errors": errors})
			return
		}
		output, err := fn(ctx, *input)
		if err != nil {
			ctx.String(500, err.Error())
			return
		}
		res := any(output)
		if dl, ok := res.(*download); ok {
			defer dl.Close()
			ctx.DataFromReader(200, int64(dl.size), dl.mimetype, dl, map[string]string{})
			return
		}
		ctx.JSON(200, output)
	}
}

type download struct {
	name     string
	size     int
	mimetype string
	io.ReadCloser
}
