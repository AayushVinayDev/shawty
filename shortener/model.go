package shortener

type Redirect struct {
	Code      string `json:"code" bson:"code" msgpack:"code" protobuf:"bytes,1,opt,name=code"`
	URL       string `json:"url" bson:"url" msgpack:"url" validate:"empty=false & format=url" protobuf:"bytes,2,opt,name=url"`
	CreatedAt int64  `json:"created_at" bson:"created_at" msgpack:"created_at" protobuf:"varint,3,opt,name=created_at,json=createdAt"`
}
