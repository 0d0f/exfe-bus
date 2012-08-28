package tokenmanager

import (
	"crypto/md5"
	"fmt"
	"github.com/googollee/go-mysql"
	"io"
	"math"
	"math/rand"
	"time"
)

const TOKEN_RANDOM_LENGTH = 16

const (
	CREATE        = "CREATE TABLE `tokens` (`id` SERIAL NOT NULL, `token` CHAR(64) NOT NULL, `created_at` DATETIME NOT NULL, `expire_at` DATETIME, `resource` TEXT)"
	INSERT        = "INSERT INTO `tokens` VALUES (null, '%s', '%s', '%s', '%s')"
	SELECT        = "SELECT expire_at, resource FROM `tokens` WHERE token='%s'"
	UPDATE_EXPIRE = "UPDATE `tokens` SET expire_at='%s' WHERE token='%s'"
	DELETE        = "DELETE FROM `tokens` WHERE token='%s'"
)

var ExpiredError = fmt.Errorf("token expired")
var NeverExpire = time.Duration(-1)

type TokenManager struct {
	db *mysql.Client
	r  *rand.Rand
}

func New(db *mysql.Client) *TokenManager {
	return &TokenManager{
		db: db,
		r:  rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (m *TokenManager) GenerateToken(resource string, expire_duration time.Duration) (string, error) {
	now := time.Now().UTC()
	stamp := fmt.Sprintf("%s-%d", resource, now.Unix())
	hash := md5.New()
	io.WriteString(hash, stamp)
	token := fmt.Sprintf("%x%x", hash.Sum(nil), m.randBytes(TOKEN_RANDOM_LENGTH))

	sql := fmt.Sprintf(INSERT, token, now.Format(time.RFC3339), m.getExpireStr(expire_duration), resource)
	err := m.db.Query(sql)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (m *TokenManager) GetResource(token string) (string, error) {
	sql := fmt.Sprintf(SELECT, token)
	err := m.db.Query(sql)
	if err != nil {
		return "", err
	}

	res, err := m.db.StoreResult()
	if err != nil {
		return "", err
	}
	defer res.Free()

	row := res.FetchRow()
	if row == nil {
		return "", fmt.Errorf("no token")
	}

	expire_at_str, resource := string(row[0].([]byte)), string(row[1].([]byte))
	if expire_at_str == "0000-00-00 00:00:00" {
		return resource, nil
	}

	expire_at, err := time.Parse("2006-01-02 15:04:05", expire_at_str)
	if err != nil {
		return "", err
	}

	if expire_at.Sub(time.Now()) <= 0 {
		err = ExpiredError
	}
	return resource, err
}

func (m *TokenManager) VerifyToken(token, resource string) (bool, error) {
	r, err := m.GetResource(token)
	if err != nil && err != ExpiredError {
		return false, err
	}
	return r == resource, err
}

func (m *TokenManager) DeleteToken(token string) error {
	sql := fmt.Sprintf(DELETE, token)
	err := m.db.Query(sql)
	return err
}

func (m *TokenManager) RefreshToken(token string, duration time.Duration) error {
	sql := fmt.Sprintf(UPDATE_EXPIRE, m.getExpireStr(duration), token)
	err := m.db.Query(sql)
	return err
}

func (m *TokenManager) ExpireToken(token string) error {
	return m.RefreshToken(token, 0)
}

func (m *TokenManager) getExpireStr(duration time.Duration) string {
	if duration == NeverExpire {
		return "0000-00-00 00:00:00"
	}
	expire := time.Now().Add(duration)
	return expire.UTC().Format(time.RFC3339)
}

func (m *TokenManager) randBytes(length int) []byte {
	ret := make([]byte, length, length)
	for i := range ret {
		ret[i] = byte(m.r.Int31n(math.MaxInt8))
	}
	return ret
}
