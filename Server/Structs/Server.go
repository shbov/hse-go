package Structs

import (
	"Server/Funcs"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"math/rand"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

const shutdownTimeout = 20 * time.Second

type Server struct {
	config     *Config
	httpServer *http.Server
}

func NewServer(config *Config) *Server {
	return &Server{
		config: config,
	}
}

func (server *Server) GetHardOp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// имитируем работу сложного процесса
	timeout := Funcs.RandIntInInterval(10, 20)
	time.Sleep(time.Duration(timeout) * time.Second)

	// рандомим код возврата
	if rand.Intn(2)%2 == 0 {
		response := Response{Message: "Error while processing function (rand % 2 == 0)", Code: http.StatusInternalServerError}
		response.Json(w)
		return
	}

	response := Response{Message: "Ok", Code: http.StatusOK}
	response.Json(w)
}

func (server *Server) PostDecode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var request Request

	err := decoder.Decode(&request)
	if err != nil {
		response := Response{
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		}

		response.Json(w)
		return
	}

	rawDecodedText, err := base64.StdEncoding.DecodeString(request.Message)
	if err != nil {
		response := Response{
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		}

		response.Json(w)
		return
	}

	response := Response{
		Message: string(rawDecodedText),
		Code:    http.StatusOK,
	}
	response.Json(w)
}

func (server *Server) GetVersion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// сделал общий формат ответа (outputString)
	// можно переименовать в message и будет совсем ок
	// в задании просто сказано вернуть версию
	response := Response{
		Message: server.config.Version.String(),
		Code:    http.StatusOK,
	}
	response.Json(w)
}

func (server *Server) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/version", server.GetVersion)
	mux.HandleFunc("/decode", server.PostDecode)
	mux.HandleFunc("/hard-op", server.GetHardOp)

	server.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", server.config.Port),
		Handler: mux,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	group, ctx := errgroup.WithContext(ctx)
	group.Go(func() error {
		if err := server.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}

		return nil
	})

	group.Go(func() error {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		err := server.httpServer.Shutdown(shutdownCtx)
		if err != nil {
			return err
		}

		return nil
	})

	err := group.Wait()
	if err != nil {
		log.Error(err)
	}
}
