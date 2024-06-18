package http

import (
	"context"
	"fmt"
	"io"
	"runtime/debug"
	"time"

	"animal.dev/animal/internal/api"
	"github.com/creasty/defaults"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	Config
	router *gin.Engine
}

type Config struct {
	Addr         []string
	AllowOrigins []string
	API          *api.App
}

func New(cfg Config) (*Server, error) {
	server := &Server{Config: cfg, router: gin.New()}
	server.router.MaxMultipartMemory = 15 << 20 // 15 MiB
	server.router.Use(server.middlewareLogger)
	server.router.Use(server.middlewareCrash)
	server.router.Use(server.middlewareSession)
	server.router.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Accept-Encoding", "Content-Length", "Content-Type", "Host", "Origin", "Access-Control-Allow-Credentials", "Access-Control-Allow-Headers", "Access-Control-Allow-Methods", "Access-Control-Allow-Origin"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "Content-Disposition", "Access-Control-Allow-Credentials", "Access-Control-Allow-Headers", "Access-Control-Allow-Methods", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	server.router.Handle("GET", "/dataset", endpoint(server.datasetList))
	server.router.Handle("POST", "/dataset", endpoint(server.datasetCreate))
	server.router.Handle("GET", "/dataset/:id", endpoint(server.datasetGet))
	server.router.Handle("PATCH", "/dataset/:id", endpoint(server.datasetUpdate))
	go server.router.Run(cfg.Addr...)
	return server, nil
}

func (server *Server) Close() {}

func (server *Server) middlewareLogger(ctx *gin.Context) {
	start := time.Now()
	ctx.Next()
	fmt.Printf("[%s] %d %s %s\n", time.Since(start), ctx.Writer.Status(), ctx.Request.Method, ctx.Request.URL.Path)
}

func (server *Server) middlewareCrash(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("[CRASH] %s\n%s\n", err, debug.Stack())
		}
	}()
	ctx.Next()
}

func (server *Server) middlewareSession(ctx *gin.Context) {
	var session context.Context = ctx
	ctx.Set("context", session)
}

func endpoint[IN, OUT any](fn func(api.Context, IN) (OUT, error)) gin.HandlerFunc {
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
