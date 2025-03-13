"use strict";

let loadFile = function (event) {
  var reader = new FileReader();
  reader.onload = function () {
    var output = document.getElementById("profile-img");
    if (event.target.files[0].type.match("image.*")) {
      output.src = reader.result;

      // Save the image data to local storage
      localStorage.setItem("profileImage", reader.result);
    } else {
      event.target.value = "";
      alert("please select a valid image");
    }
  };
  reader.readAsDataURL(event.target.files[0]);
};
(function () {
  let profileImage = localStorage.getItem("profileImage");
  if (profileImage) {
    document.getElementById("profile-img").src = profileImage;
  }
  let ProfileChange = document.querySelector("#profile-change");

  ProfileChange.addEventListener("change", loadFile);

  // for mail language
  const multipleCancelButton1 = new Choices("#language", {
    allowHTML: true,
    removeItemButton: true,
  });

  //To choose date
  flatpickr(".date", {});

  /*Basic Gallery */
  var lightbox = GLightbox();
  lightbox.on('open', (target) => {
  });
  var lightbox1 = GLightbox({
      selector: ".gallery"
  });
})();
