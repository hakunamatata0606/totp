package otp_test

import (
	"example/totp/otp"
	"example/totp/util"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func simpleVerify(t *testing.T, i int, om otp.OtpManagerIf, clientSecret string, interval int64) {
	time.Sleep(time.Duration(i) * time.Second)
	nowRoundedUtc := util.RoundTimeUTC(time.Now(), time.Duration(interval)*time.Second)

	elapsed := int64(time.Since(nowRoundedUtc).Seconds())
	if elapsed >= (interval - 1) {
		time.Sleep(time.Duration((elapsed - interval + 1)) * time.Second)
	}

	signed1 := om.GenerateOtp([]byte(clientSecret))

	time.Sleep(time.Duration(1) * time.Second)
	signed2 := om.GenerateOtp([]byte(clientSecret))
	require.Equal(t, signed1, signed2)

	time.Sleep(time.Duration(interval) * time.Second)
	signed3 := om.GenerateOtp([]byte(clientSecret))
	require.NotEqual(t, signed1, signed3)
}

func TestGenerateOtp(t *testing.T) {
	interval := 5
	om := otp.New([]byte("aloha"), int64(interval), 10)
	clientSecret := "client_secret"
	var wg sync.WaitGroup

	for i := range 10 {
		wg.Add(1)
		cSecret := clientSecret + fmt.Sprintf("%d", i)
		go func() {
			defer wg.Done()
			simpleVerify(t, i, om, cSecret, int64(interval))
		}()
	}
	wg.Wait()
}

// func TestGenerateOtpLibc(t *testing.T) {
// 	interval := 5
// 	om := otp.NewLibc([]byte("aloha"), int64(interval), 10)
// 	clientSecret := "client_secret"
// 	var wg sync.WaitGroup

// 	for i := range 3 {
// 		wg.Add(1)
// 		cSecret := clientSecret + fmt.Sprintf("%d", i)
// 		go func() {
// 			defer wg.Done()
// 			simpleVerify(t, i, om, cSecret, int64(interval))
// 		}()
// 	}
// 	wg.Wait()
// }
