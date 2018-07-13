import React from 'react'
import TableRow from './TableRow'
import './Table.css'

export default (props) => 
    <table>
        <thead>
            <tr>
                <th>Name</th>
                <th>Email</th>
                <th>Description</th>
                <th>Ações</th>
            </tr>
        </thead>
        <tbody>
            {
                props.lista.map( (elem, index) => {
                    return ( 
                        <TableRow key={index} 
                        {...elem} onDelete={props.onDelete}
                        onUpdate={props.onUpdate} index={index}/>
                    )
                })
            }
        </tbody>
    </table>