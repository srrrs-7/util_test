package model

type Key string

type Value []byte

type CacheModel struct {
	Key   Key   `json:"key"`
	Value Value `json:"value"`
}
