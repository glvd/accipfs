package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"hash"
	"reflect"
)

// Key ...
const key = "trfs_hash"

// DefaultOption ...
var DefaultOption = &Options{
	Hash: hmac.New(func() hash.Hash {
		return sha256.New()
	}, []byte(key)),
	TagName: "hash",
	ZeroNil: false,
}

// ErrNotStringer is returned when there's an error with hash:"string"
type ErrNotStringer struct {
	Field string
}

// Options are options that are available for hashing.
type Options struct {
	// Encoder is the hash function to use. If this isn't set, it will
	// default to sha256.
	Hash hash.Hash

	// TagName is the struct tag to look at when hashing the structure.
	// By default this is "hash".
	TagName string

	// ZeroNil is flag determining if nil pointer should be treated equal
	// to a zero value of pointed type. By default this is false.
	ZeroNil bool
}

type walker struct {
	h       hash.Hash
	tag     string
	zeronil bool
}

type visitOpts struct {
	// Flags are a bitmask of flags to affect behavior of this visit
	Flags visitFlag

	// Information about the struct containing this field
	Struct      interface{}
	StructField string
}

// Error implements error for ErrNotStringer
func (ens *ErrNotStringer) Error() string {
	return fmt.Sprintf("hashstructure: %s has hash:\"string\" set, but does not implement fmt.Stringer", ens.Field)
}

// Sum returns the hash value of an arbitrary value.
//
// If opts is nil, then default options will be used. See HashOptions
// for the default values. The same *HashOptions value cannot be used
// concurrently. None of the values within a *HashOptions struct are
// safe to read/write while hashing is being done.
//
// Notes on the value:
//
//   * Unexported fields on structs are ignored and do not affect the
//     hash value.
//
//   * Adding an exported field to a struct with the zero value will change
//     the hash value.
//
// For structs, the hashing can be controlled using tags. For example:
//
//    struct {
//        ID string
//        UUID string `hash:"ignore"`
//    }
//
// The available tag values are:
//
//   * "ignore" or "-" - The field will be ignored and not affect the hash code.
//
//   * "set" - The field will be treated as a set, where ordering doesn't
//             affect the hash code. This only works for slices.
//
//   * "string" - The field will be hashed as a string, only works when the
//                field implements fmt.Stringer
//
func Sum(v interface{}) ([]byte, error) {
	// Reset the hash
	DefaultOption.Hash.Reset()

	// Create our walker and walk the structure
	w := &walker{
		h:       DefaultOption.Hash,
		tag:     DefaultOption.TagName,
		zeronil: DefaultOption.ZeroNil,
	}
	return w.visit(reflect.ValueOf(v), nil)
}

