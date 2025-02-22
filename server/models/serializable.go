package models

type Serializable interface {
	ToJson() string
	FromJson(jsonData string) Serializable
}
