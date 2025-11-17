function confirmListAction(identifier, action) {
  document.getElementById('list-action-name').value = action;
  document.getElementById('list-action-item-id').value = identifier;
  document.getElementById('list-action-desc').innerHTML = `Do you really want to <b>${action}</b> page <b>${identifier}</b>?`;
  document.getElementById('list-action-modal').showModal();
}

document.addEventListener("htmx:afterRequest", (e) => {
  document.getElementById('list-action-modal').close();
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
