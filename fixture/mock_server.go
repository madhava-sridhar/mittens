package fixture

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/test/grpc_testing"
)

type PathResponseHandler struct {
	Path string
	PathHandlerFunc func(rw http.ResponseWriter, r *http.Request) 
}

// StartGrpcTargetTestServer starts a gRPC server on provided port
// It uses the test.proto from grpc-testing: https://github.com/grpc/grpc-go/blob/40a879c23a0dc77234d17e0699d074d5fd151bd0/test/grpc_testing/test.proto
func StartGrpcTargetTestServer(port int) *grpc.Server {
	svr := grpc.NewServer()
	grpc_testing.RegisterTestServiceServer(svr, &grpc_testing.UnimplementedTestServiceServer{})
	reflection.Register(svr)
	uri := ":" + fmt.Sprint(port)
	l, _ := net.Listen("tcp", uri)
	go func(){
		err := svr.Serve(l)
		if(err != nil ){
			log.Fatal("Server failed : " ,err)
		}
	}()

	return svr
}

func StartHttpTargetTestServer(port int, path_handlers []PathResponseHandler, disable_default_path bool) *http.Server {
	if(!disable_default_path) {
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		})

		http.HandleFunc("/delay", func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(time.Millisecond * 100)
			w.WriteHeader(http.StatusNoContent)
		})
	}

	for _, path_handler := range path_handlers {
		http.HandleFunc(path_handler.Path, path_handler.PathHandlerFunc)
	}

	base_url := ":" + fmt.Sprint(port)
	server := &http.Server{Addr: base_url}

	go func(){
		err := server.ListenAndServe()
		if(err != nil ){
			log.Fatal("Server failed : " ,err)
		}
	}()
	return server
}