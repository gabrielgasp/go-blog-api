package user

import "github.com/gin-gonic/gin"

func LoadUserRoutes(app *gin.RouterGroup) {
	app.POST("/signup", signupHandler)
	app.POST("/login", loginHandler)
	app.GET("/list", listUsersHandler)
	app.GET("/:id", getUserByIdHandler)
}
