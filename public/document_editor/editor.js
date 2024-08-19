document.addEventListener('DOMContentLoaded', function() {
    function addNewPage() {
        var editorWrapper = document.getElementById('editorwrapper');
        var newPageContainer = document.createElement('div');
        newPageContainer.className = 'bg-white shadow-lg';
        newPageContainer.style.width = '794px';
        newPageContainer.style.height = '1123px';

        var newEditor = document.createElement('div');
        newEditor.className = 'page';
        newEditor.contentEditable = 'true';
        newEditor.style.width = '100%';
        newEditor.style.height = '100%';
        newEditor.style.boxSizing = 'border-box';

        newPageContainer.appendChild(newEditor);
        editorWrapper.appendChild(newPageContainer);

        // Optionnel: mettre le focus sur la nouvelle page
        newEditor.focus();
    }

    function isPageFull(editor) {
        // Vérifie si l'utilisateur est arrivé en bas de la page
        return editor.scrollHeight > editor.clientHeight;
    }

    document.getElementById('editorwrapper').addEventListener('keydown', function(event) {
        if (event.key === 'Enter') {
            event.preventDefault(); // Empêche le comportement par défaut (ajout d'un <br>)

            var selection = window.getSelection();
            var range = selection.getRangeAt(0);
            var editor = event.target;

            // Crée un nouveau paragraphe
            var newParagraph = document.createElement('p');
            newParagraph.innerHTML = '&#8203;'; // Utilise un espace non sécable pour rendre le paragraphe visible

            // Insère le nouveau paragraphe à la fin du conteneur #editor
            editor.appendChild(newParagraph);

            // Positionne le curseur dans le nouveau paragraphe
            range.setStart(newParagraph, 0);
            range.collapse(true);
            selection.removeAllRanges();
            selection.addRange(range);

            // Vérifie si une nouvelle page doit être ajoutée
            var pages = document.querySelectorAll('#editorwrapper .bg-white.shadow-lg');
            var currentPage = pages[pages.length - 1];
            var currentEditor = currentPage.querySelector('.page');

            // Si la page actuelle est pleine, ajoute une nouvelle page
            if (isPageFull(currentEditor)) {
                addNewPage();
            }
        }
    });

    window.addEventListener('load', function() {
        var firstPage = document.querySelector('#editorwrapper .page');
        if (firstPage) {
            firstPage.focus(); // Met le focus sur la première page
        }
    });
});
