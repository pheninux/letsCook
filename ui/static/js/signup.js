(function () {
    let image = document.getElementById('avatar');
    image.src = "../static/img/defaultAvatar.png"

})();

function uploadFile(event) {
    let image = document.getElementById('avatar');
    image.src = URL.createObjectURL(event.target.files[0]);
}
