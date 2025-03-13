/* To choose date */
flatpickr("#date", {});

const dateInputs = document.querySelectorAll(".date");

dateInputs.forEach((input) => {
  flatpickr(input, {});
});
