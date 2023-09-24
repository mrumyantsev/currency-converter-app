package fsops

import (
	"io"
	"os"

	"github.com/mrumyantsev/currency-converter/internal/pkg/config"
	"github.com/mrumyantsev/currency-converter/internal/pkg/utils"
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
		return nil, utils.DecorateError("cannot open file", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, utils.DecorateError("cannot read all data from file", err)
	}

	return data, nil
}

func (f *FsOps) OverwriteCurrencyDataFile(data []byte) error {
	err := os.WriteFile(f.config.CurrencySourceFile, data, 0644)
	if err != nil {
		return utils.DecorateError("cannot write file", err)
	}

	return nil
}
