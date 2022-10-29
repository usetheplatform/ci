package observer

import (
	"os"

	. "github.com/usetheplatform/ci-system/pkg/common"
)

type Store struct {
	path string
}

func NewStore(path string) Store {
	return Store{path: path}
}

func (s *Store) Exists() (bool, error) {
	if _, err := os.Stat(s.path); err == nil || os.IsExist(err) {
		return true, err
	} else {
		return false, err
	}
}

func (s *Store) Write(content string) error {
	if exists, err := s.Exists(); err == nil && exists == false {
		f, err := os.Create(s.path)
		defer f.Close()
		CheckIfError(err)
	}

	return os.WriteFile(s.path, []byte(content), 0644)
}

func (s *Store) Read() (*string, error) {
	if exists, err := s.Exists(); err == nil && exists == true {
		parts, err := os.ReadFile(s.path)
		content := string(parts)

		return &content, err
	} else {
		return nil, err
	}
}

func (s *Store) Clear() error {
	if exists, err := s.Exists(); err == nil && exists == true {
		return os.Remove(s.path)
	}

	return nil
}
