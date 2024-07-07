package files

import (
	"bot-adviser/lib/e"
	"bot-adviser/storage"
	"encoding/gob"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

type Storage struct {
	basePath string
}

const defaultPerm = 0774

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Save(page *storage.Page) (err error) {
	defer func() { err = e.WrapIfErr("can't save", err) }()

	fPath := filepath.Join(s.basePath, page.UserName)

	if err := os.MkdirAll(fPath, defaultPerm); err != nil {
		return err
	}

	fName, err := fileName(page)
	if err != nil {
		return err
	}

	fPath = filepath.Join(fPath, fName)

	file, err := os.Create(fPath)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()
	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}
	return nil

}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() { err = e.WrapIfErr("can't save", err) }()

	fPath := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(fPath)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, storage.ErrNoSavedPages
	}

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(files))
	file := files[n]

	return s.decodePage(filepath.Join(fPath, file.Name()))
	// open decode
}

func (s Storage) Remove(p *storage.Page) error {
	fileName, err := fileName(p)
	if err != nil {
		return e.Wrap("can't remove file", err)
	}
	path := filepath.Join(s.basePath, p.UserName, fileName)

	if err := os.Remove(path); err != nil {
		msg := fmt.Sprintf("can't remove file:\n %s", path)
		return e.Wrap(msg, err)
	}
	return nil
}

func (s Storage) IsExists(p *storage.Page) (bool, error) {
	fileName, err := fileName(p)
	if err != nil {
		return false, e.Wrap("can't check file", err)
	}
	path := filepath.Join(s.basePath, p.UserName, fileName)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		msg := fmt.Sprintf("file does not exists : %s", path)
		return false, e.Wrap(msg, err)
	}

	return true, nil

}

func (s Storage) decodePage(filepath string) (*storage.Page, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, e.Wrap("can't decode page", err)
	}
	defer func() { _ = f.Close() }()
	var p storage.Page
	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, e.Wrap("can't decode page", err)
	}
	return &p, nil

}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
