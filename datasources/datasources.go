package datasources

import (
	"fmt"
)


type DatasourceType uint64

const (
	Price	DatasourceType = 1
	Event  	DatasourceType = 2
)


type Datasource interface {
	Id() uint64
	DsType() DatasourceType
	Name() string
	Description() string
	Value() (uint64, error)
	Interval() uint64
}

func GetAllDatasources() []Datasource {
	var datasources []Datasource
	datasources = append(datasources, &UsdBtcRoundedRandom{})
	datasources = append(datasources, &EurBtcRounded{})
	datasources = append(datasources, &TxConfirmation{})
	return datasources
}

func GetDatasource(id uint64) (Datasource, error) {
	switch id {
	case 1:
		return &UsdBtcRoundedRandom{}, nil
	case 2:
		return &EurBtcRounded{}, nil
	case 3:
		return &TxConfirmation{}, nil
	default:
		return nil, fmt.Errorf("Data source with ID %d not known", id)
	}
}

func HasDatasource(id uint64) bool {
	return (id <= 2)
}
