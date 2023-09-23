package storage

// import (
// 	"io/fs"
// 	"io/ioutil"
// 	"os"
// )

// const (
// 	MIGRATION_PERFORMED        = "performed"
// 	MIGRATION_NOT_PERFORMED    = "not " + MIGRATION_PERFORMED
// 	MIGRATION_STATUS_FILE_NAME = "migration-status.txt"
// )

// func (s *Storage) Migrate() error {
// 	var (
// 		statusFilePath = s.config.StorageMigrationDir + MIGRATION_STATUS_FILE_NAME
// 		fileContent    []byte
// 		dirEntries     []fs.DirEntry
// 		err            error
// 	)

// 	fileContent, err = readFile(&statusFilePath)
// 	if err != nil {
// 		return err
// 	}

// 	if string(fileContent) == MIGRATION_PERFORMED {
// 		return nil
// 	}

// 	dirEntries, err = os.ReadDir(s.config.StorageMigrationDir)

// 	for _, dirEntry := range dirEntries {
// 		if dirEntry.IsDir() || (dirEntry.Name() == MIGRATION_STATUS_FILE_NAME) {
// 			continue
// 		}

// 		// s.ExecuteQueryFile(s.config.StorageMigrationDir + dirEntry.Name())
// 	}

// 	err = overwriteFile(&statusFilePath, []byte(MIGRATION_PERFORMED))
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func readFile(filePath *string) ([]byte, error) {
// 	var (
// 		f           *os.File
// 		fileContent []byte
// 		err         error
// 	)

// 	f, err = os.Open(*filePath)
// 	if err != nil {
// 		return nil, err
// 	}

// 	fileContent, err = ioutil.ReadAll(f)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = f.Close()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return fileContent, nil
// }

// func overwriteFile(filePath *string, newContent []byte) error {
// 	var (
// 		f   *os.File
// 		err error
// 	)

// 	f, err = os.OpenFile(*filePath, os.O_TRUNC|os.O_WRONLY, 0)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = f.Write(newContent)
// 	if err != nil {
// 		return err
// 	}

// 	err = f.Close()
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
