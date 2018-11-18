import React from 'react'
import { Table } from 'semantic-ui-react'

function createHistoryRows(history) {
    var rows = [];
    var i = 0;

    history.forEach(function(h) {
        const row = (
            <Table.Row key={i}>
                <Table.Cell>{h.path}</Table.Cell>
                <Table.Cell>{h.date}</Table.Cell>
                <Table.Cell>{h.user_id}</Table.Cell>
            </Table.Row>
        );
        rows.push(row);
        i = i + 1;
    })

    return rows;
}

function HistoryTable(props) {
    const ready = props.ready;
    if (!ready) {
        return (
            <div>Loading...</div>
        );
    }

    const historyRows = createHistoryRows(props.history);
    return (
        <div>
            <Table celled>
                <Table.Header>
                    <Table.Row>
                        <Table.HeaderCell>Path</Table.HeaderCell>
                        <Table.HeaderCell>Date</Table.HeaderCell>
                        <Table.HeaderCell>User ID</Table.HeaderCell>
                    </Table.Row>
                </Table.Header>

                <Table.Body>
                    {historyRows}
                </Table.Body>
            </Table>
        </div>
    )
}

export default HistoryTable;
