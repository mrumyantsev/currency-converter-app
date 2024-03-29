package fsops

import (
	"errors"
	"io"
	"os"
	"path"

	"github.com/mrumyantsev/currency-converter-app/internal/pkg/config"
	"github.com/mrumyantsev/currency-converter-app/pkg/lib/e"
)

const (
	saveDir  = "./save"
	filePerm = 0644
	dirPerm  = 0755
)

type FsOps struct {
	config *config.Config
}

func New(cfg *config.Config) *FsOps {
	return &FsOps{config: cfg}
}

func (f *FsOps) CurrencyData() ([]byte, error) {
	err := makeDirIfNotExist(saveDir)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(path.Join(saveDir, f.config.CurrencySourceFile))
	if err != nil {
		return nil, e.Wrap("could not open file", err)
	}
	defer func() { _ = file.Close() }()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, e.Wrap("could not read all data from file", err)
	}

	return data, nil
}

func (f *FsOps) OverwriteCurrencyDataFile(data []byte) error {
	err := makeDirIfNotExist(saveDir)
	if err != nil {
		return err
	}

	err = os.WriteFile(
		path.Join(saveDir, f.config.CurrencySourceFile),
		data,
		filePerm,
	)
	if err != nil {
		return e.Wrap("could not write file", err)
	}

	return nil
}

func makeDirIfNotExist(path string) error {
	_, err := os.Stat(path)
	if !errors.Is(err, os.ErrNotExist) {
		return e.Wrap("could not check for save directory existence", err)
	}

	if err = os.Mkdir(path, dirPerm); err != nil {
		return e.Wrap("could not make save directory", err)
	}

	return nil
}
