package service

import (
	"errors"
	"io"
	"io/fs"
	"os"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/HardDie/mmr_boost_server/internal/logger"
)

var (
	swaggerCache []byte
)

type System struct {
}

func NewSystem() *System {
	return &System{}
}

func (s *System) GetSwagger() ([]byte, error) {
	// Check cache
	if swaggerCache != nil {
		return swaggerCache, nil
	}

	// Open file
	file, err := os.Open("api.swagger.yaml")
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, status.Error(codes.NotFound, "can't find swagger.yaml file")
		}
		logger.Error.Println("error opening swagger.yaml file:", err.Error())
		return nil, status.Error(codes.Internal, "internal")
	}
	defer file.Close()

	// Read data from file
	data, err := io.ReadAll(file)
	if err != nil {
		logger.Error.Println("error reading swagger.yaml file:", err.Error())
		return nil, status.Error(codes.Internal, "internal")
	}

	// Save cache and return result
	swaggerCache = data
	return data, nil
}
