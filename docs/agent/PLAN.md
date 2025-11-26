# Development Plan for the upcoming tasks

This document outlines the remaining tasks to complete.

## F001: Admin Interface

1. **Admin Route:**
 - Establish a dedicated, top-level `/admin` route for all administrative functionalities to improve URL structure and organization. Move the current endpoints under `/admin`.

## F002: Additional Actions & Refinements

1. **Preview Functionality:**
 - Create a new non-Admin route and controller to render a preview of a page. 
 - The page data is loaded and transformed to JSON-LDThis
 - The JSON-LD data is passed to the new endpoint with POST method.
 - Create a separate package what transforms the JSON-LD data to dynamic HTML output.

## F003: Deployment

1. **Raspberry PI Zero packaging:**
 - Create Makefile task to cross-compile application for Raspberry PI Zero.
 - Ensure all necessary dependencies are included.
 - Define the build and run commands for the application.

2. **Android Packaging:**
 - Package the application into a Golang library.
 - Create an Android application what imports this Golang library.
 - Develop a build process to package the application for Android platforms.
 - This process will embed admin templates and static assets into the final binary, creating a self-contained and easily deployable application.
