package data

import (
	"testing"

	"github.com/ElrondNetwork/elrond-go/core/check"
	"github.com/stretchr/testify/assert"
)

func TestInitialSmartContract_OwnerBytes(t *testing.T) {
	t.Parallel()

	ownerBytes := []byte("owner bytes")
	isc := &InitialSmartContract{}
	isc.SetOwnerBytes(ownerBytes)
	recovered := isc.OwnerBytes()

	assert.Equal(t, ownerBytes, recovered)
}

func TestInitialSmartContract_VmTypeBytes(t *testing.T) {
	t.Parallel()

	vmTypeBytes := []byte("vm type bytes")
	isc := &InitialSmartContract{}
	isc.SetVmTypeBytes(vmTypeBytes)
	recovered := isc.VmTypeBytes()

	assert.Equal(t, vmTypeBytes, recovered)
}

func TestInitialSmartContract_Getters(t *testing.T) {
	t.Parallel()

	owner := "owner"
	filename := "filename"
	vmType := "vm type"
	initParams := "init parameters"
	scType := "type"
	isc := &InitialSmartContract{
		Owner:          owner,
		Filename:       filename,
		VmType:         vmType,
		Type:           scType,
		InitParameters: initParams,
	}

	assert.False(t, check.IfNil(isc))
	assert.Equal(t, owner, isc.GetOwner())
	assert.Equal(t, filename, isc.GetFilename())
	assert.Equal(t, vmType, isc.GetVmType())
	assert.Equal(t, scType, isc.GetType())
	assert.Equal(t, initParams, isc.GetInitParameters())
}