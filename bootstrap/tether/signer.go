package tether

import (
	"io"

	"golang.org/x/crypto/ssh"
)

type ContainerSigner struct {
	id string
}

func (c *ContainerSigner) PublicKey() ssh.PublicKey {
	return *c
}

// we're going to ignore everything for the moment as we're repurposing the host key for the id.
// later we may use a genuine host key and an SSH out-of-band request to get the container id.
func (c *ContainerSigner) Sign(rand io.Reader, data []byte) (*ssh.Signature, error) {
	return &ssh.Signature{
		Format: "container-id",
		Blob:   []byte{},
	}, nil
}

func (c ContainerSigner) Type() string {
	return "container-id"
}

func (c ContainerSigner) Marshal() []byte {
	return []byte(c.id)
}

func (c ContainerSigner) Verify(data []byte, sig *ssh.Signature) error {
	return nil
}

func NewSigner(id string) *ContainerSigner {
	signer := &ContainerSigner{
		id: id,
	}

	return signer
}
