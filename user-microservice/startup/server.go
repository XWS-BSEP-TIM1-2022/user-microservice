package startup

import (
	"fmt"
	userService "github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/user"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/token"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/tracer"
	otgo "github.com/opentracing/opentracing-go"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"user-microservice/application"
	"user-microservice/infrastructure/api"
	"user-microservice/infrastructure/persistance"
	"user-microservice/model"
	"user-microservice/startup/config"
)

type Server struct {
	config     *config.Config
	tracer     otgo.Tracer
	closer     io.Closer
	jwtManager *token.JwtManager
}

func NewServer(config *config.Config) *Server {
	tracer, closer := tracer.Init(config.UserServiceName)
	otgo.SetGlobalTracer(tracer)
	jwtManager := token.NewJwtManagerDislinkt(config.ExpiresIn)
	return &Server{
		config:     config,
		tracer:     tracer,
		closer:     closer,
		jwtManager: jwtManager,
	}
}

func (server *Server) GetTracer() otgo.Tracer {
	return server.tracer
}

func (server *Server) GetCloser() io.Closer {
	return server.closer
}

func (server *Server) Start() {
	mongoClient := server.initMongoClient()
	userStore := server.initUserStore(mongoClient)
	userService := server.initUserService(userStore)
	authService := server.initAuthService(userStore)
	userHandler := server.initUserHandler(userService, authService)

	server.startGrpcServer(userHandler)
}

func (server *Server) initMongoClient() *mongo.Client {
	client, err := persistance.GetClient(server.config.UserDBHost, server.config.UserDBPort)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func (server *Server) startGrpcServer(userHandler *api.UserHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	log.Println(fmt.Sprintf("started grpc server on localhost:%s", server.config.Port))
	userService.RegisterUserServiceServer(grpcServer, userHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

func (server *Server) initUserStore(client *mongo.Client) model.UserStore {
	store := persistance.NewUserMongoDBStore(client)
	/*store.DeleteAll()
	for _, user := range users {
		_, err := store.Create(user)
		if err != nil {
			log.Fatal(err)
		}
	}*/
	return store
}

func (server *Server) initUserService(store model.UserStore) *application.UserService {
	return application.NewUserService(store)
}

func (server *Server) initUserHandler(service *application.UserService, authService *application.AuthService) *api.UserHandler {
	return api.NewUserHandler(service, authService)
}

func (server *Server) initAuthService(store model.UserStore) *application.AuthService {
	return application.NewAuthService(store, server.jwtManager)
}
