package fsops

import (
	"io"
	"os"

	"github.com/mrumyantsev/currency-converter/internal/pkg/config"

	"github.com/mrumyantsev/fastlog"
)

type FsOps struct {
	config *config.Config
}

func New(cfg *config.Config) *FsOps {
	return &FsOps{
		config: cfg,
	}
}

func (f *FsOps) GetCurrencyData() []byte {
	fastlog.Debug("getting data from local file")

	file, err := os.Open(f.config.CurrencySourceFile)
	if err != nil {
		fastlog.Fatal("cannot open file", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		fastlog.Fatal("cannot read all data from file", err)
	}

	return data
}

func (f *FsOps) OverwriteCurrencyDataFile(data []byte) {
	err := os.WriteFile(f.config.CurrencySourceFile, data, 0644)
	if err != nil {
		fastlog.Fatal("cannot write file", err)
	}
}
