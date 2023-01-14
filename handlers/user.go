package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gabrielgaspar447/go-blog-api/errs"
	"github.com/gabrielgaspar447/go-blog-api/models"
	"github.com/gabrielgaspar447/go-blog-api/services"
	"github.com/gabrielgaspar447/go-blog-api/utils"
	"github.com/gin-gonic/gin"
)

func UserSignup(c *gin.Context) {
	var input models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		handleUserErrors(c, err)
		return
	}

	token, err := services.UserSignup(&input)

	if err != nil {
		handleUserErrors(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": token})
}

func UserLogin(c *gin.Context) {
	var input models.LoginDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		handleUserErrors(c, err)
		return
	}

	token, err := services.UserLogin(&input)

	if err != nil {
		handleUserErrors(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": token})
}

func UserList(c *gin.Context) {
	includePosts := c.Query("posts") == "true"

	var users []models.User

	err := services.UserList(&users, includePosts)

	if err != nil {
		handleUserErrors(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
}

func UserGetById(c *gin.Context) {
	includePosts := c.Query("posts") == "true"
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errs.ErrInvalidId})
		return
	}

	var user models.User

	err = services.UserGetById(&user, uint(id), includePosts)

	if err != nil {
		handleUserErrors(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func UserDeleteSelf(c *gin.Context) {
	id := c.GetUint("userId")

	err := services.UserDeleteSelf(id)

	if err != nil {
		handleUserErrors(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func handleUserErrors(c *gin.Context, err error) {
	if valErrs := utils.GetValidationErrors(err); valErrs != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": valErrs})
		return
	}

	switch err {
	case errs.ErrUserAlreadyExists:
		c.JSON(http.StatusConflict, gin.H{"error": errs.ErrUserAlreadyExists.Error()})
		return
	case errs.ErrUserNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": errs.ErrUserNotFound.Error()})
		return
	case errs.ErrInvalidPassword:
		fmt.Println("invalid password")
		c.JSON(http.StatusUnauthorized, gin.H{"error": errs.ErrInvalidPassword.Error()})
		return
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": errs.ErrUnknown.Error()})
		return
	}
}