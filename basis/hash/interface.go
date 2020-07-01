package hash

// Encoder is an interface that can optionally be implemented by
// a struct. It will be called for each field in the struct to check whether
// it should be included in the hash.
type Encoder interface {
	Encode(field string, v interface{}) (bool, error)
}

// MapEncoder is an interface that can optionally be implemented by
// a struct. It will be called when a map-type field is found to ask the
// struct if the map item should be included in the hash.
type MapEncoder interface {
	EncodeMap(field string, k, v interface{}) (bool, error)
}
