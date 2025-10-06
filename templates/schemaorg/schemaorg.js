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
    window.location.href = `/schema/${schemaClass}/edit`;
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
    propertyOrderList.innerHTML = "";
    const showAll = showHiddenToggle.checked;
    const selectedIdentifier =
      document.getElementById("loaded-identifier").value;
    const selectedSecondaryIdentifier = document.getElementById(
      "loaded-secondary-identifier",
    ).value;

    // 1. Clear dynamic lists
    const propertyOrders = [];
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
      const propertyOrder = item.dataset.propertyOrder;
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
        listItem.dataset.propertyOrder = propertyOrder;
        listItem.innerHTML = `<div class="handle w-full"><i class="fa-solid fa-sort w-5 h-5 mr-2"></i><span>${propertyName}</span></div>`;
        propertyOrders.push(listItem);

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

    propertyOrders.sort((a, b) => {
      return a.dataset.propertyOrder - b.dataset.propertyOrder;
    });

    propertyOrders.forEach((item) => {
      propertyOrderList.appendChild(item);
    });
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
  let sortable = new Sortable(propertyOrderList, {
    animation: 150,
    ghostClass: "sortable-ghost",
    handle: ".handle",
    dataIdAttr: "data-property-name",
  });

  const form = document.querySelector("#edit-schema-form");
  if (form) {
    form.addEventListener("submit", (e) => {
      document.getElementById("property-order").value = sortable.toArray();
    });
  }
}
