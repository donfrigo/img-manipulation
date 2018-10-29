package helpers

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func GetCurrentFolder () string{
	ex, err := os.Executable()
	CheckErr(err)
	return filepath.Dir(ex)
}

func GenerateRandomFolderName () string {
	randBytes := make([]byte, 8)
	rand.Read(randBytes)

	return hex.EncodeToString(randBytes)
}

// returns hash of file found at filepath
func GetFileHash (filePath string) string {
	var returnMD5String string
	file, err := os.Open(filePath)
	if err != nil {
		return returnMD5String
	}

	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String
	}

	hashInBytes := hash.Sum(nil)[:16]
	returnMD5String = hex.EncodeToString(hashInBytes)
	return returnMD5String

}

// convert a float number to a string with arbitrary precision
func FloatToString(inputNum float64, precision int) string {
	return strconv.FormatFloat(inputNum, 'f', precision, 64)
}

func RetryHttpGet(attempts int, sleep time.Duration, callBackUrl string) {
	var err error

	for i := 0; ; i++ {
		_, err = http.Get(callBackUrl)
		if err == nil {
			return
		}

		if i >= (attempts - 1) {
			break
		}

		time.Sleep(sleep)
	}
	log.Printf("Unsuccessful %d Attempts, Last Error: %s", attempts, err)
}
