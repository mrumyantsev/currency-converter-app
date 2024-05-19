package fsops

import (
	"errors"
	"io"
	"os"
	"path"

	"github.com/mrumyantsev/currency-converter-app/internal/pkg/config"
	"github.com/mrumyantsev/currency-converter-app/pkg/lib/errlib"
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
		return nil, errlib.Wrap(err, "could not open file")
	}
	defer func() { _ = file.Close() }()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, errlib.Wrap(err, "could not read all data from file")
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
		return errlib.Wrap(err, "could not write file")
	}

	return nil
}

func makeDirIfNotExist(path string) error {
	_, err := os.Stat(path)
	if !errors.Is(err, os.ErrNotExist) {
		return errlib.Wrap(err, "could not check for save directory existence")
	}

	if err = os.Mkdir(path, dirPerm); err != nil {
		return errlib.Wrap(err, "could not make save directory")
	}

	return nil
}
