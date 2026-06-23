# Define schema

Once the Schema.org schema is selected, the next step is to define the class.
In this step we can define:
- what properties we need from the class
- the type for the property (in Schema.org a property can be represented by different types, what can be simple or complex data types)
- the backing HTML component for the property (how it should be rendered)
- the order of the properties
- if the properties should be searchable (returning results during a search operation)
- if the properties should appear on the listing page (listable)
- if they will be mandatory for a page
- the Identifier field what will be a generated unique ID of the page
- the Secondary Identifier what will be a human readable unique name, like the title of the page what would be used in URL slugs and search references

## Validation rules

Identifier and Secondary Identifier must be selected.

## Data mapping

schema name
can't be updated
Identifier: identifier
Secondary Identifier: secondary-identifier


The **Schema Properties** have the following attributes:
- Order: the position of the property
    - HTML: `property-order` name
    - DTO: `Order` from `adminschema.schemaPropDto` 
    - domain: `Order` from `schema.Property` 
    - SQL: `order` from `schema_meta_properties` 


order: property-order

property: property name
hide: isLoadedClass?
mandatory: property-mandatory
searchable: property-searchable
listable: property-listable

type: property-type
HTML component: property-component


`schema_meta` SQL table:
- id: table row identifier
- name: Schema.org name of the schema e.g. BlogPosting
- identifier: name of the property selected as identifier
- secondary_identifier: name of the property selected as secondary identifier

`schema.SchemaMeta` domain:
- Name: mapped from `schema_meta.name` DB column
- Identifier: mapped from `schema_meta.identifier` DB column
- SecondaryIdentifier: mapped from `schema_meta.secondary_identifer` DB column
- Properties: mapped from `schema_meta_properties` DB table

`adminschema.schemaDto` DTO:
- IsLoaded: TODO
- Name: mapped from `SchemaMeta.Name` domain field
- Description: description of the Schema.org field
- CanonicalURL: TODO
- Properties: mapped from `SchemaMeta.Properties` domain field
- Identifier: TODO
- SecondaryIdentifier: TODO

