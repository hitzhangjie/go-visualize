package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Req struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
	Address   string `json:"address"`
	privacy   string `json:"privacy"`
}

func TestJSONMarshaler(t *testing.T) {
	m := Marshalers[JSON]
	req := &Req{
		FirstName: "jie",
		LastName:  "zhang",
		Age:       30,
		Address:   "China",
		privacy:   "something",
	}
	dat, err := m.Marshal(req)
	t.Logf("marshal data: %s", string(dat))
	assert.Nil(t, err)
	assert.NotEmpty(t, dat)

	req2 := Req{}
	err = m.Unmarshal(dat, &req2)
	t.Logf("unmarshal data: %+v", req2)
	assert.Nil(t, err)
	assert.NotEmpty(t, &req2)
}

func TestFORMMarshaler(t *testing.T) {
	m := Marshalers[FORM]
	req := Req{
		FirstName: "jie",
		LastName:  "zhang",
		Age:       30,
		Address:   "China",
		privacy:   "something",
	}
	dat, err := m.Marshal(&req)
	t.Logf("marshal data: %s", string(dat))
	assert.Nil(t, err)
	assert.NotEmpty(t, dat)

	req2 := Req{}
	err = m.Unmarshal(dat, &req2)
	t.Logf("unmarshal data: %+v", req2)
	assert.Nil(t, err)
	assert.NotEmpty(t, &req2)
}
