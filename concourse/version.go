package concourse

import (
	"crypto/sha1"
	"fmt"
)

type Version struct {
	CredentialsSha1 string `json:"credentials_sha1"`
	Server          string `json:"server"`
}

func NewVersion(bytesToSha1 []byte, server string) Version {
	return Version{
		CredentialsSha1: fmt.Sprintf("%x", sha1.Sum(bytesToSha1)),
		Server:          server,
	}
}
