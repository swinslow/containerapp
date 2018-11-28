import React from 'react'
import { Table } from 'semantic-ui-react'

function createUserRows(users) {
    var rows = [];
    var i = 0;

    users.forEach(function(u) {
        const isAdmin = (u.is_admin) ? "ADMIN" : "" 
        const row = (
            <Table.Row key={i}>
                <Table.Cell>{u.id}</Table.Cell>
                <Table.Cell>{u.email}</Table.Cell>
                <Table.Cell>{u.name}</Table.Cell>
                <Table.Cell>{isAdmin}</Table.Cell>
            </Table.Row>
        );
        rows.push(row);
        i = i + 1;
    })

    return rows;
}

function UsersTable(props) {
    if (props.users === null) {
        return (
            <div>Loading...</div>
        );
    }

    const userRows = createUserRows(props.users);
    return (
        <div>
            <Table celled>
                <Table.Header>
                    <Table.Row>
                        <Table.HeaderCell>ID</Table.HeaderCell>
                        <Table.HeaderCell>Email</Table.HeaderCell>
                        <Table.HeaderCell>Name</Table.HeaderCell>
                        <Table.HeaderCell>Admin</Table.HeaderCell>
                    </Table.Row>
                </Table.Header>

                <Table.Body>
                    {userRows}
                </Table.Body>
            </Table>
        </div>
    )
}

export default UsersTable;
