package qr

import (

	"github.com/skip2/go-qrcode"
)

func GenerateQR(link string) ([]byte, error) {
	return qrcode.Encode(link, qrcode.Medium, 256)
}