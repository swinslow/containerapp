import React, { Component } from 'react';
import { BrowserRouter as Router, Route, Switch } from "react-router-dom";
import axios from 'axios';
import 'semantic-ui-css/semantic.min.css';

import LoginPane from './LoginPane.js';
import PathInput from './PathInput';
import TokenManager from './TokenManager';
import AdminPane from './AdminPane';
import './App.css';

const APIROOT = 'http://localhost:3005'

class App extends Component {
  constructor(props) {
    super(props);

    this.state = {
      history: null,
      emailInputContents: "",
      pathInputContents: "",
      lastPathResponse: null,
      jwtToken: null,
      myself: {
        isKnownUser: false,
        name: "",
        email: "",
        isAdmin: false
      }
    }

    // set bindings
    this.handleRefresh = this.handleRefresh.bind(this)
    this.handlePathInputChange = this.handlePathInputChange.bind(this)
    this.handlePathInputSubmit = this.handlePathInputSubmit.bind(this)
    this.handleLogout = this.handleLogout.bind(this)
    this.isLoggedIn = this.isLoggedIn.bind(this)
    this.setToken = this.setToken.bind(this)
    this.resetMyself = this.resetMyself.bind(this)
    this.setMyself = this.setMyself.bind(this)
    this.refreshHistory = this.refreshHistory.bind(this)
  }

  componentDidMount() {
    // create token manager and fetch token
    this.tokenManager = new TokenManager(APIROOT, this.setToken, this.setMyself)

    // load token from local storage if it's already there
    const checkToken = localStorage.getItem("apitoken")
    if (checkToken) {
      this.setToken(checkToken)
    }
  }

  isLoggedIn() {
    if (this.state === undefined) {
      return false;
    }
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
    const config = {
      headers: {
          "Authorization": "Bearer " + this.state.jwtToken,
          "Content-Type": "application/json"
      }
    }
    axios.get(requestedEndpoint, config)
      .then(res => {
        this.setState({
          lastPathResponse: res.data,
          pathInputContents: ""
        });
      })
      .catch(err => {
        // const errorFlag = true;
        // const errorMsg = "Couldn't load list of projects";
        // this.setState({errorFlag, errorMsg});
      });
  }

  setToken = (token) => {
    this.setState({jwtToken: token});
    if (token === null) {
      localStorage.removeItem("apitoken");
    } else {
      localStorage.setItem("apitoken", token);
      this.tokenManager.fetchLoginInfo(token);
    }
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
    const config = {
      headers: {
          "Authorization": "Bearer " + this.state.jwtToken,
          "Content-Type": "application/json"
      }
    }
    axios.get(historyEndpoint, config)
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

  handleLogout = (e) => {
    e.preventDefault();
    this.setState({lastPathResponse: null});
    this.resetMyself();
    this.setToken(null);
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
                    <PathInput myself={this.state.myself}
                              pathInputValue={this.state.pathInputContents}
                              lastPathResponse={this.state.lastPathResponse}
                              isLoggedIn={this.isLoggedIn}
                              onChange={this.handlePathInputChange}
                              onSubmit={this.handlePathInputSubmit}
                              onLogout={this.handleLogout}
                    />
                  </div>
                </header>
              </div>
            </Route>
            <Route path="/admin">
              <AdminPane onRefresh={this.refreshHistory}
                         isLoggedIn={this.isLoggedIn}
                         history={this.state.history} />
            </Route>
          </Switch>
        </div>
      </Router>
    );
  }
}

export default App;
