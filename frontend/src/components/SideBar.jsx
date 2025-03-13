import React from "react";
import {Link, useLocation} from "react-router-dom";
import TopBar from "./TopBar.jsx";

const SideBar = () => {
    const location = useLocation();

    return (
        <aside className="app-sidebar sticky" id="sidebar">

            <div className="main-sidebar-header">
                <Link to="/" className="h4">Jobs Dashboard</Link>
                <Link to="/" className="h5" hidden>Jobs</Link>
            </div>

            <div className="main-sidebar" id="sidebar-scroll">

                <nav className="main-menu-container nav nav-pills flex-column sub-open">
                    <ul className="main-menu">

                        <li className="slide">
                            <Link to="/channels" className="side-menu__item">
                                <i className="fe fe-settings side-menu__icon"></i>
                                <span className="side-menu__label">Channels</span>
                            </Link>
                        </li>

                        <li className="slide">
                            <Link to="/imports" className="side-menu__item">
                                <i className="fe fe-layers side-menu__icon"></i>
                                <span className="side-menu__label">Imports</span>
                            </Link>
                        </li>
                    </ul>
                </nav>

            </div>

        </aside>
    );
}

export default SideBar;
