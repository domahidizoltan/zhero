document.addEventListener("DOMContentLoaded", initScripts);

function initScripts(evt) {
  const searchIcon = document.querySelector(".search-icon");
  const searchPopup = document.getElementById("searchPopup");
  const closeSearch = document.getElementById("closeSearch");
  const searchInput = document.getElementById("searchInput");
  const searchButton = document.getElementById("searchButton");

  searchIcon.addEventListener("click", (e) => {
    e.preventDefault();
    searchPopup.classList.add("active");
    searchInput.focus();
  });

  closeSearch.addEventListener("click", () => {
    searchPopup.classList.remove("active");
  });

  searchButton.addEventListener("click", () => {
    const query = searchInput.value;
    if (query) {
      alert(`Searching for: ${query}`);
    }
  });

  searchPopup.addEventListener("click", (e) => {
    if (e.target === searchPopup) {
      searchPopup.classList.remove("active");
    }
  });

  document.addEventListener("keydown", (e) => {
    if (e.key === "Escape" && searchPopup.classList.contains("active")) {
      searchPopup.classList.remove("active");
    }
  });

  showCurtain();
}

const pageKey = "zhero-ageConfirmed";
const period24h = 86400000;
function showCurtain() {
  const modal = document.getElementById("age-modal");
  if (modal) {
    const confirmed = sessionStorage.getItem(pageKey);
    if (confirmed && Date.now() - parseInt(confirmed) < period24h) return;

    modal.showModal();
    document.getElementById("age-confirm-btn").onclick = function () {
      sessionStorage.setItem(pageKey, Date.now());
      modal.close();
    };
  }
}
