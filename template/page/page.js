function confirmListAction(identifier, action) {
  list_action_name.value = action;
  list_action_item_id.value = identifier;
  list_action_desc.innerHTML = `Do you really want to <b>${action}</b> page <b>${identifier}</b>?`;
  list_action_modal.showModal();
}

document.addEventListener("htmx:afterRequest", (e) => {
  list_action_modal.close();
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
