import React from 'react';
import logo from './logo.svg';
import './App.css';
import { BrowserRouter as Router } from "react-router-dom";
import {
  Routes,
  Route,
  Navigate
} from "react-router";

const X = () => {
  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <p>
          Edit <code>src/App.tsx</code> and save to reload.
        </p>
        <a
          className="App-link"
          href="https://reactjs.org"
          target="_blank"
          rel="noopener noreferrer"
        >
          Learn React
        </a>
      </header>
    </div >
  )
}

const TODO = () => {
  return (
    <div>TODO</div>
  )
}

const App = () => {
  return (
    <Router>
      <Routes>
        <Route path="/login" element={<X />} />
        <Route path="/sign-up" element={<TODO />} />
        <Route path="/chat" element={<TODO />} />
        <Route path="*" element={<Navigate to="/login" />} />
      </Routes>
    </Router >
  );
}

export default App;
