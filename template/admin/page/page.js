function confirmListAction(identifier, action) {
  document.getElementById("list-action-name").value = action;
  document.getElementById("list-action-item-id").value = identifier;
  document.getElementById("list-action-desc").innerHTML =
    `Do you really want to <b>${action}</b> page <b>${identifier}</b>?`;
  document.getElementById("list-action-modal").showModal();
}

document.addEventListener("htmx:afterRequest", (e) => {
  const listActionModal = document.getElementById("list-action-modal");
  if (!listActionModal) {
    return;
  }

  listActionModal.close();
  const triggerHeader = e.detail.xhr.getResponseHeader("HX-Trigger");
  if (!triggerHeader) {
    return;
  }

  try {
    const triggerData = JSON.parse(triggerHeader);
    if (triggerData.showError) {
      popup(triggerData.showError);
    }
  } catch (err) {
    console.error(err);
  }
});

const previewHost =
  window.location.protocol + "//" + window.location.hostname + ":" + 8080;

function submitPreview(cls) {
  const form = document.getElementById("edit-page-form");
  const restoreAction = form.action;
  const restoreTarget = form.target;
  try {
    form.action = `${previewHost}/preview/${cls}`;
    form.target = "_blank";
    form.submit();
  } finally {
    form.action = restoreAction;
    form.target = restoreTarget;
  }
}

document.addEventListener("DOMContentLoaded", function () {
  document.querySelectorAll("[data-soft-limit]").forEach(function (el) {
    var limit = parseInt(el.getAttribute("data-soft-limit"));
    var check = function () {
      if (el.value.length > limit) {
        el.classList.add("bg-yellow-100");
      } else {
        el.classList.remove("bg-yellow-100");
      }
    };
    check();
    el.addEventListener("input", check);
  });
});

let refPattern = /#ZHERO#([^#]+)#\{([^}]*)\}#/gim;
let modalContainerFor = "";
let selectionPos = 0;
document.addEventListener("DOMContentLoaded", function () {
  let openRefModal = function (e) {
    if (e.key && e.key !== "#") {
      return;
    }

    modalContainerFor = e.target.name;
    if (selectionPos == 0) {
      selectionPos = e.target.selectionStart;
    }

    let fieldType = e.target.getAttribute("data-field-type");
    if (e.target.tagName === "TEXTAREA") {
      fieldType = "";
    }
    [typeID, linkText, altText] = ["", "", ""];
    if (e.target.value.length > 0) {
      while ((match = refPattern.exec(e.target.value))) {
        if (refPattern.lastIndex < e.target.selectionStart) {
          continue;
        }

        [fieldType, typeID] = match[1].split("/");
        const meta = JSON.parse(`{${match[2]}}`);
        linkText = meta.linkText || "";
        altText = meta.altText || "";
      }
    }

    htmx.ajax(
      "GET",
      `/admin/page/reference:search?type=${fieldType}&id=${typeID}&linkText=${linkText}&altText=${altText}`,
      {
        target: "#modal-container",
        push: false,
      },
    );
  };

  document.addEventListener("keyup", function (e) {
    if (e.key === "Escape" && modalContainerFor !== "") {
      closeSearchReferenceModal();
    }
  });

  document.querySelectorAll("[ref-hook=true]").forEach((item) => {
    item.addEventListener("dblclick", openRefModal);
    item.addEventListener("keyup", openRefModal);
  });
});

function initTypeSelect() {
  if (document.getElementById("select-ref")) {
    let select = new TomSelect("#select-ref", {
      create: false,
      sortField: { field: "text", direction: "asc" },
      render: {
        option: function (data, escape) {
          let cls = "";
          if (data.enabled == "false") {
            cls = 'class="text-neutral-content"';
          }
          return `<div ${cls}>${data.text}</div>`;
        },
        item: function (item, escape) {
          return `<div>${item.text}</div>`;
        },
      },
    });

    select.on("change", function (e) {
      const label = select.getItem(e).innerHTML;
      document.getElementsByName("link-text")[0].value = label;
      document.getElementsByName("alt-text")[0].value = label;
    });
  } else if (document.getElementById("select-ref-type")) {
    new TomSelect("#select-ref-type", {
      create: false,
      sortField: { field: "text", direction: "asc" },
    });
  }
}

