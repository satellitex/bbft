package convertor

import (
	"github.com/satellitex/bbft/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBlock_GetHash(t *testing.T) {
	blocks := make([]model.Block, 20)
	for id, _ := range blocks {
		blocks[id] = randomBlock(t)
	}
	for id, a := range blocks {
		for jd, b := range blocks {
			if id != jd {
				assert.NotEqual(t, getHash(t, a), getHash(t, b))
			} else {
				assert.Equal(t, getHash(t, a), getHash(t, b))
			}
		}
	}
}

func TestBlock_Verify(t *testing.T) {

}
