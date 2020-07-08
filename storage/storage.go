package storage

type Storage interface {
	Start() error
	Stop() error
	Write() error
	Reader() (Reader, error)
}

// Reader interface
type Reader interface {
	// When the key doesn't exist, return nil for the value
	GetCF(cf string, key []byte) ([]byte, error)
	//IterCF(cf string) DBIterator
	Close()
}
