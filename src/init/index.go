package init

import (
	"api/src/dal/otp"
	"api/src/dal/user"
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	customerpb "stash.bms.bz/turf/generic-proto-files.git/customer/v1"
	facilitypb "stash.bms.bz/turf/generic-proto-files.git/facility/v1"
	merchantpb "stash.bms.bz/turf/generic-proto-files.git/merchant/v1"
	venuepb "stash.bms.bz/turf/generic-proto-files.git/venue/v1"

	"google.golang.org/grpc/reflection"
	"stash.bms.bz/bms/gologger.git"

	facilityService "api/src/service/facility"
	merchantService "api/src/service/merchant"
	otpService "api/src/service/otp"
	uploadService "api/src/service/upload"
	userService "api/src/service/user"
	venueService "api/src/service/venue"

	facilitydal "api/src/dal/facility"
	merchantdal "api/src/dal/merchant"
	userdal "api/src/dal/user"
	venuedal "api/src/dal/venue"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"api/src/dal"
	io "api/src/models"
	facility "api/src/modules/facilities"
	merchant "api/src/modules/merchants"
	venue "api/src/modules/venues"
	utile "api/src/utils"
)

const grpcPort = ":5001"
const httpPort = ":80"

var db user.Repository
var logger *gologger.Logger

func Initialize() {
	logger = gologger.New("api", true)
	dbConnections := dal.Initialize()
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		logger.Log("fatal", "datastore", "ab-101", "error listening to grpcport	", "connectionError", "", map[string]interface{}{}, true)
	}

	db = user.NewUserRepo(logger, dbConnections)
	s := grpc.NewServer(grpc.UnaryInterceptor(serverInterceptor))
	//s := grpc.NewServer()
	athUserService := userService.NewBasicAthUserService(logger, db)
	athMerchantService := merchantService.NewBasicAthMerchantService(logger, merchantdal.NewMerchantRepo(logger, dbConnections), userdal.NewUserRepo(logger, dbConnections))
	athOtpService := otpService.NewBasicAthOtpService(logger, otp.NewOTPRepo(logger, dbConnections), db)
	athVenueService := venueService.NewBasicAthVenueService(logger, venuedal.NewVenueRepo(logger, dbConnections))
	athFacilityService := facilityService.NewBasicAthFacilityService(logger, facilitydal.NewFacilityRepo(logger, dbConnections))
	//athCustomerService := customerService.NewBasicAthCustomerService(logger, db)
	athUploadService := uploadService.NewBasicAthUploadService(logger, merchantdal.NewMerchantRepo(logger, dbConnections), venuedal.NewVenueRepo(logger, dbConnections))

	merchantpb.RegisterMerchantServer(s, merchant.NewMerchantHandler(logger, merchantdal.NewMerchantRepo(logger, dbConnections), athMerchantService, athUserService, athOtpService, athUploadService))
	//customerpb.RegisterCustomerServer(s, customer.NewCustomerHandler(logger, athCustomerService, athOtpService))
	facilitypb.RegisterFacilityServer(s, facility.NewFacilityHandler(logger, facilitydal.NewFacilityRepo(logger, dbConnections), athFacilityService))
	venuepb.RegisterVenueServer(s, venue.NewVenueHandler(logger, venuedal.NewVenueRepo(logger, dbConnections), athVenueService, athUploadService))
	// Register reflection service on gRPC server.
	reflection.Register(s)

	log.Println("grrpc server has started on", grpcPort)

	// Start the gRPC server in goroutine
	go s.Serve(lis)

	// Start the HTTP server for Rest
	log.Println("Starting HTTP server on port " + httpPort)
	run()

}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: true, EmitDefaults: true}))

	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := merchantpb.RegisterMerchantHandlerFromEndpoint(ctx, mux, "localhost"+grpcPort, opts); err != nil {
		logger.Log("fatal", "datastore", "ab-101", "failed to start HTTP gateway", "connectionError", "", map[string]interface{}{"err": err.Error()}, true)
	}
	if err := customerpb.RegisterCustomerHandlerFromEndpoint(ctx, mux, "localhost"+grpcPort, opts); err != nil {
		logger.Log("fatal", "datastore", "ab-101", "failed to start HTTP gateway", "connectionError", "", map[string]interface{}{"err": err.Error()}, true)
	}
	if err := facilitypb.RegisterFacilityHandlerFromEndpoint(ctx, mux, "localhost"+grpcPort, opts); err != nil {
		logger.Log("fatal", "datastore", "ab-101", "failed to start HTTP gateway", "connectionError", "", map[string]interface{}{"err": err.Error()}, true)
	}
	if err := venuepb.RegisterVenueHandlerFromEndpoint(ctx, mux, "localhost"+grpcPort, opts); err != nil {
		logger.Log("fatal", "datastore", "ab-101", "failed to start HTTP gateway", "connectionError", "", map[string]interface{}{"err": err.Error()}, true)
	}
	//http.HandleFunc("/swagger/", serveSwagger)
	return http.ListenAndServe(httpPort, http.StripPrefix("/api", addClientID(mux)))
	//return http.ListenAndServe(httpPort, http.StripPrefix("/api", mux))
}

