package es

const MappingTpl = `
{
	"settings":{
		"number_of_shards":1,
		"number_of_replicas":0
	},
	"mappings":{
		"properties":{
			"name":{
				"type":"keyword",
				"analyzer": "ik_max_word"
			},
			"keywords":{
				"type":"text",
				"analyzer": "ik_max_word"
			}
		}
	}
}
`
