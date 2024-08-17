document.addEventListener('DOMContentLoaded', function () {
    const searchInput = document.querySelector('input[name="q"]');
    const searchResults = document.getElementById('search-results');

    searchInput.addEventListener('focus', function () {
        searchResults.classList.remove('hidden');
    });

    searchInput.addEventListener('blur', function () {
        // On utilise setTimeout pour permettre au clic sur un résultat de se produire avant de cacher les résultats
        setTimeout(() => {
            searchResults.classList.add('hidden');
        }, 100);
    });
});
