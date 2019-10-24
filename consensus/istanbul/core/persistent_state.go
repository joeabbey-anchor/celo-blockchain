package core

import (
	"fmt"

	"github.com/ethereum/go-ethereum/log"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
)

func (c *core) savePrepareMessageToDisk(
	messageType uint64,
	roundNumber *big.Int,
	sequenceNumber *big.Int,
	msg []byte) error {
	fileName, err := c.generateFileName(messageType, roundNumber, sequenceNumber)
	if err != nil {
		return err
	}
	err2 := writeToDisk(fileName, msg)
	log.Debug("savePrepareMessageToDisk/wrote file to the disk", "file", fileName, "error", err2)
	return err2
}

// getPreprepareMessageFromDisk returns the prepared message
// If the file does not exist, it returns (nil, nil).
// It returns an error for all other failure cases.
func (c *core) getPreprepareMessageFromDisk(
	messageType uint64,
	roundNumber *big.Int,
	sequenceNumber *big.Int) ([]byte, error) {
	fileName, err := c.generateFileName(messageType, roundNumber, sequenceNumber)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadFile(fileName)
	if os.IsNotExist(err) {
		log.Debug("getPreprepareMessageFromDisk/file does not exist", "file", fileName)
		return nil, nil
	}
	log.Debug("getPreprepareMessageFromDisk/file found on the disk", "file", fileName)
	return data, err
}

// deleteMessageFromDisk deletes all files from this round and the sequence number `sequenceNumber`
func (c *core) deleteMessageFromDisk(
	roundNumber *big.Int, sequenceNumber *big.Int) error {
	dir := c.backend.GetDataDir()
	// This pattern must be  similar to the filenames generated by
	// generateFileName function.
	filePattern := filepath.Join(dir,
		fmt.Sprintf("geth_istanbul_sequence_%s_round_%s_type_.*",
			sequenceNumber.String(), roundNumber.String()))
	files, err := filepath.Glob(filePattern)
	if err != nil {
		panic("File pattern is bad: " + filePattern)
	}
	for i, file := range files {
		log.Debug("Deleting file", "file", file, "index", i, "total", len(files))
		err := os.Remove(file)
		if err == nil {
			log.Debug("Deleted file", "file", file)
		} else {
			log.Error("Failed to delete file", "file", file)
		}
	}
	return nil
}

func (c *core) generateFileName(
	messageType uint64,
	roundNumber *big.Int,
	sequenceNumber *big.Int) (string, error) {
	dir := c.backend.GetDataDir()
	fileName := fmt.Sprintf("geth_istanbul_sequence_%s_round_%s_type_%d",
		sequenceNumber.String(), roundNumber.String(), messageType)
	return filepath.Join(dir, fileName), nil
}

func writeToDisk(filePath string, data []byte) error {
	fp, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer fp.Close()

	_, err2 := fp.Write(data)
	if err2 != nil {
		return err2
	}
	err3 := fp.Sync()
	if err3 != nil {
		return err3
	}
	fpDir, err4 := os.Open(filepath.Dir(filePath))
	if err4 != nil {
		return err4
	}
	log.Debug("Syncing dir %s\n", fpDir.Name())
	err5 := fpDir.Sync()
	if err5 != nil {
		return err5
	}

	err6 := fp.Close()
	return err6
}