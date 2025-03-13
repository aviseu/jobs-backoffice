import React from "react";
import {Link, useLocation} from "react-router-dom";

const TopBar = () => {

    const toggleTheme = () => {
        let html = document.querySelector("html");
        if (html.getAttribute("data-theme-mode") === "dark") {
            html.setAttribute("data-theme-mode", "light");
            html.removeAttribute("data-bg-theme");
            if (!localStorage.getItem("primaryRGB")) {
                html.setAttribute("style", "");
            }
            document
                .querySelector("html")
                .style.removeProperty("--body-bg-rgb", localStorage.bodyBgRGB);
            html.style.removeProperty("--body-bg-rgb2");
            html.style.removeProperty("--light-rgb");
            html.style.removeProperty("--form-control-bg");
            localStorage.removeItem("roxlistdarktheme");
            localStorage.removeItem("bodylightRGB");
            localStorage.removeItem("bodyBgRGB");
        } else {
            html.setAttribute("data-theme-mode", "dark");
            if (!localStorage.getItem("primaryRGB")) {
                html.setAttribute("style", "");
            }
            localStorage.setItem("roxlistdarktheme", "true");
            localStorage.removeItem("bodylightRGB");
            localStorage.removeItem("bodyBgRGB");
        }
    };

    const toggleSidebar = () => {
        let html = document.querySelector("html")
        if (html.getAttribute("data-toggled") === "close") {
            html.setAttribute("data-vertical-style", "overlay");
            html.setAttribute("data-toggled", "icon-overlay-close");
        } else {
            html.removeAttribute("data-vertical-style");
            html.setAttribute("data-toggled", "close");
        }
    }

    return (
        <header className="app-header">
            <div className="main-header-container container-fluid">

                <div className="header-content-left">
                    <div className="header-element">
                        <div className="horizontal-logo">
                            <Link className="header-logo" to="/">
                                <h1>Jobs Dashboard</h1>
                            </Link>
                        </div>
                    </div>

                    <div className="header-element">
                        <button aria-label="Hide Sidebar"
                           className="sidemenu-toggle header-link animated-arrow hor-toggle horizontal-navtoggle text-decoration-none btn btn-link"
                           data-bs-toggle="sidebar" onClick={(event) => toggleSidebar()}>
                            <i className="header-link-icon fe fe-menu"></i>
                        </button>
                    </div>

                </div>


                <div className="header-content-right">

                    <div className="header-element">
                        <a href="javascript:void(0);" className="header-link text-decoration-none" data-bs-toggle="modal"
                           data-bs-target="#locationModal">
                            <i className="fe fe-globe header-link-icon"></i>
                        </a>
                    </div>

                    <div className="header-element header-theme-mode">
                        <button onClick={(event) => toggleTheme()} className=" header-link layout-setting text-decoration-none btn btn-link">
                            <i className="fe fe-moon header-link-icon light-layout">
                            </i>

                            <i className="fe fe-sun header-link-icon dark-layout"></i>
                        </button>
                    </div>
                </div>

            </div>
        </header>
    );
};

export default TopBar;