import React, { Component } from 'react';
import ButtonExampleEmphasis from './ButtonExampleEmphasis';
import HistoryTable from './HistoryTable';
import logo from './logo.svg';
import 'semantic-ui-css/semantic.min.css';
import './App.css';

class App extends Component {
  constructor(props) {
    super(props);

    this.state = {
      ready: true,
      history: [
        {"path": "/abc", "date": "2018-11-17T21:55:00Z"},
        {"path": "/def", "date": "2018-11-16T15:22:00Z"},
        {"path": "/ghi", "date": "2018-11-15T00:00:00Z"}
      ]
    }
  }

  render() {
    return (
      <div className="App">
        <header className="App-header">
          <img src={logo} className="App-logo" alt="logo" />
          <h1>Hello!</h1>
          <HistoryTable ready={this.state.ready} history={this.state.history} />
          <a
            className="App-link"
            href="https://reactjs.org"
            target="_blank"
            rel="noopener noreferrer"
          >
            Learn React
          </a><br />
          <ButtonExampleEmphasis />
        </header>
      </div>
    );
  }
}

export default App;
