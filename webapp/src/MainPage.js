import React, { Component } from 'react'
import { Form, Label } from 'semantic-ui-react'

import PathInput from './PathInput'

class MainPage extends Component {

    render() {
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
}

export default MainPage;
