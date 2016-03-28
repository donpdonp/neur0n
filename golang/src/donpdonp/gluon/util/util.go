package util

import (
  "encoding/base64"
  "crypto/sha1"

  "gopkg.in/satori/go.uuid.v0"
)


func Sha1Base64(word string) string {
  word_sha := sha1.Sum([]byte(word))
  word_b64 := base64.StdEncoding.EncodeToString(word_sha[:])
  return word_b64
}

func Snowflake() string {
  return uuid.NewV4().String()[0:8]
}