package wallet

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const defaultWalletFile = "./tmp/wallets_%s.data"

type Wallets struct {
	Wallets     map[string]*Wallet
	WalletsFile string
}

func CreateWallets(nodeId, walletsFilePath string) (*Wallets, error) {
	wallets := Wallets{
		Wallets:     make(map[string]*Wallet),
		WalletsFile: walletsFilePath,
	}
	if wallets.WalletsFile == "" {
		wallets.WalletsFile = defaultWalletFile
	}
	err := wallets.LoadFile(nodeId)

	return &wallets, err
}

func (ws *Wallets) AddWallet() (string, error) {
	wallet, err := NewWallet()
	if err != nil {
		return "", err
	}
	address := string(wallet.Address())

	ws.Wallets[address] = wallet

	return address, nil
}

func (ws *Wallets) GetAllAddresses() []string {
	var addresses []string

	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}

	return addresses
}

func (ws Wallets) GetWallet(address string) Wallet {
	return *ws.Wallets[address]
}

func (ws *Wallets) LoadFile(nodeId string) error {
	walletFile := fmt.Sprintf(ws.WalletsFile, nodeId)
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}

	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		return err
	}

	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	var wallets Wallets
	err = decoder.Decode(&(wallets.Wallets))
	if err != nil {
		return err
	}

	ws.Wallets = wallets.Wallets

	return nil
}

func (ws *Wallets) SaveFile(nodeId string) {
	var content bytes.Buffer
	walletFile := fmt.Sprintf(defaultWalletFile, nodeId)

	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws.Wallets)
	if err != nil {
		log.Panic(err)
	}

	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}
