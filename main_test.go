package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	want := "test data"
	testkey := "key"
	store.data[testkey] = want
	got := store.data[testkey]
	assert.Equal(t, want, got)
}

func TestFormat(t *testing.T) {
	key := "receipt"
	s := "{\"register_id\":3001,\"receipt_code\":{{receipt.code}},\"customer_id\":10077,\"item_code\":\"978020137962\",\"price\":{{receipt.total.price}}}"
	store.data[key] = "{\"code\":\"30012021012717353301\",\"itemlist\":[],\"total\":{\"discount_amount\":0,\"price\":15,\"taxes\":[]}}"
	want := "{\"register_id\":3001,\"receipt_code\":30012021012717353301,\"customer_id\":10077,\"item_code\":\"978020137962\",\"price\":15}"
	got, err := format(s)
	if assert.NoError(t, err) {
		assert.Equal(t, want, got)
	}
}

func TestFormatExpectation(t *testing.T) {
	expect := "{\"base\":135,\"id\":1,\"product_id\":1,\"purchase\":95,\"tax_group_id\":1}"
	res := "{\"base\":135,\"id\":1,\"product_id\":1,\"purchase\":95,\"tax_group_id\":1}"

	assert.True(t, testExpectation(res, expect))
}
