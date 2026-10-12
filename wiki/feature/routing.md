# Routing

When a page is requested, first we check the URL. There is a routing table where the SEO compatible URL slugs are stored. When an URL is found here, we check if it is the latest version. When we have a newer version, the page is redirect to the latest version with a 301 status. This is important, so old visitor links can still work even if the Admins are changing the URLS.  
When the URL is not found in the routing table, we try to load the page by matching the schema name and the identifier from the URL.  
If there is no match, a 404 Not Found page is returned.  


## Validation rules

When not empty, it must be unique in the routing table.

## Data mapping

- **Route**: the current route assigned to a page    
  type | location | mapping | note
  -|-|-|-
  HTML | `template/admin/page/edit.hbs`| `route` name 
  DTO | `adminpage.pageDto` | `Route` 
  domain | `page.Page` | `Route` | 
  SQL | `route` | `route` | the latest route version

- **Routing Table**: this SQL table stores all the route versions in the `route` table
  - `route` column: the unique route across all pages
  - `page` column: composite page identifier in format `<schema_name>/<identifier>`
  - version: 1 based index of the route version of a given `page`

