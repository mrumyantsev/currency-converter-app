package fsops

import (
	"io"
	"os"

	"github.com/mrumyantsev/currency-converter-app/internal/pkg/config"
	"github.com/mrumyantsev/currency-converter-app/pkg/lib"
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
		return nil, lib.DecorateError("cannot open file", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, lib.DecorateError("cannot read all data from file", err)
	}

	return data, nil
}

func (f *FsOps) OverwriteCurrencyDataFile(data []byte) error {
	err := os.WriteFile(f.config.CurrencySourceFile, data, 0644)
	if err != nil {
		return lib.DecorateError("cannot write file", err)
	}

	return nil
}
