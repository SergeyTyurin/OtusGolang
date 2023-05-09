package hw10programoptimization

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	easyjson "github.com/mailru/easyjson"
)

//easyjson:json
type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return
	}
	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		var user User
		if err = easyjson.Unmarshal([]byte(line), &user); err != nil {
			return
		}
		result[i] = user
	}
	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)
	reg, err := regexp.Compile("\\." + domain)
	if err != nil {
		return nil, err
	}
	size := len(u)
	for i := 0; i < size; i++ {
		if reg.MatchString(u[i].Email) {
			result[strings.ToLower(strings.SplitN(u[i].Email, "@", 2)[1])]++
		}
	}
	return result, nil
}
