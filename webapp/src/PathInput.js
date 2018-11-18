import React from 'react'
import { Form } from 'semantic-ui-react'

function PathInput(props) {
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
