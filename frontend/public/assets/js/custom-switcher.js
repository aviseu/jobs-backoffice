"use strict";

let mainContent;
(function () {

    let html = document.querySelector('html');

    //RTL
    if (!localStorage.getItem("roxlistrtl")) {
        // html.setAttribute("dir" , "rtl") // for rtl version
    }

    //Light Theme Style
    if (!localStorage.getItem("roxlistlighttheme")) {
        // html.setAttribute("data-theme-mode" , "light") // for light theme
    }

    //Dark Theme Style
    if (!localStorage.getItem("roxlistdarktheme")) {
        // html.setAttribute("data-theme-mode" , "dark") // for dark theme
    }

    /*RTL Start*/
    if (html.getAttribute('dir') === "rtl") {
        rtlFn();
    }
    /*RTL End*/

    if (document.querySelector("#switcher-canvas")) {
        localStorageBackup();
    }
    localStorageBackup();
    switcherClick();
    checkOptions();

})();

function switcherClick() {
    let ltrBtn, rtlBtn, lightBtn, darkBtn,  resetBtn,
    primaryDefaultColor1Btn, primaryDefaultColor2Btn, primaryDefaultColor3Btn, primaryDefaultColor4Btn, primaryDefaultColor5Btn,
    bgDefaultColor1Btn, bgDefaultColor2Btn, bgDefaultColor3Btn, bgDefaultColor4Btn, bgDefaultColor5Btn,  ResetAll;
    let html = document.querySelector('html');
    lightBtn = document.querySelector('#switcher-light-theme');
    darkBtn = document.querySelector('#switcher-dark-theme');
    ltrBtn = document.querySelector('#switcher-ltr');
    rtlBtn = document.querySelector('#switcher-rtl');
    resetBtn = document.querySelector('#resetbtn');
    primaryDefaultColor1Btn = document.querySelector('#switcher-primary1');
    primaryDefaultColor2Btn = document.querySelector('#switcher-primary2');
    primaryDefaultColor3Btn = document.querySelector('#switcher-primary3');
    primaryDefaultColor4Btn = document.querySelector('#switcher-primary4');
    primaryDefaultColor5Btn = document.querySelector('#switcher-primary5');
    bgDefaultColor1Btn = document.querySelector('#switcher-background1');
    bgDefaultColor2Btn = document.querySelector('#switcher-background2');
    bgDefaultColor3Btn = document.querySelector('#switcher-background3');
    bgDefaultColor4Btn = document.querySelector('#switcher-background4');
    bgDefaultColor5Btn = document.querySelector('#switcher-background5');
    ResetAll = document.querySelector('#reset-all');


    // primary theme
    if(primaryDefaultColor1Btn){
        primaryDefaultColor1Btn.addEventListener('click', () => {
            localStorage.setItem("primaryRGB", "21, 62, 172");
            html.style.setProperty('--primary-rgb', `21, 62, 172`);
            document.querySelector("#switcher-primary1").checked = true
            updateColors();
        })
    }
    if(primaryDefaultColor2Btn){
        primaryDefaultColor2Btn.addEventListener('click', () => {
            localStorage.setItem("primaryRGB", "85, 27, 130");
            html.style.setProperty('--primary-rgb', `85, 27, 130`);
            document.querySelector("#switcher-primary2").checked = true
            updateColors();
        })
    }
    if(primaryDefaultColor3Btn){
        primaryDefaultColor3Btn.addEventListener('click', () => {
            localStorage.setItem("primaryRGB", "135, 90, 25");
            html.style.setProperty('--primary-rgb', `135, 90, 25`);
            document.querySelector("#switcher-primary3").checked = true
            updateColors();
        })
    }
    if(primaryDefaultColor4Btn){
        primaryDefaultColor4Btn.addEventListener('click', () => {
            localStorage.setItem("primaryRGB", "25, 135, 59");
            html.style.setProperty('--primary-rgb', `25, 135, 59`);
            document.querySelector("#switcher-primary4").checked = true
            updateColors();
        })
    }
    if(primaryDefaultColor5Btn){
        primaryDefaultColor5Btn.addEventListener('click', () => {
            localStorage.setItem("primaryRGB", "135, 20, 99");
            html.style.setProperty('--primary-rgb', `135, 20, 99`);
            document.querySelector("#switcher-primary5").checked = true
            updateColors();
        })
    }

    // Background theme
    if(bgDefaultColor1Btn){
        bgDefaultColor1Btn.addEventListener('click', () => {
           localStorage.setItem('bodyBgRGB', "12, 39, 95");
           localStorage.setItem('bodylightRGB', "29, 55, 109");
           html.setAttribute('data-theme-mode', 'dark');
           document.querySelector('html').style.setProperty('--body-bg-rgb', localStorage.bodyBgRGB);
           document.querySelector('html').style.setProperty('--body-bg-rgb2', localStorage.bodylightRGB);
           document.querySelector('html').style.setProperty('--light-rgb', "29, 55, 109");
           document.querySelector('html').style.setProperty('--form-control-bg', "rgb(29, 55, 109)");
           document.querySelector('#switcher-dark-theme').checked = true;
           document.querySelector("#switcher-background1").checked = true
       })
    }
    if(bgDefaultColor2Btn){
        bgDefaultColor2Btn.addEventListener('click', () => {
           localStorage.setItem('bodyBgRGB', "46, 10, 74");
           localStorage.setItem('bodylightRGB', "72, 31, 104");
           html.setAttribute('data-theme-mode', 'dark');
           document.querySelector('html').style.setProperty('--body-bg-rgb', localStorage.bodyBgRGB);
           document.querySelector('html').style.setProperty('--body-bg-rgb2', localStorage.bodylightRGB);
           document.querySelector('html').style.setProperty('--light-rgb', "72, 31, 104");
           document.querySelector('html').style.setProperty('--form-control-bg', "rgb(72, 31, 104)");
           document.querySelector('#switcher-dark-theme').checked = true;
           document.querySelector("#switcher-background2").checked = true
       })
    }
    if(bgDefaultColor3Btn){
        bgDefaultColor3Btn.addEventListener('click', () => {
           localStorage.setItem('bodyBgRGB', "106, 65, 4");
           localStorage.setItem('bodylightRGB', "120, 82, 25");
           html.setAttribute('data-theme-mode', 'dark');
           document.querySelector('html').style.setProperty('--body-bg-rgb', localStorage.bodyBgRGB);
           document.querySelector('html').style.setProperty('--body-bg-rgb2', localStorage.bodylightRGB);
           document.querySelector('html').style.setProperty('--light-rgb', "120, 82, 25");
           document.querySelector('html').style.setProperty('--form-control-bg', "rgb(120, 82, 25)");
           document.querySelector('#switcher-dark-theme').checked = true;
           document.querySelector("#switcher-background3").checked = true
       })
    }
    if(bgDefaultColor4Btn){
        bgDefaultColor4Btn.addEventListener('click', () => {
           localStorage.setItem('bodyBgRGB', "7, 86, 35");
           localStorage.setItem('bodylightRGB', "25, 109, 55");
           html.setAttribute('data-theme-mode', 'dark');
           document.querySelector('html').style.setProperty('--body-bg-rgb', localStorage.bodyBgRGB);
           document.querySelector('html').style.setProperty('--body-bg-rgb2', localStorage.bodylightRGB);
           document.querySelector('html').style.setProperty('--light-rgb', "25, 109, 55");
           document.querySelector('html').style.setProperty('--form-control-bg', "rgb(25, 109, 55)");
           document.querySelector('#switcher-dark-theme').checked = true;
           document.querySelector("#switcher-background4").checked = true
       })
    }
    if(bgDefaultColor5Btn){
        bgDefaultColor5Btn.addEventListener('click', () => {
           localStorage.setItem('bodyBgRGB', "86, 7, 58");
           localStorage.setItem('bodylightRGB', "121, 26, 88");
           html.setAttribute('data-theme-mode', 'dark');
           document.querySelector('html').style.setProperty('--body-bg-rgb', localStorage.bodyBgRGB);
           document.querySelector('html').style.setProperty('--body-bg-rgb2', localStorage.bodylightRGB);
           document.querySelector('html').style.setProperty('--light-rgb', "121, 26, 88");
           document.querySelector('html').style.setProperty('--form-control-bg', "rgb(121, 26, 88)");
           document.querySelector('#switcher-dark-theme').checked = true;
           document.querySelector("#switcher-background5").checked = true
       })
    }

    /* Light Layout Start */
    if(lightBtn){
        lightBtn.addEventListener('click', () => {
            lightFn();
            localStorage.removeItem("bodylightRGB")
            localStorage.removeItem("bodyBgRGB")
        })
    }
    /* Light Layout End */

    /* Dark Layout Start */
    if(darkBtn){
        darkBtn.addEventListener('click', () => {
            darkFn();
        });
    }
    /* Dark Layout End */

    /* rtl start */
    if(rtlBtn){
        rtlBtn.addEventListener('click', () => {
            localStorage.setItem("roxlistrtl", true);
            localStorage.removeItem("roxlistltr");
            rtlFn();
            if (document.querySelector(".noUi-target")) {
                console.log("working");
                document.querySelectorAll(".noUi-origin").forEach((e) => {
                    e.classList.add("transform-none");
                });
            }
        });
        /* rtl end */
    }

    /* ltr start */
    if(ltrBtn){
        ltrBtn.addEventListener('click', () => {
            //    local storage
            localStorage.setItem("roxlistltr", true);
            localStorage.removeItem("roxlistrtl");
            ltrFn();
            if (document.querySelector(".noUi-target")) {
                document.querySelectorAll(".noUi-origin").forEach((e) => {
                    e.classList.remove("transform-none");
                });
            }
        });
    }
    /* ltr end */

    // reset all start
    if(ResetAll){
        ResetAll.addEventListener('click', () => {
            ResetAllFn();
            if (document.querySelector(".noUi-target")) {
                console.log("working");
                document.querySelectorAll(".noUi-origin").forEach((e) => {
                    e.classList.add("transform-none");
                });
            }
        })
    }
    // reset all start
}

