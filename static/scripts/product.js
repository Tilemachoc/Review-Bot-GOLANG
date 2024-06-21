document.addEventListener("DOMContentLoaded", function() {
    document.getElementById("product-form").addEventListener("submit", function(event) {
        event.preventDefault();

        const selectedProduct = document.getElementById("product").value;

        window.location.href = `/?product=${selectedProduct}`
    })
})