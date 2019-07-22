package userclient

type Cache interface {
	Get(key string, obj interface{}) error
	Set(key string, obj interface{}) error
	Delete(key string) error
}
