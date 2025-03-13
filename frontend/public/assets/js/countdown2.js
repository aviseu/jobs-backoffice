
// Get all elements with the specified class name
const countdownElements = document.querySelectorAll('.hour');
const countdownElements2 = document.querySelectorAll('.hour2');

// Update the countdown for each element
function updateCountdown() {
  // Get the current date and time
  const now = new Date().getTime();

  // Set the target date and time (e.g., 1 January 2024 00:00:00)
  const target = new Date('2025-01-01T00:00:00').getTime();

  // Calculate the time remaining
  let timeRemaining = target - now;

  // If the countdown has reached zero or is negative, reset it to the next target date
  if (timeRemaining <= 0) {
    const nextTarget = new Date(target + (24 * 60 * 60 * 1000)); // Add 24 hours to the target time
    timeRemaining = nextTarget - now;
  }

  // Calculate days, hours, minutes, and seconds
  const days = Math.floor(timeRemaining / (1000 * 60 * 60 * 24));
  const hours = Math.floor((timeRemaining % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
  const minutes = Math.floor((timeRemaining % (1000 * 60 * 60)) / (1000 * 60));
  const seconds = Math.floor((timeRemaining % (1000 * 60)) / 1000);

  // Update the content of each countdown element
  Array.from(countdownElements).forEach(element => {
    element.innerHTML = `${days} days : ${hours} hours : ${minutes} minutes : ${seconds} seconds`;
  });

  // Repeat the update every second
  setTimeout(updateCountdown, 1000);
}
function updateCountdown2() {
  // Get the current date and time
  const now = new Date().getTime();

  // Set the target date and time (e.g., 1 January 2024 00:00:00)
  const target = new Date('2025-01-01T00:00:00').getTime();

  // Calculate the time remaining
  let timeRemaining = target - now;

  // If the countdown has reached zero or is negative, reset it to the next target date
  if (timeRemaining <= 0) {
    const nextTarget = new Date(target + (24 * 60 * 60 * 1000)); // Add 24 hours to the target time
    timeRemaining = nextTarget - now;
  }

  // Calculate days, hours, minutes, and seconds
  const hours = Math.floor((timeRemaining % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
  const minutes = Math.floor((timeRemaining % (1000 * 60 * 60)) / (1000 * 60));
  const seconds = Math.floor((timeRemaining % (1000 * 60)) / 1000);

  // Update the content of each countdown element
  Array.from(countdownElements2).forEach(element => {
    element.innerHTML = `${hours} hours : ${minutes} minutes : ${seconds} seconds`;
  });

  // Repeat the update every second
  setTimeout(updateCountdown2, 1000);
}

countdownElements.forEach(a =>{
    updateCountdown();
})

countdownElements2.forEach(a =>{
    updateCountdown2();
})
// Start the countdown