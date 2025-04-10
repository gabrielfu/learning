package main

type ResponseHeaderV0 struct {
	CorrelationID int32
}

type ResponseHeaderV1 struct {
	CorrelationID int32
	TagBuffer     byte
}

type Response struct {
	Header interface{}
	Body   interface{}
}

// Implements BinaryMarshaler interface
func (r Response) MarshalBinary() ([]byte, error) {
	out, err := marshal(r.Header)
	if err != nil {
		return nil, err
	}

	if r.Body != nil {
		body, err := marshal(r.Body)
		if err != nil {
			return nil, err
		}
		out = append(out, body...)
	}
	messageSize := marshalInt32BigEndian(int32(len(out)))
	out = append(messageSize, out...)
	return out, nil
}

type ApiVersionsV4ResponseBody struct {
	ErrorCode      int16
	ApiKeys        Array[ApiVersionsV4ApiKey]
	ThrottleTimeMs int32
	TagBuffer      byte
}

type ApiVersionsV4ApiKey struct {
	ApiKey     int16
	MinVersion int16
	MaxVersion int16
	TagBuffer  byte
}

type DescribeTopicPartitionsV0ResponseBody struct {
	ThrottleTimeMs int32
	Topics         Array[DescribeTopicPartitionsV0Topic]
	NextCursor     byte
	TagBuffer      byte
}

type DescribeTopicPartitionsV0Topic struct {
	ErrorCode                 int16
	TopicName                 CompactString
	TopicID                   UUID
	IsInternal                bool
	Partitions                byte // 1 indicates empty array. TODO: change to array
	TopicAuthorizedOperations int32
	TagBuffer                 byte
}
