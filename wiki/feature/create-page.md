# Create page

Pages are listed in the Admin panel in a table with filtering and paging capabilities. In this `template/admin/page/list.hbs` template we can also start creating a new page.

When we create a page the Identifier is not editable. It is generated once the page is created, and it can't be updated.  
During page create or update we can:
- edit page property values
- edit SEO values
- define custom URL route
- upload files
- enable or disable a page
- preview the page (even before the first save or when it is disabled)


## Validation rules

All the page properties marked as mandatory must have a value on save.

## Data mapping

The **Page** has the following attributes during creation:
- **Enabled**: enables or disables a page  
  type | location | mapping | note
  -|-|-|-
  HTML | `template/admin/page/edit.hbs`| `is-enabled` name 
  DTO | `adminpage.pageDto` | `IsEnabled` 
  domain | `page.Page` | `IsEnabled` | 
  SQL | `page` | `enabled`

- **Identifier**: the unique identifier of the page which is automatically generated on create and not updatable  
  Note: Here we store the value of the identifier, not the property name  
  type | location | mapping | note
  -|-|-|-
  HTML | `template/admin/page/edit.hbs`| `identifier` name | hidden field used on form submit
  DTO | `adminpage.pageDto` | `Identifier`
  domain | `page.Page` | `Identifier`
  SQL | `page` | `identifier`

### Meta section

This section stores the SEO meta values and page route.  
Route mapping details are covered in `wiki/feature/routing.md`. 
- **Meta**:  
  type | location | mapping | note
  -|-|-|-
  HTML | `template/admin/page/edit.hbs`| names with `meta-` prefix  
  DTO | `adminpage.pageDto` | `Meta`
  domain | `page.Page` | `Meta` 
  SQL | `page` | `meta`

`Meta` fields are serialized to JSON object to database, to speed up loading with page data without using SQL joins.  
On the edit page the fields from the JSON object are prefixed with `meta-` and they are dynamically processed on submit.  

The `Meta` object has the following fields:  
name | HTML component name | `adminpage.pageMeta` DTO field | `page.PageMeta` domain field | JSON key | note 
-|-|-|-|-|-
Meta Title | `meta-title` | `Title` | `Title` | `title`
Meta Description | `meta-description` | `Description` | `Description` | `description`
robots directive | `meta-robots-noindex` and `meta-robots-nofollow` | `Robots` | `Robots` | `robots` | an array with checked values, e.g. `["noindex","nofollow"]`
Adult content | `meta-rating-adult` | `Rating` | `Rating` | `rating` | value `adult` indicates that the page content is age restricted
OG Title | `meta-og-title` | `OGTitle` | `OGTitle` | `ogTitle` | OpenGraph title when sharing page on social media
OG Description | `meta-og-description` | `OGDescription` | `OGDescription` | `ogDescription` | OpenGraph description when sharing page on social media
Field Meta | `alt-text-<field_name>` | `FieldMeta` | `FieldMeta` | `fieldMeta` | key-value map to store extra SEO values for some fields

Field meta has compound keys, starting with the given field name followed by the meta property name.
Field meta can have the following keys:
- `<field_name>:altText`: "alt" tag property of links and images 

### Data section

This section enumerates the fields from the schema properties defined as the blueprint.  
The field positions are indicated by the Order attribute.  
The fields have the following visual marks:
- red star: is mandatory
- list icon: is listable
- magnifying glass icon: is searchable  

- **Fields data**:  
The data fields are serialized to JSON object to database, and also used as map in the domain.  
On the edit page the fields from the JSON object are prefixed with `field-` and they are dynamically processed on submit.  
Field data is defined in `adminpage.fieldDto` and it contains all the details needed to render the field, including the Type and Component.  
  type | location | mapping | note
  -|-|-|-
  HTML | `template/admin/page/edit.hbs`| names with `field-` prefix  
  DTO | `adminpage.pageDto` | `Fields`
  domain | `page.Page` | `Data` 
  SQL | `page` | `data`


- **Secondary Identifier**:  
The Secondary Identifier is extracted from the fields.  
  type | location | mapping | note
  -|-|-|-
  HTML | `template/admin/page/edit.hbs`| `field-<name>` from `adminschema.schemaDto#SecondaryIdentifier`, e.g. `field-author`
  DTO | `adminpage.pageDto` | `SecondaryIdentifier` for field name and `SecondaryIdentifierValue` for field value
  domain | `page.Page` | `SecondaryIdentifier` | field value
  SQL | `page` | `secondary_identifier` | field value

- **Data**:  
Zhero prefers JSON-LD as the primary read-data representation of the pages, so it could be returned on API and SEO headers. HTML could also be rendered from this.
The `data` column of the `page` SQL table is mapped to a JSON-LD object, where the object is a key-value map. The key is a field name and the value is the field value. For the moment this map has 2 extra fields:
- `@id`: the Identifier (duplicated)
- `@type`: the schema name from Schema.org

- **Listable data**:  
Listable data is represented as JSON in the database level, where the keys are the listable field names.  
  type | location | mapping | note
  -|-|-|-
  HTML | `template/admin/page/edit.hbs`| | only marked with `fa-list` when the field is listable
  DTO | `adminpage.fieldDto` | `ListableData` | key-value map
  domain | `page.Page` | `ListableData` | key-value map
  SQL | `page` | `listable_data`


- **Searchable fields**:  
To improve full text search we allow only 5 search slots (defined in `model.MaxSearchVals`). 
These search slots are indexed with 0-4, and the `Searchable` field values are stored in these slots, and also their values in `data` and `listable_data` are removed. These values are replaced back to the model and dto when the page data is read from the database.  
  type | location | mapping | note
  -|-|-|-
  HTML | `template/admin/page/edit.hbs`| | only marked with `fa-magnify-icon` when the field is searchable
  DTO | `adminpage.fieldDto` | `IsSearchable`
  domain | `page.Page` | `SearchVals` | array of searchable field values
  SQL | `page_search` | `col0`..`col4` 
`secondary_identifier` column in DB table `page` is also considered searchable, but it's not added to the search slots.


- **References**:  
Page fields with component type `ReferenceSearch` are linkable references, what can be internal or external links.
These links are represented in the following format `#ZHERO#<schema_name>/<identifier>#<ref_meta>#`, where:
  - `<schema_name>`: the Schema.org name of the schema
  - `<identifier>`: the unique Identifier of the page
  - `<ref_meta>`: some meta data of the reference represented in key-value JSON object  
Possible values for `<ref_meta>`:
  - `linkText`: visible text of the anchor link
  - `altText`: SEO "alt" tag of the link or image
The references are extracted to a deduplicated list, and the values in page data are replaced with `ref<idx>`, where `<idx>` is the index counter for the occurrence of the reference.  
  type | location | mapping | note
  -|-|-|-
  HTML | `template/admin/page/edit.hbs`| | visible in HTML input and text areas
  DTO | `adminpage.fieldDto` | `References`
  domain | `page.Page` | `References`
  SQL | `page` | `references`

