import React, { Component } from 'react'
import { Link, Redirect } from "react-router-dom";
import { Button, Form, Label } from 'semantic-ui-react'
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
                    <div>
                        <Form onSubmit={this.props.onNewUserSubmit}>
                            <Form.Field>
                                <Label>Create New User</Label>
                                <Form.Input type="text" name="name"
                                            value={this.props.newUserName}
                                            onChange={this.props.onNewUserChange} />
                                <Form.Input type="text" name="email"
                                            value={this.props.newUserEmail}
                                            onChange={this.props.onNewUserChange} />
                                <Form.Input type="checkbox" name="isAdmin"
                                            label="Is user an admin?"
                                            checked={this.props.newUserIsAdmin}
                                            onChange={this.props.onNewUserChange} />
                                <Form.Button icon="arrow right" />
                            </Form.Field>
                        </Form>
                    </div>
                </header>
            </div>

        );
    }
}

export default AdminPane;
