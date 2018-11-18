import React, { Component } from 'react';
import axios from 'axios';
import 'semantic-ui-css/semantic.min.css';
import { Button } from 'semantic-ui-react';

import HistoryTable from './HistoryTable';
import './App.css';

const APIROOT = 'http://localhost:3005'

class App extends Component {
  constructor(props) {
    super(props);

    this.state = {
      ready: false,
      history: []
    }
  }

  componentDidMount() {
    // set bindings
    this.handleRefresh = this.handleRefresh.bind(this)

    // load history for the first time
    this.refreshHistory()
  }

  handleRefresh = (e) => {
    this.refreshHistory();
  }

  refreshHistory() {
    // call to retrieve JSON and update state
    const historyEndpoint = APIROOT + '/history';
    axios.get(historyEndpoint)
      .then(res => {
        const history = res.data;
        const ready = true;
        this.setState({ready, history});
      })
      .catch(err => {
        // const errorFlag = true;
        // const errorMsg = "Couldn't load list of projects";
        // this.setState({errorFlag, errorMsg});
        const ready = false;
        this.setState({ready});
      });
  }

  render() {
    return (
      <div className="App">
        <header className="App-header">
          <div>History <Button onClick={this.handleRefresh} icon='refresh' /></div>
          <HistoryTable ready={this.state.ready} history={this.state.history} />
        </header>
      </div>
    );
  }
}

export default App;
