package library

import (
	"crypto/rsa"
	"io/ioutil"

	jwt "github.com/dgrijalva/jwt-go"
)

//ReadPrivateKey Read RSA key from file
func ReadPrivateKey(savePrivateFileTo string) (key *rsa.PrivateKey, err error) {
	privateKey, err := ioutil.ReadFile(savePrivateFileTo)
	return jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	//privateKey,err := jwt.ParsePKCS1PrivateKey(keyByte)
}

//ReadPrivateKeyByte Read RSA key from file in byte
func ReadPrivateKeyByte(savePrivateFileTo string) (key []byte, err error) {
	return ioutil.ReadFile(savePrivateFileTo)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	//privateKey,err := jwt.ParsePKCS1PrivateKey(keyByte)
}

//ReadPublicKey Read RSA public key from file
func ReadPublicKey(savePrivateFileTo string) (key *rsa.PublicKey, err error) {
	publicKey, err := ioutil.ReadFile(savePrivateFileTo)
	return jwt.ParseRSAPublicKeyFromPEM(publicKey)
	// privateKey,err := x509.ParsePKCS1PrivateKey(keyByte)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	//privateKey,err := jwt.ParsePKCS1PrivateKey(keyByte)
}
