import React, { Component } from 'react';
import { BrowserRouter as Router, Route, Switch } from "react-router-dom";
import axios from 'axios';
import 'semantic-ui-css/semantic.min.css';

import LoginPane from './LoginPane.js';
import PathInput from './PathInput';
import TokenManager from './TokenManager';
import FetchManager from './FetchManager';
import './App.css';

const APIROOT = 'http://localhost:3005'

class App extends Component {
  constructor(props) {
    super(props);

    this.state = {
      history: null,
      emailInputContents: "",
      pathInputContents: "",
      jwtToken: null,
      myself: {
        isKnownUser: false,
        name: "",
        email: "",
        isAdmin: false
      }
    }
  }

  componentDidMount() {
    // set bindings
    this.handleRefresh = this.handleRefresh.bind(this)
    this.handlePathInputChange = this.handlePathInputChange.bind(this)
    this.handlePathInputSubmit = this.handlePathInputSubmit.bind(this)
    this.isLoggedIn = this.isLoggedIn.bind(this)
    this.setToken = this.setToken.bind(this)
    this.resetMyself = this.resetMyself.bind(this)
    this.setMyself = this.setMyself.bind(this)

    // load history for the first time
    this.refreshHistory()

    // create token manager and fetch token
    this.tokenManager = new TokenManager(APIROOT, this.setToken, this.setMyself)

    // create fetch manager
    this.fetchManager = new FetchManager()
  }

  isLoggedIn() {
    if (this.state === undefined) {
      return false;
    }
    //console.log("this.state.jwtToken = " + this.state.jwtToken);
    return this.state.jwtToken !== null && this.state.jwtToken !== undefined
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

  setToken = (token) => {
    this.setState({jwtToken: token});
    this.tokenManager.fetchLoginInfo(this.state.jwtToken);
  }

  resetMyself = () => {
    this.setState({myself: {
      isKnownUser: false,
      name: "",
      email: "",
      isAdmin: false
    }});
  }

  setMyself = (userdata) => {
    if (userdata == null) {
      this.resetMyself();
    } else {
      this.setState({myself: userdata})
    }
  }

  refreshHistory() {
    // call to retrieve JSON and update state
    const historyEndpoint = APIROOT + '/admin/history';
    axios.get(historyEndpoint)
      .then(res => {
        const history = res.data;
        this.setState({history});
      })
      .catch(err => {
        const history = null;
        this.setState({history});
      });
  }

  handleEmailInputChange = (e) => {
    this.setState({emailInputContents: e.target.value});
  }

  handleEmailInputSubmit = (e) => {
    this.tokenManager.fetchToken(this.state.emailInputContents);
    this.setState({emailInputContents: ""});
  }

  render() {
    return (
      <Router>
        <div>
          <Switch>
            <Route exact path="/" render={(props) => <LoginPane {...props}
                  isLoggedIn={this.isLoggedIn}
                  emailInputValue={this.state.emailInputContents}
                  onChange={this.handleEmailInputChange}
                  onSubmit={this.handleEmailInputSubmit}
                />
              }/>
            <Route path="/app">
              <div className="App">
                <header className="App-header">
                  <div>
                    <PathInput pathInputValue={this.state.pathInputContents}
                              isLoggedIn={this.isLoggedIn}
                              onChange={this.handlePathInputChange}
                              onSubmit={this.handlePathInputSubmit}
                    />
                  </div>
                </header>
              </div>
            </Route>
          </Switch>
        </div>
      </Router>
    );
  }
}

export default App;
