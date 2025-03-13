/*Testinominal Carouel End*/
var swiper = new Swiper(".thumb-2", {
  spaceBetween: 10,
  slidesPerView: 5,
  freeMode: true,
  loop: true,
  watchSlidesProgress: true,
  autoplay: {
    delay: 2500,
    disableOnInteraction: false,
  },
});
var swiper12 = new Swiper(".thumb-1", {
  spaceBetween: 10,
  autoplay: {
    delay: 2500,
    disableOnInteraction: false,
  },
  loop: true,
  navigation: {
    nextEl: ".swiper-button-next",
    prevEl: ".swiper-button-prev",
  },
  thumbs: {
    swiper: swiper,
  },
});
var swiper3 = new Swiper(".related", {
  loop: true,
  slidesPerView: 1,
  spaceBetween: 30,
  autoplay: {
    delay: 1500,
    disableOnInteraction: false,
  },
  breakpoints: {
    768: {
      slidesPerView: 2,
    },
    1024: {
      slidesPerView: 2,
    },
    1400: {
      slidesPerView: 3,
    },
  },
});
var swiper13 = new Swiper(".related2", {
  loop: true,
  slidesPerView: 1,
  spaceBetween: 30,
  autoplay: {
    delay: 1500,
    disableOnInteraction: false,
  },
  breakpoints: {
    768: {
      slidesPerView: 2,
    },
  },
});
var swiper14 = new Swiper(".related3", {
  loop: true,
  slidesPerView: 1,
  spaceBetween: 30,
  autoplay: {
    delay: 1500,
    disableOnInteraction: false,
  },
  breakpoints: {
    576: {
      slidesPerView: 2,
      spaceBetween: 10,
    },
    992: {
      slidesPerView: 3,
    },
    1400: {
      slidesPerView: 4,
      spaceBetween: 20,
    },
  },
});

(function () {
  "use strict";
  const lightboxElements = document.querySelectorAll(".glightbox");
  if (lightboxElements.length > 0) {
    const lightboxVideo = GLightbox({
      selector: ".glightbox",
    });

    lightboxVideo.on("slide_changed", ({ prev, current }) => {
      const { slideIndex, slideNode, slideConfig, player } = current;
      // Your code to handle slide changes goes here
    });
  }
  // Select all elements with the specified class
  var elements = document.querySelectorAll(".rating-stars2");

  // Loop through each element and create an instance of raterJs
  elements.forEach(function (element) {
    var starRating = raterJs({
      starSize: 24,
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
})();
