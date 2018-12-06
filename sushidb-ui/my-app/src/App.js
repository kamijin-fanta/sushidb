import React, { Component } from "react";

import { Route } from "react-router";
import { BrowserRouter } from "react-router-dom";
import { Header } from "./Header";
import { Home } from "./pages/Home";
import "./App.css";

class App extends Component {
  render() {
    return (
      <BrowserRouter>
        <div className="app">
          <Header />
          <Route path="/" component={Home} />
        </div>
      </BrowserRouter>
    );
  }
}

export default App;
