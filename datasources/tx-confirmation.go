package datasources

import(
	"github.com/gertjaap/dlcoracle/gcfg"
	"errors"
)

type TxConfirmation struct {

	txids [][]byte

} 



func (ds *TxConfirmation) Id() uint64 {
	return 3
}


func (ds *TxConfirmation) DsType() DatasourceType {
	return Event
}


func (ds *TxConfirmation) Name() string {
	return "Tx Confirmation"
}


func (ds *TxConfirmation) Description() string {
	return "Publishes 0 if the transaction has not yet been confirmed, 1 in other case"
}


func (ds *TxConfirmation) Interval() uint64 {
	return gcfg.Interval
}



func (ds *TxConfirmation) Value() (uint64, error) {


	return 1, nil
}


func (ds *TxConfirmation) GetTxs() ([][]byte, error) {

	if ds.DsType() != Event{
		return nil, errors.New("GetTx: datasource typme must to be Event")
	}

	return ds.txids, nil
}

func (ds *TxConfirmation) SetTx(txid []byte) error {

	ds.txids = append(ds.txids, txid)

	return nil
}


func (ds *TxConfirmation) IsTxConfirmed(txid []byte) bool {

	return true
}

