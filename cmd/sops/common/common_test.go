package common

import (
	"github.com/getsops/sops/v3"
	"github.com/getsops/sops/v3/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDiffShamirThresholdWhenIntroducingGroupsAndThresholdBelowGroupSize(t *testing.T) {
	metadata := sops.Metadata{ShamirThreshold: 0}
	conf := MakeConfig(2, 3)
	diff := DiffShamirThreshold(metadata, &conf)

	assert.Equal(t, 0, diff.Old)
	assert.Equal(t, 2, diff.New)
}

func TestDiffShamirThresholdWhenIntroducingGroupsAndThresholdAboveGroupSize(t *testing.T) {
	metadata := sops.Metadata{ShamirThreshold: 0}
	conf := MakeConfig(4, 3)
	diff := DiffShamirThreshold(metadata, &conf)

	assert.Equal(t, 0, diff.Old)
	assert.Equal(t, 3, diff.New)
}

func TestDiffShamirThresholdWhenIntroducingGroupsAndNoThresholdIsConfigured(t *testing.T) {
	metadata := sops.Metadata{ShamirThreshold: 0}
	conf := MakeConfig(0, 3)
	diff := DiffShamirThreshold(metadata, &conf)

	assert.Equal(t, 0, diff.Old)
	assert.Equal(t, 3, diff.New)
}

func TestDiffShamirThresholdWhenReducingThreshold(t *testing.T) {
	metadata := sops.Metadata{ShamirThreshold: 3}
	conf := MakeConfig(2, 3)
	diff := DiffShamirThreshold(metadata, &conf)

	assert.Equal(t, 3, diff.Old)
	assert.Equal(t, 2, diff.New)
}

func TestDiffShamirThresholdWhenIncreasingThreshold(t *testing.T) {
	metadata := sops.Metadata{ShamirThreshold: 2}
	conf := MakeConfig(3, 4)
	diff := DiffShamirThreshold(metadata, &conf)

	assert.Equal(t, 2, diff.Old)
	assert.Equal(t, 3, diff.New)
}

func TestDiffShamirThresholdWhenRemovingThresholdConfiguration(t *testing.T) {
	metadata := sops.Metadata{ShamirThreshold: 2}
	conf := MakeConfig(0, 4)
	diff := DiffShamirThreshold(metadata, &conf)

	assert.Equal(t, 2, diff.Old)
	assert.Equal(t, 4, diff.New)
}

func TestDiffShamirThresholdWhenIntroducingSingleGroup(t *testing.T) {
	metadata := sops.Metadata{ShamirThreshold: 0}
	conf := MakeConfig(2, 1)
	diff := DiffShamirThreshold(metadata, &conf)

	assert.Equal(t, 0, diff.Old)
	assert.Equal(t, 0, diff.New)
}

func TestDiffShamirThresholdWhenReducingToSingleGroup(t *testing.T) {
	metadata := sops.Metadata{ShamirThreshold: 3}
	conf := MakeConfig(2, 1)
	diff := DiffShamirThreshold(metadata, &conf)

	assert.Equal(t, 3, diff.Old)
	assert.Equal(t, 0, diff.New)
}

func MakeConfig(shamirThreshold int, keygroups int) config.Config {
	return config.Config{
		ShamirThreshold: shamirThreshold,
		KeyGroups:       make([]sops.KeyGroup, keygroups),
	}
}