function ltrFn() {
    let html = document.querySelector('html')
    document.querySelector("#style")?.setAttribute("href", "/assets/libs/bootstrap/css/bootstrap.min.css");
    html.setAttribute("dir", "ltr");
    document.querySelector('#switcher-ltr').checked = true;
    checkOptions();
}

function rtlFn() {
    let html = document.querySelector('html');
    html.setAttribute("dir", "rtl");
    document.querySelector("#style")?.setAttribute("href", "/assets/libs/bootstrap/css/bootstrap.rtl.min.css");
    checkOptions();
}

function lightFn() {
    let html = document.querySelector('html');
    html.setAttribute('data-theme-mode', 'light');
    if(document.querySelector('#switcher-light-theme')){
        document.querySelector('#switcher-light-theme').checked = true;
    }
    localStorage.removeItem("roxlistdarktheme");
    localStorage.removeItem("roxlistbgColor");
    localStorage.removeItem("roxlistbgwhite");
    localStorage.removeItem("bodyBgRGB");
    localStorage.removeItem("bodylightRGB");
    checkOptions();
    html.style.removeProperty('--body-bg-rgb');
    html.style.removeProperty('--body-bg-rgb2');
    html.style.removeProperty("--light-rgb");
    html.style.removeProperty("--form-control-bg");
    updateColors();
}

