package utile

import (
	io "api/src/models"
	"bytes"
	"context"
	"errors"
	"math/rand"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	// CurrentUser is the key to extract current logged in user
	CurrentUser = "CurrentUser"
)

// get user id from context
func GetUserIDFromContext(ctx context.Context) (int, error) {
	ctxUser := ctx.Value(CurrentUser)
	if ctxUser == nil {
		return 0, status.Error(codes.Unauthenticated, "user must be logged in")
	}
	currentUser, ok := ctxUser.(*io.AthUser)
	if !ok {
		return 0, status.Error(codes.Unauthenticated, "invalid user")
	}
	return currentUser.UserID, nil
}

// RandomString ... Generates random string
func RandomString(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

// HashAndSalt ... Password encryption
func HashAndSalt(pwd []byte) (string, error) {

	// Use GenerateFromPassword to hash & salt pwd
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash), nil
}

// ComparePasswords ...
func ComparePasswords(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		return false
	}

	return true
}

// Hash ...
func Hash(pword string) (hash string, err error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pword), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("error in generating hash of password")
	}
	return string(b), nil
}

// ValidPassword .. check if password is valid
func ValidPassword(pword, hash string) bool {
	return nil == bcrypt.CompareHashAndPassword(
		[]byte(hash), []byte(pword))
}

func ValidateBase64(base64Str string, entityType string) (base64StrValid string, contentType string, valid bool) {
	contentType = ""
	base64Plain := ""
	imgPrefix := "data:image/jpeg;base64,"
	imgPrefix1 := "data:image/png;base64,"
	pdfPrefix := "data:application/pdf;base64,"
	var base64Psplit []string
	if entityType == "merchant" {
		if strings.HasPrefix(base64Str, imgPrefix) || strings.HasPrefix(base64Str, pdfPrefix) {
			base64Psplit = strings.Split(base64Str, ";base64,")
		}
		if len(base64Psplit) > 0 {
			base64Plain = base64Psplit[1]
			contentTypeData := strings.Split(base64Psplit[0], ":")
			if len(contentTypeData) > 0 {
				contentType = contentTypeData[1]
				return base64Plain, contentType, true
			}
			return base64Plain, contentType, false
		}
	} else if entityType == "venue" {
		if strings.HasPrefix(base64Str, imgPrefix) || strings.HasPrefix(base64Str, imgPrefix1) {
			base64Psplit = strings.Split(base64Str, ";base64,")
		}
		if len(base64Psplit) > 0 {
			base64Plain = base64Psplit[1]
			contentTypeData := strings.Split(base64Psplit[0], ":")
			if len(contentTypeData) > 0 {
				contentType = contentTypeData[1]
				return base64Plain, contentType, true
			}
			return base64Plain, contentType, false
		}
	}
	return base64Plain, contentType, false
}

func UploadToS3(blobImg []byte, contentType string, fileName string) error {

	awsAccessKey := "GOOG1EOAWLCKVIQ2LZPPB5TRGTPA5EUME7CFP44JY4AV6L32YZMO22JPWYOOY"
	awsSecret := "yEES6J/Fx3rnp9I9ABoPogEWV2aPEID93tlhN4ew"
	endpoint := "https://storage.googleapis.com"
	bucket := "sit-turf"
	region := "auto"

	creds := credentials.NewStaticCredentials(awsAccessKey, awsSecret, "")
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
		Endpoint:    aws.String(endpoint),
		LogLevel:    aws.LogLevel(aws.LogDebugWithHTTPBody),
	}))
	s3Svc := s3.New(sess)

	path := "/account/" + fileName
	params := &s3.PutObjectInput{
		Bucket:          aws.String(bucket),
		Key:             aws.String(path),
		Body:            bytes.NewReader(blobImg),
		ContentEncoding: aws.String("base64"),
		ContentType:     aws.String(contentType),
	}

	// upload object
	_, err := s3Svc.PutObject(params)
	return err
}

func VenueUploadToS3(blobImg []byte, contentType string, fileName string) error {

	awsAccessKey := "GOOG1E3AWXTVKDSAHJGQZ4LGJ6P4H5BA66ZBGP3KBKJTVQC64GZBILKBEVS5Q"
	awsSecret := "15S6NedVpGRrsOJGY72IQLMCBwj6L3vv/Wt+PMHe"
	endpoint := "https://storage.googleapis.com"
	bucket := "sit-turf-images"
	region := "auto"

	creds := credentials.NewStaticCredentials(awsAccessKey, awsSecret, "")
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
		Endpoint:    aws.String(endpoint),
		LogLevel:    aws.LogLevel(aws.LogDebugWithHTTPBody),
	}))
	s3Svc := s3.New(sess)

	path := "/venue/" + fileName
	params := &s3.PutObjectInput{
		Bucket:          aws.String(bucket),
		Key:             aws.String(path),
		Body:            bytes.NewReader(blobImg),
		ContentEncoding: aws.String("base64"),
		ContentType:     aws.String(contentType),
	}

	// upload object
	_, err := s3Svc.PutObject(params)
	return err
}
