document.addEventListener("DOMContentLoaded", initScripts);

function initScripts(evt) {
  setTimeout(function () {
    document.querySelectorAll(".alert-success").forEach((el) => {
      el.classList.add("hide");
    });
  }, 7000);
}

function popup(text) {
  document.getElementById("popup-text").innerHTML = text;
  document.getElementById("popup-modal").showModal();
}
