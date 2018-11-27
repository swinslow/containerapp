import React from 'react'
import { Link, Redirect } from "react-router-dom";
import { Form } from 'semantic-ui-react'

function PathInput(props) {
    if (props.isLoggedIn() === false) {
        return <Redirect to='/' />
    }

    if (props.myself !== null && props.email !== "" && !props.myself.isKnownUser) {
        return <div>
            <h2>Unknown user {props.myself.email}</h2>
            <button onClick={props.onLogout}>Go to Login page</button>
        </div>
    }

    const adminJSX = (props.myself.isAdmin) ?
        <span> (<Link to="/admin">ADMIN</Link>)</span> :
        <span> </span>;
    const lastPathString = (props.lastPathResponse !== null) ?
        "(last response: " + JSON.stringify(props.lastPathResponse) + ")" :
        "";

    return (
        <div>
            <h4>Logged in as {props.myself.email} {adminJSX}</h4>
            <div>
                <Form onSubmit={props.onSubmit}>
                    <Form.Group>
                        <Form.Input type="text"
                                    value={props.pathInputValue}
                                    onChange={props.onChange} />
                        <Form.Button icon="arrow right" />
                    </Form.Group>
                </Form>
            </div>
            <div>{lastPathString}</div>
            <button onClick={props.onLogout}>Logout</button>
        </div>
    )
}

export default PathInput;
