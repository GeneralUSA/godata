package datastore

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strings"
)

type HStore struct {
	data map[string]string

	changed []string
	removed []string
}

func NewHStore() *HStore {
	return &HStore{
		data:    make(map[string]string),
		changed: make([]string, 0),
		removed: make([]string, 0),
	}
}

func (h *HStore) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	as_bytes := src.([]byte)
	if len(as_bytes) == 0 {
		return nil
	}
	parts := bytes.Split(as_bytes, []byte(", "))

	h.data = make(map[string]string)
	for i := 0; i < len(parts); i++ {
		sub_parts := bytes.Split(parts[i], []byte("=>"))
		key := string(bytes.Trim(sub_parts[0], "\""))
		val := string(bytes.Trim(sub_parts[1], "\""))
		h.data[key] = val
	}

	return nil
}

func (h HStore) Value() (driver.Value, error) {
	pairs := make([]string, 0, len(h.data))
	for k, v := range h.data {
		pairs = append(pairs, fmt.Sprintf(`"%s"=>"%s"`, k, v))
	}
	return strings.Join(pairs, ","), nil
}

func (h HStore) Get(key string) (string, bool) {
	d, ok := h.data[key]
	return d, ok
}

func (h *HStore) Set(key string, value string) {
	h.changed = append(h.changed, key)
	h.data[key] = value
}

func (h *HStore) Remove(key string) {
	delete(h.data, key)
	h.removed = append(h.removed, key)
}

func (h HStore) Keys() []string {
	keys := make([]string, len(h.data))
	i := 0
	for k := range h.data {
		keys[i] = k
		i++
	}
	return keys
}

func (h HStore) update(id int, update *sql.Stmt, delete *sql.Stmt) {
	tx, _ := db.Begin()

	updateStatement := tx.Stmt(update)
	deleteStatement := tx.Stmt(delete)

	for k := range h.changed {
		updateStatement.Exec(h.changed[k], h.data[h.changed[k]], id)
	}

	for k := range h.removed {
		fmt.Println(deleteStatement.Exec(h.removed[k], id))
	}

	tx.Commit()
	h.clearChanges()
}

func (h HStore) clearChanges() {
	h.changed = make([]string, 0)
	h.removed = make([]string, 0)
}