function closeSearchReferenceModal() {
  let modal = document.getElementById("search-reference");
  if (modal) {
    modal.close();
    modal.classList.remove("modal-open");
  }
  modalContainerFor = "";
  selectionPos = 0;
}

function setReferenceType() {
  let refType = document.getElementById("select-ref-type");
  const fieldType = refType.value;

  htmx.ajax("GET", `/admin/page/reference:search?type=${fieldType}`, {
    target: "#modal-container",
    push: false,
  });
}

function setReference(type) {
  const form = document.getElementById("set-reference-form");
  const formData = new FormData(form);
  const formProps = Object.fromEntries(formData);

  let identifier;
  if (formProps["page"]) {
    identifier = formProps["page"];
  }

  let meta = {};
  if (formProps["link-text"]) {
    meta["linkText"] = formProps["link-text"];
  }
  if (formProps["alt-text"]) {
    meta["altText"] = formProps["alt-text"];
  }

  const val = `#ZHERO#${type}/${identifier}#${JSON.stringify(meta)}#`;
  let cont = document.getElementById(modalContainerFor);

  if (selectionPos > 0 && cont.value.indexOf("#", selectionPos - 1) > -1) {
    selectionPos--;
  }

  let replaced = false;
  if (cont.value.length == 0) {
    replaced = true;
    cont.value = val;
  } else {
    while ((match = refPattern.exec(cont.value))) {
      const pos = match.index;
      const len = match[0].length;
      if (pos + len < selectionPos) {
        continue;
      }

      cont.value = cont.value.slice(0, pos) + val + cont.value.slice(pos + len);
      replaced = true;
      break;
    }
  }

  if (cont.value.length != 0 && !replaced) {
    cont.value =
      cont.value.slice(0, selectionPos) +
      val +
      cont.value.slice(selectionPos + val.length);
  }

  closeSearchReferenceModal();
}

document.addEventListener("htmx:afterRequest", (e) => {
  if (!e.detail.successful) return;
  const triggerHeader = e.detail.xhr.getResponseHeader("HX-Trigger");
  if (!triggerHeader) return;
  try {
    const triggerData = JSON.parse(triggerHeader);
    if (triggerData.fileUploaded) {
      const { field, fileName } = triggerData.fileUploaded;
      const hiddenField = document.getElementById(`field-${field}`);
      if (hiddenField) hiddenField.value = fileName;

      const thumbnail = document.getElementById("field-thumbnail");
      if (thumbnail) {
        thumbnail.value = fileName.replaceAll(".", "_thumb.");
      }
    }
  } catch (err) {
    console.error("Failed to parse file upload response:", err);
  }
});

document.addEventListener("htmx:responseError", (e) => {
  const container = e.target.closest(".file-upload-container");
  if (container) {
    popup(
      "<i class='fa-solid fa-circle-xmark text-error'></i> " +
        e.detail.xhr.responseText,
    );
    const input = container.querySelector('input[type="file"]');
    if (input) input.value = "";
  }
});

let deleteFilePattern = /'([^']+)'/gim;
document.addEventListener("htmx:confirm", function (e) {
  if (!e.detail.question) return;

  e.preventDefault();
  let matches = e.detail.question.match(deleteFilePattern);

  document.getElementById("file-delete-description").textContent =
    e.detail.question;
  const fileDeleteModal = document.getElementById("file-delete-modal");
  document
    .getElementById("confirm-file-delete-button")
    .addEventListener("click", (_) => {
      e.detail.issueRequest(true);
      fileDeleteModal.close();

      if (matches != null && matches.length > 0) {
        const fieldName = matches[0].replaceAll("'", "");
        const fields = ["file-input", "alt-text", "field"];
        fields.forEach((prefix) => {
          document.getElementById(prefix + "-" + fieldName).value = "";
        });
        document.getElementById("field-" + fieldName).value = "";

        const thumbnail = document.getElementById("field-thumbnail");
        if (thumbnail) {
          thumbnail.value = "";
          document.getElementById("alt-text-thumbnail").value = "";
        }
      }
    });
  fileDeleteModal.showModal();
});
