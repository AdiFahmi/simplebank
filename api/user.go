package api

import (
	"log"
	"net/http"
	"time"

	"github.com/VividCortex/mysqlerr"
	db "github.com/adifahmi/simplebank/db/sqlc"
	"github.com/adifahmi/simplebank/util"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type createUserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPass, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPass,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	res, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if sqlErr, ok := err.(*mysql.MySQLError); ok {
			log.Printf("Err %s with number %d", err.Error(), sqlErr.Number)
			if sqlErr.Number == mysqlerr.ER_DUP_ENTRY {
				ctx.JSON(http.StatusBadRequest, errorStringResponse("Username/email already exists"))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	createdUserID, err := res.LastInsertId()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	user, err := server.store.GetUserByID(ctx, createdUserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := createUserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		CreatedAt:         user.CreatedAt,
		PasswordChangedAt: user.PasswordChangedAt,
	}

	ctx.JSON(http.StatusOK, resp)
}