func (w *walker) visit(v reflect.Value, opts *visitOpts) ([]byte, error) {
	t := reflect.TypeOf(0)

	// Loop since these can be wrapped in multiple layers of pointers
	// and interfaces.
	for {
		// If we have an interface, dereference it. We have to do this up
		// here because it might be a nil in there and the check below must
		// catch that.
		if v.Kind() == reflect.Interface {
			v = v.Elem()
			continue
		}

		if v.Kind() == reflect.Ptr {
			if w.zeronil {
				t = v.Type().Elem()
			}
			v = reflect.Indirect(v)
			continue
		}

		break
	}

	// If it is nil, treat it like a zero.
	if !v.IsValid() {
		v = reflect.Zero(t)
	}

	// Binary writing can use raw ints, we have to convert to
	// a sized-int, we'll choose the largest...
	switch v.Kind() {
	case reflect.Int:
		v = reflect.ValueOf(int64(v.Int()))
	case reflect.Uint:
		v = reflect.ValueOf(uint64(v.Uint()))
	case reflect.Bool:
		var tmp int8
		if v.Bool() {
			tmp = 1
		}
		v = reflect.ValueOf(tmp)
	}

	k := v.Kind()

	// We can shortcut numeric values by directly binary writing them
	if k >= reflect.Int && k <= reflect.Complex64 {
		// A direct hash calculation
		w.h.Reset()
		err := binary.Write(w.h, binary.LittleEndian, v.Interface())
		return w.h.Sum(nil), err
	}

	switch k {
	case reflect.Array:
		var h []byte
		l := v.Len()
		for i := 0; i < l; i++ {
			current, err := w.visit(v.Index(i), nil)
			if err != nil {
				return nil, err
			}

			h = hashUpdateOrdered(w.h, h, current)
		}

		return h, nil

	case reflect.Map:
		var includeMap MapEncoder
		if opts != nil && opts.Struct != nil {
			if v, ok := opts.Struct.(MapEncoder); ok {
				includeMap = v
			}
		}

		// Build the hash for the map. We do this by XOR-ing all the key
		// and value hashes. This makes it deterministic despite ordering.
		var h []byte
		for _, k := range v.MapKeys() {
			v := v.MapIndex(k)
			if includeMap != nil {
				incl, err := includeMap.EncodeMap(
					opts.StructField, k.Interface(), v.Interface())
				if err != nil {
					return nil, err
				}
				if !incl {
					continue
				}
			}

			kh, err := w.visit(k, nil)
			if err != nil {
				return nil, err
			}
			vh, err := w.visit(v, nil)
			if err != nil {
				return nil, err
			}

			fieldHash := hashUpdateOrdered(w.h, kh, vh)
			h = hashUpdateUnordered(h, fieldHash)
		}

		return h, nil

	case reflect.Struct:
		parent := v.Interface()
		var include Encoder
		if impl, ok := parent.(Encoder); ok {
			include = impl
		}

		t := v.Type()
		h, err := w.visit(reflect.ValueOf(t.Name()), nil)
		if err != nil {
			return nil, err
		}

		l := v.NumField()
		for i := 0; i < l; i++ {
			if innerV := v.Field(i); v.CanSet() || t.Field(i).Name != "_" {
				var f visitFlag
				fieldType := t.Field(i)
				if fieldType.PkgPath != "" {
					// Unexported
					continue
				}

				tag := fieldType.Tag.Get(w.tag)
				if tag == "ignore" || tag == "-" {
					// Ignore this field
					continue
				}

				// if string is set, use the string value
				if tag == "string" {
					if impl, ok := innerV.Interface().(fmt.Stringer); ok {
						innerV = reflect.ValueOf(impl.String())
					} else {
						return nil, &ErrNotStringer{
							Field: v.Type().Field(i).Name,
						}
					}
				}

				// Check if we implement encoder and check it
				if include != nil {
					incl, err := include.Encode(fieldType.Name, innerV)
					if err != nil {
						return nil, err
					}
					if !incl {
						continue
					}
				}

				switch tag {
				case "set":
					f |= visitFlagSet
				}

				kh, err := w.visit(reflect.ValueOf(fieldType.Name), nil)
				if err != nil {
					return nil, err
				}

				vh, err := w.visit(innerV, &visitOpts{
					Flags:       f,
					Struct:      parent,
					StructField: fieldType.Name,
				})
				if err != nil {
					return nil, err
				}

				fieldHash := hashUpdateOrdered(w.h, kh, vh)
				h = hashUpdateUnordered(h, fieldHash)
			}
		}

		return h, nil

	case reflect.Slice:
		// We have two behaviors here. If it isn't a set, then we just
		// visit all the elements. If it is a set, then we do a deterministic
		// hash code.
		var h []byte
		var set bool
		if opts != nil {
			set = (opts.Flags & visitFlagSet) != 0
		}
		l := v.Len()
		for i := 0; i < l; i++ {
			current, err := w.visit(v.Index(i), nil)
			if err != nil {
				return nil, err
			}

			if set {
				h = hashUpdateUnordered(h, current)
			} else {
				h = hashUpdateOrdered(w.h, h, current)
			}
		}

		return h, nil

	case reflect.String:
		// Directly hash
		w.h.Reset()
		_, err := w.h.Write([]byte(v.String()))
		return w.h.Sum(nil), err

	default:
		return nil, fmt.Errorf("unknown kind to hash: %s", k)
	}

}

func hashUpdateOrdered(h hash.Hash, a, b []byte) []byte {
	// For ordered updates, use a real hash function
	h.Reset()

	// We just panic if the binary writes fail because we are writing
	// an int64 which should never be fail-able.
	e1 := binary.Write(h, binary.LittleEndian, a)
	e2 := binary.Write(h, binary.LittleEndian, b)
	if e1 != nil {
		panic(e1)
	}
	if e2 != nil {
		panic(e2)
	}

	return h.Sum(nil)
}

func hashUpdateUnordered(a, b []byte) []byte {
	if len(a) < len(b) {
		tmp := b
		b = a
		a = tmp
	}
	for i := range a {
		if i >= len(b) {
			a[i] = a[i] ^ 0
			continue
		}
		a[i] = a[i] ^ b[i]
	}
	return a
}

// visitFlag is used as a bitmask for affecting visit behavior
type visitFlag uint

const (
	visitFlagInvalid visitFlag = iota
	visitFlagSet               = iota << 1
)
