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
	CREATE        = "CREATE TABLE `%s` (`id` SERIAL NOT NULL, `token` CHAR(64) NOT NULL, `created_at` DATETIME NOT NULL, `expire_at` DATETIME NOT NULL, `resource` TEXT NOT NULL, `data` TEXT NOT NULL)"
	INSERT        = "INSERT INTO `%s` VALUES (null, '%%s', '%%s', '%%s', '%%s', '')"
	SELECT        = "SELECT expire_at, resource, data FROM `%s` WHERE token='%%s'"
	UPDATE_EXPIRE = "UPDATE `%s` SET expire_at='%%s' WHERE token='%%s'"
	UPDATE_DATA   = "UPDATE `%s` SET data='%%s' WHERE token='%%s'"
	DELETE        = "DELETE FROM `%s` WHERE token='%%s'"
)

var ExpiredError = fmt.Errorf("token expired")
var NeverExpire = time.Duration(-1)

type TokenManager struct {
	db      *mysql.Client
	r       *rand.Rand
	create  string
	insert  string
	select_ string
	update  string
	delete  string
}

func New(db *mysql.Client, tableName string) *TokenManager {
	return &TokenManager{
		db:      db,
		r:       rand.New(rand.NewSource(time.Now().UnixNano())),
		create:  fmt.Sprintf(CREATE, tableName),
		insert:  fmt.Sprintf(INSERT, tableName),
		select_: fmt.Sprintf(SELECT, tableName),
		update:  fmt.Sprintf(UPDATE_EXPIRE, tableName),
		delete:  fmt.Sprintf(DELETE, tableName),
	}
}

func (m *TokenManager) GenerateToken(resource string, expireAfterSecond time.Duration) (string, error) {
	now := time.Now().UTC()
	stamp := fmt.Sprintf("%s-%d", resource, now.Unix())
	hash := md5.New()
	io.WriteString(hash, stamp)
	token := fmt.Sprintf("%x%x", hash.Sum(nil), m.randBytes(TOKEN_RANDOM_LENGTH))

	sql := fmt.Sprintf(m.insert, token, now.Format(time.RFC3339), m.getExpireStr(expireAfterSecond), resource)
	err := m.db.Query(sql)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (m *TokenManager) GetResource(token string) (resource, data string, err error) {
	sql := fmt.Sprintf(m.select_, token)
	err = m.db.Query(sql)
	if err != nil {
		return
	}

	var res *mysql.Result
	res, err = m.db.StoreResult()
	if err != nil {
		return
	}
	defer res.Free()

	row := res.FetchRow()
	if row == nil {
		err = fmt.Errorf("no token")
		return
	}

	expire_at_str := string(row[0].([]byte))
	if row[1] != nil {
		resource = string(row[1].([]byte))
	}
	if row[2] != nil {
		data = string(row[2].([]byte))
	}
	if expire_at_str == "0000-00-00 00:00:00" {
		return
	}

	expire_at, err := time.Parse("2006-01-02 15:04:05", expire_at_str)
	if err != nil {
		return
	}

	if expire_at.Sub(time.Now()) <= 0 {
		err = ExpiredError
	}
	return
}

func (m *TokenManager) VerifyToken(token, resource string) (bool, string, error) {
	r, data, err := m.GetResource(token)
	if err != nil && err != ExpiredError {
		return false, "", err
	}
	if r != resource {
		return false, "", err
	}
	return true, data, err
}

func (m *TokenManager) DeleteToken(token string) error {
	sql := fmt.Sprintf(m.delete, token)
	err := m.db.Query(sql)
	return err
}

func (m *TokenManager) RefreshToken(token string, duration time.Duration) error {
	sql := fmt.Sprintf(m.update, m.getExpireStr(duration), token)
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
