package readfile

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// ! data type any: but must pass pointer of data access json
func ReadJSON(filePath string, data interface{}) error {
	jsonFile, err := os.Open(filePath)
	// if we os.Open returns an error then handle it
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	// read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, data)
	if err != nil {
		return err
	}
	return nil
}
