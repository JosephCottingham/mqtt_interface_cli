package brokerManager

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// Stores the credentials required to connect to a MQTT broker
type Broker struct {
	Uuid       string `json:"uuid"`
	Name       string `json:"name"`
	ClientName string `json:"clientName"`
	Uri        string `json:"uri"`
	Port       int    `json:"port"`
	Username   string `json:"username"`
	Password   string `json:"password"`
}

// Storages the list of mqtt broker credentials
type BrokerData struct {
	Brokers []Broker `json:"brokers"`
	Hash    []byte   `json:"hash"`
}

// Get the path to the users broker storage file.
func getBrokerDataFilepath() string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
	}

	brokerDataFilepath := (dirname + "/.mqttBrokerData")
	// If the file does not exsist then create a new file
	_, err = os.Stat(brokerDataFilepath)
	if errors.Is(err, os.ErrNotExist) {
		os.WriteFile(brokerDataFilepath, []byte{}, 0644)
	}
	return brokerDataFilepath
}

// Get the path to the file that stores the user's password hash.
// There is defintly a better way to this but for now this is how we are checking that the user password is valid and therefore the data can be parsed into the broker json.
func getBrokerDataPassFilepath() string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
	}

	brokerDataFilepath := (dirname + "/.mqttBrokerDataPass")

	_, err = os.Stat(brokerDataFilepath)
	if errors.Is(err, os.ErrNotExist) {
		os.WriteFile(brokerDataFilepath, getHash([]byte("mqtt")), 0644)
	}
	return brokerDataFilepath
}

func ReadBrokerData(password string) (BrokerData, error) {
	brokerDataFilepath := getBrokerDataFilepath()
	brokerData := BrokerData{}

	if checkPasswordHash(password) == false {
		return brokerData, errors.New("Invalid Password")
	}
	dat, err := os.ReadFile(brokerDataFilepath)
	if err != nil {
		fmt.Println(err)
	}
	brokerDataJsonBytesEncrypted := []byte(dat)
	brokerDataJsonBytes, _ := decrypt([]byte(password), brokerDataJsonBytesEncrypted)

	_ = json.Unmarshal(brokerDataJsonBytes, &brokerData)

	return brokerData, nil
}

func WriteBrokerData(brokerData BrokerData, password string) {
	brokerDataFilepath := getBrokerDataFilepath()
	brokerDataJson, _ := json.Marshal(&brokerData)

	brokerDataJsonBytes := []byte(string(brokerDataJson))
	brokerDataJsonBytesEncrypted, err := encrypt([]byte(password), brokerDataJsonBytes)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(brokerDataFilepath, brokerDataJsonBytesEncrypted, 0644)
	if err != nil {
		panic(err)
	}
}

func RemoveBroker(brokerData BrokerData, broker Broker, password string) BrokerData {
	index := getBrokerIndexByUuid(brokerData, broker.Uuid)
	brokerData.Brokers = remove(brokerData.Brokers, index)
	WriteBrokerData(brokerData, password)
	return brokerData
}

func AddBroker(brokerData BrokerData, broker Broker, password string) BrokerData {
	broker.Uuid = createUuid(brokerData)
	brokerData.Brokers = append(brokerData.Brokers, broker)
	WriteBrokerData(brokerData, password)
	return brokerData
}

func getBrokerIndexByUuid(brokerData BrokerData, uuid string) int {
	for i, b := range brokerData.Brokers {
		if b.Uuid == uuid {
			return i
		}
	}
	return -1
}

func createUuid(brokerData BrokerData) string {
	uuid := ""
	for true {
		uuid = randomString(10)
		index := getBrokerIndexByUuid(brokerData, uuid)
		if index == -1 {
			break
		}
	}
	return uuid
}

// Generates a random string of a given length
func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}

// Removes slice element at index(s) and returns new slice
func remove[T any](slice []T, s int) []T {
	return append(slice[:s], slice[s+1:]...)
}

func addKeySpacing(key []byte) []byte {
	spacing_len := 16 - (len(key) % 16)
	spacing := make([]byte, spacing_len)
	return append(key, spacing...)
}

func encrypt(key, data []byte) ([]byte, error) {
	if len(data) <= 0 {
		return data, nil
	}

	key = addKeySpacing(key)

	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	return ciphertext, nil
}

func decrypt(key, data []byte) ([]byte, error) {
	if len(data) <= 0 {
		return data, nil
	}

	key = addKeySpacing(key)

	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func checkPasswordHash(password string) bool {
	brokerDataPassFilepath := getBrokerDataPassFilepath()

	dat, err := os.ReadFile(brokerDataPassFilepath)
	if err != nil {
		fmt.Println(err)
	}

	passwordHash := []byte(dat)
	if equal(getHash([]byte(password)), passwordHash) {
		return true
	}
	return false
}

func getHash(data []byte) []byte {
	hash := sha256.New()
	return hash.Sum(data)
}

func equal(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
