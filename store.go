package session

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

// 中间件实现session功能
type Store interface {
	Options(option Option)
	Get(*session) *session
	Save()
}

type store struct {
	keyPairs []byte
	Option
	session *session
}

func NewStore(keyPairs []byte) Store {
	return &store{
		keyPairs: keyPairs,
	}
}

func (s *store) Options(option Option) {
	if option.SetHeader == "" {
		option.SetHeader = defaultSetHeader
	}
	if option.Header == "" {
		option.Header = defaultHeader
	}
	s.Option = option
	option.Key = s.keyPairs
}

func (s *store) Get(se *session) *session {
	s.session = se // 添加session结构体
	auth := se.request.Header.Get(s.Option.Header)
	// 未找到 auth
	if auth == "" {
		return s.ReturnNullValues()
	}
	encoded := Authorization(se.name, auth)       // 前端返回的sessionID
	decryption := Decryption(encoded, s.keyPairs) // 解密
	decryptionSlice := strings.Split(decryption, ";")
	if len(decryptionSlice) != 2 {
		return s.ReturnNullValues()
	}

	valueStr := decryptionSlice[0]
	maxAgeTimeStamp := decryptionSlice[1]
	nowTime := time.Now().Unix() // 当前时间戳s
	atoi, _ := strconv.Atoi(maxAgeTimeStamp)
	maxAgeTimeStampInt64 := int64(atoi)

	if maxAgeTimeStampInt64-nowTime > 0 {
		// 没过期
		split := strings.Split(valueStr, ":PartingLine:")
		s.session.values = map[string]interface{}{
			split[0]: split[1],
		}
		return s.session
	}

	// 过期了
	return s.ReturnNullValues()
}

func (s *store) Save() {
	auth := s.SpliceAuthorization(s.session.name, s.session.values)
	s.session.writer.Header().Add(s.Option.SetHeader, auth)
}

// 返回session中空的values
func (s *store) ReturnNullValues() *session {
	s.session.values = map[string]interface{}{}
	return s.session
}

/*
	拼接返回header Authorization
*/
func (s *store) SpliceAuthorization(name string, value map[string]interface{}) string {
	var valueName string
	var valueValue interface{}
	for k, v := range value {
		valueName = k
		valueValue = v
	}

	var valueValueStr string
	res, ok := valueValue.(string)
	if ok {
		valueValueStr = res
	} else {
		marshal, _ := json.Marshal(valueValue)
		valueValueStr = string(marshal)
	}
	valueStr := valueName + ":PartingLine:" + valueValueStr

	var expires time.Time
	maxAge := s.Option.MaxAge
	if s.Option.MaxAge > 0 {
		d := time.Duration(maxAge) * time.Second
		expires = time.Now().Add(d)
	} else if maxAge < 0 {
		// Set it to the past to expire now.
		expires = time.Unix(1, 0)
	}

	// 过期时的时间戳(毫秒)
	maxAgeTimeStamp := expires.Unix()
	maxAgeStr := strconv.Itoa(maxAge)
	expiresStr := expires.String()

	// 返回编码之后的sessionID包含过期的时间戳
	encodedBefore := valueStr + ";" + strconv.FormatInt(maxAgeTimeStamp, 10)
	encoded := Encrypt(encodedBefore, s.keyPairs) // 加密
	return name + "=" + encoded + ";" + "Max-Age=" + maxAgeStr + ";" + "Expires=" + expiresStr
}

/*
	header Authorization中的切割出name中的value
	header Authorization格式:
		Authorization: SESSIONID=YG6kI9DnlamYDFbgyUw==; Expires=Thu, 17 Aug 2023 09:15:50 GMT; Max-Age=864000; Secure; SameSite=None
	切出来之后value: YG6kI9DnlamYDFbgyUw==
*/
func Authorization(name, getAuthHeader string) string {
	authorizationslice := strings.Split(getAuthHeader, ";")
	for _, j := range authorizationslice {
		if strings.Contains(j, name+"=") {
			jSlice := strings.Split(j, name+"=")
			if len(jSlice) == 2 {
				return jSlice[1]
			}
		}
	}
	return ""
}
