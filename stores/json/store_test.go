package json

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mozilla.org/sops/v3"
)

func TestDecodeJSON(t *testing.T) {
	in := `
{
   "glossary":{
      "title":"example glossary",
      "GlossDiv":{
         "title":"S",
         "GlossList":{
            "GlossEntry":{
               "ID":"SGML",
               "SortAs":"SGML",
               "GlossTerm":"Standard Generalized Markup Language",
               "Acronym":"SGML",
               "Abbrev":"ISO 8879:1986",
               "GlossDef":{
                  "para":"A meta-markup language, used to create markup languages such as DocBook.",
                  "GlossSeeAlso":[
                     "GML",
                     "XML"
                  ]
               },
               "GlossSee":"markup"
            }
         }
      }
   }
}`
	expected := sops.TreeBranch{
		sops.TreeItem{
			Key: "glossary",
			Value: sops.TreeBranch{
				sops.TreeItem{
					Key:   "title",
					Value: "example glossary",
				},
				sops.TreeItem{
					Key: "GlossDiv",
					Value: sops.TreeBranch{
						sops.TreeItem{
							Key:   "title",
							Value: "S",
						},
						sops.TreeItem{
							Key: "GlossList",
							Value: sops.TreeBranch{
								sops.TreeItem{
									Key: "GlossEntry",
									Value: sops.TreeBranch{
										sops.TreeItem{
											Key:   "ID",
											Value: "SGML",
										},
										sops.TreeItem{
											Key:   "SortAs",
											Value: "SGML",
										},
										sops.TreeItem{
											Key:   "GlossTerm",
											Value: "Standard Generalized Markup Language",
										},
										sops.TreeItem{
											Key:   "Acronym",
											Value: "SGML",
										},
										sops.TreeItem{
											Key:   "Abbrev",
											Value: "ISO 8879:1986",
										},
										sops.TreeItem{
											Key: "GlossDef",
											Value: sops.TreeBranch{
												sops.TreeItem{
													Key:   "para",
													Value: "A meta-markup language, used to create markup languages such as DocBook.",
												},
												sops.TreeItem{
													Key: "GlossSeeAlso",
													Value: []interface{}{
														"GML",
														"XML",
													},
												},
											},
										},
										sops.TreeItem{
											Key:   "GlossSee",
											Value: "markup",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	branch, err := Store{}.treeBranchFromJSON([]byte(in))
	assert.Nil(t, err)
	assert.Equal(t, expected, branch)
}

func TestDecodeSimpleJSONObject(t *testing.T) {
	in := `{"foo": "bar", "baz": 2}`
	expected := sops.TreeBranch{
		sops.TreeItem{
			Key:   "foo",
			Value: "bar",
		},
		sops.TreeItem{
			Key:   "baz",
			Value: 2.0,
		},
	}
	branch, err := Store{}.treeBranchFromJSON([]byte(in))
	assert.Nil(t, err)
	assert.Equal(t, expected, branch)
}

func TestDecodeNumber(t *testing.T) {
	in := `42`
	_, err := Store{}.treeBranchFromJSON([]byte(in))
	assert.NotNil(t, err)
}

func TestDecodeNestedJSONObject(t *testing.T) {
	in := `{"foo": {"foo": "bar"}}`
	expected := sops.TreeBranch{
		sops.TreeItem{
			Key: "foo",
			Value: sops.TreeBranch{
				sops.TreeItem{
					Key:   "foo",
					Value: "bar",
				},
			},
		},
	}
	branch, err := Store{}.treeBranchFromJSON([]byte(in))
	assert.Nil(t, err)
	assert.Equal(t, expected, branch)
}

func TestDecodeJSONWithArray(t *testing.T) {
	in := `{"foo": {"foo": [1, 2, 3]}, "bar": "baz"}`
	expected := sops.TreeBranch{
		sops.TreeItem{
			Key: "foo",
			Value: sops.TreeBranch{
				sops.TreeItem{
					Key:   "foo",
					Value: []interface{}{1.0, 2.0, 3.0},
				},
			},
		},
		sops.TreeItem{
			Key:   "bar",
			Value: "baz",
		},
	}
	branch, err := Store{}.treeBranchFromJSON([]byte(in))
	assert.Nil(t, err)
	assert.Equal(t, expected, branch)
}

func TestDecodeJSONArrayOfObjects(t *testing.T) {
	in := `{"foo": [{"bar": "foo"}, {"foo": "bar"}]}`
	expected := sops.TreeBranch{
		sops.TreeItem{
			Key: "foo",
			Value: []interface{}{
				sops.TreeBranch{
					sops.TreeItem{
						Key:   "bar",
						Value: "foo",
					},
				},
				sops.TreeBranch{
					sops.TreeItem{
						Key:   "foo",
						Value: "bar",
					},
				},
			},
		},
	}
	branch, err := Store{}.treeBranchFromJSON([]byte(in))
	assert.Nil(t, err)
	assert.Equal(t, expected, branch)
}

func TestDecodeJSONArrayOfArrays(t *testing.T) {
	in := `{"foo": [[["foo", {"bar": "foo"}]]]}`
	expected := sops.TreeBranch{
		sops.TreeItem{
			Key: "foo",
			Value: []interface{}{
				[]interface{}{
					[]interface{}{
						"foo",
						sops.TreeBranch{
							sops.TreeItem{
								Key:   "bar",
								Value: "foo",
							},
						},
					},
				},
			},
		},
	}
	branch, err := Store{}.treeBranchFromJSON([]byte(in))
	assert.Nil(t, err)
	assert.Equal(t, expected, branch)
}

func TestEncodeSimpleJSON(t *testing.T) {
	branch := sops.TreeBranch{
		sops.TreeItem{
			Key:   "foo",
			Value: "bar",
		},
		sops.TreeItem{
			Key:   "foo",
			Value: 3.0,
		},
		sops.TreeItem{
			Key:   "bar",
			Value: false,
		},
	}
	out, err := Store{}.jsonFromTreeBranch(branch)
	assert.Nil(t, err)
	expected, _ := Store{}.treeBranchFromJSON(out)
	assert.Equal(t, expected, branch)
}

func TestEncodeJSONWithEscaping(t *testing.T) {
	branch := sops.TreeBranch{
		sops.TreeItem{
			Key:   "foo\\bar",
			Value: "value",
		},
		sops.TreeItem{
			Key:   "a_key_with\"quotes\"",
			Value: 4.0,
		},
		sops.TreeItem{
			Key:   "baz\\\\foo",
			Value: 2.0,
		},
	}
	out, err := Store{}.jsonFromTreeBranch(branch)
	assert.Nil(t, err)
	expected, _ := Store{}.treeBranchFromJSON(out)
	assert.Equal(t, expected, branch)
}

func TestEncodeJSONArrayOfObjects(t *testing.T) {
	tree := sops.Tree{
		Branches: sops.TreeBranches{
			sops.TreeBranch{
				sops.TreeItem{
					Key: "foo",
					Value: []interface{}{
						sops.TreeBranch{
							sops.TreeItem{
								Key:   "foo",
								Value: 3,
							},
							sops.TreeItem{
								Key:   "bar",
								Value: false,
							},
						},
						2,
					},
				},
			},
		},
	}
	expected := `{
	"foo": [
		{
			"foo": 3,
			"bar": false
		},
		2
	]
}`
	store := Store{}
	out, err := store.EmitPlainFile(tree.Branches)
	assert.Nil(t, err)
	assert.Equal(t, expected, string(out))
}

func TestUnmarshalMetadataFromNonSOPSFile(t *testing.T) {
	data := []byte(`{"hello": 2}`)
	store := Store{}
	_, err := store.LoadEncryptedFile(data)
	assert.Equal(t, sops.MetadataNotFound, err)
}

func TestLoadJSONFormattedBinaryFile(t *testing.T) {
	// This is JSON data, but we want SOPS to interpret it as binary,
	// e.g. because the --input-type binary flag was provided.
	data := []byte(`{"hello": 2}`)
	store := BinaryStore{}
	branches, err := store.LoadPlainFile(data)
	assert.Nil(t, err)
	assert.Equal(t, "data", branches[0][0].Key)
}

func TestEmitValueString(t *testing.T) {
	bytes, err := (&Store{}).EmitValue("hello")
	assert.Nil(t, err)
	assert.Equal(t, []byte("\"hello\""), bytes)
}
