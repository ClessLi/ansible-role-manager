package store

var client Factory

type Factory interface {
	Inventory() InventoryStore
	//Close() error
}

func Client() Factory {
	return client
}

func SetClient(factory Factory) {
	client = factory
}
