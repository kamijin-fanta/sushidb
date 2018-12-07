import React from "react";

import { NavLink } from "react-router-dom";

import "./Header.css";

export function Header() {
  return (
    <header className="module-header">
      <div className="title">
        <NavLink to="/">SushiDB</NavLink>
      </div>
      <div className="links">
        <NavLink to="/" exact>SushiDB</NavLink>
        <NavLink to="/cluster/store">Cluster</NavLink>
        <NavLink to="/keys">View</NavLink>
        <a href="#">Write</a>
      </div>
    </header>
  );
}
