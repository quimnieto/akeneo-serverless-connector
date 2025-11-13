package akeneo

// EventTypes contains all available Akeneo event types
var EventTypes = []string{
	// Product events
	"com.akeneo.pim.v1.product.created",
	"com.akeneo.pim.v1.product.updated",
	"com.akeneo.pim.v1.product.deleted",

	// Product model events
	"com.akeneo.pim.v1.product-model.created",
	"com.akeneo.pim.v1.product-model.updated",
	"com.akeneo.pim.v1.product-model.deleted",

	// Category events
	"com.akeneo.pim.v1.category.created",
	"com.akeneo.pim.v1.category.updated",
	"com.akeneo.pim.v1.category.deleted",

	// Attribute events
	"com.akeneo.pim.v1.attribute.created",
	"com.akeneo.pim.v1.attribute.updated",
	"com.akeneo.pim.v1.attribute.deleted",

	// Attribute option events
	"com.akeneo.pim.v1.attribute-option.created",
	"com.akeneo.pim.v1.attribute-option.updated",
	"com.akeneo.pim.v1.attribute-option.deleted",

	// Attribute group events
	"com.akeneo.pim.v1.attribute-group.created",
	"com.akeneo.pim.v1.attribute-group.updated",
	"com.akeneo.pim.v1.attribute-group.deleted",

	// Family events
	"com.akeneo.pim.v1.family.created",
	"com.akeneo.pim.v1.family.updated",
	"com.akeneo.pim.v1.family.deleted",

	// Family variant events
	"com.akeneo.pim.v1.family-variant.created",
	"com.akeneo.pim.v1.family-variant.updated",
	"com.akeneo.pim.v1.family-variant.deleted",

	// Reference entity record events
	"com.akeneo.pim.v1.reference-entity-record.created",
	"com.akeneo.pim.v1.reference-entity-record.updated",
	"com.akeneo.pim.v1.reference-entity-record.deleted",

	// Delta events (only changes)
	"com.akeneo.pim.v1.product.updated.delta",
	"com.akeneo.pim.v1.product-model.updated.delta",
}

// GetEventTypes returns the list of available event types
func GetEventTypes() []string {
	return EventTypes
}
