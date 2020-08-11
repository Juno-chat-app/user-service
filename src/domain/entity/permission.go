package entity

type Permission struct {
	Key   string `bson:"key"`
	Value string `bson:"value"`
}
