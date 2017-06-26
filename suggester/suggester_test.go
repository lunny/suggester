package suggester

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSearch(t *testing.T) {
	var plates = []string{
		"沪B20176",
		"皖F1A020",
		"鄂B01345",
		"沪F09124",
	}

	client := New("http://127.0.0.1:9191/")
	var unitID int64 = 1
	var prefix = "k"
	for i, p := range plates {
		err := client.AddIndex(prefix, p, int64(i), unitID)
		if err != nil {
			t.Fatal(err)
		}
	}

	var searchKases = map[string][]Result{
		"沪B": []Result{
			{0, "沪B20176"},
		},
		"皖": []Result{
			{1, "皖F1A020"},
		},
		"皖f": []Result{
			{1, "皖F1A020"},
		},
		"沪": []Result{
			{0, "沪B20176"},
			{3, "沪F09124"},
		},
		"沪b": []Result{
			{0, "沪B20176"},
		},
		"鄂": []Result{
			{2, "鄂B01345"},
		},
	}

	for k, res := range searchKases {
		start := time.Now()
		results, err := client.Search(prefix, k, unitID, 10)
		assert.NoError(t, err)
		assert.Equal(t, res, results)
		fmt.Println("search", k, " escaped", time.Now().Sub(start))
	}

	assert.NoError(t, client.DelIndex(prefix, "沪B20176", 0, unitID))
	results, err := client.Search(prefix, "沪B", unitID, 10)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(results))
}
