package main

func handleRequest(request Request) Response {
	if request.Header.RequestAPIKey == 18 {
		return handleApiVersions(request)
	}
	if request.Header.RequestAPIKey == 75 {
		return handleDescribeTopicPartitions(request)
	}
	return Response{
		Header: ResponseHeaderV0{
			CorrelationID: request.Header.CorrelationID,
		},
		Body: nil,
	}
}

func handleApiVersions(request Request) Response {
	if request.Header.RequestAPIVersion < 0 || request.Header.RequestAPIVersion > 4 {
		return Response{
			Header: ResponseHeaderV0{
				CorrelationID: request.Header.CorrelationID,
			},
			Body: ApiVersionsV4ResponseBody{
				ErrorCode: 35,
			},
		}
	}

	return Response{
		Header: ResponseHeaderV0{
			CorrelationID: request.Header.CorrelationID,
		},
		Body: ApiVersionsV4ResponseBody{
			ErrorCode: 0,
			ApiKeys: []ApiVersionsV4ApiKey{
				{
					ApiKey:     18,
					MinVersion: 0,
					MaxVersion: 4,
				},
				{
					ApiKey:     75,
					MinVersion: 0,
					MaxVersion: 0,
				},
			},
			ThrottleTimeMs: 0,
		},
	}
}

func handleDescribeTopicPartitions(request Request) Response {
	requestBody := request.Body.(*DescribeTopicPartitionsV0RequestBody)
	topics := make([]DescribeTopicPartitionsV0Topic, len(requestBody.Topics))
	for i, topic := range requestBody.Topics {
		topics[i] = DescribeTopicPartitionsV0Topic{
			ErrorCode:                 3,
			TopicName:                 topic.Name,
			TopicID:                   UUID{},
			IsInternal:                false,
			Partitions:                1,
			TopicAuthorizedOperations: 3576,
		}
	}
	return Response{
		Header: ResponseHeaderV1{
			CorrelationID: request.Header.CorrelationID,
		},
		Body: DescribeTopicPartitionsV0ResponseBody{
			ThrottleTimeMs: 0,
			Topics:         topics,
			NextCursor:     255,
		},
	}
}
