//
// thin layer manage structured value set on cache layer
package redis

type ThinLayer struct {
	cache ICache
}

func NewThinLayer(cache ICache) ThinLayer {
	thinsLayer := ThinLayer{
		cache: cache,
	}

	return thinsLayer
}

// todo :: implement structured cache layer
