package otp

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"example/totp/util"
	"fmt"
	"log"
	"math"
	"time"
	"unsafe"
)

/*
#cgo LDFLAGS: -L./build -lotp
#include "../libc/otp.h"
#include <stdlib.h>

// Declare C function
uint64_t generate_otp(const char* secret, size_t secret_size, const char* client_secret, size_t client_secret_size, uint64_t interval, uint8_t digit);
*/
import "C"

const (
	timeFormat = "%d-%02d-%02d %02d:%02d:%02d"
)

type OtpManagerIf interface {
	GenerateOtp(clientSecret []byte) uint64
}

type otpManagerImpl struct {
	secret   []byte
	interval int64
	digit    uint
}

type otpManagerLibcImpl struct {
	secret   []byte
	interval int64
	digit    uint
}

func New(secret []byte, interval int64, digit uint) OtpManagerIf {
	return &otpManagerImpl{
		secret:   secret,
		interval: interval,
		digit:    digit,
	}
}

func NewLibc(secret []byte, interval int64, digit uint) OtpManagerIf {
	return &otpManagerLibcImpl{
		secret:   secret,
		interval: interval,
		digit:    digit,
	}
}

func (om *otpManagerImpl) GenerateOtp(clientSecret []byte) uint64 {
	now := util.RoundTimeUTC(time.Now(), time.Duration(om.interval)*time.Second)
	nowStr := fmt.Sprintf(
		timeFormat,
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(),
	)
	h := hmac.New(sha256.New, om.secret)
	payload := append([]byte(nowStr), []byte(clientSecret)...)
	_, err := h.Write(payload)
	if err != nil {
		log.Fatal("Failed to signed client secret: ", err)
	}
	hashed := h.Sum([]byte(nil))[0:4]
	hashed[3] = 0
	truncated := binary.LittleEndian.Uint32(hashed)
	return uint64(truncated) % uint64(math.Pow10(int(om.digit)))
}

func (om *otpManagerLibcImpl) GenerateOtp(clientSecret []byte) uint64 {
	secretSize := C.size_t(len(om.secret))
	clientSecretSize := C.size_t(len(clientSecret))
	secretCString := (*C.char)(C.CBytes(om.secret))
	clientSecretCString := (*C.char)(C.CBytes(clientSecret))
	interval := C.uint64_t(om.interval)
	digit := C.uint8_t(om.digit)

	defer C.free(unsafe.Pointer(secretCString))
	defer C.free(unsafe.Pointer(clientSecretCString))

	result := uint64(C.generate_otp(secretCString, secretSize, clientSecretCString, clientSecretSize, interval, digit))
	if result == 0 {
		log.Println("Failed to generate otp")
	}
	return result
}
