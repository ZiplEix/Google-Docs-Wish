function generateUUID() {
    // Génère un UUID v4
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
        var r = Math.random() * 16 | 0,
            v = c === 'x' ? r : (r & 0x3 | 0x8);
        return v.toString(16);
    });
}

document.addEventListener('DOMContentLoaded', function() {
    function addNewPage(afterPage) {
        var editorWrapper = document.getElementById('editorwrapper');
        var newPageContainer = document.createElement('div');
        var pageId = generateUUID(); // Générer un UUID pour la page
        newPageContainer.className = 'bg-white shadow-lg';
        newPageContainer.style.width = '794px';
        newPageContainer.style.height = '1123px';
        newPageContainer.id = 'page-' + pageId; // Ajouter l'ID à la page

        var newEditor = document.createElement('div');
        newEditor.className = 'page';
        newEditor.contentEditable = 'true';
        newEditor.style.width = '100%';
        newEditor.style.height = '100%';
        newEditor.style.boxSizing = 'border-box';
        newEditor.id = 'editor-' + pageId; // Ajouter l'ID à l'éditeur

        newPageContainer.appendChild(newEditor);

        // Si 'afterPage' est fourni, insère la nouvelle page après la page spécifiée
        if (afterPage) {
            editorWrapper.insertBefore(newPageContainer, afterPage.nextSibling);
        } else {
            // Sinon, ajoute la nouvelle page à la fin
            editorWrapper.appendChild(newPageContainer);
        }

        // Optionnel: mettre le focus sur la nouvelle page
        newEditor.focus();
    }

    function isPageFull(editor) {
        // Vérifie si l'utilisateur est arrivé en bas de la page
        return editor.scrollHeight > editor.clientHeight;
    }

    function addParagraphToNextPage(nextPage) {
        var firstParagraph = document.createElement('p');
        // Utilisez un espace non sécable pour rendre le paragraphe visible
        firstParagraph.innerHTML = '&#8203;';
        firstParagraph.id = generateUUID(); // Ajouter un UUID au paragraphe
        nextPage.insertBefore(firstParagraph, nextPage.firstChild); // Insère le paragraphe au début de la page
        return firstParagraph;
    }

    document.getElementById('editorwrapper').addEventListener('keydown', function(event) {
        if (event.key === 'Enter') {
            event.preventDefault(); // Empêche le comportement par défaut (ajout d'un <br>)

            var selection = window.getSelection();
            var range = selection.getRangeAt(0);
            var editor = event.target;

            // Crée un nouveau paragraphe
            var newParagraph = document.createElement('p');
            // Utilisez un espace non sécable pour rendre le paragraphe visible
            newParagraph.innerHTML = '&#8203;';
            newParagraph.id = generateUUID(); // Ajouter un UUID au paragraphe

            // Trouve la page actuelle
            var currentPage = editor.closest('.bg-white.shadow-lg');
            var nextPage = currentPage.nextElementSibling ? currentPage.nextElementSibling.querySelector('.page') : null;

            if (isPageFull(editor)) {
                // Si la page actuelle est pleine
                if (nextPage) {
                    // S'il y a une page suivante, ajoute le paragraphe au début de la page suivante
                    var firstParagraph = addParagraphToNextPage(nextPage);
                    // Positionne le curseur dans le nouveau paragraphe sur la page suivante
                    range.setStart(firstParagraph, 0);
                    range.collapse(true);
                    selection.removeAllRanges();
                    selection.addRange(range);
                } else {
                    // Sinon, ajoute une nouvelle page après la page actuelle
                    addNewPage(currentPage);
                }
            } else {
                // Si la page actuelle n'est pas pleine
                // Si le curseur est au début du paragraphe, insère le nouveau paragraphe avant celui-ci
                if (range.startOffset === 0 && range.endOffset === 0) {
                    // Insère le nouveau paragraphe avant le paragraphe actuel
                    var currentParagraph = range.startContainer;
                    currentParagraph.parentNode.insertBefore(newParagraph, currentParagraph);
                } else {
                    // Sinon, ajoute le paragraphe à la fin
                    editor.appendChild(newParagraph);
                }
                // Positionne le curseur dans le nouveau paragraphe
                range.setStart(newParagraph, 0);
                range.collapse(true);
                selection.removeAllRanges();
                selection.addRange(range);
            }
        }
    });

    window.addEventListener('load', function() {
        var firstPage = document.querySelector('#editorwrapper .page');
        if (firstPage) {
            firstPage.focus(); // Met le focus sur la première page
            // add the first paragraph
            var firstParagraph = document.createElement('p');
            firstParagraph.innerHTML = '&#8203;';
            firstParagraph.id = generateUUID(); // Ajouter un UUID au paragraphe
            firstPage.appendChild(firstParagraph);
        }
    });
});
