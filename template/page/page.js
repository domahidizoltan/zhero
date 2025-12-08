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

function submitPreview(cls) {
  const form = document.getElementById("edit-page-form");
  const restoreAction = form.action;
  const restoreTarget = form.target;
  const host =
    window.location.protocol + "//" + window.location.hostname + ":" + 8080;
  try {
    form.action = host + "/preview?class=" + cls;
    form.target = "_blank";
    form.submit();
  } finally {
    form.action = restoreAction;
    form.target = restoreTarget;
  }
}
