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
	userCfg  UsersConfig
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
	Username string `xml:"username,attr" json:"username"`
	Password string `xml:"password,attr" json:"password"`
	Roles    Roles  `xml:"roles" json:"roles"`
}

// UserClaim authorization claim
type UserClaim struct {
	Username string `json:"username"`
	Roles    Roles
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

	userCfg = UsersConfig{}
	xml.Unmarshal(data, &userCfg)

	for _, u := range userCfg.Users {
		userList[u.Username] = &u
	}
}

func createToken(user User, expire int) (string, error) {
	claims := UserClaim{
		user.Username,
		user.Roles,
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

// ParseClaimFromRequest parse claims from user request
func ParseClaimFromRequest(req *http.Request) (*jwt.Token, UserClaim, error) {
	_tok := req.Header.Get("Authorization")
	if len(_tok) > 6 && strings.ToUpper(_tok[0:7]) == "BEARER " {
		_tok = _tok[7:]
	}

	usercm := UserClaim{}
	token, err := jwt.ParseWithClaims(_tok, &usercm, keyLookupFn)
	if err != nil {
		return nil, usercm, err
	}
	return token, usercm, nil
}

// Authenticate authenticates user and returns token
func Authenticate(reqUser User, expire int) (string, bool) {
	user, found := userList[reqUser.Username]
	if found == false || (user.Password != reqUser.Password) {
		return "", false
	}

	tokenString, err := createToken(*user, expire)
	if err != nil {
		return "", false
	}

	return tokenString, true
}
