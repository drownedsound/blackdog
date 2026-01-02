document.addEventListener("DOMContentLoaded", () => {
    const btn = document.getElementById("submitButton");
    if (btn) {
        btn.addEventListener("click", (event) => {
            window.alert("application submitted")
        })
    }
})
