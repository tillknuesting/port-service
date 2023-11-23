package database_test

import (
	"context"
	"ports-service/internal/adapters/database"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	testCases := []struct {
		name    string
		memDB   database.MemDB[string]
		key     string
		value   string
		wantErr bool
	}{
		{
			name:    "NewValue",
			memDB:   database.MemDB[string]{DB: make(map[string]string)},
			key:     "newKey",
			value:   "newValue",
			wantErr: false,
		},
		{
			name:    "UpdateValue",
			memDB:   database.MemDB[string]{DB: make(map[string]string)},
			key:     "existingKey",
			value:   "updatedValue",
			wantErr: false,
		},
		{
			name:    "BooleanValue",
			memDB:   database.MemDB[string]{DB: make(map[string]string)},
			key:     "boolKey",
			value:   "true",
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.memDB.Set(context.Background(), tc.key, tc.value)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.value, tc.memDB.DB[tc.key])
			}
		})
	}
}
