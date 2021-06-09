package config

type StorageConfig interface {
	GetConnStr() string 
	GetTimeout() int
}
