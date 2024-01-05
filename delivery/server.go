package delivery

import (
	"database/sql"
	"fmt"

	"catatan-keuangan-app/config"
	"catatan-keuangan-app/delivery/controller"
	"catatan-keuangan-app/delivery/middleware"
	"catatan-keuangan-app/repository"
	"catatan-keuangan-app/shared/service"
	"catatan-keuangan-app/usecase"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

type Server struct {
	userUC     usecase.UserUseCase
	authUC     usecase.AuthUseCase
	expenseUC  usecase.ExpenseUseCase
	engine     *gin.Engine
	jwtService service.JwtService
	host       string
}

func (s *Server) initRoute() {
	rg := s.engine.Group(config.ApiGroup)
	authMid := middleware.NewAuthMiddleware(s.jwtService)
	controller.NewUserController(s.userUC, rg, authMid).Route()
	controller.NewAuthController(s.authUC, rg).Route()
	controller.NewExpenseController(s.expenseUC, rg, authMid).Route()
}

func (s *Server) Run() {
	s.initRoute()
	if err := s.engine.Run(s.host); err != nil {
		panic(fmt.Errorf("server not running on host %s, becauce error %v", s.host, err.Error()))
	}
}

func NewServer() *Server {
	cfg, _ := config.NewConfig()
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)
	db, err := sql.Open(cfg.Driver, dsn)
	if err != nil {
		panic("connection error")
	}
	userRepo := repository.NewUserRepository(db)
	expenseRepo := repository.NewExpenseRepository(db)

	userUC := usecase.NewUserUseCase(userRepo)
	taskUC := usecase.NewExpenseUseCase(expenseRepo)
	jwtService := service.NewJwtService(cfg.TokenConfig)
	authUC := usecase.NewAuthUseCase(userUC, jwtService)
	engine := gin.Default()
	host := fmt.Sprintf(":%s", cfg.ApiPort)
	return &Server{
		userUC:     userUC,
		authUC:     authUC,
		expenseUC:  taskUC,
		engine:     engine,
		host:       host,
		jwtService: jwtService,
	}
}
