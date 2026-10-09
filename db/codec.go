package db

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/abuabdillatief/sample/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
)

var (
	CodecRegistry *bsoncodec.Registry
	defaultKey    = "secureencryptionkeynoonewillknow"
)

func init() {
	CodecRegistry = setCodecRegistry(bson.NewRegistry())
}

func setCodecRegistry(rb *bsoncodec.Registry) *bsoncodec.Registry {
	// Register time.Time codec
	rb.RegisterTypeEncoder(reflect.TypeOf(time.Time{}), &timeCodec{})
	rb.RegisterTypeDecoder(reflect.TypeOf(time.Time{}), &timeCodec{})

	// Register Passport codec
	rb.RegisterTypeEncoder(reflect.TypeOf(models.Passport{}), &passportCodec{})
	rb.RegisterTypeDecoder(reflect.TypeOf(models.Passport{}), &passportCodec{})

	return rb
}

type timeCodec struct{}

func (tc *timeCodec) EncodeValue(ectx bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	if !val.IsValid() || val.IsZero() {
		return vw.WriteNull()
	}

	t := val.Interface().(time.Time)
	timestamp := t.Unix() // Convert to Unix timestamp

	return vw.WriteInt64(timestamp)
}

func (tc *timeCodec) DecodeValue(dctx bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.CanSet() {
		return errors.New("value must be settable")
	}

	timestamp, err := vr.ReadInt64()
	if err != nil {
		return err
	}

	t := time.Unix(timestamp, 0) // Convert Unix timestamp back to time.Time
	val.Set(reflect.ValueOf(t))

	return nil
}

type passportCodec struct{}

func (tc *passportCodec) EncodeValue(ectx bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	if !val.IsValid() || val.IsZero() {
		return vw.WriteNull()
	}

	passport := val.Interface().(models.Passport)
	b, err := json.Marshal(passport)
	if err != nil {
		log.Println("Failed to marshal passport data", err)

		return err
	}

	key := os.Getenv("ENCRYPTION_KEY")
	if key == "" {
		key = defaultKey
	}

	encrpyted, err := Encrypt(string(b), key)
	if err != nil {
		log.Println("Failed to encrypt passport data", err)

		return err
	}

	return vw.WriteString(encrpyted)
}

func (tc *passportCodec) DecodeValue(dctx bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.CanSet() {
		return errors.New("value must be settable")
	}

	passportAsString, err := vr.ReadString()
	if err != nil {
		log.Println("Failed to read passport data", err)

		return err
	}

	key := os.Getenv("ENCRYPTION_KEY")
	if key == "" {
		key = defaultKey
	}

	decrypted, err := Decrypt(passportAsString, key)
	if err != nil {
		log.Println("Failed to decrypt passport data", err)
		return err
	}

	var p models.Passport
	if err := json.Unmarshal([]byte(decrypted), &p); err != nil {
		log.Println("Failed to unmarshal passport data", err)
		return err
	}

	val.Set(reflect.ValueOf(p))

	return nil
}
