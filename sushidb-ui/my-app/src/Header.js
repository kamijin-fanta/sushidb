import React from "react";

import { Link } from "react-router-dom";

import "./Header.css";

export function Header() {
  return (
    <div className="header">
      <div className="title">
        <Link to="/">SushiDB</Link>
      </div>
      <div className="links">
        <div>
          <Link to="/">SushiDB</Link>
        </div>
        <div>
          <a href="#">View</a>
        </div>
        <div>
          <a href="#">Write</a>
        </div>
      </div>
    </div>
  );
}
