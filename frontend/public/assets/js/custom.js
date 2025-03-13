"use strict";

/* .. tooltip .. */
const tooltipTriggerList = document.querySelectorAll(
  '[data-bs-toggle="tooltip"]'
);
const tooltipList = [...tooltipTriggerList].map(
  (tooltipTriggerEl) => new bootstrap.Tooltip(tooltipTriggerEl)
);

/* .. popover .. */
const popoverTriggerList = document.querySelectorAll(
  '[data-bs-toggle="popover"]'
);
const popoverList = [...popoverTriggerList].map(
  (popoverTriggerEl) => new bootstrap.Popover(popoverTriggerEl)
);

/* Start::back-to-top */
 const scrollToTop = document.querySelector(".scrollToTop");
 const $rootElement = document.documentElement;
 const $body = document.body;
 window.onscroll = () => {
   const scrollTop = window.scrollY || window.pageYOffset;
   const clientHt = $rootElement.scrollHeight - $rootElement.clientHeight;
   if (window.scrollY > 100) {
     scrollToTop.style.display = "flex";
   } else {
     scrollToTop.style.display = "none";
   }
 };
 scrollToTop.onclick = () => {
   window.scrollTo(0, 0);
 };
/* End::back-to-top */

// reveal items on scroll
function reveal() {
  var reveals = document.querySelectorAll(".reveal");

  if (reveals) {
    for (var i = 0; i < reveals.length; i++) {
      var windowHeight = window.innerHeight;
      var cardTop = reveals[i].getBoundingClientRect().top;
      var cardRevealPoint = 130;

      if (cardTop < windowHeight - cardRevealPoint) {
        reveals[i].classList.add("reveal-active");
      } else {
        reveals[i].classList.remove("reveal-active");
      }
    }
  }
}
window.addEventListener("scroll", reveal);
reveal(); //Run

var coverImg = document.querySelectorAll(".cover-image");

coverImg.forEach((ele) => {
  var attr = ele.getAttribute("data-bs-image-src");
  if (attr && typeof attr !== typeof undefined && attr !== false) {
    ele.style.background = `url(${attr}) center center`;
  }
});

/* Choices JS */
document.addEventListener("DOMContentLoaded", function () {
  var genericExamples = document.querySelectorAll("[data-trigger]");
  for (let i = 0; i < genericExamples.length; ++i) {
    var element = genericExamples[i];
    new Choices(element, {
      allowHTML: true,
      placeholderValue: "This is a placeholder set in the config",
      searchPlaceholderValue: "Search Here",
    });
  }
});
/* Choices JS */

// Select all elements with the specified class
var elements = document.querySelectorAll(".rating-stars");

// Loop through each element and create an instance of raterJs
elements.forEach(function (element) {
  var starRating = raterJs({
    starSize: 14,
    max: 5,
    rating: 4,
    element: element,
    disableText: "Custom disable text!",
    ratingText: "My custom rating text {rating}",
    showToolTip: true,
    rateCallback: function ratingCallback(rating, done) {
      starRating.setRating(rating);
      starRating.disable();
      done();
    },
  });
});
/* Custom messages */

var links = document.getElementsByClassName("location-link");
const modal = document.getElementById("locationModal");
for (var i = 0; i < links.length; i++) {
  links[i].addEventListener("click", function (event) {
    event.preventDefault();
    // closeModal(); // Close the modal
    // openModal()
    addLinkToInput(this);
    document.querySelector(".btn-close").click();
  });
}

function addLinkToInput(link) {
  var inputBox = document.getElementById("sale-location");
  inputBox.value = link.innerText;
}

/*login-popup*/
if (document.querySelector("#popup")) {
  var myModal = new bootstrap.Modal(document.getElementById("popup"), {});
  myModal.show();
}



if (document.querySelector("#product-price-range")) {
  // products price range slider
  var mergingTooltipSlider = document.getElementById("product-price-range");
  noUiSlider.create(mergingTooltipSlider, {
    start: [10000, 100000],
    connect: true,
    tooltips: [true, true],
    range: {
      min: 10000,
      max: 100000,
    },
  });
}
if (document.querySelector("#price-range")) {
  // products price range slider
  var mergingTooltipSlider = document.getElementById("price-range");
  var input0 = document.getElementById("input-with-keypress-0");
  var input1 = document.getElementById("input-with-keypress-1");
  var inputs = [input0, input1];
  noUiSlider.create(mergingTooltipSlider, {
    start: [0, 100000],
    connect: true,
    tooltips: [true, true],
    range: {
      min: 10000,
      max: 100000,
    },
  });
  mergingTooltipSlider.noUiSlider.on("update", function (values, handle) {
    inputs[handle].value = values[handle];
  });
}

if (document.querySelector("#vscroll ")) {
  /* header dropdowns scroll */
  var myHeaderShortcut = document.getElementById("vscroll");
  new SimpleBar(myHeaderShortcut, { autoHide: true });
}

/* header theme toggle */
function toggleTheme() {
  let html = document.querySelector("html");
  if (html.getAttribute("data-theme-mode") === "dark") {
    html.setAttribute("data-theme-mode", "light");
    html.removeAttribute("data-bg-theme");
    if (!localStorage.getItem("primaryRGB")) {
      html.setAttribute("style", "");
    }
    document.querySelector("#switcher-light-theme").checked = true;
    document
      .querySelector("html")
      .style.removeProperty("--body-bg-rgb", localStorage.bodyBgRGB);
    checkOptions();
    html.style.removeProperty("--body-bg-rgb2");
    html.style.removeProperty("--light-rgb");
    html.style.removeProperty("--form-control-bg");
    document.querySelector("#switcher-light-theme").checked = true;
    document.querySelector("#switcher-background5").checked = false;
    document.querySelector("#switcher-background4").checked = false;
    document.querySelector("#switcher-background3").checked = false;
    document.querySelector("#switcher-background2").checked = false;
    document.querySelector("#switcher-background1").checked = false;
    localStorage.removeItem("roxlistdarktheme");
    localStorage.removeItem("bodylightRGB");
    localStorage.removeItem("bodyBgRGB");
  } else {
    html.setAttribute("data-theme-mode", "dark");
    document.querySelector("#switcher-dark-theme").checked = true;
    if (!localStorage.getItem("primaryRGB")) {
      html.setAttribute("style", "");
    }
    checkOptions();
    document.querySelector("#switcher-dark-theme").checked = true;
    document.querySelector("#switcher-background5").checked = false;
    document.querySelector("#switcher-background4").checked = false;
    document.querySelector("#switcher-background3").checked = false;
    document.querySelector("#switcher-background2").checked = false;
    document.querySelector("#switcher-background1").checked = false;
    localStorage.setItem("roxlistdarktheme", "true");
    localStorage.removeItem("bodylightRGB");
    localStorage.removeItem("bodyBgRGB");
  }
}
let layoutSetting = document.querySelector(".layout-setting");
if (document.querySelector(".layout-setting")) {
  layoutSetting.addEventListener("click", toggleTheme);
}
/* header theme toggle */