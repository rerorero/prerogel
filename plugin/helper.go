package plugin

import (
	"encoding/binary"
	"fmt"
	"hash/fnv"

	"github.com/gogo/protobuf/types"
)

// HashPartition calculates hash of id then mod
func HashPartition(id VertexID, nrOfPartitions uint64) (uint64, error) {
	h := fnv.New64()
	if _, err := h.Write([]byte(string(id))); err != nil {
		return 0, err
	}
	return h.Sum64() % nrOfPartitions, nil
}

// ConvertStringToAny converts string to any
func ConvertStringToAny(val interface{}) (*types.Any, error) {
	if val == nil {
		return nil, nil
	}
	s, ok := val.(string)
	if !ok {
		return nil, fmt.Errorf("not string message: %#v", val)
	}
	return &types.Any{Value: []byte(s)}, nil
}

// ConvertAntToString converts any to string
func ConvertAntToString(pb *types.Any) (string, error) {
	if pb == nil {
		return "", nil
	}
	return string(pb.Value), nil
}

// ConvertUint32ToAny converts interface as a uint32 to any
func ConvertUint32ToAny(val interface{}) (*types.Any, error) {
	n, ok := val.(uint32)
	if !ok {
		return nil, fmt.Errorf("not uint32 value: %#v", val)
	}
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, n)
	return &types.Any{Value: b}, nil
}

// ConvertAnyToUint32 converts any to uint32
func ConvertAnyToUint32(pb *types.Any) (uint32, error) {
	if pb == nil {
		return 0, nil
	}
	if len(pb.Value) != 4 {
		return 0, fmt.Errorf("invalid uint32 message buffer length: %d", len(pb.Value))
	}
	return binary.BigEndian.Uint32(pb.Value), nil
}
