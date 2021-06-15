package db

const (
	base62Members = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func getShortenURL(id int64) string {
	var res []byte
	N := int64(len(base62Members))
	for ; id > 0; id /= N {
		res = append(res, base62Members[id%N])
	}
	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
	}

	return string(res)
}
