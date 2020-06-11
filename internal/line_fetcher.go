package topfew

type LineFetcher interface {
	Next() ([]byte, error)
}
