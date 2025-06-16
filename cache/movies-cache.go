package cache

type DataCache interface {
	Set(key string, value any)
	Get(key string, dest any) bool
}