import React, { Component } from 'react'
import { Link, Redirect } from "react-router-dom";
import { Button } from 'semantic-ui-react'
import HistoryTable from './HistoryTable';
import UsersTable from './UsersTable';

class AdminPane extends Component {
    componentDidMount() {
        if (this.props.onRefresh) {
            this.props.onRefresh()
        }
    }

    render() {
        if (this.props.isLoggedIn() !== true) {
            return <Redirect to='/' />
        }

        return (
            <div className="App">
                <header className="App-header">
                    <h4><Link to="/app">Back to app</Link></h4>
                    <div>History <Button onClick={this.props.onRefresh} icon='refresh' /></div>
                    <HistoryTable history={this.props.history} />
                    <br />
                    <div>Users <Button onClick={this.props.onRefresh} icon='refresh' /></div>
                    <UsersTable users={this.props.users} />
                </header>
            </div>

        );
    }
}

export default AdminPane;