function darkFn() {
    let html = document.querySelector('html');
    html.setAttribute('data-theme-mode', 'dark');
    updateColors();
    localStorage.setItem("roxlistdarktheme", true);
    localStorage.removeItem("roxlistlighttheme");
    localStorage.removeItem("bodyBgRGB");
    localStorage.removeItem("roxlistbgColor");
    localStorage.removeItem("roxlistbgwhite");
    checkOptions();
}

function ResetAllFn() {
    // imageColorChange();
    let html = document.querySelector('html');
    checkOptions();

    // clearing localstorage
    localStorage.clear();

    // reseting to light
    
    // clear primary & bg color
    html.style.removeProperty(`--primary-rgb`);
    html.style.removeProperty(`--body-bg-rgb`);

    // reseting to ltr
    ltrFn();
    
    lightFn();
    // reseting chart colors
    updateColors();

    // resetting theme primary
    document.querySelector("#switcher-primary1").checked = false
    document.querySelector("#switcher-primary2").checked = false
    document.querySelector("#switcher-primary3").checked = false
    document.querySelector("#switcher-primary4").checked = false
    document.querySelector("#switcher-primary5").checked = false

    // resetting theme background
    document.querySelector("#switcher-background1").checked = false
    document.querySelector("#switcher-background2").checked = false
    document.querySelector("#switcher-background3").checked = false
    document.querySelector("#switcher-background4").checked = false
    document.querySelector("#switcher-background5").checked = false

}

