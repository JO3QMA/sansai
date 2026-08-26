package mercari

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/jo3qma/sansai/internal/model"
)

func TestSearchIntegration(t *testing.T) {
	if os.Getenv("SANSAI_INTEGRATION") == "" {
		t.Skip("set SANSAI_INTEGRATION=1 to run")
	}
	c := &Client{}
	res, err := c.Search(context.Background(), "nintendo switch", model.SearchOptions{Limit: 2})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Items) == 0 {
		t.Fatal("expected items")
	}
	b, _ := json.MarshalIndent(res, "", "  ")
	t.Log(string(b))
}
