import React from 'react'
import { Redirect } from "react-router-dom";
import { Form, Label } from 'semantic-ui-react'

function LoginPane(props) {
    if (props.isLoggedIn() === true) {
        return <Redirect to='/app' />
    }

    return (
        <div className="App">
            <header className="App-header">
                <div>
                    <Form onSubmit={props.onSubmit}>
                        <Form.Field>
                            <Label>Log in with Email</Label>
                            <Form.Input type="text"
                                        value={props.emailInputValue}
                                        onChange={props.onChange} />
                            <Form.Button icon="arrow right" />
                        </Form.Field>
                    </Form>
                </div>
            </header>
        </div>
    );
}

export default LoginPane;
