package transmission

// Request is a generic request for transmission
type Request struct {
	Method    string    `json:"method"`
	Arguments Arguments `json:"arguments"`
}

// Arguments are a list of required or optional fields for our requests
type Arguments struct {
	Fields Fields `json:"fields"`
	IDs    IDs    `json:"ids"`
}

// IDs is a list of ids to filter out requests to
type IDs struct {
	IDs []int `json:"ids"`
}

// Fields is a list of fields to filter our requests to
type Fields struct {
	Fields []string `json:"fields"`
}

// TorrentsResponse is a response struct for torrent-get
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
