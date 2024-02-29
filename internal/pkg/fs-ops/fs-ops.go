package fsops

import (
	"errors"
	"io"
	"log"
	"os"

	"github.com/mrumyantsev/currency-converter-app/internal/pkg/config"
	"github.com/mrumyantsev/currency-converter-app/pkg/lib/e"
)

type FsOps struct {
	config *config.Config
}

func New(cfg *config.Config) *FsOps {
	return &FsOps{
		config: cfg,
	}
}

func (f *FsOps) GetCurrencyData() ([]byte, error) {
	file, err := os.Open(f.config.CurrencySourceFile)
	if err != nil {
		return nil, e.Wrap("cannot open file", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, e.Wrap("cannot read all data from file", err)
	}

	return data, nil
}

func (f *FsOps) OverwriteCurrencyDataFile(data []byte) error {
	path := "sample"
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}

	err := os.WriteFile(f.config.CurrencySourceFile, data, 0777)
	if err != nil {
		return e.Wrap("cannot write file", err)
	}

	return nil
}
