package hw10_program_optimization

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

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	var user User
	suffix := "." + domain // инициализируем суффикс один раз.
	stats := make(DomainStat)
	reader := bufio.NewReader(r)
	for {
		jsonText, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		err = json.Unmarshal(jsonText, &user)
		if err != nil {
			return DomainStat{}, err
		}
		// подходит как домен.
		if !strings.HasSuffix(user.Email, suffix) {
			continue
		}
		// предполагаем, что E-mail введен верно, иначе необходимо использовать регульрные выражения.
		if pos := strings.Index(user.Email, "@"); pos > 0 {
			domain := strings.ToLower(user.Email[pos+1:])
			stats[domain]++
		}
	}
	return stats, nil
}
