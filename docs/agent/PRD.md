# 1. Project Overview
Zhero is a Schema.org-based, SEO-first, headless Content Management System (CMS). 
It is designed for non-technical users to create, manage, and deliver structured content effortlessly. 
The system leverages Schema.org's vocabulary to ensure all content is machine-readable and optimized for search engines, generating JSON-LD for rich results. 
The application is a single self-contained executable, making it easy to deploy on various platforms, including low-resource hardware (e.g. Raspberry Pi Zero, or Android phones).

# 2. Target Users
- **Primary Users:** Non-technical content creators, small business owners, and marketers who need to publish structured data for SEO purposes but lack the technical skills for database management or programming.
- **User Needs:**
    - A simple, web-based admin interface to define content models (schemas).
    - An intuitive way to perform CRUD operations on content entries (pages).
    - Automatically generated, SEO-friendly structured data output.

# 3. Key Features

## 3.1. Schema Management
- **Schema Discovery:** Users can search the entire Schema.org class hierarchy to find a base for their content models.
- **Schema Definition:**
    - Create new schemas based on a selected Schema.org class.
    - Inherit all properties from the parent Schema.org class and its ancestors.
    - Customize properties:
        - Show/hide properties from the content entry form.
        - Set properties as mandatory.
        - Mark properties as searchable.
        - Define the data type and HTML component for data entry.
        - Re-order properties for the entry form.
    - **Identifiers:** Define a unique primary and a secondary identifier for content entries from the available properties.

## 3.2. Content Management (Pages)
- **Content Listing:** View all content entries for a specific schema, with searching and sorting capabilities.
- **Content Creation & Editing:** A dynamically generated form based on the schema definition for creating and updating content entries.
- **Data Persistence:** Content is stored in a local SQLite database.

# 4. Architecture & Design
- **Layered Architecture:** The backend follows a clean architecture pattern with distinct layers for controllers, services (domain logic), and repositories (data access).
- **Dependency Injection:** Dependencies are explicitly passed via constructors to promote modularity and testability.
- **Single Executable:** The final application is a single binary with embedded assets, requiring no external dependencies for deployment.
- **Web-based interface:** Administrators and users can manage data over a web-based graphical interface.

# 5. Success Metrics
- **Ease of Use:** Time taken for a non-technical user to create a new schema and publish their first content page.
- **SEO Performance:** Validity of the generated structured data as measured by tools like Google's Rich Results Test.
- **Performance:** Low resource consumption (CPU, Memory) of the running application.