// Authorization unary interceptor function to handle authorize per RPC call
func serverInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	if info.FullMethod != "/Merchant.Merchant/SignupMerchant" &&
		info.FullMethod != "/Merchant.Merchant/PhoneVerifyMerchant" &&
		info.FullMethod != "/Merchant.Merchant/EmailVerifyMerchant" &&
		info.FullMethod != "/Merchant.Merchant/PhoneVerifyTeam" &&
		info.FullMethod != "/Merchant.Merchant/ForgotPasswordMerchant" &&
		info.FullMethod != "/Merchant.Merchant/ResetPasswordMerchant" &&
		info.FullMethod != "/Merchant.Merchant/ResendCode" &&
		info.FullMethod != "/Merchant.Merchant/UploadDoc" &&
		info.FullMethod != "/Merchant.Merchant/LoginMerchant" {
		if user, err := authorize(ctx); err != nil {
			return nil, err
		} else {
			ctx = context.WithValue(ctx, utile.CurrentUser, user)
		}
	}
	// Calls the handler
	h, err := handler(ctx, req)
	// Logging with grpclog (grpclog.LoggerV2)
	log.Printf("Request - Method:%s\tDuration:%s\tError:%v\n",
		info.FullMethod,
		time.Since(start),
		err)
	return h, err
}

// authorize function authorizes the token received from Metadata
func authorize(ctx context.Context) (*io.AthUser, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "Retrieving metadata is failed")
	}
	authHeader, ok := md["authorization"]
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "Authorization token is not supplied")
	}
	token := authHeader[0]
	parts := strings.Split(token, ":")
	if len(parts) < 2 {
		return nil, status.Errorf(codes.Unauthenticated, "invalid authorization token passed")
	}
	userIDStr := parts[1]
	token = parts[0]
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid authorization token passed")
	}
	masterDB, err := db.GetDB(ctx, dal.MasterDBConnectionName)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	user, err := db.GetUserByID(context.Background(), masterDB, userID, false)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	// if time.Now().After(user.TokenExpireAt) {
	// 	return nil, status.Errorf(codes.Unauthenticated, "authorization token expired")
	// }
	if !utile.ValidPassword(token, user.TokenHash) {
		return nil, status.Errorf(codes.Unauthenticated, "invalid authorization token passed")
	}
	return &user, nil
}

// addClientID is a middleware that injects a client/User ID into the context of each
// request.
func addClientID(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")
		if (*r).Method == "OPTIONS" {
			return
		}
		if r.URL.Path == "/v1/merchants/signup" ||
			r.URL.Path == "/v1/merchants/verify/phone" ||
			r.URL.Path == "/v1/merchants/verify/email" ||
			r.URL.Path == "/v1/merchants/resend-code" ||
			r.URL.Path == "/v1/merchants/upload" ||
			r.URL.Path == "/v1/merchants/forgot-password" ||
			r.URL.Path == "/v1/merchants/login" ||
			r.URL.Path == "/v1/merchants/verify/phone/team" ||
			r.URL.Path == "/v1/merchants/reset-password" {
			h.ServeHTTP(w, r)
			return
		}
		errorMessage := map[string]string{}
		auth := r.Header.Get("Authorization")
		if auth == "" {
			errorMessage["error"] = "authorization token is not supplied"
			b, _ := json.Marshal(errorMessage)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(b)
			return
		}
		parts := strings.Split(auth, ":")
		if len(parts) < 2 {
			errorMessage["error"] = "invalid authorization token passed"
			b, _ := json.Marshal(errorMessage)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(b)
			return
		}
		userIDStr := parts[1]
		token := parts[0]
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			errorMessage["error"] = "invalid authorization token passed"
			b, _ := json.Marshal(errorMessage)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(b)
			return
		}
		masterDB, err := db.GetDB(context.Background(), dal.MasterDBConnectionName)
		if err != nil {
			errorMessage["error"] = "user not found"
			b, _ := json.Marshal(errorMessage)
			w.WriteHeader(http.StatusNotFound)
			w.Write(b)
			return
		}
		user, err := db.GetUserByID(context.Background(), masterDB, userID, false)
		if err != nil {
			errorMessage["error"] = "user not found"
			b, _ := json.Marshal(errorMessage)
			w.WriteHeader(http.StatusNotFound)
			w.Write(b)
			return
		}
		// if time.Now().After(user.TokenExpireAt) {
		// 	errorMessage["error"] = "authorization token expired"
		// 	b, _ := json.Marshal(errorMessage)
		// 	w.WriteHeader(http.StatusUnauthorized)
		// 	w.Write(b)
		// 	return
		// }
		if !utile.ValidPassword(token, user.TokenHash) {
			errorMessage["error"] = "invalid authorization token passed"
			b, _ := json.Marshal(errorMessage)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(b)
			return
		}
		h.ServeHTTP(w, r)
	})
}
