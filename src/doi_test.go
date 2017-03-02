package ginDoi

import (
	"testing"
	"reflect"
)

func TestDoiFileMarshalling(t *testing.T) {
	var data = `
title: test
uri: test
authors:
  - testme
  - testme2
`
	mockstruc := CBerry{Title: "test",URI: "test",
		Authors:           []string{"testme","testme2"}}
	val, info := validDoiFile([]byte(data))
	if ! val && reflect.DeepEqual(info,mockstruc){
		t.Logf("Did not properly unmarshal doi data:%+v", info)
		t.Fail()
	}
	t.Log("[OK] Unmarschaling doi files works")
}

func TestGinDataSource(t *testing.T) {
	ds := GinDataSource{GinURL: "https://repo.gin.g-node.org"}
		_,err := ds.GetDoiFile("git@gin.g-node.org:testi/test")
	if err!=nil{
		t.Log(err)
		t.Fail()
	}else{
		t.Log("[OK] can get Doi Files")
	}
}