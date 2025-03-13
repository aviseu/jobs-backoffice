(function () {
  "use strict";

  /*Featured Carouel Start*/
  var swiper01 = new Swiper(".mySwipermain", {
    loop: true,
    slidesPerView: 6,
    spaceBetween: 10,
    freeMode: true,
  });
  /*Featured Carouel End*/

  /*Featured Carouel Start*/
  var swiper1 = new Swiper(".mySwiper", {
    autoplay: {
      enabled: false,
    },

    loop: true,
    slidesPerView: 1,
    spaceBetween: 10,
    freeMode: true,
    navigation: {
      nextEl: ".swiper-button-next",
      prevEl: ".swiper-button-prev",
    },
    breakpoints: {
      576: {
        slidesPerView: 2,
        spaceBetween: 10,
      },
      992: {
        slidesPerView: 3,
        spaceBetween: 50,
      },
      1400: {
        slidesPerView: 4,
        spaceBetween: 20,
      },
    },
  });
  /*Featured Carouel End*/

  /*latest Ads Carouel Start*/
  var swiper2 = new Swiper(".mySwiper2", {
    autoplay: {
      delay: 2500,
      disableOnInteraction: false,
    },
    loop: true,
    slidesPerView: 1,
    spaceBetween: 10,
    freeMode: true,
    navigation: {
      nextEl: ".swiper-button-next",
      prevEl: ".swiper-button-prev",
    },
    breakpoints: {
      768: {
        slidesPerView: 2,
        spaceBetween: 40,
      },
      1024: {
        slidesPerView: 3,
        spaceBetween: 50,
      },
      1400: {
        slidesPerView: 4,
        spaceBetween: 20,
      },
    },
  });
  /*latest Ads Carouel End*/

  /*Testinominal Carouel Start*/
  var swiper3 = new Swiper(".mySwiper3", {
    loop: true,
    slidesPerView: 1,
    spaceBetween: 10,
    navigation: {
      nextEl: ".swiper-button-next",
      prevEl: ".swiper-button-prev",
    },
    breakpoints: {
      768: {
        slidesPerView: 2,
        spaceBetween: 40,
      },
      1024: {
        slidesPerView: 2,
        spaceBetween: 50,
      },
      1400: {
        slidesPerView: 3,
        spaceBetween: 20,
      },
    },
  });
  /*Testinominal Carouel End*/
  /*Testinominal Carouel Start*/
  var swiper4 = new Swiper(".mySwiper4", {
    autoplay: {
      delay: 2500,
      disableOnInteraction: false,
    },
    loop: true,
    slidesPerView: 1,
    spaceBetween: 10,
    navigation: {
      nextEl: ".swiper-button-next",
      prevEl: ".swiper-button-prev",
    },
    breakpoints: {
      768: {
        slidesPerView: 2,
        spaceBetween: 40,
      },
      1024: {
        slidesPerView: 2,
        spaceBetween: 50,
      },
      1400: {
        slidesPerView: 4,
        spaceBetween: 20,
      },
    },
  });
  /*Testinominal Carouel End*/
  var swiper = new Swiper(".thumb-2", {
    spaceBetween: 10,
    slidesPerView: 4,
    freeMode: true,
    loop: true,
    watchSlidesProgress: true,
  });
  var swiper12 = new Swiper(".thumb-1", {
    spaceBetween: 10,
    loop: true,
    navigation: {
      nextEl: ".swiper-button-next",
      prevEl: ".swiper-button-prev",
    },
    thumbs: {
      swiper: swiper,
    },
  });

  /*Into Page Start*/
  var swiper5 = new Swiper(".introswipe", {
    loop: true,
    autoplay: {
      delay: 2500,
      disableOnInteraction: false,
    },
    pagination: {
      el: ".swiper-pagination",
      clickable: true,
    },
  });
  /*Into Page End*/

  /*Home Page Start*/
  var swiper5 = new Swiper(".homeSwiper", {
    loop: true,
    autoplay: {
      delay: 2500,
      disableOnInteraction: false,
    },
    pagination: {
      el: ".swiper-pagination",
      clickable: true,
    },
  });
  /*Home Page End*/
})();
