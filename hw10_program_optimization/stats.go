package hw10programoptimization

import (
	"bufio"
	"errors"
	"io"
	"strings"

	json "github.com/mailru/easyjson"
)

//easyjson:json
type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	var user User
	suffix := "." + domain // инициализируем суффикс один раз.
	stats := make(DomainStat)
	reader := bufio.NewReader(r)
	for {
		jsonText, _, err := reader.ReadLine()
		if errors.Is(err, io.EOF) {
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
