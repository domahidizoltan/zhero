function loadSchemaClasses(evt) {
  htmx.ajax("GET", "/class-hierarchy", "#classHierarchy");
  document
    .getElementById("classHierarchy")
    .addEventListener("htmx:afterSwap", function (evt) {
      let classHierarchy = JSON.parse(evt.target.textContent);
      fillSchemaClasses(classHierarchy);
    });
}

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

function initPropertyOrderWidget(evt) {
  const showHiddenToggle = document.getElementById(
    "show-hidden-properties-toggle",
  );
  const propertyItems = document.querySelectorAll(".property-item");
  const propertyOrderList = document.getElementById("property-order-list");
  const identifierSelect = document.getElementById("identifier");
  const secondaryIdentifierSelect = document.getElementById(
    "secondary-identifier",
  );

  const updateUI = () => {
    const showAll = showHiddenToggle.checked;
    const selectedIdentifier = identifierSelect.value;
    const selectedSecondaryIdentifier = secondaryIdentifierSelect.value;

    // 1. Clear dynamic lists
    propertyOrderList.innerHTML = "";
    while (identifierSelect.options.length > 1) {
      identifierSelect.remove(1);
    }

    while (secondaryIdentifierSelect.options.length > 1) {
      secondaryIdentifierSelect.remove(1);
    }

    // 2. Iterate over each property to  update its state
    propertyItems.forEach((item) => {
      const hideToggle = item.querySelector(".hide-toggle");
      const propertyName = item.dataset.propertyName;
      if (!hideToggle || !propertyName) return;

      // Update visibility based on its own toggle and the main toggle
      if (hideToggle.checked && !showAll) {
        item.classList.add("is-hidden");
      } else {
        item.classList.remove("is-hidden");
      }

      // Update order list and identifier  lists if property is NOT hidden by its own toggle
      if (!hideToggle.checked) {
        // Add to Property Order list
        const listItem = document.createElement("div");
        listItem.classList.add(
          "bg-base-100",
          "border",
          "border-base-300",
          "p-2",
          "rounded-md",
          "flex",
          "items-center",
        );
        listItem.dataset.propertyName = propertyName;
        listItem.innerHTML = `<div class="handle w-full"><i class="fa-solid fa-sort w-5 h-5 mr-2"></i><span>${propertyName}</span></div>`;
        propertyOrderList.appendChild(listItem);

        // Add to Identifier dropdowns
        const option = document.createElement("option");
        option.value = propertyName;
        option.textContent = propertyName;
        identifierSelect.appendChild(option.cloneNode(true));
        secondaryIdentifierSelect.appendChild(option.cloneNode(true));
      }
    });

    // 3.  Restore selections in dropdowns
    identifierSelect.value = selectedIdentifier;
    secondaryIdentifierSelect.value = selectedSecondaryIdentifier;
  };

  // Attach  event listeners
  showHiddenToggle.addEventListener("change", updateUI);
  propertyItems.forEach((item) => {
    const hideToggle = item.querySelector(".hide-toggle");
    if (hideToggle) {
      hideToggle.addEventListener("change", updateUI);
    }
  });

  // Initial UI setup on page load
  updateUI();

  // Initialize SortableJS
  new Sortable(propertyOrderList, {
    animation: 150,
    ghostClass: "sortable-ghost",
    handle: ".handle",
  });
}
