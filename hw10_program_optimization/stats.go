package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"

	"github.com/mailru/easyjson"
)

//easyjson:json
type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	stats := DomainStat{}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		user, err := getUser(scanner.Bytes())
		if err != nil {
			return nil, err
		}

		if validDomain := getEmailDomain(user, domain); validDomain != "" {
			stats[validDomain]++
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return stats, nil
}

func getUser(line []byte) (*User, error) {
	var user User

	if err := easyjson.Unmarshal(line, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func getEmailDomain(user *User, domain string) string {
	sp := strings.Split(user.Email, "@")
	if len(sp) != 2 {
		return ""
	}

	emailDomain := sp[1]
	if strings.Contains(emailDomain, "."+domain) {
		return strings.ToLower(emailDomain)
	}

	return ""
}
