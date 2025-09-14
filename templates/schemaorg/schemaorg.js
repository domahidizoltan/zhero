document.addEventListener("DOMContentLoaded", function (evt) {
  htmx.ajax("GET", "/class-hierarchy", "#classHierarchy");
  document
    .getElementById("classHierarchy")
    .addEventListener("htmx:afterSwap", function (evt) {
      let classHierarchy = JSON.parse(evt.target.textContent);
      fillSchemaClasses(classHierarchy);
    });
});

function fillSchemaClasses(classHierarchy) {
  let opts = Array();
  classHierarchy.forEach((c) => {
    let cls = c[c.length - 1];
    let breadcrumb = "";
    if (c.length > 1) {
      breadcrumb = c.slice(0, c.length - 1).join(" > ") + " > ";
    }
    opts.push({ class: cls, breadcrumb: breadcrumb + cls });
  });

  let id = "schema-class";
  htmx.removeClass(htmx.find(`#${id}`), "hidden");

  new TomSelect(`#${id}`, {
    valueField: "class",
    searchField: "breadcrumb",
    options: opts,
    maxOptions: opts.length,
    maxItesm: 1,
    render: {
      option: function (data, escape) {
        return `<div>
          <div class="text-base">${data.class}</div>
          <div class="test-xs text-gray-500">${escape(data.breadcrumb)}</div>
          </div>`;
      },
      item: function (data, escape) {
        return `<div title="${escape(data.breadcrumb)}">${data.class}</div>`;
      },
    },
  });

  htmx.remove(htmx.find(`#${id}-loading`));
}

function navigateToEditSchemaPage() {
  let schemaClass = document.getElementById("schema-class").value;
  if (!schemaClass) {
    popup("No schema selected!");
  } else {
    let slug = schemaClass
      .replace(/([A-Z])/g, "-$1")
      .toLowerCase()
      .substring(1);
    window.location.href = `/schema/${slug}/edit`;
  }
}
