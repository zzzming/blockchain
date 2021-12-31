package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"

	"github.com/zzzming/mbt/src/util"
	"golang.org/x/crypto/ripemd160"
)

const (
	checksumLength = 4
	version        = byte(0x00)
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func (w Wallet) Address() []byte {
	pubHash := PublicKeyHash(w.PublicKey)

	versionedHash := append([]byte{version}, pubHash...)
	checksum := Checksum(versionedHash)

	fullHash := append(versionedHash, checksum...)
	// Reference from the original base58.h https://github.com/bitcoin/bitcoin/blob/master/src/base58.h#L6
	// Why base-58 instead of standard base-64 encoding?
	// - Don't want 0OIl characters that look the same in some fonts and
	//      could be used to create visually identical looking account numbers.
	// - A string with non-alphanumeric characters is not as easily accepted as an account number.
	// - E-mail usually won't line-break if there's no punctuation to break at.
	// - Doubleclicking selects the whole number as one word if it's all alphanumeric.
	address := util.Base58Encode(fullHash)

	return address
}

// NewKeyPair implements public and private keys generation based on Elliptic Curve Digital Signature Algorithm or ECDSA.
// ECDSA is a cryptographic algorithm used by Bitcoin to ensure that funds can only be spent by their rightful owners.
// It is dependent on the curve order and hash function used. For bitcoin these are Secp256k1 and SHA256(SHA256()) respectively.
// reference https://en.bitcoin.it/wiki/Elliptic_Curve_Digital_Signature_Algorithm
func NewKeyPair() (ecdsa.PrivateKey, []byte, error) {
	curve := elliptic.P256()

	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return *private, nil, err
	}

	pub := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, pub, nil
}

func NewWallet() (*Wallet, error) {
	private, public, err := NewKeyPair()
	if err != nil {
		return nil, err
	}
	wallet := Wallet{private, public}

	return &wallet, nil
}

// PublicKeyHash implements the original Bitcoin hashing 160-bit of sha 256 hashing of the address' public key
// ripemd160 is used to create a short hash also provide more security over the first hashing
// reference https://en.bitcoin.it/wiki/Technical_background_of_version_1_Bitcoin_addresses
func PublicKeyHash(pubKey []byte) []byte {
	pubHash := sha256.Sum256(pubKey)

	hasher := ripemd160.New()
	_, err := hasher.Write(pubHash[:])
	if err != nil {
		log.Panic(err)
	}

	publicRipMD := hasher.Sum(nil)

	return publicRipMD
}

func Checksum(payload []byte) []byte {
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:checksumLength]
}

func ValidateAddress(address string) bool {
	pubKeyHash := util.Base58Decode([]byte(address))
	actualChecksum := pubKeyHash[len(pubKeyHash)-checksumLength:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-checksumLength]
	targetChecksum := Checksum(append([]byte{version}, pubKeyHash...))

	return bytes.Equal(actualChecksum, targetChecksum)
}
