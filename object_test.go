package notionapi_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/conduitio-labs/notionapi"
)

func TestDate(t *testing.T) {
	t.Run(".UnmarshalText", func(t *testing.T) {
		var d notionapi.Date

		t.Run("OK datetime with timezone", func(t *testing.T) {
			data := []byte("1987-02-13T00:00:00.000+01:00")
			err := d.UnmarshalText(data)
			if err != nil {
				t.Fatal(err)
			}
		})
		t.Run("OK date", func(t *testing.T) {
			data := []byte("1985-01-02")
			err := d.UnmarshalText(data)
			if err != nil {
				t.Fatal(err)
			}
		})
		t.Run("NOK", func(t *testing.T) {
			data := []byte("1985")
			err := d.UnmarshalText(data)
			if err == nil {
				t.Fatalf("expected an error, got none")
			}
		})
	})
}

func TestColor_MarshalText(t *testing.T) {
	type Foo struct {
		Test notionapi.Color `json:"test"`
	}

	t.Run("marshall to color if color is not empty", func(t *testing.T) {
		f := Foo{Test: notionapi.ColorGreen}
		r, err := json.Marshal(f)
		if err != nil {
			t.Fatal(err)
		}
		want := []byte(`{"test":"green"}`)
		if !reflect.DeepEqual(r, want) {
			t.Errorf("Color.MarshallText error() got = %v, want %v", r, want)
		}
	})

	t.Run("marshall to default color if color is empty", func(t *testing.T) {
		f := Foo{}
		r, err := json.Marshal(f)
		if err != nil {
			t.Fatal(err)
		}
		want := []byte(`{"test":"default"}`)
		if !reflect.DeepEqual(r, want) {
			t.Errorf("Color.MarshallText error() got = %v, want %v", r, want)
		}
	})
}
