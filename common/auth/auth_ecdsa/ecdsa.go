package auth_ecdsa

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	r "math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pavlo67/workshop/common/identity"

	"github.com/btcsuite/btcutil/base58"
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common/auth"
	"github.com/pavlo67/workshop/common/libraries/encrlib"
	"github.com/pavlo67/workshop/common/libraries/strlib"
)

const Cryptype encrlib.Cryptype = "ecdsa"
const Proto = "ecdsa"

var _ auth.Operator = &authECDSA{}

var errWrongSignature = errors.New("wrong signature")
var errWrongNumber = errors.New("wrong user's number")
var errEmptyPublicKeyAddress = errors.New("empty public Key address")
var errEmptyPrivateKeyGenerated = errors.New("empty private key generated")

var cnt uint32

type Session struct {
	IP        string
	StartedAt time.Time
}

type authECDSA struct {
	sessions map[uint64]Session
	mutex    *sync.Mutex

	maxSessionDuration time.Duration
	numbersLimit       int

	acceptableIDs []string
}

func New(numbersLimit int, maxSessionDuration time.Duration, acceptableIDs []string) (auth.Operator, error) {
	r.Seed(time.Now().UnixNano())

	is := &authECDSA{
		sessions: map[uint64]Session{},
		mutex:    &sync.Mutex{},

		maxSessionDuration: maxSessionDuration,
		numbersLimit:       numbersLimit,

		acceptableIDs: acceptableIDs,
	}

	return is, nil
}

// 	SetCreds creates either session-generated key or new "BTC identity" and returns it
func (is *authECDSA) SetCreds(userKey identity.Key, creds auth.Creds, toSet auth.CredsType) (identity.Key, *auth.Creds, error) {
	if toSet == auth.CredsKeyToSignature {
		now := time.Now()

		is.mutex.Lock() // Lock() -----------------------------------------------------

		if is.numbersLimit > 0 && len(is.sessions) >= is.numbersLimit {
			for n, s := range is.sessions {
				if now.Sub(s.StartedAt) >= is.maxSessionDuration {
					delete(is.sessions, n)
				}
			}
		}

		cnt++
		numberToSend := uint64(cnt)<<32 + uint64(r.Uint32())

		is.sessions[numberToSend] = Session{
			IP:        creds.Values[auth.CredsIP], // TODO??? check if IP isn't empty
			StartedAt: now,
		}

		is.mutex.Unlock() // Unlock() -------------------------------------------------

		return "", &auth.Creds{
			Cryptype: encrlib.NoCrypt,
			Values:   auth.Values{auth.CredsKeyToSignature: strconv.FormatUint(numberToSend, 10)},
		}, nil
	}

	// TODO: modify acceptableIDs if it's necessary

	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", nil, err
	} else if privKey == nil {
		return "", nil, errEmptyPrivateKeyGenerated
	}

	privKeyBytes, err := encrlib.ECDSASerialize(*privKey)
	if err != nil {
		return "", nil, err
	}

	publKeyBase58 := base58.Encode(encrlib.ECDSAPublicKey(*privKey))
	nickname := publKeyBase58
	if creds.Values[auth.CredsNickname] != "" {
		nickname = creds.Values[auth.CredsNickname]
	}

	credsNew := &auth.Creds{
		Values: map[auth.CredsType]string{
			auth.CredsNickname:          nickname,
			auth.CredsPrivateKey:        string(privKeyBytes),
			auth.CredsPublicKeyBase58:   publKeyBase58,
			auth.CredsPublicKeyEncoding: Proto,
		},
	}

	return identity.Key(Proto + "://" + publKeyBase58), credsNew, nil
}

func (is *authECDSA) Authorize(toAuth auth.Creds) (*auth.User, error) {
	if toAuth.Values[auth.CredsPublicKeyEncoding] != Proto {
		return nil, auth.ErrEncryptionType
	}

	publKeyBase58 := toAuth.Values[auth.CredsPublicKeyBase58]
	if len(publKeyBase58) < 1 {
		return nil, errEmptyPublicKeyAddress
	}
	publKey := base58.Decode(publKeyBase58)

	// TODO: use mutex is is.acceptableIDs can be modified using .SetCreds or somehow else
	if is.acceptableIDs != nil && !strlib.In(is.acceptableIDs, publKeyBase58) {
		return nil, nil
	}

	keyToSignature := toAuth.Values[auth.CredsKeyToSignature]
	numberToSend, err := strconv.ParseUint(keyToSignature, 10, 64)
	if err != nil {
		return nil, errors.Wrap(auth.ErrSignaturedKey, "not a number!")
	}

	is.mutex.Lock() // Lock() -----------------------------------------------------

	session, ok := is.sessions[numberToSend]
	if !ok {
		return nil, errors.Wrap(auth.ErrSignaturedKey, "no appropriate session")
		is.mutex.Unlock() // Unlock() ---------------------------------------------
	}
	delete(is.sessions, numberToSend)

	is.mutex.Unlock() // Unlock() -------------------------------------------------

	if time.Now().Sub(session.StartedAt) > is.maxSessionDuration {
		return nil, errors.Wrap(auth.ErrAuthSession, "session is expired")
	}

	if session.IP != toAuth.Values[auth.CredsIP] {
		return nil, auth.ErrIP
	}

	signature := []byte(toAuth.Values[auth.CredsSignature])

	if !encrlib.ECDSAVerify(keyToSignature, publKey, signature) {
		return nil, errWrongSignature
	}

	var nickname = publKeyBase58
	if nicknameReceived := toAuth.Values[auth.CredsNickname]; strings.TrimSpace(nicknameReceived) != "" {
		nickname = nicknameReceived
	}

	return &auth.User{
		Key: identity.Key(Proto + "://" + publKeyBase58),
		Creds: auth.Creds{
			Cryptype: encrlib.NoCrypt,
			Values: auth.Values{
				auth.CredsNickname: nickname,
			},
		},
	}, nil
}

//
//func (*authECDSA) Accepts() ([]auth.CredsType, error) {
//	return []auth.CredsType{auth.CredsSignature}, nil
//}
