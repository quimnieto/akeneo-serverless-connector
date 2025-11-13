package akeneo

// EventTypes contains all available Akeneo event types
var EventTypes = []string{
	// Product events
	"product.created",
	"product.updated",
	"product.removed",

	// Product model events
	"product_model.created",
	"product_model.updated",
	"product_model.removed",

	// Category events
	"category.created",
	"category.updated",
	"category.removed",

	// Attribute events
	"attribute.created",
	"attribute.updated",
	"attribute.removed",

	// Attribute option events
	"attribute_option.created",
	"attribute_option.updated",
	"attribute_option.removed",

	// Attribute group events
	"attribute_group.created",
	"attribute_group.updated",
	"attribute_group.removed",

	// Family events
	"family.created",
	"family.updated",
	"family.removed",

	// Family variant events
	"family_variant.created",
	"family_variant.updated",
	"family_variant.removed",

	// Reference entity record events
	"reference_entity_record.created",
	"reference_entity_record.updated",
	"reference_entity_record.removed",

	// Delta events (only changes)
	"com.akeneo.pim.v1.product.updated.delta",
	"com.akeneo.pim.v1.product-model.updated.delta",
}

// GetEventTypes returns the list of available event types
func GetEventTypes() []string {
	return EventTypes
}
