package stablediffusion

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStableDiffusion(t *testing.T) {

	sd := NewStableDiffusion("127.0.0.1:19000")

	answer, err := sd.Do(context.Background(), "a girl")
	assert.NoError(t, err)

	assert.NotEmpty(t, answer)

	data, err := base64.StdEncoding.DecodeString(answer)
	assert.NoError(t, err)

	assert.NoError(t,
		ioutil.WriteFile("1.jpg", data, 0644))
}
