package bimap

type Bimap[K, V comparable] struct {
	fwd map[K]V
	inv map[V]K
}

func New[K, V comparable]() *Bimap[K, V] {
	return &Bimap[K, V]{
		fwd: make(map[K]V),
		inv: make(map[V]K),
	}
}

func NewFromMap[K, V comparable](fwd map[K]V) *Bimap[K, V] {
	inv := make(map[V]K, len(fwd))
	for k, v := range fwd {
		inv[v] = k
	}

	return &Bimap[K, V]{
		fwd: fwd,
		inv: inv,
	}
}

func (b *Bimap[K, V]) Get(key K) (V, bool) {
	v, ok := b.fwd[key]
	return v, ok
}

func (b *Bimap[K, V]) GetInverse(key V) (K, bool) {
	v, ok := b.inv[key]
	return v, ok
}

func (b *Bimap[K, V]) Put(key K, value V) {
	b.fwd[key] = value
	b.inv[value] = key
}
