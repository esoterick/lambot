package transmission

// Request is a generic request for transmission
type Request struct {
	Method    string    `json:"method`
	Arguments Arguments `json:"arguments"`
}

type Arguments struct {
	Fields Fields `json:"fields"`
	IDs
}

type IDs struct {
	IDs []int `json:"ids"`
}

type Fields struct {
	Fields []string `json:"fields"`
}

type TorrentsResponse struct {
	Arguments []interface{} `json:"arguments"`
	Result    string        `json:"result"`
}

// {
// 	"method": "torrent-get",
// 	"arguments": {
// 			"fields": [
// 		"id",
// 		"name",
// 		"status",
// 		"rateDownload",
// 		"rateUpload"
// 		]
// 	}
// }

// {
//     "arguments": {
//         "torrents": []
//     },
//     "result": "success"
// }
