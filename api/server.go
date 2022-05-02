package api

import (
	"fmt"

	db "github.com/adifahmi/simplebank/db/sqlc"
	"github.com/adifahmi/simplebank/token"
	"github.com/adifahmi/simplebank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	server.setUpRouter()
	return server, nil
}

func (server *Server) setUpRouter() {
	router := gin.Default()

	router.GET("/ping", server.ping)
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)

	router.POST("/transfers", server.createTransfer)

	router.POST("/users", server.createUser)
	router.POST("/login", server.loginuser)
	server.router = router
}

func (server *Server) Start(addr string) error {
	return server.router.Run(addr)
}

func (server *Server) ping(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "pong",
	})
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}

func errorStringResponse(err string) gin.H {
	return gin.H{
		"error": err,
	}
}
