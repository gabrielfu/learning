package main

import "encoding"

type RequestHeaderV2 struct {
	RequestAPIKey     int16
	RequestAPIVersion int16
	CorrelationID     int32
	ClientID          String
	TagBuffer         byte
}

func (h *RequestHeaderV2) Unmarshal(data []byte) (int, error) {
	cursor := 0
	requestApiKey, err := unmarshalInt16BigEndian(data[cursor : cursor+2])
	if err != nil {
		return 0, err
	}
	cursor += 2
	requestApiVersion, err := unmarshalInt16BigEndian(data[cursor : cursor+2])
	if err != nil {
		return 0, err
	}
	cursor += 2
	correlationID, err := unmarshalInt32BigEndian(data[cursor : cursor+4])
	if err != nil {
		return 0, err
	}
	cursor += 4
	clientIDLen, err := unmarshalInt16BigEndian(data[cursor : cursor+2])
	if err != nil {
		return 0, err
	}
	cursor += 2
	clientID := String(data[cursor : cursor+int(clientIDLen)])
	cursor += int(clientIDLen)
	tagBuffer := data[cursor]
	cursor++
	h.RequestAPIKey = requestApiKey
	h.RequestAPIVersion = requestApiVersion
	h.CorrelationID = correlationID
	h.ClientID = clientID
	h.TagBuffer = tagBuffer
	return cursor, nil
}

type Request struct {
	MessageSize int32
	Header      RequestHeaderV2
	Body        encoding.BinaryUnmarshaler
}

// Implements BinaryUnmarshaler interface
func (r *Request) UnmarshalBinary(data []byte) error {
	cursor := 0
	messageSize, err := unmarshalInt32BigEndian(data[:4])
	if err != nil {
		return err
	}
	r.MessageSize = messageSize
	cursor += 4

	n, err := r.Header.Unmarshal(data[cursor:])
	if err != nil {
		return err
	}
	cursor += n

	if r.Header.RequestAPIKey == 75 {
		var body DescribeTopicPartitionsV0RequestBody
		err = body.UnmarshalBinary(data[cursor:])
		if err != nil {
			return err
		}
		r.Body = &body
	}
	return nil
}

type DescribeTopicPartitionsV0RequestBody struct {
	Topics                 []DescribeTopicPartitionsV0RequestTopic
	ResponsePartitionLimit int32
	NextCursor             byte
	TagBuffer              byte
}

func (b *DescribeTopicPartitionsV0RequestBody) UnmarshalBinary(data []byte) error {
	cursor := 0
	topicArrayLen := int(data[cursor]) - 1
	cursor++

	for i := 0; i < topicArrayLen; i++ {
		nameLen := int(data[cursor]) - 1
		cursor++
		name := CompactString(data[cursor : cursor+nameLen])
		cursor += nameLen
		tagBuffer := data[cursor]
		cursor++
		b.Topics = append(b.Topics, DescribeTopicPartitionsV0RequestTopic{
			Name:      name,
			TagBuffer: tagBuffer,
		})
	}
	responsePartitionLimit, err := unmarshalInt32BigEndian(data[cursor : cursor+4])
	if err != nil {
		return err
	}
	b.ResponsePartitionLimit = responsePartitionLimit
	cursor += 4

	b.NextCursor = data[cursor]
	cursor++

	b.TagBuffer = data[cursor]
	cursor++
	return nil
}

type DescribeTopicPartitionsV0RequestTopic struct {
	Name      CompactString
	TagBuffer byte
}

func (b *DescribeTopicPartitionsV0RequestTopic) UnmarshalBinary(data []byte) error {
	l := int(data[0])
	b.Name = CompactString(data[1:l])
	b.TagBuffer = data[l]
	return nil
}
