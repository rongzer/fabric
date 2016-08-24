// util
package util

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"
)

func GenerateUUID() string {
	uuid := GenerateBytesUUID()
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

func GenerateBytesUUID() []byte {
	uuid := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, uuid)
	if err != nil {
		panic(fmt.Sprintf("Error generating UUID: %s", err))
	}

	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80

	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40

	return uuid
}

func GenerateUUIDFromCert(b []byte) string {
	h := md5.New()
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}

func GetAffiliation(b []byte) ([]string, error) {
	cert, err := x509.ParseCertificate(b)
	if err != nil {
		return nil, err
	}
	var s string
	if cert != nil && cert.Subject.CommonName != "" {
		s = cert.Subject.CommonName
	}
	ss := strings.Split(s, "\\")
	if len(ss) > 2 {
		return ss, nil
	} else {
		return nil, errors.New("证书的CommonName格式不正确")
	}
}
func GetAffiliationFromString(scert string) (string, string, error) {
	cert, err := hex.DecodeString(scert)
	if err != nil {
		return "", "", err
	}
	affs, err := GetAffiliation(cert)
	if err != nil {
		return "", "", err
	}
	uuid := affs[0]
	aff := affs[1]
	if uuid == "" {
		return "", "", errors.New("uuid is nil")
	}
	return uuid, aff, nil
}
func PageRow(pagesize, pagenum, length int64) (int64, int64) {
	var begin, end int64
	total := length/pagesize + 1
	if pagenum < total {
		begin = (pagenum - 1) * pagesize
		end = pagenum * pagesize
	} else if pagenum == total {
		begin = (pagenum - 1) * pagesize
		end = length
	} else {
		begin = 0
		end = 0
	}
	return begin, end
}
