package util

import (
	"crypto/md5"
	"errors"
	"fmt"
	_ "net/http/pprof"
	"os"
	"strings"

	wphash "github.com/GerardSoleCa/wordpress-hash-go"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	uuid "github.com/gofrs/uuid"
	"github.com/lithammer/shortuuid"
	"golang.org/x/crypto/bcrypt"
)

type (
	NullLogger struct{}
	Country    struct {
		Name string
		ISO2 string
	}
)

func ShortID() string {
	return shortuuid.New()
}

func ShortIDS(number int) (result []string) {
	for index := 0; index < number; index++ {
		result = append(result, shortuuid.New())
	}
	return
}

func Hash(password string) string {
	hashedPW, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPW)
}

func GenerateUUID() (result string, err error) {
	newUuid, err := uuid.NewV4()
	if err != nil {
		return
	}
	result = newUuid.String()
	return
}

func MD5(input string) (out string) {
	return fmt.Sprintf("%x", md5.Sum([]byte(input)))
}

func GenerateUUIDS(number int) (result []string, err error) {
	for index := 0; index < number; index++ {
		newUuid, err2 := uuid.NewV4()
		if err = err2; err != nil {
			return
		}
		result = append(result, newUuid.String())
	}

	return
}

func ValidateEmail(email string) (sanitized string, err error) {
	if err = validation.Validate(email, is.Email); err != nil {
		return
	}

	sanitized = email

	//This check does not guard against custom domains on G-Suite
	if !strings.HasSuffix(email, "gmail.com") {
		return
	}

	emailParts := strings.Split(email, "@")
	if len(emailParts) < 2 {
		return "", errors.New("Invalid email")
	}

	//remove periods and get only strings before plus sign
	emailParts = strings.Split(strings.Replace(emailParts[0], ".", "", -1), "+")
	if len(emailParts) >= 1 {
		return emailParts[0] + "@gmail.com", nil
	}

	return email, nil
}

func CheckWordpressPassword(password, hash string) bool {
	return wphash.CheckPassword(password, hash)
}

func HashWordpressPassword(password string) string {
	return wphash.HashPassword(password)
}

func (NullLogger) Print(...interface{}) {}

func EnsureDir(path string, perms os.FileMode) (err error) {
	if err = CheckPath(path); err == nil {
		return
	}

	return makeDir(path, perms)

}

func EnsureDirs(paths []string, perms os.FileMode) (err error) {
	for _, path := range paths {
		// fmt.Printf("Ensuring %s exists\n", path)
		if err = EnsureDir(path, perms); err != nil {
			return
		}
	}
	return
}

func CheckPath(path string) (err error) {
	if _, err = os.Stat(path); os.IsNotExist(err) {
		return errors.New(fmt.Sprintf("%s does not exist", path))
	}
	return nil
}

func makeDir(dir string, perms os.FileMode) (err error) {
	if err = os.MkdirAll(dir, perms); err != nil {
		return errors.New("Could not create directory:" + dir)
	}

	return
}

func GetCountry(name string) (country Country, slug string, err error) {
	country, ok := Countries[name]
	if ok {
		return
	}

	code, ok := countryTranslations[name]
	if !ok {
		return Country{}, "", fmt.Errorf("País %s não encontrado", name)
	}

	if country, ok = Countries[code]; !ok {
		return Country{}, "", fmt.Errorf("País %s não encontrado", name)
	}

	return
}
