package controller

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/teris-io/shortid"
)

const (
	GenerationTaskIdPrefix = "qry"
)

type PrefixIdGenerator struct {
	sid *shortid.Shortid
}

func NewPrefixedIdGenerator() *PrefixIdGenerator {
	sid, err := shortid.New(1, shortid.DefaultABC, 1)

	if err != nil {
		log.Fatalf("unable to initialize id generator: %v", err)
		return nil
	}
	return &PrefixIdGenerator{
		sid: sid,
	}
}

func (generator *PrefixIdGenerator) Generator(prefix string) (string, error) {
	if len(prefix) == 0 {
		return "", fmt.Errorf("prefix must not be empty")
	}

	str, err := generator.sid.Generate()
	if err != nil {
		return "", fmt.Errorf("unable to generate id for prefix %s: %v", prefix, err)
	}

	return fmt.Sprintf("%s_%s", prefix, str), nil
}
