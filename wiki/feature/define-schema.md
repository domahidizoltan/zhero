# Define schema

Once the Schema.org schema is selected, the next step is to define the class.
In this step we can define:
- what properties we need from the class
- the property type: in Schema.org a property can be represented by different types, what can be a simple or complex data type
- the property backing HTML component: how the property should be rendered
- if the property is searchable: it returns results during a search operation
- if the property is listable: it should appear on the listing page
- if the property is mandatory for a page
- if the property is hidden on page creation
- the property position (order)
- the Identifier field: a generated unique ID of the page
- the Secondary Identifier: a human readable unique name, like the title of the page what would be used in URL slugs and search references
The schema name can't be updated. It is taken from the Schema.org database, and it is only used as a reference.

## Validation rules

Identifier and Secondary Identifier must be selected.

## Data mapping

The **Schema** is organized in 2 database tables:
- `schema_meta` is responsible for the minimal schema information:
  - `id`: database identifier of the record
  - `name`: the Schema.org name of the schema
  - `identifier`: the schema property name used as Identifier
  - `secondary_identifier`: the schema property name used as Secondary Identifier

- `schema_meta_properties` stores the properties of each defined class:
  - `schema_name`: the Schema.org name of the schema
  - `name`: the Schema.org name of the property
  - `mandatory`: shows if the property is mandatory
  - `searchable`: shows if the property value is available for search
  - `type`: the selected Schema.org type of the property
  - `component`: the selected HTML Component of the property
  - `order`: the position of the property
  - `listable`: shows if the property appears on the listing pages

The **Schema Properties** have the following attributes:
- **Name**: the Schema.org name of the property  
  type | location | mapping | note
  -|-|-|-
  HTML | `template/admin/schemaorg/edit-property.partial.hbs`| `property-name` name 
  DTO | `adminschema.schemaPropDto` | `Name` 
  domain | `schema.Property` | `Name` | 
  SQL | `schema_meta_properties` | `name`

- **Hidden**: shows is the property is hidden on the Admin edit page, and will not be used in pages  
  type | location | mapping | note
  -|-|-|-
  HTML | `template/admin/schemaorg/edit-property.partial.hbs`| | checkbox is checked when `NotUsed` is true for `adminschema.schemaPropDto` 
  DTO | `adminschema.schemaPropDto` | `NotUsed`
  domain | | | loaded only when it is not marked as hidden
  SQL | | | saved in `schema_meta_properties` only when it is not marked as hidden

- **Mandatory**: shows if the property is mandatory for the page  
  type | location | mapping | note
  -|-|-|-
  HTML | `template/admin/schemaorg/edit-property.partial.hbs`| `property-mandatory` name
  DTO | `adminschema.schemaPropDto` | `Mandatory`
  domain | `schema.Property` | `Mandatory` | 
  SQL | `schema_meta_properties` | `mandatory`

- **Searchable**: shows if the property value is available for search  
  type | location | mapping | note
  -|-|-|-
  HTML | `template/admin/schemaorg/edit-property.partial.hbs`| `property-searchable` name 
  DTO | `adminschema.schemaPropDto` | `Searchable`
  domain | `schema.Property` | `Searchable` | 
  SQL | `schema_meta_properties` | `searchable`

- **Listable**: shows if the property appears on the listing pages  
  type | location | mapping | note
  -|-|-|-
  HTML | `template/admin/schemaorg/edit-property.partial.hbs`| `property-listable` name 
  DTO | `adminschema.schemaPropDto` | `Listable`
  domain | `schema.Property` | `Listable` | 
  SQL | `schema_meta_properties` | `listable`

- **Type**: the Schema.org type of the property  
  type | location | mapping | note
  -|-|-|-
  HTML | `template/admin/schemaorg/edit-property.partial.hbs`| `property-type` name | a select list of possible types from `PossibleTypes` of `adminschema.schemaPropDto`
  DTO | `adminschema.schemaPropDto` | `SelectedType`
  domain | `schema.Property` | `Type` | 
  SQL | `schema_meta_properties` | `type`

- **HTML Component**: how to render the property as HTML  
  type | location | mapping | note
  -|-|-|-
  HTML | `template/admin/schemaorg/edit-property.partial.hbs`| `property-component` name | a select list of possible components
  DTO | `adminschema.schemaPropDto` | `SelectedComponent`
  domain | `schema.Property` | `Component` | 
  SQL | `schema_meta_properties` | `component`

- **Order**: the position of the property  
  type | location | mapping | note
  -|-|-|-
  HTML | `template/admin/schemaorg/edit.hbs` | `property-order` name | this is an ordered comma separated list of property names
  DTO | `adminschema.schemaPropDto`| `Order` | 
  domain | `schema.Property` | `Order` | 
  SQL | `schema_meta_properties` | `order` | 

- **Identifier**: the property used as Identifier  
  type | location | mapping | note
  -|-|-|-
  HTML | `template/admin/schemaorg/edit.hbs`| `identifier` name
  DTO | `adminschema.schemaDto` | `Identifier`
  domain | `schema.SchemaMeta` | `Identifier` | 
  SQL | `schema_meta` | `secondary_identifier`

- **Secondary Identifier**: the property used as Secondary Identifier  
  type | location | mapping | note
  -|-|-|-
  HTML | `template/admin/schemaorg/edit.hbs` | `secondary-identifier` name
  DTO | `adminschema.schemaDto` | `SecondaryIdentifier`
  domain | `schema.SchemaMeta` | `SecondaryIdentifier` | 
  SQL | `schema_meta` | `identifier`

