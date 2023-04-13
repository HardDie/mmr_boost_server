package service

import (
	"errors"
	"io"
	"os"

	"github.com/HardDie/mmr_boost_server/internal/logger"
)

var (
	swaggerCache []byte
)

type system struct {
}

func newSystem() system {
	return system{}
}

func (s *system) SystemGetSwagger() ([]byte, error) {
	// Check cache
	if swaggerCache != nil {
		return swaggerCache, nil
	}

	// Open file
	file, err := os.Open("swagger.yaml")
	if err != nil {
		logger.Error.Println("error opening swagger.yaml file:", err.Error())
		return nil, errors.New("can't find swagger.yaml file")
	}
	defer file.Close()

	// Read data from file
	data, err := io.ReadAll(file)
	if err != nil {
		logger.Error.Println("error reading swagger.yaml file:", err.Error())
		return nil, errors.New("error reading swagger.yaml file")
	}

	// Save cache and return result
	swaggerCache = data
	return data, nil
}
