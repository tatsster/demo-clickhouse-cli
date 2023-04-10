package entity

import "github.com/tikivn/clickhousectl/internal/utils/encrypt"

type UserStatus string

const (
	UserStatus_ACTIVE   = "active"
	UserStatus_INACTIVE = "inactive"

	secretPasswordKey = "trt2352bfd"
)

type User struct {
	Username       string     `json:"username"`
	Password       string     `json:"-"`
	AllowDatabases []string   `json:"allow_databases"`
	Status         UserStatus `json:"status"`
}

func (u *User) SetPassword(password string) error {
	secret, err := encrypt.Encrypt([]byte(password), secretPasswordKey)
	if err != nil {
		return err
	}
	u.Password = secret
	return nil
}

func (u *User) doubleSha1Password() (string, error) {
	if u.Password == "" {
		return "", nil
	}

	data, err := encrypt.Decrypt([]byte(u.Password), secretPasswordKey)
	if err != nil {
		return "", err
	}

	password := encrypt.DoubleSha1(data)
	return password, nil
}
