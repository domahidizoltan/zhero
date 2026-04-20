function loadSchemaClasses(evt) {
  htmx.ajax("GET", "/admin/schema/class-hierarchy", "#class-hierarchy");
  document
    .getElementById("class-hierarchy")
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
          <div class="text-base text-base-content">${data.class}</div>
          <div class="test-xs text-base-content/60">${escape(data.breadcrumb)}</div>
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
    window.location.href = `/admin/schema/edit/${schemaClass}`;
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
        item.classList.add("hide");
      } else {
        item.classList.remove("hide");
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
    form.noValidate = true;
    form.addEventListener("submit", (e) => {
      e.preventDefault();
      if (form.checkValidity()) {
        document.getElementById("property-order").value = sortable.toArray();
        document.getElementsByName("identifiers-fieldset")[0].disabled = false;
        const idName = document.getElementById("loaded-identifier").value;
        if (idName != "") {
          document.getElementsByName(idName + "-fieldset")[0].disabled = false;
        }
        const secIdName = document.getElementById(
          "loaded-secondary-identifier",
        ).value;
        if (secIdName != "") {
          document.getElementsByName(secIdName + "-fieldset")[0].disabled =
            false;
        }
        form.submit();
      }
    });
  }
}

function countListableProperties(secondaryIdentifierName) {
  return function (e) {
    if (e.target.matches('input[name="property-listable"]')) {
      let max = 3;
      max += document.querySelector(
        `input[name="property-listable"][value="${secondaryIdentifierName}"]`,
      )
        ? 1
        : 0;
      max += document.querySelector(
        `input[name="property-listable"][value="thumbnail"]`,
      )
        ? 1
        : document.querySelector(
              `input[name="property-listable"][value="image"]`,
            )
          ? 1
          : 0;

      const checkedCount = document.querySelectorAll(
        'input[name="property-listable"]:checked',
      ).length;

      if (checkedCount > max) {
        alert(
          `Maximum ${max} additional listable properties allowed (plus SecondaryIdentifier)`,
        );
        e.target.checked = false;
      }
    }
  };
}

function getTypeComponents(targetID, selectedComponent, typ) {
  let target = document.getElementById("prop-component-" + targetID);
  if (!target) {
    return;
  }

  let components = [];
  const t = typ.toLowerCase();
  if (t.includes("image") || t.includes("thumbnail")) {
    components = ["URL", "File"];
  } else if (t.includes("phone") || t.includes("fax")) {
    components = ["Tel", "TextInput"];
  } else if (t.includes("email")) {
    components = ["Email", "TextInput"];
  } else if (t.includes("color")) {
    components = ["Color", "TextInput"];
  } else if (t.includes("url") || t.includes("web")) {
    components = ["URL", "TextInput"];
  } else {
    switch (typ) {
      case "Boolean":
        components = ["Checkbox"];
        break;
      case ("Date", "DateTime", "Number", "Quantity", "Time"):
        components = [typ];
        break;
      case "Text":
        components = [
          "TextInput",
          "TextArea",
          "Color",
          "Email",
          "Tel",
          "URL",
          "Select",
          "File",
        ];
        break;
      default:
        components = ["TextInput", "URL", "ReferenceSearch"];
        break;
    }
  }

  while (target.firstChild) {
    target.removeChild(target.lastChild);
  }
  components.forEach((item) => {
    target.appendChild(new Option(item, item, true, selectedComponent == item));
  });
}

document.addEventListener("DOMContentLoaded", function () {
  document.querySelectorAll("[data-selected-component]").forEach(function (el) {
    let vals = el.getAttribute("data-selected-component").split("=");
    let name = vals[0];
    let selectedComponent = vals[1];
    getTypeComponents(name, selectedComponent, el.value);

    el.addEventListener("change", function (e) {
      getTypeComponents(name, selectedComponent, el.value);
    });
  });
});
