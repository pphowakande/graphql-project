## mockgen: generate mock
.PHONY : mockgen
mockgen:
	@cd ./src/service/user && mockgen -package=mock -destination=./mock/user_mock.go api/src/service/user AthUserService
	@cd ./src/service/otp && mockgen -package=mock -destination=./mock/otp_mock.go api/src/service/otp AthOtpService
	@cd ./src/service/upload && mockgen -package=mock -destination=./mock/upload_mock.go api/src/service/upload AthUploadService
	@cd ./src/service/merchant && mockgen -package=mock -destination=./mock/merchant_mock.go api/src/service/merchant AthMerchantService
	@cd ./src/service/facility && mockgen -package=mock -destination=./mock/facility_mock.go api/src/service/facility AthFacilityService
	@cd ./src/service/venue && mockgen -package=mock -destination=./mock/venue_mock.go api/src/service/venue AthVenueService

## run_tests: Runs all the tests
.PHONY: run_tests
run_tests:
	@go test -v ./...

## install: installs the proto dependencies
.PHONY: install
install:
	@go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	@go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
	@go get -u github.com/golang/protobuf/protoc-gen-go