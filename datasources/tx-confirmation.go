package datasources

import(
	"github.com/gertjaap/dlcoracle/gcfg"
)

type TxConfirmation struct {
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

