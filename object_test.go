package notionapi_test

import (
	"testing"

	"github.com/jomei/notionapi"
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
