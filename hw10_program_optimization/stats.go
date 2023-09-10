package hw10programoptimization

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	domainRe, err := regexp.Compile(fmt.Sprintf(`\w+@(\w+\.%s)`, domain))
	if err != nil {
		return nil, err
	}

	stats := DomainStat{}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		user, err := getUser(scanner.Bytes())
		if err != nil {
			return nil, err
		}

		if validDomain := getEmailDomain(user, domainRe); validDomain != "" {
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

	if err := json.Unmarshal(line, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func getEmailDomain(user *User, domainRe *regexp.Regexp) string {
	if matches := domainRe.FindStringSubmatch(user.Email); len(matches) != 0 {
		return strings.ToLower(matches[1])
	}

	return ""
}
