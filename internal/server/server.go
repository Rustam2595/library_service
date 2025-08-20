package server

type Server struct {
	Host   string
	Router *gin.Engine
}

func New(host string) *Server {
	r := gin.Default()
	userGroup := r.Group("/user")
	{
		userGroup.POST("/register")
		userGroup.POST("/auth")
		userGroup.GET("/")
	}
	bookGroup := r.Group("/book")
	{
		bookGroup.GET("/all-books")
		bookGroup.GET("/:id")
		bookGroup.POST("/add-book")
		bookGroup.DELETE("/delete/:id")
	}
	return &Server{
		Host:   host,
		Router: r}
}
func (s *Server) Run() error {
	if err := s.Router.Run(s.Host); err != nil {
		return err
	}
	return nil
}
