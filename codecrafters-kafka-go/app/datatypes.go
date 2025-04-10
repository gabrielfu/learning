package main

type String string

func (s String) MarshalBinary() ([]byte, error) {
	b := []byte(s)
	b = append(marshalInt16BigEndian(int16(len(b)+1)), b...)
	return b, nil
}

type CompactString string

func (s CompactString) MarshalBinary() ([]byte, error) {
	b := []byte(s)
	b = append(marshalInt8BigEndian(int8(len(b)+1)), b...)
	return b, nil
}

type UUID [16]byte

func (u UUID) MarshalBinary() ([]byte, error) {
	b := make([]byte, 16)
	for i := 0; i < 16; i++ {
		b[i] = u[i]
	}
	return b, nil
}

type Array[T any] []T

func (a Array[T]) MarshalBinary() ([]byte, error) {
	out := marshalInt8BigEndian(int8(len(a)) + 1)
	for _, item := range a {
		b, err := marshal(item)
		if err != nil {
			return nil, err
		}
		out = append(out, b...)
	}
	return out, nil
}
