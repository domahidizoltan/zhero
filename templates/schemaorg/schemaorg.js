// Function to handle item selection and trigger HTMX load
function selectItem(itemId) {
  console.log("Selected item: " + itemId);
  // Trigger HTMX request to load details for the selected item
  htmx.ajax("GET", `/items/${itemId}`, "#detail-view");
}
// Add event listener for keyup on the search input to trigger page load on Enter key
document.getElementById("search").addEventListener("keyup", function (event) {
  if (event.key === "Enter") {
    const searchTerm = event.target.value;
    if (searchTerm) {
      // Trigger HTMX request to load search results
      htmx.ajax("GET", `/search?q=${searchTerm}`, "#results");
    }
  }
});
