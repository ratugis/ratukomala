package pasetobackend

import (
	"fmt"
	"testing"

	"github.com/aiteung/atdb"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"
)

// PASETO
func TestGeneratePrivateKeyPaseto(t *testing.T) {
	privateKey, publicKey := watoken.GenerateKey()
	fmt.Println(privateKey)
	fmt.Println(publicKey)
	hasil, err := watoken.Encode("nia", privateKey)
	fmt.Println(hasil, err)
}

// func TestValidateToken(t *testing.T) {
// 	tokenstring := "v4.public.eyJleHAiOiIyMDIzLTEwLTI2VDAwOjU5OjQyKzA3OjAwIiwiaWF0IjoiMjAyMy0xMC0yNVQyMjo1OTo0MiswNzowMCIsImlkIjoiZGFmZmEiLCJuYmYiOiIyMDIzLTEwLTI1VDIyOjU5OjQyKzA3OjAwIn0vVwSRywIoGrJC0wrQY-0IXISkClJ" // Gantilah dengan token PASETO yang sesuai
// 	publicKey := "0cdc615519e70986f9821233fabe8e91ba3cc36486242a19801581d7c83fda5e"
// 	payload, _err := watoken.Decode(publicKey, tokenstring)
// 	if _err != nil {
// 		fmt.Println("expired token", _err)
// 	} else {
// 		fmt.Println("ID: ", payload.Id)
// 		fmt.Println("Di mulai: ", payload.Nbf)
// 		fmt.Println("Di buat: ", payload.Iat)
// 		fmt.Println("Expired: ", payload.Exp)
// 	}
// }

// Hash Pass
func TestGeneratePasswordHash(t *testing.T) {
	password := "kepoah"
	hash, _ := HashPassword(password) // ignore error for the sake of simplicity
	fmt.Println("Password:", password)
	fmt.Println("Hash:    ", hash)

	match := CheckPasswordHash(password, hash)
	fmt.Println("Match:   ", match)
}

func TestHashFunction(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "petapedia")
	var userdata User
	userdata.Username = "nia"
	userdata.Password = "males"

	filter := bson.M{"username": userdata.Username}
	res := atdb.GetOneDoc[User](mconn, "gis", filter)
	fmt.Println("Mongo User Result: ", res)
	hash, _ := HashPassword(userdata.Password)
	fmt.Println("Hash Password : ", hash)
	match := CheckPasswordHash(userdata.Password, res.Password)
	fmt.Println("Match:   ", match)

}

func TestIsPasswordValid(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "petapedia")
	var userdata User
	userdata.Username = "nia"
	userdata.Password = "males"

	anu := IsPasswordValid(mconn, "gis", userdata)
	fmt.Println(anu)
}

// User
func TestInsertUser(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "petapedia")
	var userdata User
	userdata.Username = "nia"
	userdata.Password = "males"

	nama := InsertUser(mconn, "gis", userdata)
	fmt.Println(nama)

}
