package main

type MockSessionStore struct {
	Session *Session
}

func (store MockSessionStore) Find(string) (*Session, error) {
	return store.Session, nil
}

func (store MockSessionStore) Save(*Session) error {
	return nil
}

func (store MockSessionStore) Delete(*Session) error {
	return nil
}
