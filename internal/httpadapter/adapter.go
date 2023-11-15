package httpadapter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/toshi0607/chi-prometheus"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc"
	"log"

	"github.com/juju/zaputil/zapctx"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/riandyrn/otelchi"
	"go.opentelemetry.io/otel"
	"net/http"
	"time"

	"github.com/anonimpopov/hw4/internal/docs" // go:generate
	"github.com/anonimpopov/hw4/internal/model"
	"github.com/anonimpopov/hw4/internal/service"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"moul.io/chizap"
)

var (
	TracerName = "demo_service"
	tracer     = otel.Tracer(TracerName)
)

// @title Auth API
// @version 1.0
// @description This is a simple auth server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:9000
// @BasePath /api/v1
// @query.collection.format multi

// @securityDefinitions.basic BasicAuth

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @securitydefinitions.oauth2.application OAuth2Application
// @tokenUrl https://example.com/oauth/token
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.implicit OAuth2Implicit
// @authorizationurl https://example.com/oauthorize
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.password OAuth2Password
// @tokenUrl /v1/login
// @scope.read Grants read access
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.accessCode OAuth2AccessCode
// @tokenUrl https://example.com/oauth/token
// @authorizationurl https://example.com/oauthorize
// @scope.admin Grants read and write access to administrative information

// @x-extension-openapi {"example": "value on a json format"}

type adapter struct {
	config      *Config
	authService service.Auth
	server      *http.Server
}

func testFunc(ctx context.Context) {
	_, span := tracer.Start(ctx, "test func")
	defer span.End()
}

func (a *adapter) Error(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "error")
	defer span.End()

	testFunc(r.Context())
	err := errors.New("oops... one more error")
	if err != nil {
		writeError(w, err)
	}
}

// Login Auth godoc
// @Summary authorize login and password
// @Description authorize user by login and password
// @Accept json
// @Param credentials body Credentials{} false "user credentials"
// @Success 200 {object} model.TokenPair
// @Failure 403 {object} Error
// @Failure 404 {object} Error
// @Failure 500 {object} Error
// @Router /login [post]
func (a *adapter) Login(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "login")
	defer span.End()

	var credentials Credentials

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		writeError(w, err)
		return
	}

	tokenPair, err := a.authService.Login(r.Context(), credentials.Login, credentials.Password)
	if err != nil {
		writeError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     a.config.AccessTokenCookie,
		Value:    tokenPair.AccessToken,
		Path:     "/",
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     a.config.RefreshTokenCookie,
		Value:    tokenPair.RefreshToken,
		Path:     "/",
		HttpOnly: true,
	})

	writeJSONResponse(w, http.StatusOK, tokenPair)
}

// Auth godoc
// @Summary validate authorization
// @Description validate authorization
// @Success 200 {object} model.TokenPair
// @Failure 403 {object} Error
// @Failure 500 {object} Error
// @Router /validate [post]
func (a *adapter) Validate(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "validate")
	defer span.End()

	accessToken, err := r.Cookie(a.config.AccessTokenCookie)
	if err != nil {
		writeError(w, fmt.Errorf("%w: %s", service.ErrForbidden, err))
		return
	}

	refreshToken, _ := r.Cookie(a.config.RefreshTokenCookie)
	if err != nil {
		writeError(w, fmt.Errorf("%w: %s", service.ErrForbidden, err))
		return
	}

	tokenPair := &model.TokenPair{
		AccessToken:  accessToken.Value,
		RefreshToken: refreshToken.Value,
	}

	tokenPair, err = a.authService.ValidateAndRefresh(r.Context(), tokenPair)
	if err != nil {
		writeError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     a.config.AccessTokenCookie,
		Value:    tokenPair.AccessToken,
		Path:     "/",
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     a.config.RefreshTokenCookie,
		Value:    tokenPair.RefreshToken,
		Path:     "/",
		HttpOnly: true,
	})

	writeJSONResponse(w, http.StatusOK, tokenPair)
}

// Auth godoc
// @Summary logout user
// @Description logout user
// @Success 200
// @Router /logout [post]
func (a *adapter) Logout(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "logout")
	defer span.End()

	http.SetCookie(w, &http.Cookie{
		Name:     a.config.AccessTokenCookie,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
	})

	http.SetCookie(w, &http.Cookie{
		Name:     a.config.RefreshTokenCookie,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
	})
}

func (a *adapter) Serve(ctx context.Context) error {
	lg := zapctx.Logger(ctx)

	shut := initTracerProvider()
	defer shut()

	r := chi.NewRouter()
	apiRouter := chi.NewRouter()

	apiRouter.Use(otelchi.Middleware("my-server", otelchi.WithChiRoutes(r)))

	m := chiprometheus.New("main")
	m.MustRegisterDefault()
	apiRouter.Use(m.Handler)

	apiRouter.Use(chizap.New(lg, &chizap.Opts{
		WithReferer:   true,
		WithUserAgent: true,
	}))

	apiRouter.Post("/login", http.HandlerFunc(a.Login))
	apiRouter.Post("/validate", http.HandlerFunc(a.Validate))
	apiRouter.Post("/logout", http.HandlerFunc(a.Logout))
	apiRouter.Get("/error", http.HandlerFunc(a.Error))

	// установка маршрута для документации
	// Адрес, по которому будет доступен doc.json
	apiRouter.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(fmt.Sprintf("%s/swagger/doc.json", a.config.BasePath))))

	r.Mount(a.config.BasePath, apiRouter)

	a.server = &http.Server{Addr: a.config.ServeAddress, Handler: r}

	if a.config.UseTLS {
		return a.server.ListenAndServeTLS(a.config.TLSCrtFile, a.config.TLSKeyFile)
	}

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		err := http.ListenAndServe(":9000", nil)
		if err != nil {
			lg.Fatal(err.Error())
		}
	}()

	return a.server.ListenAndServe()
}

func (a *adapter) Shutdown(ctx context.Context) {
	_ = a.server.Shutdown(ctx)
}

func New(
	config *Config,
	authorizer service.Auth) Adapter {

	if config.SwaggerAddress != "" {
		docs.SwaggerInfo.Host = config.SwaggerAddress
	} else {
		docs.SwaggerInfo.Host = config.ServeAddress
	}

	docs.SwaggerInfo.BasePath = config.BasePath

	return &adapter{
		config:      config,
		authService: authorizer,
	}
}

func initTracerProvider() func() {
	ctx := context.Background()
	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String("task 4"),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	traceClient := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint("docker.for.mac.host.internal:4317"),
		otlptracegrpc.WithDialOption(grpc.WithBlock()))
	sctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	traceExp, err := otlptrace.New(sctx, traceClient)
	if err != nil {
		log.Fatal(err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExp)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(tracerProvider)

	return func() {
		cxt, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		if err := traceExp.Shutdown(cxt); err != nil {
			otel.Handle(err)
		}
	}
}
