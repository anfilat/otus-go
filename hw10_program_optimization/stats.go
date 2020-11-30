package hw10_program_optimization //nolint:golint,stylecheck

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

// очень хочется убрать здесь все поля кроме Email, но кажется это будет читерство

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

// из условий непонятно, можно ли менять функцию GetDomainStat, поэтому она остается неизменной

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type NextUser = func() (*User, bool, error)

// да error здесь всегда nil, но надо вписаться в существующую GetDomainStat
//nolint:unparam
func getUsers(r io.Reader) (NextUser, error) {
	var auser User
	s := bufio.NewScanner(r)

	nextUser := func() (user *User, ok bool, err error) {
		ok = s.Scan()

		if !ok {
			err = s.Err()
			if err != nil {
				err = fmt.Errorf("error with reading data: %w", err)
			}
			return
		}

		if err = jsoniter.ConfigFastest.Unmarshal(s.Bytes(), &auser); err != nil {
			err = fmt.Errorf("error with reading user: %w", err)
			return
		}

		user = &auser
		return
	}

	return nextUser, nil
}

func countDomains(nextUser NextUser, domain string) (DomainStat, error) {
	result := make(DomainStat)
	domainMask := "." + domain

	for {
		user, ok, err := nextUser()
		if err != nil {
			return nil, err
		}
		if !ok {
			break
		}

		email := strings.ToLower(user.Email)
		matched := strings.HasSuffix(email, domainMask)

		if matched {
			n := strings.LastIndex(email, "@")
			if n == -1 {
				return nil, fmt.Errorf("wrong email: %s", email)
			}
			result[email[n+1:]]++
		}
	}

	return result, nil
}
