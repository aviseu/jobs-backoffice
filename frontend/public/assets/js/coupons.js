
class CouponManager {
    constructor() {
      this.cpnBtns = document.querySelectorAll(".coupon-btn");
      this.cpnCodes = document.querySelectorAll(".coupon-code");

      this.cpnBtns.forEach((btn) => {
        btn.addEventListener("click", this.copyCode.bind(this, btn));
      });
    }

    copyCode(btn) {
      const index = Array.from(this.cpnBtns).indexOf(btn);
      const cpnCode = this.cpnCodes[index];

      navigator.clipboard.writeText(cpnCode.innerHTML);
      btn.innerHTML = "COPIED";

      setTimeout(() => {
        btn.innerHTML = "COPY CODE";
      }, 3000);
    }
  }

  const couponManager = new CouponManager();