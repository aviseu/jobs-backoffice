document.addEventListener('DOMContentLoaded', function () {
    'use strict';
    var form = document.querySelector('.needs-validation');
    
    form.addEventListener('submit', function (event) {
        event.preventDefault(); // Prevent the default form submission behavior
        
        if (form.checkValidity()) {
            form.classList.add('was-validated');
            window.location.href = 'index.html'; // Redirect to index.html if valid
        } else {
            event.stopPropagation(); // Prevents further propagation of the current event
            form.classList.add('was-validated'); // Add validation feedback
        }
    }, false);
});