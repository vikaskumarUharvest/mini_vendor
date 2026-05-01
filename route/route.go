package route

import (
	"pgxpostgress/handler"
	"pgxpostgress/repository/postgres"
	"pgxpostgress/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupRoutes(
	r *gin.Engine,
	pool *pgxpool.Pool,
) {

	repo := postgres.NewUserRepository(pool)

	svc := service.NewUserService(repo)

	h := handler.NewUserHandler(svc)

	api := r.Group("/api")
	{
		users := api.Group("/users")
		{
			users.GET("", h.List)
			users.POST("", h.Create)
			users.GET("/:id", h.Get)
			users.PUT("/:id", h.Update)
			users.DELETE("/:id", h.Delete)
		}
	}
}