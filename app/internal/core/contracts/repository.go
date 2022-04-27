package contracts

type Repository interface {
	Put(url, key *string) error
	Get(key, url *string) error
}
