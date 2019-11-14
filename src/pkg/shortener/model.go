package shortener

import "time"

type Redirect struct {
	Code      string    `json:"code" bson:"code" msgpack:"code"`
	URL       string    `json:"url" bson:"url" msgpack:"url" validate:"empty=false & format=url"`
	CreatedAt time.Time `json:"created_at" bson:"created_at" msgpack:"created_at"`
}
