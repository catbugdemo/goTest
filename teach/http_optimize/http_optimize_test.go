package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestOptimize(t *testing.T) {
	get, e := http.Get("http://127.0.0.1:8083/optimize")
	assert.Nil(t, e)

	all, e := ioutil.ReadAll(get.Body)
	assert.Nil(t, e)
	fmt.Println(string(all))
}
