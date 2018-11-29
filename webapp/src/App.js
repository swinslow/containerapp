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
      users: null,
      emailInputContents: "",
      pathInputContents: "",
      newUserName: "",
      newUserEmail: "",
      newUserIsAdmin: false,
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
    this.handlePathInputChange = this.handlePathInputChange.bind(this)
    this.handlePathInputSubmit = this.handlePathInputSubmit.bind(this)
    this.handleNewUserChange = this.handleNewUserChange.bind(this)
    this.handleNewUserSubmit = this.handleNewUserSubmit.bind(this)
    this.handleLogout = this.handleLogout.bind(this)
    this.isLoggedIn = this.isLoggedIn.bind(this)
    this.setToken = this.setToken.bind(this)
    this.resetMyself = this.resetMyself.bind(this)
    this.setMyself = this.setMyself.bind(this)
    this.refreshHistory = this.refreshHistory.bind(this)
    this.refreshUsers = this.refreshUsers.bind(this)
    this.refreshAdmin = this.refreshAdmin.bind(this)
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
        // const errorMsg = "Couldn't get path from API";
        // this.setState({errorFlag, errorMsg});
      });
  }

  handleNewUserChange = (e) => {
    switch(e.target.name) {
      case 'name':
        this.setState({newUserName: e.target.value});
        break;
      case 'email':
        this.setState({newUserEmail: e.target.value});
        break;
      case 'isAdmin':
        this.setState({newUserIsAdmin: !this.state.newUserIsAdmin});
        break;
      default:
        console.log("unknown element name: " + e.target.name);
        break;
    }
  }

  handleNewUserSubmit = (e) => {
    e.preventDefault();
    // call to retrieve JSON and update state
    const requestedEndpoint = APIROOT + '/admin/users';
    const data = {
      name: this.state.newUserName,
      email: this.state.newUserEmail,
      is_admin: this.state.newUserIsAdmin
    }
    const config = {
      headers: {
          'X-Requested-With': 'XMLHttpRequest',
          "Authorization": "Bearer " + this.state.jwtToken,
          "Content-Type": "application/json"
      }
    }
    axios.post(requestedEndpoint, data, config)
      .then(res => {
        this.setState({
          newUserName: "",
          newUserEmail: "",
          newUserIsAdmin: false
        });
        this.refreshUsers();
      })
      .catch(err => {
        // const errorFlag = true;
        // const errorMsg = "Couldn't post new user";
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

  refreshUsers() {
    // call to retrieve JSON and update state
    const usersEndpoint = APIROOT + '/admin/users';
    const config = {
      headers: {
          "Authorization": "Bearer " + this.state.jwtToken,
          "Content-Type": "application/json"
      }
    }
    axios.get(usersEndpoint, config)
      .then(res => {
        const users = res.data;
        this.setState({users});
      })
      .catch(err => {
        const users = null;
        this.setState({users});
      });
  }

  refreshAdmin() {
    this.refreshHistory()
    this.refreshUsers()
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
              <AdminPane onRefresh={this.refreshAdmin}
                         isLoggedIn={this.isLoggedIn}
                         history={this.state.history}
                         users={this.state.users}
                         newUserName={this.state.newUserName}
                         newUserEmail={this.state.newUserEmail}
                         newUserIsAdmin={this.state.newUserIsAdmin}
                         onNewUserChange={this.handleNewUserChange}
                         onNewUserSubmit={this.handleNewUserSubmit} />
            </Route>
          </Switch>
        </div>
      </Router>
    );
  }
}

export default App;
