import React, { Component } from "react";

import { Route } from "react-router";
import { BrowserRouter } from "react-router-dom";
import { Header } from "./Header";
import "./App.css";

import { Home } from "./pages/Home";
import { Keys } from "./pages/Keys";
import { SingleMetric } from "./pages/SingleMetric";
import { MessageMetric } from "./pages/MessageMetric";

class App extends Component {
  render() {
    return (
      <BrowserRouter basename={process.env.PUBLIC_URL || "/"}>
        <div className="app">
          <Header />
          <Route path="/" component={Home} exact />
          <Route path="/keys" component={Keys} exact />
          <Route path="/metric/single/:key" component={SingleMetric} exact />
          <Route path="/metric/message/:key" component={MessageMetric} exact />
        </div>
      </BrowserRouter>
    );
  }
}

export default App;
