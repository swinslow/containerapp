import React from 'react'
import { Redirect } from "react-router-dom";
import { Form } from 'semantic-ui-react'

function PathInput(props) {
    if (props.isLoggedIn() === false) {
        return <Redirect to='/' />
    }

    return (
        <Form onSubmit={props.onSubmit}>
            <Form.Group>
                <Form.Input type="text"
                            value={props.pathInputValue}
                            onChange={props.onChange} />
                <Form.Button icon="arrow right" />
            </Form.Group>
        </Form>
    )
}

export default PathInput;
