package utils

import "strings"

type Neo4JObjectType = string
type Neo4JObjectProperties = map[string]any

type Neo4JObject struct {
	Category   Neo4JObjectType
	Properties Neo4JObjectProperties
}

func (self *Neo4JObject) AppendAsNeo4JMatch(b *strings.Builder, limits []string, queryId string) {
	if len(limits) != 2 {
		panic("There should only be two limits! Example: []string {\"[\", \"]\"}")
	}

	b.WriteString(limits[0])
	b.WriteString(queryId)
	b.WriteRune(':')
	b.WriteString(self.Category)

	propertiesCount := len(self.Properties)
	if propertiesCount > 0 {
		b.WriteString(" {")

		i := 1
		for property := range self.Properties {
			b.WriteString(property)
			b.WriteString(": ")
			b.WriteRune('$')
			b.WriteString(queryId)
			b.WriteRune('_')
			b.WriteString(property)

			if i != propertiesCount {
				b.WriteRune(',')
			}
			i++
		}
		b.WriteString("}")
	}
	b.WriteString(limits[1])
}
