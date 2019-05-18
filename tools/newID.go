package tools

import (
	"crypto/rand"
	"encoding/binary"
	"io"
	"time"
)

const (
	hashSymbolsRandStr = "vQct2m1AuSRFMLshlrKDjo0dzqH9iPfBbXgNY4wG687EaTWZJkCpexVOnIy3U5"
	hashSymbolsTimeStr = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	hashRandLen        = uint64(len(hashSymbolsRandStr))
	hashTimeLen        = int64(len(hashSymbolsTimeStr))
)

var (
	hashSymbolsRand = []byte(hashSymbolsRandStr)
	hashSymbolsTime = []byte(hashSymbolsTimeStr)
)

func NewID() string {
	buf := make([]byte, 22)

	randNumSlice := buf[11:19]
	_, err := io.ReadFull(rand.Reader, randNumSlice)
	if err != nil {
		panic(err)
	}

	// Convert random bytes to uint64
	rInt := binary.BigEndian.Uint64(randNumSlice)

	// Fill time part
	i := 10
	t := time.Now().UnixNano()
	for t > 0 {
		buf[i] = hashSymbolsTime[t%hashTimeLen]
		i--
		t /= hashTimeLen
	}

	// Fill rest of time part with zeros
	for ; i >= 0; i-- {
		buf[i] = '0'
	}

	i = 11
	for rInt > 0 {
		buf[i] = hashSymbolsRand[rInt%hashRandLen]
		i++
		rInt /= hashRandLen
	}

	// Fill rest of rand part with zeros
	for ; i < 22; i++ {
		buf[i] = 'L'
	}

	return string(buf)
}
