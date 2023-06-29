package encrypt

import (
	"testing"

	"github.com/HardDie/mmr_boost_server/internal/config"
)

func TestEncrypt_Good(t *testing.T) {
	key := config.Encrypt{Key: "somelongkeymustbe32bytesizelengt"}
	e, err := NewEncrypt(key)
	if err != nil {
		t.Fatal(err)
	}

	arg := "some string with text"
	want := arg

	encRes, err := e.Encrypt(arg)
	if err != nil {
		t.Fatal(err)
	}

	got, err := e.Decrypt(encRes)
	if err != nil {
		t.Fatal(err)
	}

	if got != want {
		t.Fatalf("got: %q, want: %q", got, want)
	}
}

func TestEncrypt_Different(t *testing.T) {
	key1 := config.Encrypt{Key: "somelongkeymustbe32bytesizeleng1"}
	key2 := config.Encrypt{Key: "somelongkeymustbe32bytesizeleng2"}

	e1, err := NewEncrypt(key1)
	if err != nil {
		t.Fatal(err)
	}
	e2, err := NewEncrypt(key2)
	if err != nil {
		t.Fatal(err)
	}

	arg := "some string with text"

	encRes, err := e1.Encrypt(arg)
	if err != nil {
		t.Fatal(err)
	}

	_, err = e2.Decrypt(encRes)
	if err == nil {
		t.Fatal("you must be not have ability decrypt with other key")
	}
}
