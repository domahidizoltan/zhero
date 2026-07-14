# Public facing features 

The public facing page is for the page visitors who will see the final form of the webpage. This application runs in a different subprocess and reachable on http://localhost:8080.  
The public facing page has a template with a header, footer, menu and search section, and they will be applied to all subpages.  

## List classes

The defined classes are dynamically listed in the menu.  
Whenever they are created and have any enabled page, it's name appears in the menu in alphabetical order.  

See [Dynamic page rendering](#dynamic-page-rendering) for details how the page is rendered.

## List pages

When the class list is selected, a listing page is rendered, where the enabled pages are listed.   
This is a minimal list, which will load only the `listable` properties, and renders a listing tile for each page.  
This page also has a paging capability.  
Each page links to the detailed representation of itself.  

See [Dynamic page rendering](#dynamic-page-rendering) for details how the page list is rendered.

## Load page

When a page is loaded we check if the page is enabled and look up for the existing URL route (see [Routing](#routing)). If these checks are failing we return a 404 Not Found page (`template/page_not_found.hbs`).  
The page is rendered by reading the header meta values and the page data. The references are translated to links and the uploaded files are displayed as links or images.  

## Routing

This feature is responsible to assign a page with a custom URL.  

See [feature/routing.md](feature/routing.md) for more details about page routing feature.  

## Dynamic page rendering

At the moment 2 kinds of page renderer interfaces are defined in `controller/controller.go`:
- only render a page
- render a page and list pages of a schema class 

```go
type UserFacingPageRenderer interface {
	Render(schemaMeta schema.SchemaMeta, data map[string]any) (string, error)
}

type UserFacingPageListRenderer interface {
	UserFacingPageRenderer
	List(schemaMeta schema.SchemaMeta, data []map[string]any, paging paging.Meta) (string, error)
}
```

These functions define how the page data should be rendered to Handlebars template.  
Page preview and simple page load now is using `DynamicPageRenderer` from `controller/pagerenderer/dynamicpagerenderer.go`. This renders page data by using a predefined and simple component decision logic and layout. The layout is not explicitly defined, but dynamically rendered by using hardcoded HTML components.
