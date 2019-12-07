package store

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"time"
	"path/filepath"

	"os"


	"github.com/adiabat/btcd/btcec"
	"github.com/boltdb/bolt"
	"github.com/gertjaap/dlcoracle/logging"
	"github.com/gertjaap/dlcoracle/gcfg"
)

var db *bolt.DB

func Init() error {
	database, err := bolt.Open(filepath.Join(gcfg.DataDir, "oracle.db"), 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		logging.Error.Fatal(err)
		return err
	}

	db = database

	// Ensure buckets exist that we need
	err = db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte("Keys"))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte("Publications"))
		return err
	})

	return err
}

func GetRPoint(datasourceId, timestamp uint64) ([33]byte, error) {


	fmt.Printf("::%s:: store.go: ====GetRPoint()==== \n", os.Args[2][len(os.Args[2])-4:])

	var pubKey [33]byte

	privKey, err := GetK(datasourceId, timestamp)

	fmt.Printf("::%s:: store.go: GetRPoint(): privKey: %x \n", os.Args[2][len(os.Args[2])-4:], privKey)

	if err != nil {
		logging.Error.Print(err)
		return pubKey, err
	}

	_, pk := btcec.PrivKeyFromBytes(btcec.S256(), privKey[:])

	copy(pubKey[:], pk.SerializeCompressed())
	
	fmt.Printf("::%s:: store.go: GetRPoint(): pubKey: %x \n", os.Args[2][len(os.Args[2])-4:], pubKey)
	fmt.Printf("::%s:: store.go: ---GetRPoint--- \n", os.Args[2][len(os.Args[2])-4:])
	
	return pubKey, nil
}

func GetK(datasourceId, timestamp uint64) ([32]byte, error) {

	fmt.Printf("::%s:: store.go: ====GetK()==== \n", os.Args[2][len(os.Args[2])-4:])

	var privKey [32]byte
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Keys"))
		key := makeStorageKey(datasourceId, timestamp)

		fmt.Printf("::%s:: store.go: GetK(): key: %x \n", os.Args[2][len(os.Args[2])-4:], key)

		priv := b.Get(key)

		fmt.Printf("::%s:: store.go: GetK(): priv: %x \n", os.Args[2][len(os.Args[2])-4:], priv)

		if priv == nil {
			_, err := rand.Read(privKey[:])
			if err != nil {
				return err
			}
			err = b.Put(key, privKey[:])
			fmt.Printf("::%s:: store.go: GetK() rand.Read(privKey): %x \n", os.Args[2][len(os.Args[2])-4:], privKey)
			return err
		} else {
			fmt.Printf("::%s:: store.go: GetK() b.Get(key): %x \n", os.Args[2][len(os.Args[2])-4:], priv)
			copy(privKey[:], priv)
		}
		return nil
	})

	if err != nil {
		logging.Error.Print(err)
		return privKey, err
	}

	fmt.Printf("::%s:: store.go: ---GetK--- \n", os.Args[2][len(os.Args[2])-4:])

	return privKey, nil
}

func Publish(rPoint [33]byte, value uint64, signature [32]byte) error {
	err := db.Update(func(tx *bolt.Tx) error {

		fmt.Printf("::%s:: store.go: ====Publish()==== \n", os.Args[2][len(os.Args[2])-4:])

		b := tx.Bucket([]byte("Publications"))

		v := b.Get(rPoint[:])
		if v != nil { // only add when not yet exists
			return fmt.Errorf("There is already a value published for this rpoint")
		}

		var buf bytes.Buffer
		binary.Write(&buf, binary.BigEndian, value)
		buf.Write(signature[:])


		fmt.Printf("::%s:: store.go: Publish() rPoint: %x, buf.Bytes(): %x \n", os.Args[2][len(os.Args[2])-4:], rPoint, buf.Bytes())

		err := b.Put(rPoint[:], buf.Bytes())
		return err
	})

	if err != nil {
		logging.Error.Print(err)
		return err
	}

	fmt.Printf("::%s:: store.go: ---Publish--- \n", os.Args[2][len(os.Args[2])-4:])

	return nil
}

func IsPublished(rPoint [33]byte) (bool, error) {
	published := false
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Publications"))
		v := b.Get(rPoint[:])
		if v != nil { // only add when not yet exists
			published = true
		}

		return nil
	})

	if err != nil {
		logging.Error.Print(err)
		return false, err
	}

	return published, nil
}

func GetPublication(rPoint [33]byte) (uint64, [32]byte, error) {


	fmt.Printf("::%s:: store.go: ====GetPublication()==== \n", os.Args[2][len(os.Args[2])-4:])

	value := uint64(0)
	signature := [32]byte{}

	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Publications"))
		v := b.Get(rPoint[:])
		if v == nil { // only add when not yet exists
			return fmt.Errorf("Publication not found")
		}

		fmt.Printf("::%s:: store.go: GetPublication(): rPoint: %x, v: %x \n", os.Args[2][len(os.Args[2])-4:], rPoint, v)

		buf := bytes.NewBuffer(v)

		err := binary.Read(buf, binary.BigEndian, &value)
		if err != nil {
			return err
		}
		copy(signature[:], buf.Next(32))

		fmt.Printf("::%s:: store.go: GetPublication(), buf.Bytes(): %x, signature: %x, value: %d \n", 
		os.Args[2][len(os.Args[2])-4:], buf.Bytes(), signature, value)	


		return nil
	})

	if err != nil {
		logging.Error.Print(err)
		return 0, [32]byte{}, err
	}


	fmt.Printf("::%s:: store.go: ---GetPublication--- \n", os.Args[2][len(os.Args[2])-4:])

	return value, signature, nil
}

func makeStorageKey(datasourceId uint64, timestamp uint64) []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, timestamp)
	binary.Write(&buf, binary.BigEndian, datasourceId)
	return buf.Bytes()
}
