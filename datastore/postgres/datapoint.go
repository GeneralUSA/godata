package postgres

import (
	"database/sql"
	"fmt"
	"github.com/YouthBuild-USA/ybdata/datastore"
)

type Datapoint struct {
	id            int
	label         string
	key           string
	subjectTypeId int
	valueTypeId   int

	isNew bool
}

type datapointStatements struct {
	load   *sql.Stmt
	create *sql.Stmt
	update *sql.Stmt
	delete *sql.Stmt
}

func newDatapointStatements(pg *PostgresProvider) datapointStatements {
	s := datapointStatements{}

	s.load = pg.mustPrepare("SELECT id, label, key, subject_type_id, value_type_id FROM datapoint WHERE id = $1")

	return s
}

func (pg PostgresProvider) CreateDatapoint(label, key string, subjectTypeId, valueTypeId int) datastore.Datapoint {
	dp := Datapoint{
		label:         label,
		key:           key,
		subjectTypeId: subjectTypeId,
		valueTypeId:   valueTypeId,
	}
	return (datastore.Datapoint)(dp)
}

func (pg PostgresProvider) LoadDatapoint(id int) (datastore.Datapoint, error) {
	row := pg.statements.datapoint.load.QueryRow(id)
	dp := Datapoint{}
	err := row.Scan(&dp.id, &dp.label, &dp.key, &dp.subjectTypeId, &dp.valueTypeId)
	if err != nil {
		return nil, err
	}
	return &dp, nil
}

func (d Datapoint) Id() int {
	return d.id
}

func (d Datapoint) Key() string {
	return d.key
}

func (d Datapoint) SetKey(key string) {
	d.key = key
}

func (d Datapoint) Label() string {
	return d.label
}

func (d Datapoint) SetLabel(label string) {
	d.label = label
}

func (dp Datapoint) ValueType() (datastore.ValueType, error) {
	return provider.LoadValueType(dp.valueTypeId)
}

func (dp Datapoint) String() string {
	vt, _ := dp.ValueType()
	s := fmt.Sprintf("%v (%v)", dp.label, vt.Name())
	var p datastore.PrimitiveType
	fields := vt.Fields()
	for i := range fields {
		p = fields[i].PrimitiveType()
		s += fmt.Sprintf("\n  %v_%v (%v)", dp.key, fields[i].Key(), p.Label())
	}
	return s
}
