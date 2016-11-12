package auth

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/SparrowDb/sparrowdb/db"
	"github.com/SparrowDb/sparrowdb/slog"

	"github.com/dgrijalva/jwt-go"
)

const (
	defaultUserFile = "user.xml"
)

var (
	userList map[string]*User
	secret   string
	letters  = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func init() {
	userList = make(map[string]*User, 0)
	secret = randSeq(256)
}

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// UsersConfig user from xml file
type UsersConfig struct {
	XMLName xml.Name `xml:"users"`
	Users   []User   `xml:"user"`
}

// User holds user info
type User struct {
	Username string `xml:"username" json:"username"`
	Password string `xml:"password" json:"password"`
}

// UserClaim authorization claim
type UserClaim struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// LoadUserConfig loads users from configuration file
func LoadUserConfig(filePath string, dbConfig *db.SparrowConfig) {
	path := filepath.Join(filePath, defaultUserFile)

	xmlFile, err := os.Open(path)
	if err != nil {
		slog.Fatalf("Could not load users definition file")
	}

	defer xmlFile.Close()

	data, _ := ioutil.ReadAll(xmlFile)

	users := UsersConfig{}
	xml.Unmarshal(data, &users)

	for _, u := range users.Users {
		userList[u.Username] = &u
	}
}

func createToken(user User, expire int) (string, error) {
	claims := UserClaim{
		user.Username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(expire) * time.Millisecond).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func keyLookupFn(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return []byte(secret), nil
}

// ValidateToken checks if token is valid
func ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, keyLookupFn)
}

// ParseFromRequest parses token from request
func ParseFromRequest(req *http.Request) (*jwt.Token, error) {
	_tok := req.Header.Get("Authorization")
	var token string

	if len(_tok) > 6 && strings.ToUpper(_tok[0:7]) == "BEARER " {
		token = _tok[7:]
	}

	return ValidateToken(token)
}

// Authenticate authenticates user and returns token
func Authenticate(reqUser User, expire int) (string, bool) {
	user, found := userList[reqUser.Username]
	if found == false || (user.Password != reqUser.Password) {
		return "", false
	}

	tokenString, err := createToken(reqUser, expire)
	if err != nil {
		return "", false
	}

	return tokenString, true
}
