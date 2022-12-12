package server

import (
	"net/http"
	"os"

	static "github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"

	"github.com/gin-contrib/gzip"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/stenic/ledger/graph"
	"github.com/stenic/ledger/graph/generated"
	auth "github.com/stenic/ledger/internal/auth"
	"github.com/stenic/ledger/internal/storage"
)

func NewServer(opts ServerOpts) *Server {
	return &Server{
		ServerOpts: opts,
	}
}

type Server struct {
	ServerOpts

	server *gin.Engine
}

func (s *Server) Listen(addr string) error {
	gin.SetMode(gin.ReleaseMode)

	engine := storage.Database{}
	engine.InitDB()
	defer engine.CloseDB()
	engine.Migrate()

	s.server = gin.New()
	logrus.Debugf("Serving static from %s", s.StaticAssetPath)
	s.server.Use(
		ginlogrus.Logger(&logrus.Logger{}, "/healthz"),
		gzip.Gzip(gzip.DefaultCompression),
		static.Serve("/", static.LocalFile(s.StaticAssetPath, true)),
		gin.Recovery(),
	)
	s.server.NoRoute(func(c *gin.Context) {
		c.File(s.StaticAssetPath + "/index.html")
		//default 404 page not found
	})

	s.server.GET("/healthz", healthHandler)
	oidcOptions := auth.ApiSecurityOptions{
		IssuerURL: s.ServerOpts.OidcIssuerURL,
		ClientID:  s.ServerOpts.OidcClientID,
	}

	if os.Getenv("DEBUG") == "true" {
		s.server.Any("/playground", gin.WrapF(playground.Handler("GraphQL playground", "/query")))
	}
	s.server.Any(
		"/query",
		auth.JwtHandler(oidcOptions),
		gin.WrapH(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))),
	)

	s.server.GET("/auth/config", func(c *gin.Context) {
		c.JSON(http.StatusOK, oidcOptions)
	})

	logrus.Info("Starting webserver")
	return s.server.Run(addr)
}
