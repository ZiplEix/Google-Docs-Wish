function generateUUID() {
    // Génère un UUID v4
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
        var r = Math.random() * 16 | 0,
            v = c === 'x' ? r : (r & 0x3 | 0x8);
        return v.toString(16);
    });
}

document.addEventListener('DOMContentLoaded', function() {
    function addNewPage() {
        var editorWrapper = document.getElementById('editorwrapper');
        var newPageContainer = document.createElement('div');
        var pageId = generateUUID(); // Générer un UUID pour la page
        newPageContainer.className = 'bg-white shadow-lg';
        newPageContainer.style.width = '794px';
        newPageContainer.style.height = '1123px';

        var newEditor = document.createElement('div');
        newEditor.className = 'page';
        newEditor.contentEditable = 'true';
        newEditor.style.width = '100%';
        newEditor.style.height = '100%';
        newEditor.style.boxSizing = 'border-box';
        newEditor.id = generateUUID(); // Ajouter l'ID à l'éditeur

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

            // Ajouter un UUID au paragraphe
            newParagraph.id = generateUUID();

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


// Voila donc l'état du code actuellement :
// HTML:
//             <main class="bg-gray-50 flex flex-col pb-8">
//                 <!-- The document content -->
//                 <div id="editorwrapper" class="flex flex-col items-center justify-center min-h-full gap-y-4">
//                     <div class="bg-white shadow-lg" style="width: 794px; height: 1123px;"> <!-- Une page -->
//                         <div contenteditable="true" class="page" style="width: 100%; height: 100%; box-sizing: border-box;">
//                             <!-- Page content here -->
//                         </div>
//                     </div>
//                 </div>
//                 <script src="/document_editor/editor.js"></script>
//             </main>

// /document_editor/editor.js:
// document.addEventListener('DOMContentLoaded', function() {
//     function addNewPage() {
//         var editorWrapper = document.getElementById('editorwrapper');
//         var newPageContainer = document.createElement('div');
//         newPageContainer.className = 'bg-white shadow-lg';
//         newPageContainer.style.width = '794px';
//         newPageContainer.style.height = '1123px';

//         var newEditor = document.createElement('div');
//         newEditor.className = 'page';
//         newEditor.contentEditable = 'true';
//         newEditor.style.width = '100%';
//         newEditor.style.height = '100%';
//         newEditor.style.boxSizing = 'border-box';

//         newPageContainer.appendChild(newEditor);
//         editorWrapper.appendChild(newPageContainer);

//         // Optionnel: mettre le focus sur la nouvelle page
//         newEditor.focus();
//     }

//     function isPageFull(editor) {
//         // Vérifie si l'utilisateur est arrivé en bas de la page
//         return editor.scrollHeight > editor.clientHeight;
//     }

//     document.getElementById('editorwrapper').addEventListener('keydown', function(event) {
//         if (event.key === 'Enter') {
//             event.preventDefault(); // Empêche le comportement par défaut (ajout d'un <br>)

//             var selection = window.getSelection();
//             var range = selection.getRangeAt(0);
//             var editor = event.target;

//             // Crée un nouveau paragraphe
//             var newParagraph = document.createElement('p');
//             newParagraph.innerHTML = '&#8203;'; // Utilise un espace non sécable pour rendre le paragraphe visible

//             // Insère le nouveau paragraphe à la fin du conteneur #editor
//             editor.appendChild(newParagraph);

//             // Positionne le curseur dans le nouveau paragraphe
//             range.setStart(newParagraph, 0);
//             range.collapse(true);
//             selection.removeAllRanges();
//             selection.addRange(range);

//             // Vérifie si une nouvelle page doit être ajoutée
//             var pages = document.querySelectorAll('#editorwrapper .bg-white.shadow-lg');
//             var currentPage = pages[pages.length - 1];
//             var currentEditor = currentPage.querySelector('.page');

//             // Si la page actuelle est pleine, ajoute une nouvelle page
//             if (isPageFull(currentEditor)) {
//                 addNewPage();
//             }
//         }
//     });

//     window.addEventListener('load', function() {
//         var firstPage = document.querySelector('#editorwrapper .page');
//         if (firstPage) {
//             firstPage.focus(); // Met le focus sur la première page
//         }
//     });
// });

// Les nouvelles page se créent parfaitement lorsque je dépasse de la dernière page.
// Par contre quelques chose se passent mal.
// J'ai crée un document vierge. Dans la première balise p j'ai marqué "test 1" puis j'ai appuyé sur la touche entré jusqu'à créer une nouvelle page. Sur cette nouvelle page j'ai écrit "Test 2" puis je suis retourné tout en bas de la première page et j'ai appuyé sur la touche entré, mais aucune nouvelle page ne s'est créé entre la page "TEST 1" et la page "TEST 2". Comment régler ca ?
