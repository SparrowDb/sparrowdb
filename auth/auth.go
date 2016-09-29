package auth

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/SparrowDb/sparrowdb/db"
	"github.com/SparrowDb/sparrowdb/slog"
	"github.com/SparrowDb/sparrowdb/util/uuid"
)

const (
	defaultUserFile = "user.xml"
)

var (
	//onlineUsers map[string]onlineUser
	userList    map[string]*User
	tokenList   map[string]string
	cleanupTime time.Duration
	muUser      sync.RWMutex
	userExpire  time.Duration
)

type UsersConfig struct {
	XMLName xml.Name `xml:"users"`
	Users   []User   `xml:"user"`
}

type User struct {
	Username string `xml:"username"`
	Password string `xml:"password"`
	expires  time.Time
	token    string
	online   bool
}

func (u *User) expired() bool {
	return u.expires.Before(time.Now())
}

func (u *User) update() {
	expiration := time.Now().Add(userExpire * time.Millisecond)
	u.expires = expiration
}

func init() {
	cleanupTime = 3000 * time.Millisecond
	userList = make(map[string]*User)
	tokenList = make(map[string]string)
	startCleanUp()
}

func startCleanUp() {
	ticker := time.Tick(cleanupTime)

	go func() {
		for {
			select {
			case <-ticker:
				for _, item := range userList {
					if item.online {
						if item.expired() {
							item.online = false
							delete(tokenList, item.token)
							item.token = ""
						}
					}
				}
			}
		}
	}()
}

// LoadUserConfig loads users from configuration file
func LoadUserConfig(filePath string, dbConfig *db.SparrowConfig) {
	path := filepath.Join(filePath, defaultUserFile)

	xmlFile, err := os.Open(path)
	if err != nil {
		slog.Fatalf("Could not load users definition file")
	}

	defer xmlFile.Close()

	userExpire = time.Duration(dbConfig.UserExpire)

	data, _ := ioutil.ReadAll(xmlFile)

	users := UsersConfig{}
	xml.Unmarshal(data, &users)

	for _, u := range users.Users {
		userList[u.Username] = &u
	}
}

// IsLogged checks if user is logged, if yes
// update expire time
func IsLogged(token string) bool {
	if v, ok := tokenList[token]; ok == true {
		user, _ := userList[v]
		if user.expired() {
			return false
		}
		user.update()
		return true
	}
	return false
}

// Authenticate authenticates user and returns token
func Authenticate(reqUser User) (bool, string) {
	user, found := userList[reqUser.Username]
	if found == false || (user.Password != reqUser.Password) {
		return false, ""
	}

	user.update()

	if user.online {
		return true, user.token
	}

	token := uuid.TimeUUID().String()
	user.token = token
	user.online = true

	tokenList[token] = user.Username

	return true, token
}
