package clogger

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_fieldCollection_add(t *testing.T) {
	fc := newFieldCollection()
	fc.add("add", "value")
	require.Equal(t, fc.retrieve("add"), "value")
	require.Empty(t, fc.retrieve("some-unknown-key"))
}

func Test_fieldCollection_addField(t *testing.T) {
	fc := newFieldCollection()
	fc.addField(field{key: "addField", value: "value"})
	require.Equal(t, fc.retrieve("addField"), "value")
	require.Empty(t, fc.retrieve("some-unknown-key"))
}

func Test_fieldCollection_fields(t *testing.T) {
	fc := newFieldCollection()
	fs := []field{
		{
			key:   "key1",
			value: "value1",
		},
		{
			key:   "key2",
			value: "value2",
		},
	}

	for _, f := range fs {
		fc.add(f.key, f.value)
	}

	require.Empty(t, fc.retrieve("add0"))

	found := fc.fields()
	for _, f := range found {
		require.Equal(t, fc.retrieve(f.key), f.value)
	}
}

func Test_fieldCollection_len(t *testing.T) {
	fc := newFieldCollection()
	fc.add("k", "v")
	fc.add("key", "val")
	require.Equal(t, 2, fc.len())
}

func Test_fieldCollection_merge(t *testing.T) {
	fc := newFieldCollection()
	fc.add("k1", "v1")
	fc.add("k2", "v2")

	ff := newFieldCollection()
	ff.add("k2", "v3") // ! when merging, identical keys will be overwritten in original fc
	ff.add("k3", "v3")
	fc.merge(ff)

	require.Equal(t, fc.retrieve("k1"), "v1")
	require.Equal(t, fc.retrieve("k2"), "v3")
	require.Equal(t, fc.retrieve("k3"), "v3")
}

func Test_newFieldCollection(t *testing.T) {
	fc := newFieldCollection()
	require.NotNil(t, fc)
	require.NotNil(t, fc.m)
	require.NotNil(t, fc.mu)
}
