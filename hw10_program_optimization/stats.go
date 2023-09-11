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
	user := User{}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		if err := getUser(&user, scanner.Bytes()); err != nil {
			return nil, err
		}

		if validDomain := getEmailDomain(&user, domain); validDomain != "" {
			stats[validDomain]++
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return stats, nil
}

func getUser(user *User, line []byte) error {
	if err := easyjson.Unmarshal(line, user); err != nil {
		return err
	}

	return nil
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
