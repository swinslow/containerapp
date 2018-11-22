import React, { Component } from 'react';
import axios from 'axios';
import 'semantic-ui-css/semantic.min.css';
import { Button } from 'semantic-ui-react';

import HistoryTable from './HistoryTable';
import PathInput from './PathInput';
import './App.css';

const APIROOT = 'http://localhost:3005'

class App extends Component {
  constructor(props) {
    super(props);

    this.state = {
      ready: false,
      history: [],
      pathInputContents: ""
    }
  }

  componentDidMount() {
    // set bindings
    this.handleRefresh = this.handleRefresh.bind(this)
    this.handlePathInputChange = this.handlePathInputChange.bind(this)
    this.handlePathInputSubmit = this.handlePathInputSubmit.bind(this)

    // load history for the first time
    this.refreshHistory()
  }

  handleRefresh = (e) => {
    this.refreshHistory();
  }

  handlePathInputChange = (e) => {
    this.setState({pathInputContents: e.target.value});
  }

  handlePathInputSubmit = (e) => {
    e.preventDefault();
    // call to retrieve JSON and update state
    const requestedEndpoint = APIROOT + '/' + this.state.pathInputContents;
    axios.get(requestedEndpoint)
      .then(res => {
        this.setState({pathInputContents: ""});
        this.refreshHistory()
      })
      .catch(err => {
        // const errorFlag = true;
        // const errorMsg = "Couldn't load list of projects";
        // this.setState({errorFlag, errorMsg});
      });
  }

  refreshHistory() {
    // call to retrieve JSON and update state
    const historyEndpoint = APIROOT + '/admin/history';
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
          <div>
            <PathInput pathInputValue={this.state.pathInputContents}
                       onChange={this.handlePathInputChange}
                       onSubmit={this.handlePathInputSubmit}
            />
          </div>
          <div>History <Button onClick={this.handleRefresh} icon='refresh' /></div>
          <HistoryTable ready={this.state.ready} history={this.state.history} />
        </header>
      </div>
    );
  }
}

export default App;
