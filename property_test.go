package notionapi

import (
	"encoding/json"
	"testing"
)

func TestProperties_UnmarshalJSON_Place(t *testing.T) {
	raw := []byte(`{
		"Place": {
			"id": "%60%40Gq",
			"type": "place",
			"place": null
		}
	}`)

	var props Properties
	if err := json.Unmarshal(raw, &props); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	p, ok := props["Place"].(*PlaceProperty)
	if !ok {
		t.Fatalf("got %T, want *PlaceProperty", props["Place"])
	}
	if p.GetID() != "%60%40Gq" {
		t.Errorf("GetID() = %q, want %q", p.GetID(), "%60%40Gq")
	}
	if p.GetType() != PropertyTypePlace {
		t.Errorf("GetType() = %q, want %q", p.GetType(), PropertyTypePlace)
	}
	if p.Place != nil {
		t.Errorf("Place = %#v, want nil", p.Place)
	}
}

func TestPropertyConfigs_UnmarshalJSON_Place(t *testing.T) {
	raw := []byte(`{
		"Place": {
			"id": "Xqz4",
			"name": "Place",
			"type": "place",
			"place": {}
		}
	}`)

	var configs PropertyConfigs
	if err := json.Unmarshal(raw, &configs); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	p, ok := configs["Place"].(*PlacePropertyConfig)
	if !ok {
		t.Fatalf("got %T, want *PlacePropertyConfig", configs["Place"])
	}
	if p.GetID() != "Xqz4" {
		t.Errorf("GetID() = %q, want %q", p.GetID(), "Xqz4")
	}
	if p.GetType() != PropertyConfigPlace {
		t.Errorf("GetType() = %q, want %q", p.GetType(), PropertyConfigPlace)
	}
}
