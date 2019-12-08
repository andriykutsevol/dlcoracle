package routes

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"net/http"
	"strconv"

	"bytes"
	"encoding/binary"	

	"github.com/gertjaap/dlcoracle/store"
	"github.com/gorilla/mux"

	"github.com/gertjaap/dlcoracle/crypto"
	"github.com/gertjaap/dlcoracle/datasources"
	"github.com/gertjaap/dlcoracle/logging"
)

type RPointResponse struct {
	R string
}

func RPointHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	datasourceId, err := strconv.ParseUint(vars["datasource"], 10, 64)
	if err != nil {
		logging.Error.Println("RPointPubKeyHandler - Invalid Datasource: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !datasources.HasDatasource(datasourceId) {
		logging.Error.Println("RPointPubKeyHandler - Invalid Datasource: ", datasourceId)
		http.Error(w, fmt.Sprintf("Invalid datasource %d", datasourceId), http.StatusInternalServerError)
		return
	}

	timestamp, err := strconv.ParseUint(vars["timestamp"], 10, 64)
	if err != nil {
		logging.Error.Println("RPointPubKeyHandler - Invalid Timestamp: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}



	rPoint, err := store.GetRPoint(datasourceId, timestamp)
	if err != nil {
		logging.Error.Println("RPointPubKeyHandler", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	publishedAlready, err := store.IsPublished(rPoint)
	if err != nil {
		logging.Error.Printf("Error determining if this is already published: %s", err.Error())
		return
	}

	if publishedAlready {
		logging.Info.Printf("Already published for data source %d and timestamp %d", datasourceId, timestamp)
		return
	}	


	var a [32]byte
	copy(a[:], crypto.RetrieveKey(crypto.KeyTypeA)[:])
	fmt.Printf("::%s:: keys.go:RPointPubKeyHandler() a: %x \n", os.Args[2][len(os.Args[2])-4:], a)

	k, err := store.GetK(datasourceId, timestamp)
	if err != nil {
		logging.Error.Printf("Could not get signing key for data source %d and timestamp %d : %s", datasourceId, timestamp, err.Error())
		return
	}
	fmt.Printf("::%s:: keys.go:RPointPubKeyHandler() k: %x \n", os.Args[2][len(os.Args[2])-4:], k)

	ds := datasources.GetAllDatasources()[0]

	valueToPublish, err := ds.Value()
	if err != nil {
		logging.Error.Printf("keys.go: Could not retrieve value for data source %d: %s", datasourceId, err.Error())
		return
	}

	// Zero pad the value before signing. Sign expects a [32]byte message
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, uint64(0))
	binary.Write(&buf, binary.BigEndian, uint64(0))
	binary.Write(&buf, binary.BigEndian, uint64(0))
	binary.Write(&buf, binary.BigEndian, valueToPublish)

	signature, err := crypto.ComputeS(a, k, buf.Bytes())
	if err != nil {
		logging.Error.Printf("Could not sign the message: %s", err.Error())
		return
	}
	
	store.Publish(rPoint, valueToPublish, signature)
	


	response := RPointResponse{
		R: hex.EncodeToString(rPoint[:]),
	}

	js, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

type PubKeyResponse struct {
	A string
}

func PubKeyHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("::%s:: keys.go: PubKeyHandler \n", os.Args[2][len(os.Args[2])-4:])

	A, err := crypto.GetPubKey(crypto.KeyTypeA)
	if err != nil {
		logging.Error.Println("PubKeyHandler", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("::%s:: keys.go: PubKeyHandler: A: %x \n", os.Args[2][len(os.Args[2])-4:], A)

	response := PubKeyResponse{
		A: hex.EncodeToString(A[:]),
	}

	js, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
