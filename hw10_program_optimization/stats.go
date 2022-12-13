package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

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

var parser = jsoniter.ConfigCompatibleWithStandardLibrary

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	var user User
	stats := make(DomainStat)
	reader := bufio.NewReader(r)
	for {
		json, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		err = parser.Unmarshal(json, &user)
		if err != nil {
			return DomainStat{}, err
		}
		// подходит как домен
		if !strings.HasSuffix(user.Email, "."+domain) {
			continue
		}
		if pos := strings.Index(user.Email, "@"); pos > 0 {
			domain := strings.ToLower(user.Email[pos+1:])
			stats[domain]++
		}
	}
	return stats, nil
}
