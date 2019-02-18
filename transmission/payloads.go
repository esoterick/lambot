package transmission

// Request is a generic request for transmission
type Request struct {
	Method    string    `json:"method"`
	Arguments Arguments `json:"arguments"`
}

// Arguments are a list of required or optional fields for our requests
type Arguments struct {
	Fields []string `json:"fields"`
	IDs    IDs      `json:"ids"`
}

// IDs is a list of ids to filter out requests to
type IDs struct {
	IDs []int `json:"ids"`
}

// TorrentsResponse is a response struct for torrent-get
type TorrentsResponse struct {
	Arguments TorrentArguments `json:"arguments"`
	Result    string           `json:"result"`
}

// TorrentArguments another layer of nonense
type TorrentArguments struct {
	Torrents []Torrent `json:"torrents"`
}

// Torrent Represents a single torrent
type Torrent struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Status       int    `json:"status"`
	RateDownload int    `json:"rateDownload"`
	RateUpload   int    `json:"rateUpload"`
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
