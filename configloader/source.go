package configloader

type Source interface {
	Get(key string) (string, bool)
}
