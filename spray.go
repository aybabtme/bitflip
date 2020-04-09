package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

func parseSprayPattern(sprayPattern string) (SprayPatternFactory, error) {
	args := strings.Split(sprayPattern, ":")
	if len(args) > 2 {
		return nil, errors.New("invalid pattern string, too many ':' symbols found")
	}
	if len(args) < 2 {
		return nil, errors.New("invalid pattern string, no ':' symbols found")
	}
	patternType := args[0]
	patternArgs := args[1]
	switch patternType {
	case "percent":
		return newSprayPercentPatternFactory(patternArgs)
	}
	return nil, fmt.Errorf("unknown spray pattern type: %q", patternType)
}

type SprayPatternFactory func(os.FileInfo) SprayPattern

type SprayPattern interface {
	Spray(io.ReadWriteSeeker, bitflipFunc) error
}

func newSprayPercentPatternFactory(args string) (SprayPatternFactory, error) {
	percent, err := strconv.ParseFloat(args, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid argument for `percent:` spray: %v", err)
	}
	return func(fi os.FileInfo) SprayPattern {
		return &sprayPercentPattern{
			r:       rand.New(rand.NewSource(time.Now().UnixNano())),
			percent: percent,
			fi:      fi,
		}
	}, nil
}

type sprayPercentPattern struct {
	r       *rand.Rand
	percent float64
	fi      os.FileInfo
}

func (sp *sprayPercentPattern) Spray(rws io.ReadWriteSeeker, flipFn bitflipFunc) error {
	toFlip := int(float64(sp.fi.Size()*8) * sp.percent / 100)
	log.Printf("randomly flipping %s out of %s (%g%% of %s) in file %q",
		humanize.SI(float64(toFlip), "bits"),
		humanize.SI(float64(sp.fi.Size()*8), "bits"),
		sp.percent,
		humanize.IBytes(uint64(sp.fi.Size())),
		sp.fi.Name(),
	)
	for i := 0; i < toFlip; i++ {
		byteOffset := sp.r.Int63n(sp.fi.Size())
		bitOffset := uint8(sp.r.Intn(8))
		if err := flipFn(rws, byteOffset, bitOffset); err != nil {
			return fmt.Errorf("spraying file with random bitflips (flip number %d): %v", i, err)
		}
	}
	return nil
}
