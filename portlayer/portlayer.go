package portlayer

import "github.com/vmware/vic/portlayer/storage"

type API interface {
	storage.ImageStorer
}
