package slashing_protectors

import (
	"fmt"
	"github.com/bloxapp/KeyVault/encryptors"
	"github.com/bloxapp/KeyVault/stores/in_memory"
	pb "github.com/wealdtech/eth2-signer-api/pb/v1"
	e2types "github.com/wealdtech/go-eth2-types/v2"
	hd "github.com/wealdtech/go-eth2-wallet-hd/v2"
	types "github.com/wealdtech/go-eth2-wallet-types/v2"
	"testing"
)

func setup() (VaultSlashingProtector, types.Account,error) {
	if err := e2types.InitBLS(); err != nil { // very important!
		return nil,nil,err
	}

	store := in_memory.NewInMemStore()
	wallet,err := hd.CreateWallet("test",[]byte(""), store, encryptors.NewPlainTextEncryptor())
	if err != nil {
		return nil,nil,err
	}
	err = wallet.Unlock([]byte(""))
	if err != nil {
		return nil,nil,err
	}

	account1,err := wallet.CreateAccount("1",[]byte(""))
	if err != nil {
		return nil,nil,err
	}

	protector := NewNormalProtection(store)
	protector.SaveAttestation(account1, &pb.SignBeaconAttestationRequest{
		Id:                   nil,
		Domain:               []byte("domain"),
		Data:                 &pb.AttestationData{
			Slot:                 30,
			CommitteeIndex:       5,
			BeaconBlockRoot:      []byte("A"),
			Source:               &pb.Checkpoint{
				Epoch:                2,
				Root:                 []byte("B"),
			},
			Target:               &pb.Checkpoint{
				Epoch:                3,
				Root:                 []byte("C"),
			},
		},
	})
	protector.SaveAttestation(account1, &pb.SignBeaconAttestationRequest{
		Id:                   nil,
		Domain:               []byte("domain"),
		Data:                 &pb.AttestationData{
			Slot:                 30,
			CommitteeIndex:       5,
			BeaconBlockRoot:      []byte("B"),
			Source:               &pb.Checkpoint{
				Epoch:                3,
				Root:                 []byte("C"),
			},
			Target:               &pb.Checkpoint{
				Epoch:                4,
				Root:                 []byte("D"),
			},
		},
	})
	return protector, account1,nil
}


func TestCorrectAttestation(t *testing.T) {
	protector,account,err := setup()
	if err != nil {
		t.Error(err)
		return
	}

	t.Run("Simple1",func(t *testing.T) {
		res, err := protector.IsSlashableAttestation(account, &pb.SignBeaconAttestationRequest{
			Id:                   nil,
			Domain:               []byte("domain"),
			Data:                 &pb.AttestationData{
				Slot:                 30,
				CommitteeIndex:       4,
				BeaconBlockRoot:      []byte("A"),
				Source:               &pb.Checkpoint{
					Epoch:                2,
					Root:                 []byte("B"),
				},
				Target:               &pb.Checkpoint{
					Epoch:                3,
					Root:                 []byte("C"),
				},
			},
		})

		if err != nil {
			t.Error(err)
		}
		if res != true {
			t.Error(fmt.Errorf("non correct attestation found no slashable"))
		}
	})

	t.Run("Simple2",func(t *testing.T) {
		res, err := protector.IsSlashableAttestation(account, &pb.SignBeaconAttestationRequest{
			Id:                   nil,
			Domain:               []byte("domain"),
			Data:                 &pb.AttestationData{
				Slot:                 30,
				CommitteeIndex:       4,
				BeaconBlockRoot:      []byte("A"),
				Source:               &pb.Checkpoint{
					Epoch:                3,
					Root:                 []byte("B"),
				},
				Target:               &pb.Checkpoint{
					Epoch:                4,
					Root:                 []byte("C"),
				},
			},
		})

		if err != nil {
			t.Error(err)
		}
		if res != true {
			t.Error(fmt.Errorf("non correct attestation found no slashable"))
		}
	})

	t.Run("Existing attestation, should not error",func(t *testing.T) {
		res, err := protector.IsSlashableAttestation(account, &pb.SignBeaconAttestationRequest{
			Id:                   nil,
			Domain:               []byte("domain"),
			Data:                 &pb.AttestationData{
				Slot:                 30,
				CommitteeIndex:       5,
				BeaconBlockRoot:      []byte("B"),
				Source:               &pb.Checkpoint{
					Epoch:                3,
					Root:                 []byte("C"),
				},
				Target:               &pb.Checkpoint{
					Epoch:                4,
					Root:                 []byte("D"),
				},
			},
		})

		if err != nil {
			t.Error(err)
		}
		if res != false {
			t.Error(fmt.Errorf("correct attestation found slashable"))
		}
	})

	t.Run("new attestation, should not error",func(t *testing.T) {
		res, err := protector.IsSlashableAttestation(account, &pb.SignBeaconAttestationRequest{
			Id:                   nil,
			Domain:               []byte("domain"),
			Data:                 &pb.AttestationData{
				Slot:                 30,
				CommitteeIndex:       5,
				BeaconBlockRoot:      []byte("E"),
				Source:               &pb.Checkpoint{
					Epoch:                6,
					Root:                 []byte("I"),
				},
				Target:               &pb.Checkpoint{
					Epoch:                7,
					Root:                 []byte("H"),
				},
			},
		})

		if err != nil {
			t.Error(err)
		}
		if res != false {
			t.Error(fmt.Errorf("correct attestation found slashable"))
		}
	})
}
