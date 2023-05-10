package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
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
	result := make(DomainStat)

	scanner := bufio.NewScanner(r)
	for {
		if !scanner.Scan() {
			break
		}

		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("get users error: %w", err)
		}

		var user User

		if err := easyjson.Unmarshal(scanner.Bytes(), &user); err != nil {
			return nil, fmt.Errorf("get users error: %w", err)
		}

		if strings.Contains(user.Email, "."+domain) {
			domainBegin := strings.Index(user.Email, "@")
			result[strings.ToLower(user.Email[domainBegin+1:])]++
		}
	}
	return result, nil
}
