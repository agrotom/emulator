package auth

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAuthEncoding(t *testing.T) {
	as := AuthService{
		TerminalIdentifier:      232,
		MobileNetworkExists:     false,
		NetworkIdentifierExists: false,
		SSRA:                    false,
		LanguageCodeExists:      false,
		InternationalMobileSubscriberIdentityExists: false,
		InternationalMobileEquipmentIdentityExists:  true,
		HomeDispatcherIdentifierExists:              true,
		HomeDispatcherIdentifier:                    34,
		InternationalMobileEquipmentIdentity:        "23213",
	}

	_, err := as.Encode()

	if err != nil {
		t.Errorf("%s", err.Error())
	}
}

func TestAuthDecoding(t *testing.T) {
	as := AuthService{
		TerminalIdentifier:      232,
		BufferSizeExists:        false,
		MobileNetworkExists:     false,
		NetworkIdentifierExists: false,
		SSRA:                    false,
		LanguageCodeExists:      false,
		InternationalMobileSubscriberIdentityExists: false,
		InternationalMobileEquipmentIdentityExists:  true,
		HomeDispatcherIdentifierExists:              true,
		HomeDispatcherIdentifier:                    34,
		InternationalMobileEquipmentIdentity:        "23213",
	}

	asBytes, err := as.Encode()

	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}

	decodedAs := AuthService{}
	err = decodedAs.Decode(asBytes)

	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}

	if diff := cmp.Diff(as, decodedAs); diff != "" {
		t.Fatal(diff)
	}

	if !reflect.DeepEqual(as, decodedAs) {
		t.Errorf("Encoded auth service and decoded one have not been coincided")
	}
}