function checkOptions() {

    // dark
    if (localStorage.getItem('roxlistdarktheme')) {
        if(document.querySelector('#switcher-dark-theme')){
            document.querySelector('#switcher-dark-theme').checked = true;
        }
    }

    //RTL
    if (localStorage.getItem('roxlistrtl')) {
        if(document.querySelector('#switcher-rtl')){
            document.querySelector('#switcher-rtl').checked = true;
        }
    }
}

// chart colors
let myVarVal,primaryRGB
function updateColors() {
    'use strict'
    primaryRGB = getComputedStyle(document.documentElement).getPropertyValue('--primary-rgb').trim();

    //get variable
    myVarVal = localStorage.getItem("primaryRGB") || primaryRGB;
}
updateColors();

function localStorageBackup() {
    // if there is a value stored, update color picker and background color
    // Used to retrive the data from local storage
    if (localStorage.primaryRGB) {
        if (document.querySelector('.theme-container-primary')) {
            document.querySelector('.theme-container-primary').value = localStorage.primaryRGB;
        }
        document.querySelector('html').style.setProperty('--primary-rgb', localStorage.primaryRGB);
        if (localStorage.primaryRGB == "21, 62, 172") {
            document.querySelector("#switcher-primary1").checked = true
        }
        if (localStorage.primaryRGB == "85, 27, 130") {
            document.querySelector("#switcher-primary2").checked = true
        }
        if (localStorage.primaryRGB == "135, 90, 25") {
            document.querySelector("#switcher-primary3").checked = true
        }
        if (localStorage.primaryRGB == "25, 135, 59") {
            document.querySelector("#switcher-primary4").checked = true
        }
        if (localStorage.primaryRGB == "135, 20, 99") {
            document.querySelector("#switcher-primary5").checked = true
        }
    }
    if (localStorage.bodyBgRGB && localStorage.bodylightRGB) {
        if (document.querySelector('.theme-container-background')) {
            document.querySelector('.theme-container-background').value = localStorage.bodyBgRGB;
        }
        document.querySelector('html').style.setProperty('--body-bg-rgb', localStorage.bodyBgRGB);
        let html = document.querySelector('html');
        html.setAttribute('data-theme-mode', 'dark');
        document.querySelector('#switcher-dark-theme').checked = true;
        if (localStorage.bodyBgRGB == "12, 39, 95") {
            document.querySelector("#switcher-background1").checked = true
        }
        if (localStorage.bodyBgRGB == "46, 10, 74") {
            document.querySelector("#switcher-background2").checked = true
        }
        if (localStorage.bodyBgRGB == "106, 65, 4") {
            document.querySelector("#switcher-background3").checked = true
        }
        if (localStorage.bodyBgRGB == "7, 86, 35") {
            document.querySelector("#switcher-background4").checked = true
        }
        if (localStorage.bodyBgRGB == "86, 7, 58") {
            document.querySelector("#switcher-background5").checked = true
        }
    }
    if (localStorage.roxlistdarktheme) {
        let html = document.querySelector('html');
        html.setAttribute('data-theme-mode', 'dark');
    }
    if (localStorage.roxlistrtl) {
        let html = document.querySelector('html');
        html.setAttribute('dir', 'rtl');
        rtlFn();
    }
    if(localStorage.roxlistltr){
        ltrFn()
    }
}